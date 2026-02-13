package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"demo/backend-go/internal/domain"
)

type JoinToCreateRepository struct {
	db *sql.DB
}

func NewJoinToCreateRepository(db *sql.DB) *JoinToCreateRepository {
	return &JoinToCreateRepository{db: db}
}

func (r *JoinToCreateRepository) Create(ctx context.Context, input domain.JoinToCreateConfigCreateInput) (*domain.JoinToCreateConfig, error) {
	query := `
		INSERT INTO join_to_create_configs (
			guild_id, user_id, join_channel_id, category_id,
			channel_name_prefix, user_limit, private_channel, enabled
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, guild_id, user_id, join_channel_id, category_id,
		          channel_name_prefix, user_limit, private_channel, enabled,
		          created_at, updated_at
	`

	var config domain.JoinToCreateConfig
	err := r.db.QueryRowContext(
		ctx, query,
		input.GuildID,
		input.UserID,
		input.JoinChannelID,
		input.CategoryID,
		input.ChannelNamePrefix,
		input.UserLimit,
		input.PrivateChannel,
		input.Enabled,
	).Scan(
		&config.ID,
		&config.GuildID,
		&config.UserID,
		&config.JoinChannelID,
		&config.CategoryID,
		&config.ChannelNamePrefix,
		&config.UserLimit,
		&config.PrivateChannel,
		&config.Enabled,
		&config.CreatedAt,
		&config.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create join-to-create config: %w", err)
	}

	return &config, nil
}

func (r *JoinToCreateRepository) GetByID(ctx context.Context, id int64) (*domain.JoinToCreateConfig, error) {
	query := `
		SELECT id, guild_id, user_id, join_channel_id, category_id,
		       channel_name_prefix, user_limit, private_channel, enabled,
		       created_at, updated_at
		FROM join_to_create_configs
		WHERE id = $1
	`

	var config domain.JoinToCreateConfig
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&config.ID,
		&config.GuildID,
		&config.UserID,
		&config.JoinChannelID,
		&config.CategoryID,
		&config.ChannelNamePrefix,
		&config.UserLimit,
		&config.PrivateChannel,
		&config.Enabled,
		&config.CreatedAt,
		&config.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get join-to-create config: %w", err)
	}

	return &config, nil
}

func (r *JoinToCreateRepository) GetByGuild(ctx context.Context, guildID int64) ([]domain.JoinToCreateConfig, error) {
	query := `
		SELECT id, guild_id, user_id, join_channel_id, category_id,
		       channel_name_prefix, user_limit, private_channel, enabled,
		       created_at, updated_at
		FROM join_to_create_configs
		WHERE guild_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, guildID)
	if err != nil {
		return nil, fmt.Errorf("failed to get configs by guild: %w", err)
	}
	defer rows.Close()

	var configs []domain.JoinToCreateConfig
	for rows.Next() {
		var config domain.JoinToCreateConfig
		err := rows.Scan(
			&config.ID,
			&config.GuildID,
			&config.UserID,
			&config.JoinChannelID,
			&config.CategoryID,
			&config.ChannelNamePrefix,
			&config.UserLimit,
			&config.PrivateChannel,
			&config.Enabled,
			&config.CreatedAt,
			&config.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan config: %w", err)
		}
		configs = append(configs, config)
	}

	return configs, nil
}

func (r *JoinToCreateRepository) GetByUser(ctx context.Context, userID int64) ([]domain.JoinToCreateConfig, error) {
	query := `
		SELECT c.id, c.guild_id, c.user_id, c.join_channel_id, c.category_id,
		       c.channel_name_prefix, c.user_limit, c.private_channel, c.enabled,
		       c.created_at, c.updated_at
		FROM join_to_create_configs c
		INNER JOIN discord_guilds g ON c.guild_id = g.id
		WHERE c.user_id = $1 AND g.is_active = true
		ORDER BY c.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get configs by user: %w", err)
	}
	defer rows.Close()

	var configs []domain.JoinToCreateConfig
	for rows.Next() {
		var config domain.JoinToCreateConfig
		err := rows.Scan(
			&config.ID,
			&config.GuildID,
			&config.UserID,
			&config.JoinChannelID,
			&config.CategoryID,
			&config.ChannelNamePrefix,
			&config.UserLimit,
			&config.PrivateChannel,
			&config.Enabled,
			&config.CreatedAt,
			&config.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan config: %w", err)
		}
		configs = append(configs, config)
	}

	return configs, nil
}

