package service

import (
	"context"
	"fmt"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository/postgres"
)

type DiscordSettingsService struct {
	settingsRepo *postgres.DiscordGuildSettingsRepository
	guildRepo    *postgres.DiscordGuildRepository
}

func NewDiscordSettingsService(
	settingsRepo *postgres.DiscordGuildSettingsRepository,
	guildRepo *postgres.DiscordGuildRepository,
) *DiscordSettingsService {
	return &DiscordSettingsService{
		settingsRepo: settingsRepo,
		guildRepo:    guildRepo,
	}
}

func (s *DiscordSettingsService) GetSettings(ctx context.Context, guildID int64, userID string) (*domain.DiscordGuildSettings, error) {
	owns, err := s.guildRepo.UserOwnsGuild(ctx, userID, guildID)
	if err != nil {
		return nil, fmt.Errorf("failed to check guild ownership: %w", err)
	}

	if !owns {
		return nil, fmt.Errorf("user does not own this guild")
	}

	settings, err := s.settingsRepo.GetByGuildAndUser(ctx, guildID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get settings: %w", err)
	}

	if settings == nil {
		settings, err = s.CreateDefaultSettings(ctx, guildID, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to create default settings: %w", err)
		}
	}

	return settings, nil
}

func (s *DiscordSettingsService) UpdateSettings(ctx context.Context, guildID int64, userID string, input domain.DiscordGuildSettingsUpdateInput) (*domain.DiscordGuildSettings, error) {
	owns, err := s.guildRepo.UserOwnsGuild(ctx, userID, guildID)
	if err != nil {
		return nil, fmt.Errorf("failed to check guild ownership: %w", err)
	}

	if !owns {
		return nil, fmt.Errorf("user does not own this guild")
	}

	existing, err := s.settingsRepo.GetByGuildAndUser(ctx, guildID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing settings: %w", err)
	}

	if existing == nil {
		createInput := domain.DiscordGuildSettingsCreateInput{
			GuildID:                    guildID,
			UserID:                     userID,
			NotificationChannelID:      input.NotificationChannelID,
			CommandChannelID:           input.CommandChannelID,
			ActivityChannelID:          input.ActivityChannelID,
			TwitchNotificationsEnabled: input.TwitchNotificationsEnabled != nil && *input.TwitchNotificationsEnabled,
			JoinToCreateEnabled:        input.JoinToCreateEnabled != nil && *input.JoinToCreateEnabled,
			AdminRoleID:                input.AdminRoleID,
		}

		return s.settingsRepo.Create(ctx, createInput)
	}

	err = s.settingsRepo.Update(ctx, guildID, userID, input)
	if err != nil {
		return nil, fmt.Errorf("failed to update settings: %w", err)
	}

	return s.settingsRepo.GetByGuildAndUser(ctx, guildID, userID)
}

func (s *DiscordSettingsService) CreateDefaultSettings(ctx context.Context, guildID int64, userID string) (*domain.DiscordGuildSettings, error) {
	input := domain.DiscordGuildSettingsCreateInput{
		GuildID:                    guildID,
		UserID:                     userID,
		NotificationChannelID:      nil,
		CommandChannelID:           nil,
		ActivityChannelID:          nil,
		TwitchNotificationsEnabled: true,  // Enabled by default
		JoinToCreateEnabled:        true,  // Enabled by default
		AdminRoleID:                nil,
	}

	return s.settingsRepo.Create(ctx, input)
}

func (s *DiscordSettingsService) DeleteSettings(ctx context.Context, guildID int64, userID string) error {
	owns, err := s.guildRepo.UserOwnsGuild(ctx, userID, guildID)
	if err != nil {
		return fmt.Errorf("failed to check guild ownership: %w", err)
	}

	if !owns {
		return fmt.Errorf("user does not own this guild")
	}

	return s.settingsRepo.DeleteByGuildAndUser(ctx, guildID, userID)
}

func (s *DiscordSettingsService) GetNotificationChannel(ctx context.Context, userID string, guildID int64) (*int64, error) {
	return s.settingsRepo.GetNotificationChannel(ctx, userID, guildID)
}

func (s *DiscordSettingsService) GetCommandChannel(ctx context.Context, userID string, guildID int64) (*int64, error) {
	return s.settingsRepo.GetCommandChannel(ctx, userID, guildID)
}

// ValidateChannelID validates that a channel ID belongs to the guild
// TODO: Implement by querying Discord bot for guild channels
func (s *DiscordSettingsService) ValidateChannelID(ctx context.Context, guildID int64, channelID int64) (bool, error) {
	// For now, assume valid
	// In production, this should:
	// 1. Request guild channels from Discord bot
	// 2. Check if channelID is in the list
	return true, nil
}

// ValidateRoleID validates that a role ID belongs to the guild
// TODO: Implement by querying Discord bot for guild roles
func (s *DiscordSettingsService) ValidateRoleID(ctx context.Context, guildID int64, roleID int64) (bool, error) {
	// For now, assume valid
	// In production, this should:
	// 1. Request guild roles from Discord bot
	// 2. Check if roleID is in the list
	return true, nil
}