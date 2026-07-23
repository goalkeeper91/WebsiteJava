package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"demo/backend-go/internal/domain"
)

type SubathonRepository struct {
	db *sql.DB
}

func NewSubathonRepository(db *sql.DB) *SubathonRepository {
	return &SubathonRepository{db: db}
}

const subathonStateColumns = `
	user_id, time_remaining, is_running, target_timestamp,
	total_subs, total_bits, total_events,
	initial_time_minutes, sub_time_seconds, bits_time_seconds,
	event_log, created_at, updated_at
`

func scanSubathonState(scanner interface {
	Scan(dest ...interface{}) error
}) (*domain.SubathonState, error) {
	var s domain.SubathonState
	err := scanner.Scan(
		&s.UserID, &s.TimeRemaining, &s.IsRunning, &s.TargetTimestamp,
		&s.TotalSubs, &s.TotalBits, &s.TotalEvents,
		&s.InitialTimeMinutes, &s.SubTimeSeconds, &s.BitsTimeSeconds,
		&s.EventLog, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// GetByUser returns nil (no error) if the user has no subathon state yet.
func (r *SubathonRepository) GetByUser(ctx context.Context, userID string) (*domain.SubathonState, error) {
	query := fmt.Sprintf(`SELECT %s FROM subathon_timer_state WHERE user_id = $1`, subathonStateColumns)

	state, err := scanSubathonState(r.db.QueryRowContext(ctx, query, userID))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get subathon state: %w", err)
	}
	return state, nil
}

// GetOrCreate returns the user's existing state, or creates a fresh default
// row on first use (mirrors the original app's default: 30min/120s/30s).
func (r *SubathonRepository) GetOrCreate(ctx context.Context, userID string) (*domain.SubathonState, error) {
	existing, err := r.GetByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}

	query := fmt.Sprintf(`
		INSERT INTO subathon_timer_state (user_id)
		VALUES ($1)
		ON CONFLICT (user_id) DO UPDATE SET user_id = EXCLUDED.user_id
		RETURNING %s
	`, subathonStateColumns)

	state, err := scanSubathonState(r.db.QueryRowContext(ctx, query, userID))
	if err != nil {
		return nil, fmt.Errorf("failed to create subathon state: %w", err)
	}
	return state, nil
}

// Update applies a partial update and returns the resulting row.
func (r *SubathonRepository) Update(ctx context.Context, userID string, input domain.SubathonStateUpdateInput) (*domain.SubathonState, error) {
	query := fmt.Sprintf(`
		UPDATE subathon_timer_state
		SET time_remaining = COALESCE($2, time_remaining),
		    is_running = COALESCE($3, is_running),
		    target_timestamp = $4,
		    total_subs = COALESCE($5, total_subs),
		    total_bits = COALESCE($6, total_bits),
		    total_events = COALESCE($7, total_events),
		    initial_time_minutes = COALESCE($8, initial_time_minutes),
		    sub_time_seconds = COALESCE($9, sub_time_seconds),
		    bits_time_seconds = COALESCE($10, bits_time_seconds),
		    event_log = COALESCE($11, event_log),
		    updated_at = NOW()
		WHERE user_id = $1
		RETURNING %s
	`, subathonStateColumns)

	var eventLogArg interface{}
	if input.EventLog != nil {
		eventLogArg = input.EventLog
	}

	state, err := scanSubathonState(r.db.QueryRowContext(
		ctx, query,
		userID,
		input.TimeRemaining,
		input.IsRunning,
		input.TargetTimestamp,
		input.TotalSubs,
		input.TotalBits,
		input.TotalEvents,
		input.InitialTimeMinutes,
		input.SubTimeSeconds,
		input.BitsTimeSeconds,
		eventLogArg,
	))
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("subathon state not found for user %s", userID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update subathon state: %w", err)
	}
	return state, nil
}

// GetAllUserIDs lists every user with a subathon state row - the set the
// background EventSub manager keeps webhook subscriptions active for.
func (r *SubathonRepository) GetAllUserIDs(ctx context.Context) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT user_id FROM subathon_timer_state`)
	if err != nil {
		return nil, fmt.Errorf("failed to list subathon users: %w", err)
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan user id: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// MarkEventProcessed records a Twitch EventSub message ID as handled.
// Returns false if it was already recorded (a duplicate delivery - Twitch
// only guarantees at-least-once, so this is the dedup guard against
// double-counting a sub/cheer).
func (r *SubathonRepository) MarkEventProcessed(ctx context.Context, messageID, userID string) (bool, error) {
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO subathon_processed_events (message_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT (message_id) DO NOTHING
	`, messageID, userID)
	if err != nil {
		return false, fmt.Errorf("failed to mark event processed: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to get rows affected: %w", err)
	}
	return rows > 0, nil
}

// RecordFailedEvent persists an event that failed to process so it stays
// visible and can be retried, instead of a silent log line.
func (r *SubathonRepository) RecordFailedEvent(ctx context.Context, userID, messageID, eventType string, rawPayload []byte, errMsg string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO subathon_failed_events (user_id, message_id, event_type, raw_payload, error_message)
		VALUES ($1, $2, $3, $4, $5)
	`, userID, messageID, eventType, rawPayload, errMsg)
	if err != nil {
		return fmt.Errorf("failed to record failed event: %w", err)
	}
	return nil
}

// GetUnresolvedFailedEvents lists failed events still worth retrying (caps
// retries at 5 so a permanently-broken payload doesn't retry forever).
func (r *SubathonRepository) GetUnresolvedFailedEvents(ctx context.Context) ([]domain.SubathonFailedEvent, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, message_id, event_type, raw_payload, error_message,
		       retry_count, resolved, created_at, updated_at
		FROM subathon_failed_events
		WHERE resolved = false AND retry_count < 5
		ORDER BY created_at ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list failed events: %w", err)
	}
	defer rows.Close()

	var events []domain.SubathonFailedEvent
	for rows.Next() {
		var e domain.SubathonFailedEvent
		if err := rows.Scan(
			&e.ID, &e.UserID, &e.MessageID, &e.EventType, &e.RawPayload, &e.ErrorMessage,
			&e.RetryCount, &e.Resolved, &e.CreatedAt, &e.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan failed event: %w", err)
		}
		events = append(events, e)
	}
	return events, nil
}

func (r *SubathonRepository) MarkFailedEventResolved(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `UPDATE subathon_failed_events SET resolved = true, updated_at = NOW() WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to mark failed event resolved: %w", err)
	}
	return nil
}

func (r *SubathonRepository) IncrementFailedEventRetry(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `UPDATE subathon_failed_events SET retry_count = retry_count + 1, updated_at = NOW() WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to increment failed event retry count: %w", err)
	}
	return nil
}
