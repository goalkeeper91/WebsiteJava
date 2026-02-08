package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
	"demo/backend-go/internal/security"
)

var (
	ErrAuthTokenNotFound = errors.New("auth token nicht gefunden")
)

type authTokenRepository struct {
	db     *sql.DB
	crypto *security.Crypto
}

func NewAuthTokenRepository(db *sql.DB, crypto *security.Crypto) repository.AuthTokenRepository {
	return &authTokenRepository{
		db:     db,
		crypto: crypto,
	}
}

func (r *authTokenRepository) Create(ctx context.Context, token *domain.AuthToken) error {
	encryptedAccessToken, err := r.crypto.Encrypt(token.AccessToken)
	if err != nil {
		return fmt.Errorf("fehler beim Verschlüsseln des Access Tokens: %w", err)
	}

	encryptedRefreshToken, err := r.crypto.Encrypt(token.RefreshToken)
	if err != nil {
		return fmt.Errorf("fehler beim Verschlüsseln des Refresh Tokens: %w", err)
	}

	query := `
		INSERT INTO twitch_auth_tokens
		(twitch_user_id, user_name, access_token, refresh_token, token_type, expires_in, scope, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`

	err = r.db.QueryRowContext(
		ctx,
		query,
		token.TwitchUserID,
		token.Username,
		encryptedAccessToken,
		encryptedRefreshToken,
		token.TokenType,
		token.ExpiresIn,
		token.Scope,
		token.CreatedAt,
		token.UpdatedAt,
	).Scan(&token.ID)

	if err != nil {
		return fmt.Errorf("fehler beim Erstellen des Auth Tokens: %w", err)
	}

	return nil
}

func (r *authTokenRepository) GetByTwitchUserID(ctx context.Context, twitchUserID string) (*domain.AuthToken, error) {
	query := `
		SELECT id, twitch_user_id, user_name, access_token, refresh_token, token_type, expires_in, scope, created_at, updated_at
		FROM twitch_auth_tokens
		WHERE twitch_user_id = $1
	`

	token := &domain.AuthToken{}
	var encryptedAccessToken, encryptedRefreshToken string

	err := r.db.QueryRowContext(ctx, query, twitchUserID).Scan(
		&token.ID,
		&token.TwitchUserID,
		&token.Username,
		&encryptedAccessToken,
		&encryptedRefreshToken,
		&token.TokenType,
		&token.ExpiresIn,
		&token.Scope,
		&token.CreatedAt,
		&token.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrAuthTokenNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden des Auth Tokens: %w", err)
	}

	token.AccessToken, err = r.crypto.Decrypt(encryptedAccessToken)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Entschlüsseln des Access Tokens: %w", err)
	}

	token.RefreshToken, err = r.crypto.Decrypt(encryptedRefreshToken)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Entschlüsseln des Refresh Tokens: %w", err)
	}

	return token, nil
}

func (r *authTokenRepository) Update(ctx context.Context, token *domain.AuthToken) error {
	encryptedAccessToken, err := r.crypto.Encrypt(token.AccessToken)
	if err != nil {
		return fmt.Errorf("fehler beim Verschlüsseln des Access Tokens: %w", err)
	}

	encryptedRefreshToken, err := r.crypto.Encrypt(token.RefreshToken)
	if err != nil {
		return fmt.Errorf("fehler beim Verschlüsseln des Refresh Tokens: %w", err)
	}

	query := `
		UPDATE twitch_auth_tokens
		SET access_token = $2, refresh_token = $3, expires_in = $4, updated_at = $5
		WHERE twitch_user_id = $1
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		token.TwitchUserID,
		encryptedAccessToken,
		encryptedRefreshToken,
		token.ExpiresIn,
		token.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("fehler beim Update des Auth Tokens: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("fehler beim Abrufen der betroffenen Zeilen: %w", err)
	}

	if rowsAffected == 0 {
		return ErrAuthTokenNotFound
	}

	return nil
}

func (r *authTokenRepository) Delete(ctx context.Context, twitchUserID string) error {
	query := `DELETE FROM twitch_auth_tokens WHERE twitch_user_id = $1`

	result, err := r.db.ExecContext(ctx, query, twitchUserID)
	if err != nil {
		return fmt.Errorf("fehler beim Löschen des Auth Tokens: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("fehler beim Abrufen der betroffenen Zeilen: %w", err)
	}

	if rowsAffected == 0 {
		return ErrAuthTokenNotFound
	}

	return nil
}

func (r *authTokenRepository) Upsert(ctx context.Context, token *domain.AuthToken) error {
	encryptedAccessToken, err := r.crypto.Encrypt(token.AccessToken)
	if err != nil {
		return fmt.Errorf("fehler beim Verschlüsseln des Access Tokens: %w", err)
	}

	encryptedRefreshToken, err := r.crypto.Encrypt(token.RefreshToken)
	if err != nil {
		return fmt.Errorf("fehler beim Verschlüsseln des Refresh Tokens: %w", err)
	}

	query := `
		INSERT INTO twitch_auth_tokens
		(twitch_user_id, user_name, access_token, refresh_token, token_type, expires_in, scope, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (twitch_user_id)
		DO UPDATE SET
			access_token = EXCLUDED.access_token,
			refresh_token = EXCLUDED.refresh_token,
			expires_in = EXCLUDED.expires_in,
			updated_at = EXCLUDED.updated_at
		RETURNING id
	`

	err = r.db.QueryRowContext(
		ctx,
		query,
		token.TwitchUserID,
		token.Username,
		encryptedAccessToken,
		encryptedRefreshToken,
		token.TokenType,
		token.ExpiresIn,
		token.Scope,
		token.CreatedAt,
		token.UpdatedAt,
	).Scan(&token.ID)

	if err != nil {
		return fmt.Errorf("fehler beim Upsert des Auth Tokens: %w", err)
	}

	return nil
}