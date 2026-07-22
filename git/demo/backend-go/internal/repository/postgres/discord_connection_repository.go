package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/security"
)

type DiscordConnectionRepository struct {
	db     *sql.DB
	crypto *security.Crypto
}

func NewDiscordConnectionRepository(db *sql.DB, crypto *security.Crypto) *DiscordConnectionRepository {
	return &DiscordConnectionRepository{
		db:     db,
		crypto: crypto,
	}
}

func (r *DiscordConnectionRepository) Create(ctx context.Context, input domain.DiscordConnectionCreateInput) (*domain.DiscordConnection, error) {
	encryptedAccess, err := r.crypto.Encrypt(input.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt access token: %w", err)
	}

	encryptedRefresh, err := r.crypto.Encrypt(input.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt refresh token: %w", err)
	}

	query := `
		INSERT INTO discord_connections (
			user_id, discord_user_id, discord_username, discord_discriminator,
			access_token, refresh_token, token_expires_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	var conn domain.DiscordConnection
	err = r.db.QueryRowContext(
		ctx, query,
		input.UserID,
		input.DiscordUserID,
		input.DiscordUsername,
		input.DiscordDiscriminator,
		encryptedAccess,
		encryptedRefresh,
		input.TokenExpiresAt,
	).Scan(&conn.ID, &conn.CreatedAt, &conn.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create discord connection: %w", err)
	}

	conn.UserID = input.UserID
	conn.DiscordUserID = input.DiscordUserID
	conn.DiscordUsername = input.DiscordUsername
	conn.DiscordDiscriminator = input.DiscordDiscriminator
	conn.TokenExpiresAt = input.TokenExpiresAt

	return &conn, nil
}

func (r *DiscordConnectionRepository) GetByUserID(ctx context.Context, userID string) (*domain.DiscordConnection, error) {
	query := `
		SELECT id, user_id, discord_user_id, discord_username, discord_discriminator,
		       access_token, refresh_token, token_expires_at, created_at, updated_at
		FROM discord_connections
		WHERE user_id = $1
	`

	var conn domain.DiscordConnection
	var encryptedAccess, encryptedRefresh string

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&conn.ID,
		&conn.UserID,
		&conn.DiscordUserID,
		&conn.DiscordUsername,
		&conn.DiscordDiscriminator,
		&encryptedAccess,
		&encryptedRefresh,
		&conn.TokenExpiresAt,
		&conn.CreatedAt,
		&conn.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get discord connection: %w", err)
	}

	conn.AccessToken, err = r.crypto.Decrypt(encryptedAccess)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt access token: %w", err)
	}

	conn.RefreshToken, err = r.crypto.Decrypt(encryptedRefresh)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt refresh token: %w", err)
	}

	return &conn, nil
}

func (r *DiscordConnectionRepository) GetByDiscordUserID(ctx context.Context, discordUserID int64) (*domain.DiscordConnection, error) {
	query := `
		SELECT id, user_id, discord_user_id, discord_username, discord_discriminator,
		       access_token, refresh_token, token_expires_at, created_at, updated_at
		FROM discord_connections
		WHERE discord_user_id = $1
	`

	var conn domain.DiscordConnection
	var encryptedAccess, encryptedRefresh string

	err := r.db.QueryRowContext(ctx, query, discordUserID).Scan(
		&conn.ID,
		&conn.UserID,
		&conn.DiscordUserID,
		&conn.DiscordUsername,
		&conn.DiscordDiscriminator,
		&encryptedAccess,
		&encryptedRefresh,
		&conn.TokenExpiresAt,
		&conn.CreatedAt,
		&conn.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get discord connection: %w", err)
	}

	// Decrypt tokens
	conn.AccessToken, err = r.crypto.Decrypt(encryptedAccess)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt access token: %w", err)
	}

	conn.RefreshToken, err = r.crypto.Decrypt(encryptedRefresh)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt refresh token: %w", err)
	}

	return &conn, nil
}

func (r *DiscordConnectionRepository) Update(ctx context.Context, userID string, input domain.DiscordConnectionUpdateInput) error {
	query := `
		UPDATE discord_connections
		SET discord_username = COALESCE($2, discord_username),
		    discord_discriminator = COALESCE($3, discord_discriminator),
		    access_token = COALESCE($4, access_token),
		    refresh_token = COALESCE($5, refresh_token),
		    token_expires_at = COALESCE($6, token_expires_at),
		    updated_at = NOW()
		WHERE user_id = $1
	`

	var encryptedAccess, encryptedRefresh *string

	if input.AccessToken != nil {
		encrypted, err := r.crypto.Encrypt(*input.AccessToken)
		if err != nil {
			return fmt.Errorf("failed to encrypt access token: %w", err)
		}
		encryptedAccess = &encrypted
	}

	if input.RefreshToken != nil {
		encrypted, err := r.crypto.Encrypt(*input.RefreshToken)
		if err != nil {
			return fmt.Errorf("failed to encrypt refresh token: %w", err)
		}
		encryptedRefresh = &encrypted
	}

	_, err := r.db.ExecContext(
		ctx, query,
		userID,
		input.DiscordUsername,
		input.DiscordDiscriminator,
		encryptedAccess,
		encryptedRefresh,
		input.TokenExpiresAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update discord connection: %w", err)
	}

	return nil
}

func (r *DiscordConnectionRepository) Delete(ctx context.Context, userID string) error {
	query := `DELETE FROM discord_connections WHERE user_id = $1`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete discord connection: %w", err)
	}

	return nil
}

func (r *DiscordConnectionRepository) IsTokenExpired(expiresAt time.Time) bool {
	return time.Now().Add(5 * time.Minute).After(expiresAt)
}