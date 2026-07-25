package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
)

type automodRepository struct {
	db *sql.DB
}

func NewAutomodRepository(db *sql.DB) repository.AutomodRepository {
	return &automodRepository{db: db}
}

const automodSettingsColumns = `
	user_twitch_id, enabled, blocked_words, link_filter_enabled, allowed_domains,
	caps_filter_enabled, symbol_filter_enabled, emote_filter_enabled, emote_threshold, repetition_filter_enabled,
	exempt_vips, exempt_regulars, exempt_users, created_at, updated_at
`

func scanAutomodSettings(row interface{ Scan(dest ...interface{}) error }) (*domain.AutomodSettings, error) {
	s := &domain.AutomodSettings{}
	err := row.Scan(
		&s.UserTwitchID, &s.Enabled, &s.BlockedWords, &s.LinkFilterEnabled, &s.AllowedDomains,
		&s.CapsFilterEnabled, &s.SymbolFilterEnabled, &s.EmoteFilterEnabled, &s.EmoteThreshold, &s.RepetitionFilterEnabled,
		&s.ExemptVips, &s.ExemptRegulars, &s.ExemptUsers, &s.CreatedAt, &s.UpdatedAt,
	)
	return s, err
}

func (r *automodRepository) GetSettings(ctx context.Context, userTwitchID string) (*domain.AutomodSettings, error) {
	query := `SELECT ` + automodSettingsColumns + ` FROM automod_settings WHERE user_twitch_id = $1`

	s, err := scanAutomodSettings(r.db.QueryRowContext(ctx, query, userTwitchID))
	if err == sql.ErrNoRows {
		return domain.NewAutomodSettings(userTwitchID), nil
	}
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Automod-Settings: %w", err)
	}

	return s, nil
}

func (r *automodRepository) UpsertSettings(ctx context.Context, s *domain.AutomodSettings) error {
	query := `
		INSERT INTO automod_settings
		(user_twitch_id, enabled, blocked_words, link_filter_enabled, allowed_domains,
		 caps_filter_enabled, symbol_filter_enabled, emote_filter_enabled, emote_threshold, repetition_filter_enabled,
		 exempt_vips, exempt_regulars, exempt_users, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		ON CONFLICT (user_twitch_id)
		DO UPDATE SET
			enabled = EXCLUDED.enabled,
			blocked_words = EXCLUDED.blocked_words,
			link_filter_enabled = EXCLUDED.link_filter_enabled,
			allowed_domains = EXCLUDED.allowed_domains,
			caps_filter_enabled = EXCLUDED.caps_filter_enabled,
			symbol_filter_enabled = EXCLUDED.symbol_filter_enabled,
			emote_filter_enabled = EXCLUDED.emote_filter_enabled,
			emote_threshold = EXCLUDED.emote_threshold,
			repetition_filter_enabled = EXCLUDED.repetition_filter_enabled,
			exempt_vips = EXCLUDED.exempt_vips,
			exempt_regulars = EXCLUDED.exempt_regulars,
			exempt_users = EXCLUDED.exempt_users,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(
		ctx, query,
		s.UserTwitchID, s.Enabled, s.BlockedWords, s.LinkFilterEnabled, s.AllowedDomains,
		s.CapsFilterEnabled, s.SymbolFilterEnabled, s.EmoteFilterEnabled, s.EmoteThreshold, s.RepetitionFilterEnabled,
		s.ExemptVips, s.ExemptRegulars, s.ExemptUsers, s.CreatedAt, s.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("fehler beim Speichern der Automod-Settings: %w", err)
	}

	return nil
}

func (r *automodRepository) GetAllEnabledSettings(ctx context.Context) ([]*domain.AutomodSettings, error) {
	query := `SELECT ` + automodSettingsColumns + ` FROM automod_settings WHERE enabled = true`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden aktivierter Automod-Settings: %w", err)
	}
	defer rows.Close()

	settingsList := make([]*domain.AutomodSettings, 0)
	for rows.Next() {
		s, err := scanAutomodSettings(rows)
		if err != nil {
			return nil, fmt.Errorf("fehler beim Scannen der Automod-Settings: %w", err)
		}
		settingsList = append(settingsList, s)
	}

	return settingsList, nil
}

// RecordViolation is a single atomic upsert: a violation streak resets to 1
// if the offender's last violation is older than the streak TTL, otherwise
// the counter just increments - all enforced in one query so concurrent
// violations from the same offender can't race each other.
func (r *automodRepository) RecordViolation(ctx context.Context, userTwitchID, offenderTwitchID string) (int, error) {
	query := `
		INSERT INTO automod_violations (user_twitch_id, offender_twitch_id, violation_count, last_violation_at)
		VALUES ($1, $2, 1, NOW())
		ON CONFLICT (user_twitch_id, offender_twitch_id)
		DO UPDATE SET
			violation_count = CASE
				WHEN automod_violations.last_violation_at IS NULL
				     OR automod_violations.last_violation_at < NOW() - INTERVAL '24 hours'
				THEN 1
				ELSE automod_violations.violation_count + 1
			END,
			last_violation_at = NOW()
		RETURNING violation_count
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, userTwitchID, offenderTwitchID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("fehler beim Fortschreiben des Verstoß-Zählers: %w", err)
	}

	return count, nil
}

func (r *automodRepository) LogEvent(ctx context.Context, event *domain.AutomodEvent) error {
	query := `
		INSERT INTO automod_events
		(user_twitch_id, offender_twitch_id, offender_name, reason, message_excerpt, timeout_seconds, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx, query,
		event.UserTwitchID, event.OffenderTwitchID, event.OffenderName, event.Reason,
		event.MessageExcerpt, event.TimeoutSeconds, event.CreatedAt,
	).Scan(&event.ID)
	if err != nil {
		return fmt.Errorf("fehler beim Protokollieren des Automod-Events: %w", err)
	}

	return nil
}

func (r *automodRepository) GetEvents(ctx context.Context, userTwitchID string, limit, offset int) ([]*domain.AutomodEvent, int64, error) {
	var total int64
	err := r.db.QueryRowContext(
		ctx, "SELECT COUNT(*) FROM automod_events WHERE user_twitch_id = $1", userTwitchID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("fehler beim Zählen der Automod-Events: %w", err)
	}

	query := `
		SELECT id, user_twitch_id, offender_twitch_id, offender_name, reason, message_excerpt,
		       timeout_seconds, created_at
		FROM automod_events
		WHERE user_twitch_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userTwitchID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("fehler beim Laden der Automod-Events: %w", err)
	}
	defer rows.Close()

	events := make([]*domain.AutomodEvent, 0)
	for rows.Next() {
		e := &domain.AutomodEvent{}
		var excerpt sql.NullString
		if err := rows.Scan(
			&e.ID, &e.UserTwitchID, &e.OffenderTwitchID, &e.OffenderName, &e.Reason, &excerpt,
			&e.TimeoutSeconds, &e.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("fehler beim Scannen des Automod-Events: %w", err)
		}
		if excerpt.Valid {
			e.MessageExcerpt = excerpt.String
		}
		events = append(events, e)
	}

	return events, total, nil
}
