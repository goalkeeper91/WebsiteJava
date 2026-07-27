package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
)

var (
	ErrUserNotFound      = errors.New("user nicht gefunden")
	ErrUserAlreadyExists = errors.New("user existiert bereits")
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (twitch_id, username, email, discord_id, is_admin, is_bot, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.TwitchID,
		user.Username,
		user.Email,
		user.DiscordID,
		user.IsAdmin,
		user.IsBot,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("fehler beim Erstellen des Users: %w", err)
	}

	return nil
}

func (r *userRepository) GetByTwitchID(ctx context.Context, twitchID string) (*domain.User, error) {
	query := `
		SELECT twitch_id, username, email, discord_id, is_admin, is_bot, created_at, updated_at
		FROM users
		WHERE twitch_id = $1
	`

	user := &domain.User{}
	var discordID sql.NullString

	err := r.db.QueryRowContext(ctx, query, twitchID).Scan(
		&user.TwitchID,
		&user.Username,
		&user.Email,
		&discordID,
		&user.IsAdmin,
		&user.IsBot,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden des Users: %w", err)
	}

	if discordID.Valid {
		user.DiscordID = &discordID.String
	}

	return user, nil
}

func (r *userRepository) GetByDiscordID(ctx context.Context, discordID string) (*domain.User, error) {
	query := `
		SELECT twitch_id, username, email, discord_id, is_admin, is_bot, created_at, updated_at
		FROM users
		WHERE discord_id = $1
	`

	user := &domain.User{}
	var discordIDNull sql.NullString

	err := r.db.QueryRowContext(ctx, query, discordID).Scan(
		&user.TwitchID,
		&user.Username,
		&user.Email,
		&discordIDNull,
		&user.IsAdmin,
		&user.IsBot,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden des Users: %w", err)
	}

	if discordIDNull.Valid {
		user.DiscordID = &discordIDNull.String
	}

	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET username = $2, email = $3, discord_id = $4, is_admin = $5, is_bot = $6, updated_at = $7
		WHERE twitch_id = $1
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		user.TwitchID,
		user.Username,
		user.Email,
		user.DiscordID,
		user.IsAdmin,
		user.IsBot,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("fehler beim Update des Users: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("fehler beim Abrufen der betroffenen Zeilen: %w", err)
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, twitchID string) error {
	query := `DELETE FROM users WHERE twitch_id = $1`

	result, err := r.db.ExecContext(ctx, query, twitchID)
	if err != nil {
		return fmt.Errorf("fehler beim Löschen des Users: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("fehler beim Abrufen der betroffenen Zeilen: %w", err)
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *userRepository) Exists(ctx context.Context, twitchID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE twitch_id = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, twitchID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("fehler beim Prüfen der User-Existenz: %w", err)
	}

	return exists, nil
}

func (r *userRepository) GetBotUsers(ctx context.Context) ([]*domain.User, error) {
	query := `
		SELECT twitch_id, username, email, discord_id, is_admin, is_bot, created_at, updated_at
		FROM users
		WHERE is_bot = true
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Bot Users: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		var discordID sql.NullString

		err := rows.Scan(
			&user.TwitchID,
			&user.Username,
			&user.Email,
			&discordID,
			&user.IsAdmin,
			&user.IsBot,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("fehler beim Scannen: %w", err)
		}

		if discordID.Valid {
			user.DiscordID = &discordID.String
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *userRepository) GetAllNonBotUsers(ctx context.Context) ([]*domain.User, error) {
	query := `
		SELECT twitch_id, username, email, discord_id, is_admin, is_bot, created_at, updated_at
		FROM users
		WHERE is_bot = false
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Nutzer: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		var discordID sql.NullString

		err := rows.Scan(
			&user.TwitchID,
			&user.Username,
			&user.Email,
			&discordID,
			&user.IsAdmin,
			&user.IsBot,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("fehler beim Scannen: %w", err)
		}

		if discordID.Valid {
			user.DiscordID = &discordID.String
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *userRepository) GetFirstBotUser(ctx context.Context) (*domain.User, error) {
	query := `
		SELECT twitch_id, username, email, discord_id, is_admin, is_bot, created_at, updated_at
		FROM users
		WHERE is_bot = true
		ORDER BY created_at ASC
		LIMIT 1
	`

	user := &domain.User{}
	var discordID sql.NullString

	err := r.db.QueryRowContext(ctx, query).Scan(
		&user.TwitchID,
		&user.Username,
		&user.Email,
		&discordID,
		&user.IsAdmin,
		&user.IsBot,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden des Bot Users: %w", err)
	}

	if discordID.Valid {
		user.DiscordID = &discordID.String
	}

	return user, nil
}