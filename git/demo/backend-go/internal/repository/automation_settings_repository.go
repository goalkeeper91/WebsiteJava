package repository

import (
	"context"
	"demo/backend-go/internal/domain"
)

// AutomationSettingsRepository definiert alle Operationen für Automation Settings
type AutomationSettingsRepository interface {
	// Create erstellt Settings für einen User (nur einmal pro User)
	Create(ctx context.Context, settings *domain.AutomationSettings) error

	// GetByUser lädt Settings für einen User
	GetByUser(ctx context.Context, userTwitchID string) (*domain.AutomationSettings, error)

	// Update aktualisiert Settings
	Update(ctx context.Context, settings *domain.AutomationSettings) error

	// Upsert erstellt oder aktualisiert Settings
	Upsert(ctx context.Context, settings *domain.AutomationSettings) error

	// Delete löscht Settings
	Delete(ctx context.Context, userTwitchID string) error

	// GetAllEnabled lädt alle aktiven Automations (für Batch Processing)
	GetAllEnabled(ctx context.Context) ([]*domain.AutomationSettings, error)

	// GetAllEnabledWithSubscription lädt alle aktiven Automations zusammen mit
	// der jeweiligen Subscription (inkl. Tier) und Admin-Status des Nutzers -
	// fuer die laufende Tier-Durchsetzung im Clip-Detector-Poll-Loop (ein
	// einzelner JOIN statt eines Service-Aufrufs pro Kandidat).
	GetAllEnabledWithSubscription(ctx context.Context) ([]*domain.AutomationSettingsWithSubscription, error)

	// GetAllEnabledWithAutoPost lädt alle Settings mit Auto-Post enabled
	GetAllEnabledWithAutoPost(ctx context.Context) ([]*domain.AutomationSettings, error)

	// Enable aktiviert Automation für einen User
	Enable(ctx context.Context, userTwitchID string) error

	// Disable deaktiviert Automation für einen User
	Disable(ctx context.Context, userTwitchID string) error

	// UpdateTargetPlatforms aktualisiert nur die Target Platforms
	UpdateTargetPlatforms(ctx context.Context, userTwitchID string, platforms []domain.Platform) error
}