func (r *JoinToCreateRepository) GetByJoinChannel(ctx context.Context, joinChannelID int64) (*domain.JoinToCreateConfig, error) {
	query := `
		SELECT id, guild_id, user_id, join_channel_id, category_id,
		       channel_name_prefix, user_limit, private_channel, enabled,
		       created_at, updated_at
		FROM join_to_create_configs
		WHERE join_channel_id = $1 AND enabled = true
		LIMIT 1
	`

	var config domain.JoinToCreateConfig
	err := r.db.QueryRowContext(ctx, query, joinChannelID).Scan(
		&config.ID,
		&config.GuildID,
		&config.UserID,
		&config.JoinChannelID,
		&config.CategoryID,
		&config.ChannelNamePrefix,
		&config.UserLimit,
		&config.PrivateChannel,
		&config.Enabled,
		&config.CreatedAt,
		&config.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get config by join channel: %w", err)
	}

	return &config, nil
}

func (r *JoinToCreateRepository) GetEnabled(ctx context.Context) ([]domain.JoinToCreateConfig, error) {
	query := `
		SELECT c.id, c.guild_id, c.user_id, c.join_channel_id, c.category_id,
		       c.channel_name_prefix, c.user_limit, c.private_channel, c.enabled,
		       c.created_at, c.updated_at
		FROM join_to_create_configs c
		INNER JOIN discord_guilds g ON c.guild_id = g.id
		WHERE c.enabled = true AND g.is_active = true
		ORDER BY c.guild_id, c.created_at
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get enabled configs: %w", err)
	}
	defer rows.Close()

	var configs []domain.JoinToCreateConfig
	for rows.Next() {
		var config domain.JoinToCreateConfig
		err := rows.Scan(
			&config.ID,
			&config.GuildID,
			&config.UserID,
			&config.JoinChannelID,
			&config.CategoryID,
			&config.ChannelNamePrefix,
			&config.UserLimit,
			&config.PrivateChannel,
			&config.Enabled,
			&config.CreatedAt,
			&config.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan config: %w", err)
		}
		configs = append(configs, config)
	}

	return configs, nil
}

func (r *JoinToCreateRepository) Update(ctx context.Context, id int64, input domain.JoinToCreateConfigUpdateInput) error {
	query := `
		UPDATE join_to_create_configs
		SET join_channel_id = COALESCE($2, join_channel_id),
		    category_id = COALESCE($3, category_id),
		    channel_name_prefix = COALESCE($4, channel_name_prefix),
		    user_limit = COALESCE($5, user_limit),
		    private_channel = COALESCE($6, private_channel),
		    enabled = COALESCE($7, enabled),
		    updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(
		ctx, query,
		id,
		input.JoinChannelID,
		input.CategoryID,
		input.ChannelNamePrefix,
		input.UserLimit,
		input.PrivateChannel,
		input.Enabled,
	)

	if err != nil {
		return fmt.Errorf("failed to update join-to-create config: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("config with id %d not found", id)
	}

	return nil
}

func (r *JoinToCreateRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM join_to_create_configs WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete join-to-create config: %w", err)
	}

	return nil
}

func (r *JoinToCreateRepository) UserOwnsConfig(ctx context.Context, userID int64, configID int64) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM join_to_create_configs
			WHERE id = $1 AND user_id = $2
		)
	`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, configID, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check config ownership: %w", err)
	}

	return exists, nil
}

func (r *JoinToCreateRepository) GetGuildConfigCount(ctx context.Context, guildID int64) (int, error) {
	query := `
		SELECT COUNT(*) FROM join_to_create_configs
		WHERE guild_id = $1
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, guildID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get config count: %w", err)
	}

	return count, nil
}