package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
)

type scheduledMessageRepository struct {
	db *sql.DB
}

func NewScheduledMessageRepository(db *sql.DB) repository.ScheduledMessageRepository {
	return &scheduledMessageRepository{db: db}
}

func (r *scheduledMessageRepository) Create(ctx context.Context, m *domain.ScheduledMessage) error {
	query := `
		INSERT INTO scheduled_messages
		(channel_id, message, command_id, interval_seconds, enabled, only_when_live, next_send_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		m.ChannelID,
		m.Message,
		m.CommandID,
		m.IntervalSeconds,
		m.Enabled,
		m.OnlyWhenLive,
		m.NextSendAt,
		m.CreatedAt,
		m.UpdatedAt,
	).Scan(&m.ID)

	if err != nil {
		return fmt.Errorf("fehler beim Erstellen der automatisierten Nachricht: %w", err)
	}

	return nil
}

func (r *scheduledMessageRepository) GetByID(ctx context.Context, id int64, channelID string) (*domain.ScheduledMessage, error) {
	query := `
		SELECT id, channel_id, message, command_id, interval_seconds, enabled, only_when_live,
		       next_send_at, last_sent_at, created_at, updated_at
		FROM scheduled_messages
		WHERE id = $1 AND channel_id = $2
	`

	m, err := scanOneScheduledMessage(r.db.QueryRowContext(ctx, query, id, channelID))
	if err == sql.ErrNoRows {
		return nil, domain.ErrScheduledMessageNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der automatisierten Nachricht: %w", err)
	}

	return m, nil
}

func (r *scheduledMessageRepository) GetAll(ctx context.Context, channelID string, limit, offset int) ([]*domain.ScheduledMessage, int64, error) {
	var total int64
	err := r.db.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM scheduled_messages WHERE channel_id = $1",
		channelID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("fehler beim Zählen der automatisierten Nachrichten: %w", err)
	}

	query := `
		SELECT id, channel_id, message, command_id, interval_seconds, enabled, only_when_live,
		       next_send_at, last_sent_at, created_at, updated_at
		FROM scheduled_messages
		WHERE channel_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, channelID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("fehler beim Laden der automatisierten Nachrichten: %w", err)
	}
	defer rows.Close()

	messages, err := scanScheduledMessages(rows)
	return messages, total, err
}

func (r *scheduledMessageRepository) Update(ctx context.Context, m *domain.ScheduledMessage) error {
	query := `
		UPDATE scheduled_messages
		SET message = $2, interval_seconds = $3, enabled = $4, only_when_live = $5, updated_at = $6
		WHERE id = $1 AND channel_id = $7
	`

	result, err := r.db.ExecContext(
		ctx, query,
		m.ID, m.Message, m.IntervalSeconds, m.Enabled, m.OnlyWhenLive, m.UpdatedAt, m.ChannelID,
	)
	if err != nil {
		return fmt.Errorf("fehler beim Update: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("fehler beim Abrufen betroffener Zeilen: %w", err)
	}
	if rowsAffected == 0 {
		return domain.ErrScheduledMessageNotFound
	}

	return nil
}

func (r *scheduledMessageRepository) Delete(ctx context.Context, id int64, channelID string) error {
	result, err := r.db.ExecContext(ctx,
		`DELETE FROM scheduled_messages WHERE id = $1 AND channel_id = $2`, id, channelID)
	if err != nil {
		return fmt.Errorf("fehler beim Löschen: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("fehler beim Abrufen betroffener Zeilen: %w", err)
	}
	if rowsAffected == 0 {
		return domain.ErrScheduledMessageNotFound
	}

	return nil
}

func (r *scheduledMessageRepository) GetDue(ctx context.Context, limit int) ([]*domain.ScheduledMessage, error) {
	query := `
		SELECT id, channel_id, message, command_id, interval_seconds, enabled, only_when_live,
		       next_send_at, last_sent_at, created_at, updated_at
		FROM scheduled_messages
		WHERE enabled = true AND next_send_at <= NOW()
		ORDER BY next_send_at ASC
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden fälliger Nachrichten: %w", err)
	}
	defer rows.Close()

	return scanScheduledMessages(rows)
}

func (r *scheduledMessageRepository) MarkSent(ctx context.Context, id int64, sentAt, nextSendAt time.Time) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE scheduled_messages
		SET last_sent_at = $2, next_send_at = $3, updated_at = NOW()
		WHERE id = $1
	`, id, sentAt, nextSendAt)
	if err != nil {
		return fmt.Errorf("fehler beim Fortschreiben des Zeitplans: %w", err)
	}
	return nil
}

// rowScanner is satisfied by both *sql.Row and *sql.Rows, letting
// scanOneScheduledMessage share the column layout with scanScheduledMessages.
type rowScanner interface {
	Scan(dest ...interface{}) error
}

func scanOneScheduledMessage(row rowScanner) (*domain.ScheduledMessage, error) {
	m := &domain.ScheduledMessage{}
	var message sql.NullString
	var commandID sql.NullInt64

	err := row.Scan(
		&m.ID, &m.ChannelID, &message, &commandID, &m.IntervalSeconds, &m.Enabled, &m.OnlyWhenLive,
		&m.NextSendAt, &m.LastSentAt, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if message.Valid {
		m.Message = &message.String
	}
	if commandID.Valid {
		m.CommandID = &commandID.Int64
	}

	return m, nil
}

func scanScheduledMessages(rows *sql.Rows) ([]*domain.ScheduledMessage, error) {
	messages := make([]*domain.ScheduledMessage, 0)
	for rows.Next() {
		m, err := scanOneScheduledMessage(rows)
		if err != nil {
			return nil, fmt.Errorf("fehler beim Scannen der automatisierten Nachricht: %w", err)
		}
		messages = append(messages, m)
	}
	return messages, nil
}
