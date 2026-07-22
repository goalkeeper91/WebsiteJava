package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
	"demo/backend-go/internal/security"
)

var (
	ErrSocialConnectionNotFound = errors.New("social connection not found")
	ErrSocialConnectionExists   = errors.New("social connection already exists for this platform")
)

type socialConnectionRepository struct {
	db     *sql.DB
	crypto *security.Crypto
}

func NewSocialConnectionRepository(db *sql.DB, crypto *security.Crypto) repository.SocialConnectionRepository {
	return &socialConnectionRepository{
		db:     db,
		crypto: crypto,
	}
}

func (r *socialConnectionRepository) Create(ctx context.Context, connection *domain.SocialConnection) error {
	// Encrypt tokens
	encryptedAccessToken, err := r.crypto.Encrypt(connection.AccessToken)
	if err != nil {
		return fmt.Errorf("fehler beim Verschlüsseln des Access Tokens: %w", err)
	}

	encryptedRefreshToken := ""
	if connection.RefreshToken != "" {
		encryptedRefreshToken, err = r.crypto.Encrypt(connection.RefreshToken)
		if err != nil {
			return fmt.Errorf("fehler beim Verschlüsseln des Refresh Tokens: %w", err)
		}
	}

	query := `
		INSERT INTO social_connections
		(user_twitch_id, platform, platform_user_id, platform_username, 
		 access_token, refresh_token, token_type, expires_at, scope, 
		 is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id
	`

	err = r.db.QueryRowContext(
		ctx,
		query,
		connection.UserTwitchID,
		connection.Platform,
		connection.PlatformUserID,
		connection.PlatformUsername,
		encryptedAccessToken,
		encryptedRefreshToken,
		connection.TokenType,
		connection.ExpiresAt,
		connection.Scope,
		connection.IsActive,
		connection.CreatedAt,
		connection.UpdatedAt,
	).Scan(&connection.ID)

	if err != nil {
		if isUniqueViolation(err) {
			return ErrSocialConnectionExists
		}
		return fmt.Errorf("fehler beim Erstellen der Social Connection: %w", err)
	}

	return nil
}

func (r *socialConnectionRepository) GetByID(ctx context.Context, id int64) (*domain.SocialConnection, error) {
	query := `
		SELECT id, user_twitch_id, platform, platform_user_id, platform_username,
		       access_token, refresh_token, token_type, expires_at, scope,
		       is_active, created_at, updated_at, last_used_at
		FROM social_connections
		WHERE id = $1
	`

	var conn domain.SocialConnection
	var encryptedAccessToken, encryptedRefreshToken string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&conn.ID,
		&conn.UserTwitchID,
		&conn.Platform,
		&conn.PlatformUserID,
		&conn.PlatformUsername,
		&encryptedAccessToken,
		&encryptedRefreshToken,
		&conn.TokenType,
		&conn.ExpiresAt,
		&conn.Scope,
		&conn.IsActive,
		&conn.CreatedAt,
		&conn.UpdatedAt,
		&conn.LastUsedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrSocialConnectionNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Social Connection: %w", err)
	}

	// Decrypt tokens
	conn.AccessToken, err = r.crypto.Decrypt(encryptedAccessToken)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Entschlüsseln des Access Tokens: %w", err)
	}

	if encryptedRefreshToken != "" {
		conn.RefreshToken, err = r.crypto.Decrypt(encryptedRefreshToken)
		if err != nil {
			return nil, fmt.Errorf("fehler beim Entschlüsseln des Refresh Tokens: %w", err)
		}
	}

	return &conn, nil
}

