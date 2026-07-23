package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
)

type twitchChannelRepository struct {
	db *sql.DB
}

func NewTwitchChannelRepository(db *sql.DB) repository.TwitchChannelRepository {
	return &twitchChannelRepository{db: db}
}

func (r *twitchChannelRepository) Create(ctx context.Context, channel *domain.TwitchChannel) error {
	query := `
		INSERT INTO twitch_channels (twitch_user_id, user_name, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		channel.TwitchUserID,
		channel.Username,
		channel.IsActive,
		channel.CreatedAt,
		channel.UpdatedAt,
	).Scan(&channel.ID)

	if err != nil {
		return fmt.Errorf("fehler beim Erstellen des Channels: %w", err)
	}

	return nil
}

func (r *twitchChannelRepository) GetByTwitchUserID(ctx context.Context, twitchUserID string) (*domain.TwitchChannel, error) {
	query := `
		SELECT id, twitch_user_id, user_name, is_active, created_at, updated_at
		FROM twitch_channels
		WHERE twitch_user_id = $1
	`

	channel := &domain.TwitchChannel{}
	err := r.db.QueryRowContext(ctx, query, twitchUserID).Scan(
		&channel.ID,
		&channel.TwitchUserID,
		&channel.Username,
		&channel.IsActive,
		&channel.CreatedAt,
		&channel.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrChannelNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden des Channels: %w", err)
	}

	return channel, nil
}

func (r *twitchChannelRepository) GetByID(ctx context.Context, id int64) (*domain.TwitchChannel, error) {
	query := `
		SELECT id, twitch_user_id, user_name, is_active, created_at, updated_at
		FROM twitch_channels
		WHERE id = $1
	`

	channel := &domain.TwitchChannel{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&channel.ID,
		&channel.TwitchUserID,
		&channel.Username,
		&channel.IsActive,
		&channel.CreatedAt,
		&channel.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrChannelNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden des Channels: %w", err)
	}

	return channel, nil
}

func (r *twitchChannelRepository) Update(ctx context.Context, channel *domain.TwitchChannel) error {
	query := `
		UPDATE twitch_channels
		SET user_name = $2, is_active = $3, updated_at = $4
		WHERE twitch_user_id = $1
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		channel.TwitchUserID,
		channel.Username,
		channel.IsActive,
		channel.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("fehler beim Update des Channels: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("fehler beim Abrufen betroffener Zeilen: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrChannelNotFound
	}

	return nil
}

func (r *twitchChannelRepository) Upsert(ctx context.Context, channel *domain.TwitchChannel) error {
	query := `
		INSERT INTO twitch_channels (twitch_user_id, user_name, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (twitch_user_id)
		DO UPDATE SET
			user_name = EXCLUDED.user_name,
			is_active = EXCLUDED.is_active,
			updated_at = EXCLUDED.updated_at
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		channel.TwitchUserID,
		channel.Username,
		channel.IsActive,
		channel.CreatedAt,
		channel.UpdatedAt,
	).Scan(&channel.ID)

	if err != nil {
		return fmt.Errorf("fehler beim Upsert des Channels: %w", err)
	}

	return nil
}

func (r *twitchChannelRepository) Delete(ctx context.Context, twitchUserID string) error {
	query := `DELETE FROM twitch_channels WHERE twitch_user_id = $1`

	result, err := r.db.ExecContext(ctx, query, twitchUserID)
	if err != nil {
		return fmt.Errorf("fehler beim Löschen des Channels: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("fehler beim Abrufen betroffener Zeilen: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrChannelNotFound
	}

	return nil
}

func (r *twitchChannelRepository) Exists(ctx context.Context, twitchUserID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM twitch_channels WHERE twitch_user_id = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, twitchUserID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("fehler beim Prüfen der Channel-Existenz: %w", err)
	}

	return exists, nil
}