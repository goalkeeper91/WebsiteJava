package service

import (
	"context"
	"encoding/json"
	"fmt"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/infrastructure/redis"
	"demo/backend-go/internal/repository/postgres"
)

type DiscordGuildService struct {
	guildRepo    *postgres.DiscordGuildRepository
	settingsRepo *postgres.DiscordGuildSettingsRepository
	redisService *redis.RedisService
}

func NewDiscordGuildService(
	guildRepo *postgres.DiscordGuildRepository,
	settingsRepo *postgres.DiscordGuildSettingsRepository,
	redisService *redis.RedisService,
) *DiscordGuildService {
	return &DiscordGuildService{
		guildRepo:    guildRepo,
		settingsRepo: settingsRepo,
		redisService: redisService,
	}
}

func (s *DiscordGuildService) GetUserGuilds(ctx context.Context, userID string) ([]domain.DiscordGuildWithSettings, error) {
	guilds, err := s.guildRepo.GetByOwner(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user guilds: %w", err)
	}

	result := make([]domain.DiscordGuildWithSettings, 0, len(guilds))
	for _, guild := range guilds {
		settings, err := s.settingsRepo.GetByGuildAndUser(ctx, guild.ID, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to get settings for guild %d: %w", guild.ID, err)
		}

		result = append(result, domain.DiscordGuildWithSettings{
			Guild:    guild,
			Settings: settings,
		})
	}

	return result, nil
}

func (s *DiscordGuildService) GetGuild(ctx context.Context, guildID int64, userID string) (*domain.DiscordGuild, error) {
	owns, err := s.guildRepo.UserOwnsGuild(ctx, userID, guildID)
	if err != nil {
		return nil, fmt.Errorf("failed to check guild ownership: %w", err)
	}

	if !owns {
		return nil, fmt.Errorf("user does not own this guild")
	}

	guild, err := s.guildRepo.GetByID(ctx, guildID)
	if err != nil {
		return nil, fmt.Errorf("failed to get guild: %w", err)
	}

	if guild == nil {
		return nil, fmt.Errorf("guild not found")
	}

	return guild, nil
}

func (s *DiscordGuildService) GetGuildDetails(ctx context.Context, guildID int64, userID string) (*domain.DiscordGuildDetails, error) {
	owns, err := s.guildRepo.UserOwnsGuild(ctx, userID, guildID)
	if err != nil {
		return nil, fmt.Errorf("failed to check guild ownership: %w", err)
	}

	if !owns {
		return nil, fmt.Errorf("user does not own this guild")
	}

	guild, err := s.guildRepo.GetByID(ctx, guildID)
	if err != nil {
		return nil, fmt.Errorf("failed to get guild: %w", err)
	}

	if guild == nil {
		return nil, fmt.Errorf("guild not found")
	}

	settings, err := s.settingsRepo.GetByGuildAndUser(ctx, guildID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get settings: %w", err)
	}

	channels, roles, err := s.requestGuildDataFromBot(guildID)
	if err != nil {
		fmt.Printf("Warning: failed to get guild data from bot: %v\n", err)
		channels = []domain.DiscordChannel{}
		roles = []domain.DiscordRole{}
	}

	return &domain.DiscordGuildDetails{
		Guild:    *guild,
		Channels: channels,
		Roles:    roles,
		Settings: settings,
	}, nil
}

func (s *DiscordGuildService) SyncGuild(ctx context.Context, guildID int64, userID string) error {
	owns, err := s.guildRepo.UserOwnsGuild(ctx, userID, guildID)
	if err != nil {
		return fmt.Errorf("failed to check guild ownership: %w", err)
	}

	if !owns {
		return fmt.Errorf("user does not own this guild")
	}

	message := map[string]interface{}{
		"type":     "SYNC_GUILD",
		"guild_id": guildID,
	}

	return s.publishToDiscordBot(message)
}

func (s *DiscordGuildService) RemoveBot(ctx context.Context, guildID int64, userID string) error {
	owns, err := s.guildRepo.UserOwnsGuild(ctx, userID, guildID)
	if err != nil {
		return fmt.Errorf("failed to check guild ownership: %w", err)
	}

	if !owns {
		return fmt.Errorf("user does not own this guild")
	}

	message := map[string]interface{}{
		"type":     "LEAVE_GUILD",
		"guild_id": guildID,
	}

	err = s.publishToDiscordBot(message)
	if err != nil {
		return fmt.Errorf("failed to send leave request: %w", err)
	}

	// Mark guild as left in database
	err = s.guildRepo.MarkAsLeft(ctx, guildID)
	if err != nil {
		return fmt.Errorf("failed to mark guild as left: %w", err)
	}

	return nil
}

func (s *DiscordGuildService) HandleGuildJoined(ctx context.Context, guildID int64, guildName string, iconURL *string, memberCount *int, ownerDiscordID *int64) error {
	var ownerUserID *string
	// TODO: Implement lookup from discord_connections if ownerDiscordID is provided

	input := domain.DiscordGuildCreateInput{
		ID:          guildID,
		OwnerUserID: ownerUserID,
		Name:        guildName,
		IconURL:     iconURL,
		MemberCount: memberCount,
	}

	_, err := s.guildRepo.Create(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create guild: %w", err)
	}

	return nil
}

func (s *DiscordGuildService) HandleGuildLeft(ctx context.Context, guildID int64) error {
	err := s.guildRepo.MarkAsLeft(ctx, guildID)
	if err != nil {
		return fmt.Errorf("failed to mark guild as left: %w", err)
	}

	return nil
}

// LinkOwnedGuilds backfills ownership for any guild the bot already joined
// whose raw Discord owner ID matches discordUserID but that hasn't been
// linked to a platform user yet. Called right after a user connects their
// Discord account, so guilds joined before that connection existed still
// show up on the dashboard.
func (s *DiscordGuildService) LinkOwnedGuilds(ctx context.Context, discordUserID int64, userID string) error {
	linked, err := s.guildRepo.LinkOwnerByDiscordID(ctx, discordUserID, userID)
	if err != nil {
		return fmt.Errorf("failed to link owned guilds: %w", err)
	}
	if linked > 0 {
		fmt.Printf("Linked %d guild(s) to user %s via discord owner id %d\n", linked, userID, discordUserID)
	}
	return nil
}

func (s *DiscordGuildService) GetAllGuilds(ctx context.Context) ([]domain.DiscordGuild, error) {
	return s.guildRepo.GetAll(ctx)
}

func (s *DiscordGuildService) GetGuildCount(ctx context.Context) (int, error) {
	return s.guildRepo.GetGuildCount(ctx)
}

func (s *DiscordGuildService) requestGuildDataFromBot(guildID int64) ([]domain.DiscordChannel, []domain.DiscordRole, error) {
	// TODO: Implement request-response pattern with Redis
	// For now, return empty slices
	// In production, this would:
	// 1. Send request to Redis
	// 2. Wait for response (with timeout)
	// 3. Parse and return channels/roles

	return []domain.DiscordChannel{}, []domain.DiscordRole{}, nil
}

func (s *DiscordGuildService) publishToDiscordBot(message map[string]interface{}) error {
	if s.redisService == nil {
		return fmt.Errorf("redis service not available")
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