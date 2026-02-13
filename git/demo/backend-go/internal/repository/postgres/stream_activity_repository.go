package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
)

type streamActivityRepository struct {
	db *sql.DB
}

func NewStreamActivityRepository(db *sql.DB) repository.StreamActivityRepository {
	return &streamActivityRepository{db: db}
}

func (r *streamActivityRepository) Create(ctx context.Context, activity *domain.StreamActivity) error {
	query := `
		INSERT INTO stream_activities
		(twitch_user_id, type, username, display_name, viewers, bits, tier, message, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		activity.TwitchUserID,
		activity.Type,
		activity.Username,
		activity.DisplayName,
		activity.Viewers,
		activity.Bits,
		activity.Tier,
		activity.Message,
		activity.Timestamp,
	).Scan(&activity.ID)

	if err != nil {
		return fmt.Errorf("fehler beim Erstellen der Activity: %w", err)
	}

	return nil
}

func (r *streamActivityRepository) GetRecent(ctx context.Context, twitchUserID string, limit int) ([]*domain.StreamActivity, error) {
	query := `
		SELECT id, twitch_user_id, type, username, display_name, viewers, bits, tier, message, timestamp
		FROM stream_activities
		WHERE twitch_user_id = $1
		ORDER BY timestamp DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, twitchUserID, limit)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Activities: %w", err)
	}
	defer rows.Close()

	activities := make([]*domain.StreamActivity, 0)
	for rows.Next() {
		activity := &domain.StreamActivity{}
		err := rows.Scan(
			&activity.ID,
			&activity.TwitchUserID,
			&activity.Type,
			&activity.Username,
			&activity.DisplayName,
			&activity.Viewers,
			&activity.Bits,
			&activity.Tier,
			&activity.Message,
			&activity.Timestamp,
		)
		if err != nil {
			return nil, fmt.Errorf("fehler beim Scannen der Activity: %w", err)
		}
		activities = append(activities, activity)
	}

	return activities, nil
}

func (r *streamActivityRepository) GetByType(ctx context.Context, twitchUserID string, activityType domain.ActivityType, limit int) ([]*domain.StreamActivity, error) {
	query := `
		SELECT id, twitch_user_id, type, username, display_name, viewers, bits, tier, message, timestamp
		FROM stream_activities
		WHERE twitch_user_id = $1 AND type = $2
		ORDER BY timestamp DESC
		LIMIT $3
	`

	rows, err := r.db.QueryContext(ctx, query, twitchUserID, activityType, limit)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Activities: %w", err)
	}
	defer rows.Close()

	activities := make([]*domain.StreamActivity, 0)
	for rows.Next() {
		activity := &domain.StreamActivity{}
		err := rows.Scan(
			&activity.ID,
			&activity.TwitchUserID,
			&activity.Type,
			&activity.Username,
			&activity.DisplayName,
			&activity.Viewers,
			&activity.Bits,
			&activity.Tier,
			&activity.Message,
			&activity.Timestamp,
		)
		if err != nil {
			return nil, fmt.Errorf("fehler beim Scannen der Activity: %w", err)
		}
		activities = append(activities, activity)
	}

	return activities, nil
}

func (r *streamActivityRepository) DeleteOlderThan(ctx context.Context, timestamp int64) error {
	query := `DELETE FROM stream_activities WHERE EXTRACT(EPOCH FROM timestamp) < $1`

	result, err := r.db.ExecContext(ctx, query, timestamp)
	if err != nil {
		return fmt.Errorf("fehler beim Löschen alter Activities: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("fehler beim Abrufen betroffener Zeilen: %w", err)
	}

	if rowsAffected > 0 {
		fmt.Printf("✅ %d alte Activities gelöscht\n", rowsAffected)
	}

	return nil
}