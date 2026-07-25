package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
)

type loyaltyRepository struct {
	db *sql.DB
}

func NewLoyaltyRepository(db *sql.DB) repository.LoyaltyRepository {
	return &loyaltyRepository{db: db}
}

const loyaltySettingsColumns = `
	user_twitch_id, enabled, points_name, points_per_interval, interval_minutes,
	next_accrual_at, created_at, updated_at
`

func scanLoyaltySettings(row interface{ Scan(dest ...interface{}) error }) (*domain.LoyaltySettings, error) {
	s := &domain.LoyaltySettings{}
	err := row.Scan(
		&s.UserTwitchID, &s.Enabled, &s.PointsName, &s.PointsPerInterval, &s.IntervalMinutes,
		&s.NextAccrualAt, &s.CreatedAt, &s.UpdatedAt,
	)
	return s, err
}

func (r *loyaltyRepository) GetSettings(ctx context.Context, userTwitchID string) (*domain.LoyaltySettings, error) {
	query := `SELECT ` + loyaltySettingsColumns + ` FROM loyalty_settings WHERE user_twitch_id = $1`

	s, err := scanLoyaltySettings(r.db.QueryRowContext(ctx, query, userTwitchID))
	if err == sql.ErrNoRows {
		return domain.NewLoyaltySettings(userTwitchID), nil
	}
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Loyalty-Settings: %w", err)
	}

	return s, nil
}

func (r *loyaltyRepository) UpsertSettings(ctx context.Context, s *domain.LoyaltySettings) error {
	query := `
		INSERT INTO loyalty_settings
		(user_twitch_id, enabled, points_name, points_per_interval, interval_minutes,
		 next_accrual_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (user_twitch_id)
		DO UPDATE SET
			enabled = EXCLUDED.enabled,
			points_name = EXCLUDED.points_name,
			points_per_interval = EXCLUDED.points_per_interval,
			interval_minutes = EXCLUDED.interval_minutes,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(
		ctx, query,
		s.UserTwitchID, s.Enabled, s.PointsName, s.PointsPerInterval, s.IntervalMinutes,
		s.NextAccrualAt, s.CreatedAt, s.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("fehler beim Speichern der Loyalty-Settings: %w", err)
	}

	return nil
}

func (r *loyaltyRepository) GetDueForAccrual(ctx context.Context, limit int) ([]*domain.LoyaltySettings, error) {
	query := `
		SELECT ` + loyaltySettingsColumns + `
		FROM loyalty_settings
		WHERE enabled = true AND next_accrual_at <= NOW()
		ORDER BY next_accrual_at ASC
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden fälliger Loyalty-Settings: %w", err)
	}
	defer rows.Close()

	settingsList := make([]*domain.LoyaltySettings, 0)
	for rows.Next() {
		s, err := scanLoyaltySettings(rows)
		if err != nil {
			return nil, fmt.Errorf("fehler beim Scannen der Loyalty-Settings: %w", err)
		}
		settingsList = append(settingsList, s)
	}

	return settingsList, nil
}

func (r *loyaltyRepository) MarkAccrued(ctx context.Context, userTwitchID string, nextAccrualAt time.Time) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE loyalty_settings
		SET next_accrual_at = $2, updated_at = NOW()
		WHERE user_twitch_id = $1
	`, userTwitchID, nextAccrualAt)
	if err != nil {
		return fmt.Errorf("fehler beim Fortschreiben des Loyalty-Zeitplans: %w", err)
	}
	return nil
}

// CreditPoints upserts each viewer's balance in one transaction - a tick
// that fails partway through never leaves some viewers credited and others
// not for the same channel.
func (r *loyaltyRepository) CreditPoints(ctx context.Context, userTwitchID string, credits []domain.ViewerCredit) error {
	if len(credits) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("fehler beim Starten der Loyalty-Transaktion: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO loyalty_points (user_twitch_id, viewer_twitch_id, viewer_login, points, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (user_twitch_id, viewer_twitch_id)
		DO UPDATE SET
			points = loyalty_points.points + EXCLUDED.points,
			viewer_login = EXCLUDED.viewer_login,
			updated_at = NOW()
	`)
	if err != nil {
		return fmt.Errorf("fehler beim Vorbereiten der Punkte-Gutschrift: %w", err)
	}
	defer stmt.Close()

	for _, credit := range credits {
		if _, err := stmt.ExecContext(ctx, userTwitchID, credit.ViewerTwitchID, credit.ViewerLogin, credit.Amount); err != nil {
			return fmt.Errorf("fehler beim Gutschreiben von Punkten für %s: %w", credit.ViewerLogin, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("fehler beim Committen der Loyalty-Transaktion: %w", err)
	}

	return nil
}

func (r *loyaltyRepository) GetLeaderboard(ctx context.Context, userTwitchID string, limit, offset int) ([]*domain.LoyaltyPoints, int64, error) {
	var total int64
	err := r.db.QueryRowContext(
		ctx, "SELECT COUNT(*) FROM loyalty_points WHERE user_twitch_id = $1", userTwitchID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("fehler beim Zählen der Loyalty-Bestenliste: %w", err)
	}

	query := `
		SELECT user_twitch_id, viewer_twitch_id, viewer_login, points, updated_at
		FROM loyalty_points
		WHERE user_twitch_id = $1
		ORDER BY points DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userTwitchID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("fehler beim Laden der Loyalty-Bestenliste: %w", err)
	}
	defer rows.Close()

	entries := make([]*domain.LoyaltyPoints, 0)
	for rows.Next() {
		p := &domain.LoyaltyPoints{}
		if err := rows.Scan(&p.UserTwitchID, &p.ViewerTwitchID, &p.ViewerLogin, &p.Points, &p.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("fehler beim Scannen der Loyalty-Bestenliste: %w", err)
		}
		entries = append(entries, p)
	}

	return entries, total, nil
}

func (r *loyaltyRepository) GetViewerPoints(ctx context.Context, userTwitchID, viewerLogin string) (*domain.LoyaltyPoints, error) {
	query := `
		SELECT user_twitch_id, viewer_twitch_id, viewer_login, points, updated_at
		FROM loyalty_points
		WHERE user_twitch_id = $1 AND LOWER(viewer_login) = LOWER($2)
	`

	p := &domain.LoyaltyPoints{}
	err := r.db.QueryRowContext(ctx, query, userTwitchID, viewerLogin).Scan(
		&p.UserTwitchID, &p.ViewerTwitchID, &p.ViewerLogin, &p.Points, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Viewer-Punkte: %w", err)
	}

	return p, nil
}
