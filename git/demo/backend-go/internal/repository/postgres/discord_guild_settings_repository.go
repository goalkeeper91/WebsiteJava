package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"demo/backend-go/internal/domain"
)

type DiscordGuildSettingsRepository struct {
	db *sql.DB
}

func NewDiscordGuildSettingsRepository(db *sql.DB) *DiscordGuildSettingsRepository {
	return &DiscordGuildSettingsRepository{db: db}
}

func (r *DiscordGuildSettingsRepository) Create(ctx context.Context, input domain.DiscordGuildSettingsCreateInput) (*domain.DiscordGuildSettings, error) {
	query := `
		INSERT INTO discord_guild_settings (
			guild_id, user_id, notification_channel_id, command_channel_id,
			activity_channel_id, twitch_notifications_enabled,
			join_to_create_enabled, admin_role_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (guild_id, user_id) DO UPDATE SET
			notification_channel_id = EXCLUDED.notification_channel_id,
			command_channel_id = EXCLUDED.command_channel_id,
			activity_channel_id = EXCLUDED.activity_channel_id,
			twitch_notifications_enabled = EXCLUDED.twitch_notifications_enabled,
			join_to_create_enabled = EXCLUDED.join_to_create_enabled,
			admin_role_id = EXCLUDED.admin_role_id,
			updated_at = NOW()
		RETURNING id, guild_id, user_id, notification_channel_id, command_channel_id,
		          activity_channel_id, twitch_notifications_enabled, join_to_create_enabled,
		          admin_role_id, created_at, updated_at
	`

	var settings domain.DiscordGuildSettings
	err := r.db.QueryRowContext(
		ctx, query,
		input.GuildID,
		input.UserID,
		input.NotificationChannelID,
		input.CommandChannelID,
		input.ActivityChannelID,
		input.TwitchNotificationsEnabled,
		input.JoinToCreateEnabled,
		input.AdminRoleID,
	).Scan(
		&settings.ID,
		&settings.GuildID,
		&settings.UserID,
		&settings.NotificationChannelID,
		&settings.CommandChannelID,
		&settings.ActivityChannelID,
		&settings.TwitchNotificationsEnabled,
		&settings.JoinToCreateEnabled,
		&settings.AdminRoleID,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create guild settings: %w", err)
	}

	return &settings, nil
}

func (r *DiscordGuildSettingsRepository) GetByGuildAndUser(ctx context.Context, guildID int64, userID int64) (*domain.DiscordGuildSettings, error) {
	query := `
		SELECT id, guild_id, user_id, notification_channel_id, command_channel_id,
		       activity_channel_id, twitch_notifications_enabled, join_to_create_enabled,
		       admin_role_id, created_at, updated_at
		FROM discord_guild_settings
		WHERE guild_id = $1 AND user_id = $2
	`

	var settings domain.DiscordGuildSettings
	err := r.db.QueryRowContext(ctx, query, guildID, userID).Scan(
		&settings.ID,
		&settings.GuildID,
		&settings.UserID,
		&settings.NotificationChannelID,
		&settings.CommandChannelID,
		&settings.ActivityChannelID,
		&settings.TwitchNotificationsEnabled,
		&settings.JoinToCreateEnabled,
		&settings.AdminRoleID,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get guild settings: %w", err)
	}

	return &settings, nil
}

