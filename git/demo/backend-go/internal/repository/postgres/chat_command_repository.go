package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
)

type chatCommandRepository struct {
	db *sql.DB
}

func NewChatCommandRepository(db *sql.DB) repository.ChatCommandRepository {
	return &chatCommandRepository{db: db}
}

func (r *chatCommandRepository) Create(ctx context.Context, command *domain.ChatCommand) error {
	query := `
		INSERT INTO twitch_chat_commands (channel_id, trigger, response, cooldown, enabled, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		command.ChannelID,
		command.Trigger,
		command.Response,
		command.Cooldown,
		command.Enabled,
		command.CreatedAt,
		command.UpdatedAt,
	).Scan(&command.ID)

	if err != nil {
		return fmt.Errorf("fehler beim Erstellen des Commands: %w", err)
	}

	return nil
}

func (r *chatCommandRepository) GetByID(ctx context.Context, id int64, channelID string) (*domain.ChatCommand, error) {
	query := `
		SELECT id, channel_id, trigger, response, cooldown, enabled, created_at, updated_at
		FROM twitch_chat_commands
		WHERE id = $1 AND channel_id = $2
	`

	command := &domain.ChatCommand{}
	err := r.db.QueryRowContext(ctx, query, id, channelID).Scan(
		&command.ID,
		&command.ChannelID,
		&command.Trigger,
		&command.Response,
		&command.Cooldown,
		&command.Enabled,
		&command.CreatedAt,
		&command.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrCommandNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden des Commands: %w", err)
	}

	return command, nil
}

func (r *chatCommandRepository) GetByTrigger(ctx context.Context, channelID, trigger string) (*domain.ChatCommand, error) {
	query := `
		SELECT id, channel_id, trigger, response, cooldown, enabled, created_at, updated_at
		FROM twitch_chat_commands
		WHERE channel_id = $1 AND LOWER(trigger) = LOWER($2)
	`

	command := &domain.ChatCommand{}
	err := r.db.QueryRowContext(ctx, query, channelID, trigger).Scan(
		&command.ID,
		&command.ChannelID,
		&command.Trigger,
		&command.Response,
		&command.Cooldown,
		&command.Enabled,
		&command.CreatedAt,
		&command.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrCommandNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden des Commands: %w", err)
	}

	return command, nil
}

func (r *chatCommandRepository) GetAll(ctx context.Context, channelID string, limit, offset int) ([]*domain.ChatCommand, int64, error) {
	var total int64
	err := r.db.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM twitch_chat_commands WHERE channel_id = $1",
		channelID,
	).Scan(&total)

	if err != nil {
		return nil, 0, fmt.Errorf("fehler beim Zählen der Commands: %w", err)
	}

	// Commands abfragen
	query := `
		SELECT id, channel_id, trigger, response, cooldown, enabled, created_at, updated_at
		FROM twitch_chat_commands
		WHERE channel_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, channelID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("fehler beim Laden der Commands: %w", err)
	}
	defer rows.Close()

	commands := make([]*domain.ChatCommand, 0)
	for rows.Next() {
		command := &domain.ChatCommand{}
		err := rows.Scan(
			&command.ID,
			&command.ChannelID,
			&command.Trigger,
			&command.Response,
			&command.Cooldown,
			&command.Enabled,
			&command.CreatedAt,
			&command.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("fehler beim Scannen des Commands: %w", err)
		}
		commands = append(commands, command)
	}

	return commands, total, nil
}

func (r *chatCommandRepository) Search(ctx context.Context, channelID, search string, limit, offset int) ([]*domain.ChatCommand, int64, error) {
	searchPattern := "%" + strings.ToLower(search) + "%"

	// Gesamtanzahl
	var total int64
	err := r.db.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM twitch_chat_commands WHERE channel_id = $1 AND LOWER(trigger) LIKE $2",
		channelID,
		searchPattern,
	).Scan(&total)

	if err != nil {
		return nil, 0, fmt.Errorf("fehler beim Zählen: %w", err)
	}

	// Commands suchen
	query := `
		SELECT id, channel_id, trigger, response, cooldown, enabled, created_at, updated_at
		FROM twitch_chat_commands
		WHERE channel_id = $1 AND LOWER(trigger) LIKE $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.QueryContext(ctx, query, channelID, searchPattern, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("fehler beim Suchen: %w", err)
	}
	defer rows.Close()

	commands := make([]*domain.ChatCommand, 0)
	for rows.Next() {
		command := &domain.ChatCommand{}
		err := rows.Scan(
			&command.ID,
			&command.ChannelID,
			&command.Trigger,
			&command.Response,
			&command.Cooldown,
			&command.Enabled,
			&command.CreatedAt,
			&command.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("fehler beim Scannen: %w", err)
		}
		commands = append(commands, command)
	}

	return commands, total, nil
}

func (r *chatCommandRepository) GetByStatus(ctx context.Context, channelID string, enabled bool, limit, offset int) ([]*domain.ChatCommand, int64, error) {
	// Gesamtanzahl
	var total int64
	err := r.db.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM twitch_chat_commands WHERE channel_id = $1 AND enabled = $2",
		channelID,
		enabled,
	).Scan(&total)

	if err != nil {
		return nil, 0, fmt.Errorf("fehler beim Zählen: %w", err)
	}

	// Commands laden
	query := `
		SELECT id, channel_id, trigger, response, cooldown, enabled, created_at, updated_at
		FROM twitch_chat_commands
		WHERE channel_id = $1 AND enabled = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.QueryContext(ctx, query, channelID, enabled, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("fehler beim Laden: %w", err)
	}
	defer rows.Close()

	commands := make([]*domain.ChatCommand, 0)
	for rows.Next() {
		command := &domain.ChatCommand{}
		err := rows.Scan(
			&command.ID,
			&command.ChannelID,
			&command.Trigger,
			&command.Response,
			&command.Cooldown,
			&command.Enabled,
			&command.CreatedAt,
			&command.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("fehler beim Scannen: %w", err)
		}
		commands = append(commands, command)
	}

	return commands, total, nil
}

func (r *chatCommandRepository) Update(ctx context.Context, command *domain.ChatCommand) error {
	query := `
		UPDATE twitch_chat_commands
		SET trigger = $2, response = $3, cooldown = $4, enabled = $5, updated_at = $6
		WHERE id = $1 AND channel_id = $7
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		command.ID,
		command.Trigger,
		command.Response,
		command.Cooldown,
		command.Enabled,
		command.UpdatedAt,
		command.ChannelID,
	)

	if err != nil {
		return fmt.Errorf("fehler beim Update: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("fehler beim Abrufen betroffener Zeilen: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrCommandNotFound
	}

	return nil
}

func (r *chatCommandRepository) Delete(ctx context.Context, id int64, channelID string) error {
	query := `DELETE FROM twitch_chat_commands WHERE id = $1 AND channel_id = $2`

	result, err := r.db.ExecContext(ctx, query, id, channelID)
	if err != nil {
		return fmt.Errorf("fehler beim Löschen: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("fehler beim Abrufen betroffener Zeilen: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrCommandNotFound
	}

	return nil
}

func (r *chatCommandRepository) Exists(ctx context.Context, channelID, trigger string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM twitch_chat_commands WHERE channel_id = $1 AND LOWER(trigger) = LOWER($2))`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, channelID, trigger).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("fehler beim Prüfen der Existenz: %w", err)
	}

	return exists, nil
}