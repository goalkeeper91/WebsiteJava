package service

import (
	"context"
	"fmt"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/infrastructure/redis"
	"demo/backend-go/internal/repository"
)

type TwitchChannelService struct {
	channelRepo  repository.TwitchChannelRepository
	redisService *redis.RedisService
}

func NewTwitchChannelService(
	channelRepo repository.TwitchChannelRepository,
	redisService *redis.RedisService,
) *TwitchChannelService {
	return &TwitchChannelService{
		channelRepo:  channelRepo,
		redisService: redisService,
	}
}

func (s *TwitchChannelService) SyncChannel(ctx context.Context, twitchUserID, username string) error {
	channel := domain.NewTwitchChannel(twitchUserID, username)

	if err := s.channelRepo.Upsert(ctx, channel); err != nil {
		return fmt.Errorf("fehler beim Synchronisieren des Channels: %w", err)
	}

	if err := s.redisService.SendJoinChannelSignal(twitchUserID); err != nil {
		fmt.Printf("⚠️ Fehler beim Senden des Join-Signals: %v\n", err)
	}

	return nil
}

func (s *TwitchChannelService) GetChannel(ctx context.Context, twitchUserID string) (*domain.TwitchChannel, error) {
	return s.channelRepo.GetByTwitchUserID(ctx, twitchUserID)
}

func (s *TwitchChannelService) ActivateChannel(ctx context.Context, twitchUserID string) error {
	channel, err := s.channelRepo.GetByTwitchUserID(ctx, twitchUserID)
	if err != nil {
		return err
	}

	channel.Activate()

	if err := s.channelRepo.Update(ctx, channel); err != nil {
		return err
	}

	return s.redisService.SendJoinChannelSignal(twitchUserID)
}

func (s *TwitchChannelService) DeactivateChannel(ctx context.Context, twitchUserID string) error {
	channel, err := s.channelRepo.GetByTwitchUserID(ctx, twitchUserID)
	if err != nil {
		return err
	}

	channel.Deactivate()

	return s.channelRepo.Update(ctx, channel)
}