func (r *socialConnectionRepository) GetByUserAndPlatform(ctx context.Context, userTwitchID string, platform domain.Platform) (*domain.SocialConnection, error) {
	query := `
		SELECT id, user_twitch_id, platform, platform_user_id, platform_username,
		       access_token, refresh_token, token_type, expires_at, scope,
		       is_active, created_at, updated_at, last_used_at
		FROM social_connections
		WHERE user_twitch_id = $1 AND platform = $2
	`

	var conn domain.SocialConnection
	var encryptedAccessToken, encryptedRefreshToken string

	err := r.db.QueryRowContext(ctx, query, userTwitchID, platform).Scan(
		&conn.ID,
		&conn.UserTwitchID,
		&conn.Platform,
		&conn.PlatformUserID,
		&conn.PlatformUsername,
		&encryptedAccessToken,
		&encryptedRefreshToken,
		&conn.TokenType,
		&conn.ExpiresAt,
		&conn.Scope,
		&conn.IsActive,
		&conn.CreatedAt,
		&conn.UpdatedAt,
		&conn.LastUsedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrSocialConnectionNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Social Connection: %w", err)
	}

	// Decrypt
	conn.AccessToken, _ = r.crypto.Decrypt(encryptedAccessToken)
	if encryptedRefreshToken != "" {
		conn.RefreshToken, _ = r.crypto.Decrypt(encryptedRefreshToken)
	}

	return &conn, nil
}

func (r *socialConnectionRepository) GetAllByUser(ctx context.Context, userTwitchID string) ([]*domain.SocialConnection, error) {
	return r.getByUserWithFilter(ctx, userTwitchID, "")
}

func (r *socialConnectionRepository) GetActiveByUser(ctx context.Context, userTwitchID string) ([]*domain.SocialConnection, error) {
	return r.getByUserWithFilter(ctx, userTwitchID, "AND is_active = true")
}

func (r *socialConnectionRepository) getByUserWithFilter(ctx context.Context, userTwitchID, filter string) ([]*domain.SocialConnection, error) {
	query := fmt.Sprintf(`
		SELECT id, user_twitch_id, platform, platform_user_id, platform_username,
		       access_token, refresh_token, token_type, expires_at, scope,
		       is_active, created_at, updated_at, last_used_at
		FROM social_connections
		WHERE user_twitch_id = $1 %s
		ORDER BY created_at DESC
	`, filter)

	rows, err := r.db.QueryContext(ctx, query, userTwitchID)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Social Connections: %w", err)
	}
	defer rows.Close()

	var connections []*domain.SocialConnection
	for rows.Next() {
		var conn domain.SocialConnection
		var encryptedAccessToken, encryptedRefreshToken string

		err := rows.Scan(
			&conn.ID,
			&conn.UserTwitchID,
			&conn.Platform,
			&conn.PlatformUserID,
			&conn.PlatformUsername,
			&encryptedAccessToken,
			&encryptedRefreshToken,
			&conn.TokenType,
			&conn.ExpiresAt,
			&conn.Scope,
			&conn.IsActive,
			&conn.CreatedAt,
			&conn.UpdatedAt,
			&conn.LastUsedAt,
		)
		if err != nil {
			return nil, err
		}

		// Decrypt
		conn.AccessToken, _ = r.crypto.Decrypt(encryptedAccessToken)
		if encryptedRefreshToken != "" {
			conn.RefreshToken, _ = r.crypto.Decrypt(encryptedRefreshToken)
		}

		connections = append(connections, &conn)
	}

	return connections, rows.Err()
}

