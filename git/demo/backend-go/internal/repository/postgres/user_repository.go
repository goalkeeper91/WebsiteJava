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

// NewUserRepository erstellt ein neues PostgreSQL User Repository
func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (twitch_id, username, email, discord_id, is_admin, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.TwitchID,
		user.Username,
		user.Email,
		user.DiscordID,
		user.IsAdmin,
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
		SELECT twitch_id, username, email, discord_id, is_admin, created_at, updated_at
		FROM users
		WHERE twitch_id = $1
	`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, twitchID).Scan(
		&user.TwitchID,
		&user.Username,
		&user.Email,
		&user.DiscordID,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden des Users: %w", err)
	}

	return user, nil
}

func (r *userRepository) GetByDiscordID(ctx context.Context, discordID string) (*domain.User, error) {
	query := `
		SELECT twitch_id, username, email, discord_id, is_admin, created_at, updated_at
		FROM users
		WHERE discord_id = $1
	`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, discordID).Scan(
		&user.TwitchID,
		&user.Username,
		&user.Email,
		&user.DiscordID,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden des Users: %w", err)
	}

	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET username = $2, email = $3, discord_id = $4, is_admin = $5, updated_at = $6
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