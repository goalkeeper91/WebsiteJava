package service

import (
	"context"
	"encoding/json"
	"fmt"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/infrastructure/redis"
	"demo/backend-go/internal/repository/postgres"
)

type JoinToCreateService struct {
	jtcRepo      *postgres.JoinToCreateRepository
	guildRepo    *postgres.DiscordGuildRepository
	settingsRepo *postgres.DiscordGuildSettingsRepository
	redisService *redis.RedisService
}

func NewJoinToCreateService(
	jtcRepo *postgres.JoinToCreateRepository,
	guildRepo *postgres.DiscordGuildRepository,
	settingsRepo *postgres.DiscordGuildSettingsRepository,
	redisService *redis.RedisService,
) *JoinToCreateService {
	return &JoinToCreateService{
		jtcRepo:      jtcRepo,
		guildRepo:    guildRepo,
		settingsRepo: settingsRepo,
		redisService: redisService,
	}
}

func (s *JoinToCreateService) ListConfigs(ctx context.Context, userID string) ([]domain.JoinToCreateConfig, error) {
	return s.jtcRepo.GetByUser(ctx, userID)
}

func (s *JoinToCreateService) GetConfig(ctx context.Context, configID int64, userID string) (*domain.JoinToCreateConfig, error) {
	owns, err := s.jtcRepo.UserOwnsConfig(ctx, userID, configID)
	if err != nil {
		return nil, fmt.Errorf("failed to check config ownership: %w", err)
	}

	if !owns {
		return nil, fmt.Errorf("user does not own this config")
	}

	config, err := s.jtcRepo.GetByID(ctx, configID)
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	if config == nil {
		return nil, fmt.Errorf("config not found")
	}

	return config, nil
}

func (s *JoinToCreateService) CreateConfig(ctx context.Context, input domain.JoinToCreateConfigCreateInput, userID string) (*domain.JoinToCreateConfig, error) {
	if err := s.validateConfigInput(input); err != nil {
		return nil, err
	}

	owns, err := s.guildRepo.UserOwnsGuild(ctx, userID, input.GuildID)
	if err != nil {
		return nil, fmt.Errorf("failed to check guild ownership: %w", err)
	}

	if !owns {
		return nil, fmt.Errorf("user does not own this guild")
	}

	settings, err := s.settingsRepo.GetByGuildAndUser(ctx, input.GuildID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get guild settings: %w", err)
	}

	if settings != nil && !settings.JoinToCreateEnabled {
		return nil, fmt.Errorf("join-to-create is disabled for this guild")
	}

	input.UserID = userID

	config, err := s.jtcRepo.Create(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create config: %w", err)
	}

	err = s.notifyBotReloadConfigs(ctx, input.GuildID)
	if err != nil {
		fmt.Printf("Warning: failed to notify bot about config change: %v\n", err)
	}

	return config, nil
}

func (s *JoinToCreateService) UpdateConfig(ctx context.Context, configID int64, userID string, input domain.JoinToCreateConfigUpdateInput) (*domain.JoinToCreateConfig, error) {
	owns, err := s.jtcRepo.UserOwnsConfig(ctx, userID, configID)
	if err != nil {
		return nil, fmt.Errorf("failed to check config ownership: %w", err)
	}

	if !owns {
		return nil, fmt.Errorf("user does not own this config")
	}

	existing, err := s.jtcRepo.GetByID(ctx, configID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing config: %w", err)
	}

	if existing == nil {
		return nil, fmt.Errorf("config not found")
	}

	err = s.jtcRepo.Update(ctx, configID, input)
	if err != nil {
		return nil, fmt.Errorf("failed to update config: %w", err)
	}

	err = s.notifyBotReloadConfigs(ctx, existing.GuildID)
	if err != nil {
		fmt.Printf("Warning: failed to notify bot about config change: %v\n", err)
	}

	return s.jtcRepo.GetByID(ctx, configID)
}

