package service

import (
	"context"
	"fmt"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/infrastructure/redis"
	"demo/backend-go/internal/repository"
)

type AutomodService struct {
	automodRepo  repository.AutomodRepository
	redisService *redis.RedisService
}

func NewAutomodService(automodRepo repository.AutomodRepository, redisService *redis.RedisService) *AutomodService {
	return &AutomodService{
		automodRepo:  automodRepo,
		redisService: redisService,
	}
}

// GetSettings is dashboard-facing (session auth handled at the handler) -
// userTwitchID is directly the settings row's PK, no channel-ID indirection
// needed since this is a 1:1-per-channel settings table like AutomationSettings.
func (s *AutomodService) GetSettings(ctx context.Context, userTwitchID string) (*domain.AutomodSettings, error) {
	return s.automodRepo.GetSettings(ctx, userTwitchID)
}

func (s *AutomodService) UpdateSettings(ctx context.Context, userTwitchID string, input domain.AutomodSettingsUpdateInput) (*domain.AutomodSettings, error) {
	settings, err := s.automodRepo.GetSettings(ctx, userTwitchID)
	if err != nil {
		return nil, err
	}

	settings.ApplyUpdate(input)

	if err := s.automodRepo.UpsertSettings(ctx, settings); err != nil {
		return nil, fmt.Errorf("fehler beim Speichern der Automod-Settings: %w", err)
	}

	if s.redisService != nil {
		if err := s.redisService.SendRefreshAutomodSignal(userTwitchID); err != nil {
			fmt.Printf("⚠️ Fehler beim Senden des Automod-Reload-Signals: %v\n", err)
		}
	}

	return settings, nil
}

// GetAllEnabledSettingsForBot is the bot-internal endpoint's data source -
// no user-session auth at this layer, the handler gates it via the shared
// internal secret instead (same as bot_announce_handler.go).
func (s *AutomodService) GetAllEnabledSettingsForBot(ctx context.Context) ([]*domain.AutomodSettings, error) {
	return s.automodRepo.GetAllEnabledSettings(ctx)
}

// RecordViolation is called by the bot right after it detects a violation -
// it returns the timeout duration to apply for this occurrence.
func (s *AutomodService) RecordViolation(ctx context.Context, broadcasterTwitchID, offenderTwitchID string) (int, error) {
	count, err := s.automodRepo.RecordViolation(ctx, broadcasterTwitchID, offenderTwitchID)
	if err != nil {
		return 0, err
	}
	return domain.NextTimeoutSeconds(count), nil
}
