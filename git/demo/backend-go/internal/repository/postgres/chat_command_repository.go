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
		INSERT INTO twitch_chat_commands
		(channel_id, trigger, response, cooldown, enabled, command_type, n8n_workflow_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
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
		command.CommandType,
		command.N8NWorkflowID,
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
		SELECT id, channel_id, trigger, response, cooldown, enabled,
		       command_type, n8n_workflow_id, usage_count, last_used_at,
		       created_at, updated_at
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
		&command.CommandType,
		&command.N8NWorkflowID,
		&command.UsageCount,
		&command.LastUsedAt,
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
		SELECT id, channel_id, trigger, response, cooldown, enabled,
		       command_type, n8n_workflow_id, usage_count, last_used_at,
		       created_at, updated_at
		FROM twitch_chat_commands
		WHERE channel_id = $1 AND LOWER(trigger) = LOWER($2) AND enabled = true
	`

	command := &domain.ChatCommand{}
	err := r.db.QueryRowContext(ctx, query, channelID, trigger).Scan(
		&command.ID,
		&command.ChannelID,
		&command.Trigger,
		&command.Response,
		&command.Cooldown,
		&command.Enabled,
		&command.CommandType,
		&command.N8NWorkflowID,
		&command.UsageCount,
		&command.LastUsedAt,
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

	query := `
		SELECT id, channel_id, trigger, response, cooldown, enabled,
		       command_type, n8n_workflow_id, usage_count, last_used_at,
		       created_at, updated_at
		FROM twitch_chat_commands
		WHERE channel_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	return r.queryCommands(ctx, query, channelID, limit, offset)
}

func (r *chatCommandRepository) Search(ctx context.Context, channelID, search string, limit, offset int) ([]*domain.ChatCommand, int64, error) {
	searchPattern := "%" + strings.ToLower(search) + "%"

	var total int64
	err := r.db.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM twitch_chat_commands WHERE channel_id = $1 AND LOWER(trigger) LIKE $2",
		channelID, searchPattern,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("fehler beim Zählen: %w", err)
	}

	query := `
		SELECT id, channel_id, trigger, response, cooldown, enabled,
		       command_type, n8n_workflow_id, usage_count, last_used_at,
		       created_at, updated_at
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

	commands, err := scanCommands(rows)
	return commands, total, err
}

func (r *chatCommandRepository) GetByStatus(ctx context.Context, channelID string, enabled bool, limit, offset int) ([]*domain.ChatCommand, int64, error) {
	var total int64
	err := r.db.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM twitch_chat_commands WHERE channel_id = $1 AND enabled = $2",
		channelID, enabled,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("fehler beim Zählen: %w", err)
	}

	query := `
		SELECT id, channel_id, trigger, response, cooldown, enabled,
		       command_type, n8n_workflow_id, usage_count, last_used_at,
		       created_at, updated_at
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

	commands, err := scanCommands(rows)
	return commands, total, err
}

func (r *chatCommandRepository) GetByType(ctx context.Context, channelID string, commandType domain.CommandType, limit, offset int) ([]*domain.ChatCommand, int64, error) {
	var total int64
	err := r.db.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM twitch_chat_commands WHERE channel_id = $1 AND command_type = $2",
		channelID, commandType,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("fehler beim Zählen: %w", err)
	}

	query := `
		SELECT id, channel_id, trigger, response, cooldown, enabled,
		       command_type, n8n_workflow_id, usage_count, last_used_at,
		       created_at, updated_at
		FROM twitch_chat_commands
		WHERE channel_id = $1 AND command_type = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.QueryContext(ctx, query, channelID, commandType, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("fehler beim Laden nach Typ: %w", err)
	}
	defer rows.Close()

	commands, err := scanCommands(rows)
	return commands, total, err
}

func (r *chatCommandRepository) GetAdvancedCommands(ctx context.Context, channelID string) ([]*domain.ChatCommand, error) {
	query := `
		SELECT id, channel_id, trigger, response, cooldown, enabled,
		       command_type, n8n_workflow_id, usage_count, last_used_at,
		       created_at, updated_at
		FROM twitch_chat_commands
		WHERE channel_id = $1 AND command_type = 'advanced' AND enabled = true
		ORDER BY trigger ASC
	`

	rows, err := r.db.QueryContext(ctx, query, channelID)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Advanced Commands: %w", err)
	}
	defer rows.Close()

	return scanCommands(rows)
}

func (r *chatCommandRepository) Update(ctx context.Context, command *domain.ChatCommand) error {
	query := `
		UPDATE twitch_chat_commands
		SET trigger = $2, response = $3, cooldown = $4, enabled = $5,
		    command_type = $6, n8n_workflow_id = $7, updated_at = $8
		WHERE id = $1 AND channel_id = $9
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		command.ID,
		command.Trigger,
		command.Response,
		command.Cooldown,
		command.Enabled,
		command.CommandType,
		command.N8NWorkflowID,
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
	result, err := r.db.ExecContext(ctx,
		`DELETE FROM twitch_chat_commands WHERE id = $1 AND channel_id = $2`, id, channelID)
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
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM twitch_chat_commands WHERE channel_id = $1 AND LOWER(trigger) = LOWER($2))`,
		channelID, trigger,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("fehler beim Prüfen der Existenz: %w", err)
	}
	return exists, nil
}

func (r *chatCommandRepository) TrackUsage(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE twitch_chat_commands
		SET usage_count = usage_count + 1, last_used_at = NOW()
		WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("fehler beim Tracken der Nutzung: %w", err)
	}
	return nil
}

// =====================================
// HELPERS
// =====================================

func (r *chatCommandRepository) queryCommands(ctx context.Context, query string, args ...interface{}) ([]*domain.ChatCommand, int64, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("fehler beim Laden der Commands: %w", err)
	}
	defer rows.Close()

	commands, err := scanCommands(rows)
	return commands, int64(len(commands)), err
}

func scanCommands(rows *sql.Rows) ([]*domain.ChatCommand, error) {
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
			&command.CommandType,
			&command.N8NWorkflowID,
			&command.UsageCount,
			&command.LastUsedAt,
			&command.CreatedAt,
			&command.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("fehler beim Scannen des Commands: %w", err)
		}
		commands = append(commands, command)
	}
	return commands, nil
}