func (r *DiscordGuildSettingsRepository) GetByGuild(ctx context.Context, guildID int64) ([]domain.DiscordGuildSettings, error) {
	query := `
		SELECT id, guild_id, user_id, notification_channel_id, command_channel_id,
		       activity_channel_id, twitch_notifications_enabled, join_to_create_enabled,
		       admin_role_id, created_at, updated_at
		FROM discord_guild_settings
		WHERE guild_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, guildID)
	if err != nil {
		return nil, fmt.Errorf("failed to get settings by guild: %w", err)
	}
	defer rows.Close()

	var settingsList []domain.DiscordGuildSettings
	for rows.Next() {
		var settings domain.DiscordGuildSettings
		err := rows.Scan(
			&settings.ID,
			&settings.GuildID,
			&settings.UserID,
			&settings.NotificationChannelID,
			&settings.CommandChannelID,
			&settings.ActivityChannelID,
			&settings.TwitchNotificationsEnabled,
			&settings.JoinToCreateEnabled,
			&settings.AdminRoleID,
			&settings.CreatedAt,
			&settings.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan settings: %w", err)
		}
		settingsList = append(settingsList, settings)
	}

	return settingsList, nil
}

func (r *DiscordGuildSettingsRepository) GetByUser(ctx context.Context, userID int64) ([]domain.DiscordGuildSettings, error) {
	query := `
		SELECT s.id, s.guild_id, s.user_id, s.notification_channel_id, s.command_channel_id,
		       s.activity_channel_id, s.twitch_notifications_enabled, s.join_to_create_enabled,
		       s.admin_role_id, s.created_at, s.updated_at
		FROM discord_guild_settings s
		INNER JOIN discord_guilds g ON s.guild_id = g.id
		WHERE s.user_id = $1 AND g.is_active = true
		ORDER BY g.name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get settings by user: %w", err)
	}
	defer rows.Close()

	var settingsList []domain.DiscordGuildSettings
	for rows.Next() {
		var settings domain.DiscordGuildSettings
		err := rows.Scan(
			&settings.ID,
			&settings.GuildID,
			&settings.UserID,
			&settings.NotificationChannelID,
			&settings.CommandChannelID,
			&settings.ActivityChannelID,
			&settings.TwitchNotificationsEnabled,
			&settings.JoinToCreateEnabled,
			&settings.AdminRoleID,
			&settings.CreatedAt,
			&settings.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan settings: %w", err)
		}
		settingsList = append(settingsList, settings)
	}

	return settingsList, nil
}

func (r *DiscordGuildSettingsRepository) Update(ctx context.Context, guildID int64, userID int64, input domain.DiscordGuildSettingsUpdateInput) error {
	query := `
		UPDATE discord_guild_settings
		SET notification_channel_id = COALESCE($3, notification_channel_id),
		    command_channel_id = COALESCE($4, command_channel_id),
		    activity_channel_id = COALESCE($5, activity_channel_id),
		    twitch_notifications_enabled = COALESCE($6, twitch_notifications_enabled),
		    join_to_create_enabled = COALESCE($7, join_to_create_enabled),
		    admin_role_id = COALESCE($8, admin_role_id),
		    updated_at = NOW()
		WHERE guild_id = $1 AND user_id = $2
	`

	result, err := r.db.ExecContext(
		ctx, query,
		guildID,
		userID,
		input.NotificationChannelID,
		input.CommandChannelID,
		input.ActivityChannelID,
		input.TwitchNotificationsEnabled,
		input.JoinToCreateEnabled,
		input.AdminRoleID,
	)

	if err != nil {
		return fmt.Errorf("failed to update guild settings: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("settings not found for guild %d and user %d", guildID, userID)
	}

	return nil
}

func (r *DiscordGuildSettingsRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM discord_guild_settings WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete guild settings: %w", err)
	}

	return nil
}

// DeleteByGuildAndUser deletes settings for a specific guild and user
func (r *DiscordGuildSettingsRepository) DeleteByGuildAndUser(ctx context.Context, guildID int64, userID int64) error {
	query := `DELETE FROM discord_guild_settings WHERE guild_id = $1 AND user_id = $2`

	_, err := r.db.ExecContext(ctx, query, guildID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete guild settings: %w", err)
	}

	return nil
}

func (r *DiscordGuildSettingsRepository) UserCanAccessGuild(ctx context.Context, userID int64, guildID int64) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM discord_guild_settings
			WHERE user_id = $1 AND guild_id = $2
		)
	`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, userID, guildID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check guild access: %w", err)
	}

	return exists, nil
}

func (r *DiscordGuildSettingsRepository) GetNotificationChannel(ctx context.Context, userID int64, guildID int64) (*int64, error) {
	query := `
		SELECT notification_channel_id
		FROM discord_guild_settings
		WHERE user_id = $1 AND guild_id = $2 AND twitch_notifications_enabled = true
	`

	var channelID *int64
	err := r.db.QueryRowContext(ctx, query, userID, guildID).Scan(&channelID)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get notification channel: %w", err)
	}

	return channelID, nil
}

func (r *DiscordGuildSettingsRepository) GetCommandChannel(ctx context.Context, userID int64, guildID int64) (*int64, error) {
	query := `
		SELECT command_channel_id
		FROM discord_guild_settings
		WHERE user_id = $1 AND guild_id = $2 AND twitch_notifications_enabled = true
	`

	var channelID *int64
	err := r.db.QueryRowContext(ctx, query, userID, guildID).Scan(&channelID)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get command channel: %w", err)
	}

	return channelID, nil
}