package repository

import (
	"context"
	"time"

	"demo/backend-go/internal/domain"
)

// SocialConnectionRepository definiert alle Operationen für Social Connections
type SocialConnectionRepository interface {
	// Create erstellt eine neue Social Connection
	Create(ctx context.Context, connection *domain.SocialConnection) error

	// GetByID lädt eine Connection by ID
	GetByID(ctx context.Context, id int64) (*domain.SocialConnection, error)

	// GetByUserAndPlatform lädt eine Connection für einen User auf einer Platform
	GetByUserAndPlatform(ctx context.Context, userTwitchID string, platform domain.Platform) (*domain.SocialConnection, error)

	// GetAllByUser lädt alle Connections eines Users
	GetAllByUser(ctx context.Context, userTwitchID string) ([]*domain.SocialConnection, error)

	// GetActiveByUser lädt alle aktiven Connections eines Users
	GetActiveByUser(ctx context.Context, userTwitchID string) ([]*domain.SocialConnection, error)

	// Update aktualisiert eine Connection
	Update(ctx context.Context, connection *domain.SocialConnection) error

	// UpdateTokens aktualisiert nur die OAuth Tokens
	UpdateTokens(ctx context.Context, id int64, accessToken, refreshToken string, expiresAt *time.Time) error

	// Delete löscht eine Connection
	Delete(ctx context.Context, id int64) error

	// Deactivate deaktiviert eine Connection
	Deactivate(ctx context.Context, id int64) error

	// Activate aktiviert eine Connection
	Activate(ctx context.Context, id int64) error

	// GetExpiringSoon lädt alle Connections die bald ablaufen (für Token Refresh)
	GetExpiringSoon(ctx context.Context, minutes int) ([]*domain.SocialConnection, error)
}
