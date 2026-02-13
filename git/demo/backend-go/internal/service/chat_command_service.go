package service

import (
	"context"
	"fmt"
	"strconv"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/infrastructure/redis"
	"demo/backend-go/internal/repository"
)

type ChatCommandService struct {
	commandRepo repository.ChatCommandRepository
	channelRepo repository.TwitchChannelRepository
	redisService *redis.RedisService
}

func NewChatCommandService(
	commandRepo repository.ChatCommandRepository,
	channelRepo repository.TwitchChannelRepository,
	redisService *redis.RedisService,
) *ChatCommandService {
	return &ChatCommandService{
		commandRepo:  commandRepo,
		channelRepo:  channelRepo,
		redisService: redisService,
	}
}

func (s *ChatCommandService) GetCommands(ctx context.Context, twitchUserID string, limit, offset int) ([]*domain.ChatCommand, int64, error) {
	channelID, err := s.getChannelID(ctx, twitchUserID)
	if err != nil {
		return nil, 0, err
	}

	return s.commandRepo.GetAll(ctx, channelID, limit, offset)
}

func (s *ChatCommandService) SearchCommands(ctx context.Context, twitchUserID, search string, limit, offset int) ([]*domain.ChatCommand, int64, error) {
	channelID, err := s.getChannelID(ctx, twitchUserID)
	if err != nil {
		return nil, 0, err
	}

	return s.commandRepo.Search(ctx, channelID, search, limit, offset)
}

func (s *ChatCommandService) GetCommandsByStatus(ctx context.Context, twitchUserID string, enabled bool, limit, offset int) ([]*domain.ChatCommand, int64, error) {
	channelID, err := s.getChannelID(ctx, twitchUserID)
	if err != nil {
		return nil, 0, err
	}

	return s.commandRepo.GetByStatus(ctx, channelID, enabled, limit, offset)
}

func (s *ChatCommandService) GetCommandByID(ctx context.Context, twitchUserID string, id int64) (*domain.ChatCommand, error) {
	channelID, err := s.getChannelID(ctx, twitchUserID)
	if err != nil {
		return nil, err
	}

	return s.commandRepo.GetByID(ctx, id, channelID)
}

func (s *ChatCommandService) CreateCommand(ctx context.Context, twitchUserID, trigger, response string, cooldown int) (*domain.ChatCommand, error) {
	if err := domain.ValidateTrigger(trigger); err != nil {
		return nil, err
	}

	channelID, err := s.getChannelID(ctx, twitchUserID)
	if err != nil {
		return nil, err
	}

	exists, err := s.commandRepo.Exists(ctx, channelID, trigger)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Prüfen der Trigger-Existenz: %w", err)
	}
	if exists {
		return nil, domain.ErrCommandAlreadyExists
	}

	command := domain.NewChatCommand(channelID, trigger, response, cooldown)
	if err := s.commandRepo.Create(ctx, command); err != nil {
		return nil, fmt.Errorf("fehler beim Erstellen des Commands: %w", err)
	}

	if err := s.redisService.SendRefreshCommandsSignal(twitchUserID); err != nil {
		fmt.Printf("⚠️ Fehler beim Senden des Bot-Signals: %v\n", err)
	}

	return command, nil
}

func (s *ChatCommandService) UpdateCommand(
	ctx context.Context,
	twitchUserID string,
	id int64,
	trigger, response *string,
	cooldown *int,
	enabled *bool,
) (*domain.ChatCommand, error) {
	channelID, err := s.getChannelID(ctx, twitchUserID)
	if err != nil {
		return nil, err
	}

	command, err := s.commandRepo.GetByID(ctx, id, channelID)
	if err != nil {
		return nil, err
	}

	if trigger != nil && *trigger != "" {
		if err := domain.ValidateTrigger(*trigger); err != nil {
			return nil, err
		}

		normalizedTrigger := domain.NormalizeTrigger(*trigger)

		existingCmd, err := s.commandRepo.GetByTrigger(ctx, channelID, normalizedTrigger)
		if err != nil && err != domain.ErrCommandNotFound {
			return nil, err
		}
		if existingCmd != nil && existingCmd.ID != id {
			return nil, domain.ErrTriggerTaken
		}
	}

	command.Update(trigger, response, cooldown, enabled)

	if err := s.commandRepo.Update(ctx, command); err != nil {
		return nil, fmt.Errorf("fehler beim Aktualisieren des Commands: %w", err)
	}

	if err := s.redisService.SendRefreshCommandsSignal(twitchUserID); err != nil {
		fmt.Printf("⚠️ Fehler beim Senden des Bot-Signals: %v\n", err)
	}

	return command, nil
}

func (s *ChatCommandService) DeleteCommand(ctx context.Context, twitchUserID string, id int64) error {
	channelID, err := s.getChannelID(ctx, twitchUserID)
	if err != nil {
		return err
	}

	if err := s.commandRepo.Delete(ctx, id, channelID); err != nil {
		return err
	}

	if err := s.redisService.SendRefreshCommandsSignal(twitchUserID); err != nil {
		fmt.Printf("⚠️ Fehler beim Senden des Bot-Signals: %v\n", err)
	}

	return nil
}

func (s *ChatCommandService) ToggleCommand(ctx context.Context, twitchUserID string, id int64, enabled bool) (*domain.ChatCommand, error) {
	channelID, err := s.getChannelID(ctx, twitchUserID)
	if err != nil {
		return nil, err
	}

	command, err := s.commandRepo.GetByID(ctx, id, channelID)
	if err != nil {
		return nil, err
	}

	command.Toggle(enabled)

	if err := s.commandRepo.Update(ctx, command); err != nil {
		return nil, err
	}

	if err := s.redisService.SendRefreshCommandsSignal(twitchUserID); err != nil {
		fmt.Printf("⚠️ Fehler beim Senden des Bot-Signals: %v\n", err)
	}

	return command, nil
}

func (s *ChatCommandService) getChannelID(ctx context.Context, twitchUserID string) (string, error) {
	channel, err := s.channelRepo.GetByTwitchUserID(ctx, twitchUserID)
	if err != nil {
		if err == domain.ErrChannelNotFound {
			return "", domain.ErrChannelNotRegistered
		}
		return "", fmt.Errorf("fehler beim Laden des Channels: %w", err)
	}

	return strconv.FormatInt(channel.ID, 10), nil
}