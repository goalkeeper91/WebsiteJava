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

func (r *automodRepository) GetSettings(ctx context.Context, userTwitchID string) (*domain.AutomodSettings, error) {
	query := `
		SELECT user_twitch_id, enabled, blocked_words, link_filter_enabled, allowed_domains,
		       created_at, updated_at
		FROM automod_settings
		WHERE user_twitch_id = $1
	`

	s := &domain.AutomodSettings{}
	err := r.db.QueryRowContext(ctx, query, userTwitchID).Scan(
		&s.UserTwitchID, &s.Enabled, &s.BlockedWords, &s.LinkFilterEnabled, &s.AllowedDomains,
		&s.CreatedAt, &s.UpdatedAt,
	)

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
		(user_twitch_id, enabled, blocked_words, link_filter_enabled, allowed_domains, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (user_twitch_id)
		DO UPDATE SET
			enabled = EXCLUDED.enabled,
			blocked_words = EXCLUDED.blocked_words,
			link_filter_enabled = EXCLUDED.link_filter_enabled,
			allowed_domains = EXCLUDED.allowed_domains,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(
		ctx, query,
		s.UserTwitchID, s.Enabled, s.BlockedWords, s.LinkFilterEnabled, s.AllowedDomains,
		s.CreatedAt, s.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("fehler beim Speichern der Automod-Settings: %w", err)
	}

	return nil
}

func (r *automodRepository) GetAllEnabledSettings(ctx context.Context) ([]*domain.AutomodSettings, error) {
	query := `
		SELECT user_twitch_id, enabled, blocked_words, link_filter_enabled, allowed_domains,
		       created_at, updated_at
		FROM automod_settings
		WHERE enabled = true
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden aktivierter Automod-Settings: %w", err)
	}
	defer rows.Close()

	settingsList := make([]*domain.AutomodSettings, 0)
	for rows.Next() {
		s := &domain.AutomodSettings{}
		if err := rows.Scan(
			&s.UserTwitchID, &s.Enabled, &s.BlockedWords, &s.LinkFilterEnabled, &s.AllowedDomains,
			&s.CreatedAt, &s.UpdatedAt,
		); err != nil {
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
