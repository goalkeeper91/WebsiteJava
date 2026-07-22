package service

import (
	"context"
	"errors"
	"fmt"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
	"demo/backend-go/internal/repository/postgres"
)

const clipAutomationFeature = "clip_automation"

type AutomationSettingsService struct {
	settingsRepo    repository.AutomationSettingsRepository
	subscriptionSvc *SubscriptionService
}

func NewAutomationSettingsService(
	settingsRepo repository.AutomationSettingsRepository,
	subscriptionSvc *SubscriptionService,
) *AutomationSettingsService {
	return &AutomationSettingsService{
		settingsRepo:    settingsRepo,
		subscriptionSvc: subscriptionSvc,
	}
}

// GetSettings lädt die Automation Settings eines Users und legt beim ersten
// Aufruf Standard-Settings an, sofern das Feature im aktuellen Plan verfügbar ist.
func (s *AutomationSettingsService) GetSettings(ctx context.Context, userID string) (*domain.AutomationSettings, error) {
	canUse, err := s.subscriptionSvc.CanUseFeature(ctx, userID, clipAutomationFeature)
	if err != nil {
		return nil, err
	}
	if !canUse {
		return nil, domain.ErrFeatureNotAvailable
	}

	settings, err := s.settingsRepo.GetByUser(ctx, userID)
	if err != nil {
		if errors.Is(err, postgres.ErrAutomationSettingsNotFound) {
			settings = domain.NewAutomationSettings(userID)
			if err := s.settingsRepo.Create(ctx, settings); err != nil {
				return nil, fmt.Errorf("fehler beim Erstellen der Standard-Einstellungen: %w", err)
			}
			return settings, nil
		}
		return nil, fmt.Errorf("fehler beim Laden der Automation Settings: %w", err)
	}

	return settings, nil
}

// UpdateSettings aktualisiert die Automation Settings eines Users.
func (s *AutomationSettingsService) UpdateSettings(ctx context.Context, userID string, input domain.AutomationSettingsUpdateInput) (*domain.AutomationSettings, error) {
	settings, err := s.GetSettings(ctx, userID)
	if err != nil {
		return nil, err
	}

	settings.ApplyUpdate(input)

	if err := settings.Validate(); err != nil {
		return nil, err
	}

	if err := s.settingsRepo.Update(ctx, settings); err != nil {
		return nil, fmt.Errorf("fehler beim Aktualisieren der Automation Settings: %w", err)
	}

	return settings, nil
}