func (r *socialConnectionRepository) Update(ctx context.Context, connection *domain.SocialConnection) error {
	encryptedAccessToken, err := r.crypto.Encrypt(connection.AccessToken)
	if err != nil {
		return fmt.Errorf("fehler beim Verschlüsseln: %w", err)
	}

	encryptedRefreshToken := ""
	if connection.RefreshToken != "" {
		encryptedRefreshToken, err = r.crypto.Encrypt(connection.RefreshToken)
		if err != nil {
			return fmt.Errorf("fehler beim Verschlüsseln: %w", err)
		}
	}

	query := `
		UPDATE social_connections
		SET platform_username = $2,
		    access_token = $3,
		    refresh_token = $4,
		    expires_at = $5,
		    scope = $6,
		    is_active = $7,
		    updated_at = $8,
		    last_used_at = $9
		WHERE id = $1
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		connection.ID,
		connection.PlatformUsername,
		encryptedAccessToken,
		encryptedRefreshToken,
		connection.ExpiresAt,
		connection.Scope,
		connection.IsActive,
		time.Now(),
		connection.LastUsedAt,
	)

	if err != nil {
		return fmt.Errorf("fehler beim Update: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrSocialConnectionNotFound
	}

	return nil
}

func (r *socialConnectionRepository) UpdateTokens(ctx context.Context, id int64, accessToken, refreshToken string, expiresAt *time.Time) error {
	encryptedAccessToken, err := r.crypto.Encrypt(accessToken)
	if err != nil {
		return err
	}

	encryptedRefreshToken := ""
	if refreshToken != "" {
		encryptedRefreshToken, err = r.crypto.Encrypt(refreshToken)
		if err != nil {
			return err
		}
	}

	query := `
		UPDATE social_connections
		SET access_token = $2, refresh_token = $3, expires_at = $4, updated_at = $5
		WHERE id = $1
	`

	_, err = r.db.ExecContext(ctx, query, id, encryptedAccessToken, encryptedRefreshToken, expiresAt, time.Now())
	return err
}

func (r *socialConnectionRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM social_connections WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrSocialConnectionNotFound
	}

	return nil
}

func (r *socialConnectionRepository) Deactivate(ctx context.Context, id int64) error {
	return r.setActive(ctx, id, false)
}

func (r *socialConnectionRepository) Activate(ctx context.Context, id int64) error {
	return r.setActive(ctx, id, true)
}

func (r *socialConnectionRepository) setActive(ctx context.Context, id int64, active bool) error {
	query := `UPDATE social_connections SET is_active = $2, updated_at = $3 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id, active, time.Now())
	return err
}

func (r *socialConnectionRepository) GetExpiringSoon(ctx context.Context, minutes int) ([]*domain.SocialConnection, error) {
	query := `
		SELECT id, user_twitch_id, platform, platform_user_id, platform_username,
		       access_token, refresh_token, token_type, expires_at, scope,
		       is_active, created_at, updated_at, last_used_at
		FROM social_connections
		WHERE is_active = true 
		  AND expires_at IS NOT NULL
		  AND expires_at <= $1
		ORDER BY expires_at ASC
	`

	expiryThreshold := time.Now().Add(time.Duration(minutes) * time.Minute)
	rows, err := r.db.QueryContext(ctx, query, expiryThreshold)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var connections []*domain.SocialConnection
	for rows.Next() {
		var conn domain.SocialConnection
		var encryptedAccessToken, encryptedRefreshToken string

		err := rows.Scan(
			&conn.ID, &conn.UserTwitchID, &conn.Platform, &conn.PlatformUserID,
			&conn.PlatformUsername, &encryptedAccessToken, &encryptedRefreshToken,
			&conn.TokenType, &conn.ExpiresAt, &conn.Scope, &conn.IsActive,
			&conn.CreatedAt, &conn.UpdatedAt, &conn.LastUsedAt,
		)
		if err != nil {
			continue
		}

		conn.AccessToken, _ = r.crypto.Decrypt(encryptedAccessToken)
		if encryptedRefreshToken != "" {
			conn.RefreshToken, _ = r.crypto.Decrypt(encryptedRefreshToken)
		}

		connections = append(connections, &conn)
	}

	return connections, rows.Err()
}

// Helper für unique violation check
func isUniqueViolation(err error) bool {
	// PostgreSQL unique violation error code: 23505
	return err != nil && (err.Error() == "pq: duplicate key value violates unique constraint" ||
		err.Error() == "UNIQUE constraint failed")
}
