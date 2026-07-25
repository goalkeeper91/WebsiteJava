package service

import (
	"context"
	"fmt"
	"log"
	"strings"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/infrastructure/redis"
	"demo/backend-go/internal/repository"
)

type AutomodService struct {
	automodRepo  repository.AutomodRepository
	redisService *redis.RedisService
	loyaltyRepo  repository.LoyaltyRepository
}

func NewAutomodService(automodRepo repository.AutomodRepository, redisService *redis.RedisService, loyaltyRepo repository.LoyaltyRepository) *AutomodService {
	return &AutomodService{
		automodRepo:  automodRepo,
		redisService: redisService,
		loyaltyRepo:  loyaltyRepo,
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
// internal secret instead (same as bot_announce_handler.go). For channels
// with ExemptRegulars on, this also merges in the logins of viewers whose
// Loyalty points cross the channel's Regulars threshold - purely at
// response time, never persisted, so the streamer's own manually-curated
// exempt_users list in the DB is untouched.
func (s *AutomodService) GetAllEnabledSettingsForBot(ctx context.Context) ([]*domain.AutomodSettings, error) {
	settingsList, err := s.automodRepo.GetAllEnabledSettings(ctx)
	if err != nil {
		return nil, err
	}

	for _, settings := range settingsList {
		if !settings.ExemptRegulars {
			continue
		}
		regulars := s.regularsForChannel(ctx, settings.UserTwitchID)
		if len(regulars) == 0 {
			continue
		}
		settings.ExemptUsers = mergeExemptUsers(settings.ExemptUsers, regulars)
	}

	return settingsList, nil
}

// regularsForChannel returns the logins currently qualifying as Regulars for
// userTwitchID, or nil if Loyalty is disabled/has no threshold set.
func (s *AutomodService) regularsForChannel(ctx context.Context, userTwitchID string) []string {
	loyaltySettings, err := s.loyaltyRepo.GetSettings(ctx, userTwitchID)
	if err != nil {
		log.Printf("⚠️ Fehler beim Laden der Loyalty-Settings für Regulars-Merge (%s): %v", userTwitchID, err)
		return nil
	}
	if !loyaltySettings.Enabled || loyaltySettings.RegularsThreshold <= 0 {
		return nil
	}

	logins, err := s.loyaltyRepo.GetViewerLoginsAboveThreshold(ctx, userTwitchID, loyaltySettings.RegularsThreshold)
	if err != nil {
		log.Printf("⚠️ Fehler beim Laden der Regulars für %s: %v", userTwitchID, err)
		return nil
	}
	return logins
}

// mergeExemptUsers returns a new list combining existing with additions,
// deduplicated case-insensitively - existing is never mutated.
func mergeExemptUsers(existing domain.StringArray, additions []string) domain.StringArray {
	seen := make(map[string]struct{}, len(existing)+len(additions))
	merged := make(domain.StringArray, 0, len(existing)+len(additions))

	for _, login := range existing {
		key := strings.ToLower(login)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		merged = append(merged, login)
	}
	for _, login := range additions {
		key := strings.ToLower(login)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		merged = append(merged, login)
	}

	return merged
}

// RecordViolation is called by the bot right after it detects a violation -
// it advances the escalation counter, logs the event for the dashboard's
// violation history, and returns the timeout duration to apply.
func (s *AutomodService) RecordViolation(
	ctx context.Context,
	broadcasterTwitchID, offenderTwitchID, offenderName, reason, messageExcerpt string,
) (int, error) {
	count, err := s.automodRepo.RecordViolation(ctx, broadcasterTwitchID, offenderTwitchID)
	if err != nil {
		return 0, err
	}
	timeoutSeconds := domain.NextTimeoutSeconds(count)

	event := domain.NewAutomodEvent(broadcasterTwitchID, offenderTwitchID, offenderName, reason, messageExcerpt, timeoutSeconds)
	if err := s.automodRepo.LogEvent(ctx, event); err != nil {
		fmt.Printf("⚠️ Fehler beim Protokollieren des Automod-Events: %v\n", err)
	}

	return timeoutSeconds, nil
}

// GetEvents is dashboard-facing - the streamer's own violation history.
func (s *AutomodService) GetEvents(ctx context.Context, userTwitchID string, limit, offset int) ([]*domain.AutomodEvent, int64, error) {
	return s.automodRepo.GetEvents(ctx, userTwitchID, limit, offset)
}
