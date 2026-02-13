package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"demo/backend-go/internal/domain"
)

type DiscordGuildRepository struct {
	db *sql.DB
}

func NewDiscordGuildRepository(db *sql.DB) *DiscordGuildRepository {
	return &DiscordGuildRepository{db: db}
}

func (r *DiscordGuildRepository) Create(ctx context.Context, input domain.DiscordGuildCreateInput) (*domain.DiscordGuild, error) {
	now := time.Now()

	query := `
		INSERT INTO discord_guilds (
			id, owner_user_id, name, icon_url, member_count,
			bot_added_at, is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO UPDATE SET
			owner_user_id = COALESCE(EXCLUDED.owner_user_id, discord_guilds.owner_user_id),
			name = EXCLUDED.name,
			icon_url = EXCLUDED.icon_url,
			member_count = EXCLUDED.member_count,
			bot_added_at = EXCLUDED.bot_added_at,
			is_active = EXCLUDED.is_active,
			updated_at = EXCLUDED.updated_at
		RETURNING id, owner_user_id, name, icon_url, member_count,
		          bot_added_at, bot_removed_at, is_active, created_at, updated_at
	`

	var guild domain.DiscordGuild
	err := r.db.QueryRowContext(
		ctx, query,
		input.ID,
		input.OwnerUserID,
		input.Name,
		input.IconURL,
		input.MemberCount,
		now,
		true, // is_active
		now,  // created_at
		now,  // updated_at
	).Scan(
		&guild.ID,
		&guild.OwnerUserID,
		&guild.Name,
		&guild.IconURL,
		&guild.MemberCount,
		&guild.BotAddedAt,
		&guild.BotRemovedAt,
		&guild.IsActive,
		&guild.CreatedAt,
		&guild.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create discord guild: %w", err)
	}

	return &guild, nil
}

func (r *DiscordGuildRepository) GetByID(ctx context.Context, guildID int64) (*domain.DiscordGuild, error) {
	query := `
		SELECT id, owner_user_id, name, icon_url, member_count,
		       bot_added_at, bot_removed_at, is_active, created_at, updated_at
		FROM discord_guilds
		WHERE id = $1
	`

	var guild domain.DiscordGuild
	err := r.db.QueryRowContext(ctx, query, guildID).Scan(
		&guild.ID,
		&guild.OwnerUserID,
		&guild.Name,
		&guild.IconURL,
		&guild.MemberCount,
		&guild.BotAddedAt,
		&guild.BotRemovedAt,
		&guild.IsActive,
		&guild.CreatedAt,
		&guild.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get discord guild: %w", err)
	}

	return &guild, nil
}

func (r *DiscordGuildRepository) GetByOwner(ctx context.Context, userID int64) ([]domain.DiscordGuild, error) {
	query := `
		SELECT id, owner_user_id, name, icon_url, member_count,
		       bot_added_at, bot_removed_at, is_active, created_at, updated_at
		FROM discord_guilds
		WHERE owner_user_id = $1 AND is_active = true
		ORDER BY name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get guilds by owner: %w", err)
	}
	defer rows.Close()

	var guilds []domain.DiscordGuild
	for rows.Next() {
		var guild domain.DiscordGuild
		err := rows.Scan(
			&guild.ID,
			&guild.OwnerUserID,
			&guild.Name,
			&guild.IconURL,
			&guild.MemberCount,
			&guild.BotAddedAt,
			&guild.BotRemovedAt,
			&guild.IsActive,
			&guild.CreatedAt,
			&guild.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan guild: %w", err)
		}
		guilds = append(guilds, guild)
	}

	return guilds, nil
}

func (r *DiscordGuildRepository) GetAll(ctx context.Context) ([]domain.DiscordGuild, error) {
	query := `
		SELECT id, owner_user_id, name, icon_url, member_count,
		       bot_added_at, bot_removed_at, is_active, created_at, updated_at
		FROM discord_guilds
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all guilds: %w", err)
	}
	defer rows.Close()

	var guilds []domain.DiscordGuild
	for rows.Next() {
		var guild domain.DiscordGuild
		err := rows.Scan(
			&guild.ID,
			&guild.OwnerUserID,
			&guild.Name,
			&guild.IconURL,
			&guild.MemberCount,
			&guild.BotAddedAt,
			&guild.BotRemovedAt,
			&guild.IsActive,
			&guild.CreatedAt,
			&guild.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan guild: %w", err)
		}
		guilds = append(guilds, guild)
	}

	return guilds, nil
}

func (r *DiscordGuildRepository) GetActive(ctx context.Context) ([]domain.DiscordGuild, error) {
	query := `
		SELECT id, owner_user_id, name, icon_url, member_count,
		       bot_added_at, bot_removed_at, is_active, created_at, updated_at
		FROM discord_guilds
		WHERE is_active = true
		ORDER BY name ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get active guilds: %w", err)
	}
	defer rows.Close()

	var guilds []domain.DiscordGuild
	for rows.Next() {
		var guild domain.DiscordGuild
		err := rows.Scan(
			&guild.ID,
			&guild.OwnerUserID,
			&guild.Name,
			&guild.IconURL,
			&guild.MemberCount,
			&guild.BotAddedAt,
			&guild.BotRemovedAt,
			&guild.IsActive,
			&guild.CreatedAt,
			&guild.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan guild: %w", err)
		}
		guilds = append(guilds, guild)
	}

	return guilds, nil
}

func (r *DiscordGuildRepository) Update(ctx context.Context, guildID int64, input domain.DiscordGuildUpdateInput) error {
	query := `
		UPDATE discord_guilds
		SET name = COALESCE($2, name),
		    icon_url = COALESCE($3, icon_url),
		    member_count = COALESCE($4, member_count),
		    bot_added_at = COALESCE($5, bot_added_at),
		    bot_removed_at = COALESCE($6, bot_removed_at),
		    is_active = COALESCE($7, is_active),
		    updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.ExecContext(
		ctx, query,
		guildID,
		input.Name,
		input.IconURL,
		input.MemberCount,
		input.BotAddedAt,
		input.BotRemovedAt,
		input.IsActive,
	)

	if err != nil {
		return fmt.Errorf("failed to update discord guild: %w", err)
	}

	return nil
}

func (r *DiscordGuildRepository) MarkAsLeft(ctx context.Context, guildID int64) error {
	query := `
		UPDATE discord_guilds
		SET is_active = false,
		    bot_removed_at = NOW(),
		    updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, guildID)
	if err != nil {
		return fmt.Errorf("failed to mark guild as left: %w", err)
	}

	return nil
}

func (r *DiscordGuildRepository) Delete(ctx context.Context, guildID int64) error {
	query := `DELETE FROM discord_guilds WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, guildID)
	if err != nil {
		return fmt.Errorf("failed to delete discord guild: %w", err)
	}

	return nil
}

func (r *DiscordGuildRepository) UserOwnsGuild(ctx context.Context, userID int64, guildID int64) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM discord_guilds
			WHERE id = $1 AND owner_user_id = $2 AND is_active = true
		)
	`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, guildID, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check guild ownership: %w", err)
	}

	return exists, nil
}

func (r *DiscordGuildRepository) GetGuildCount(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM discord_guilds WHERE is_active = true`

	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get guild count: %w", err)
	}

	return count, nil
}