func (s *JoinToCreateService) DeleteConfig(ctx context.Context, configID int64, userID string) error {
	owns, err := s.jtcRepo.UserOwnsConfig(ctx, userID, configID)
	if err != nil {
		return fmt.Errorf("failed to check config ownership: %w", err)
	}

	if !owns {
		return fmt.Errorf("user does not own this config")
	}

	existing, err := s.jtcRepo.GetByID(ctx, configID)
	if err != nil {
		return fmt.Errorf("failed to get existing config: %w", err)
	}

	if existing == nil {
		return fmt.Errorf("config not found")
	}

	err = s.jtcRepo.Delete(ctx, configID)
	if err != nil {
		return fmt.Errorf("failed to delete config: %w", err)
	}

	err = s.notifyBotReloadConfigs(ctx, existing.GuildID)
	if err != nil {
		fmt.Printf("Warning: failed to notify bot about config deletion: %v\n", err)
	}

	return nil
}

func (s *JoinToCreateService) GetGuildConfigs(ctx context.Context, guildID int64, userID string) ([]domain.JoinToCreateConfig, error) {
	owns, err := s.guildRepo.UserOwnsGuild(ctx, userID, guildID)
	if err != nil {
		return nil, fmt.Errorf("failed to check guild ownership: %w", err)
	}

	if !owns {
		return nil, fmt.Errorf("user does not own this guild")
	}

	return s.jtcRepo.GetByGuild(ctx, guildID)
}

func (s *JoinToCreateService) GetEnabledConfigs(ctx context.Context) ([]domain.JoinToCreateConfig, error) {
	return s.jtcRepo.GetEnabled(ctx)
}

func (s *JoinToCreateService) validateConfigInput(input domain.JoinToCreateConfigCreateInput) error {
	if input.GuildID == 0 {
		return fmt.Errorf("guild_id is required")
	}

	if input.JoinChannelID == 0 {
		return fmt.Errorf("join_channel_id is required")
	}

	if input.CategoryID == 0 {
		return fmt.Errorf("category_id is required")
	}

	if input.ChannelNamePrefix == "" {
		return fmt.Errorf("channel_name_prefix is required")
	}

	if len(input.ChannelNamePrefix) > 100 {
		return fmt.Errorf("channel_name_prefix must be 100 characters or less")
	}

	if input.UserLimit < 0 {
		return fmt.Errorf("user_limit must be non-negative")
	}

	if input.UserLimit > 99 {
		return fmt.Errorf("user_limit must be 99 or less")
	}

	return nil
}

// notifyBotReloadConfigs pushes the guild's full current trigger list to the
// bot. The bot's RELOAD_CONFIGS handler does a full replace (wipes its
// in-memory config for the guild, then repopulates from "triggers") rather
// than an incremental update, so every call must include every trigger, not
// just the one that changed - omitting "triggers" entirely (the previous
// bug here) silently wiped the guild's config to empty instead of loading it.
func (s *JoinToCreateService) notifyBotReloadConfigs(ctx context.Context, guildID int64) error {
	if s.redisService == nil {
		return fmt.Errorf("redis service not available")
	}

	configs, err := s.jtcRepo.GetByGuild(ctx, guildID)
	if err != nil {
		return fmt.Errorf("failed to get guild configs: %w", err)
	}

	triggers := make([]map[string]interface{}, 0, len(configs))
	for _, c := range configs {
		triggers = append(triggers, map[string]interface{}{
			"channel_id":          c.JoinChannelID,
			"category_id":         c.CategoryID,
			"channel_name_prefix": c.ChannelNamePrefix,
			"user_limit":          c.UserLimit,
			"private_channel":     c.PrivateChannel,
			"enabled":             c.Enabled,
		})
	}

	message := map[string]interface{}{
		"type":        "RELOAD_CONFIGS",
		"guild_id":    guildID,
		"config_type": "join_to_create",
		"triggers":    triggers,
	}

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	client := s.redisService.GetClient()
	if client == nil {
		return fmt.Errorf("redis client not available")
	}

	err = client.Publish(context.Background(), "discord:events", data).Err()
	if err != nil {
		return fmt.Errorf("failed to publish to redis: %w", err)
	}

	return nil
}