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
	commandRepo     repository.ChatCommandRepository
	channelRepo     repository.TwitchChannelRepository
	redisService    *redis.RedisService
	subscriptionSvc *SubscriptionService
	n8nSvc          *N8NService
}

func NewChatCommandService(
	commandRepo repository.ChatCommandRepository,
	channelRepo repository.TwitchChannelRepository,
	redisService *redis.RedisService,
	subscriptionSvc *SubscriptionService,
	n8nSvc *N8NService,
) *ChatCommandService {
	return &ChatCommandService{
		commandRepo:     commandRepo,
		channelRepo:     channelRepo,
		redisService:    redisService,
		subscriptionSvc: subscriptionSvc,
		n8nSvc:          n8nSvc,
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

func (s *ChatCommandService) GetCommandsByType(ctx context.Context, twitchUserID string, cmdType domain.CommandType, limit, offset int) ([]*domain.ChatCommand, int64, error) {
	channelID, err := s.getChannelID(ctx, twitchUserID)
	if err != nil {
		return nil, 0, err
	}
	return s.commandRepo.GetByType(ctx, channelID, cmdType, limit, offset)
}

func (s *ChatCommandService) GetCommandByID(ctx context.Context, twitchUserID string, id int64) (*domain.ChatCommand, error) {
	channelID, err := s.getChannelID(ctx, twitchUserID)
	if err != nil {
		return nil, err
	}
	return s.commandRepo.GetByID(ctx, id, channelID)
}

// CreateCommand creates a simple text command (free)
func (s *ChatCommandService) CreateCommand(ctx context.Context, twitchUserID, trigger, response string, cooldown int) (*domain.ChatCommand, error) {
	if err := domain.ValidateTrigger(trigger); err != nil {
		return nil, err
	}

	channelID, err := s.getChannelID(ctx, twitchUserID)
	if err != nil {
		return nil, err
	}

	// Check command limit based on subscription
	commands, total, err := s.commandRepo.GetAll(ctx, channelID, 1, 0)
	_ = commands
	if err != nil {
		return nil, fmt.Errorf("fehler beim Zählen der Commands: %w", err)
	}
	if err := s.subscriptionSvc.CheckCommandLimit(ctx, twitchUserID, int(total)); err != nil {
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

	s.notifyBot(twitchUserID)
	return command, nil
}

// CreateAdvancedCommand creates an n8n-powered command (pro/premium)
func (s *ChatCommandService) CreateAdvancedCommand(ctx context.Context, twitchUserID, trigger, workflowID string, cooldown int) (*domain.ChatCommand, error) {
	if err := domain.ValidateTrigger(trigger); err != nil {
		return nil, err
	}

	// Check n8n subscription
	canUse, err := s.subscriptionSvc.CanUseN8N(ctx, twitchUserID)
	if err != nil {
		return nil, err
	}
	if !canUse {
		return nil, domain.ErrFeatureNotAvailable
	}

	// Check n8n is configured
	ready, err := s.n8nSvc.IsReady(ctx, twitchUserID)
	if err != nil {
		return nil, err
	}
	if !ready {
		return nil, domain.ErrN8NIntegrationNotReady
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

	command := domain.NewAdvancedChatCommand(channelID, trigger, workflowID, cooldown)
	if err := s.commandRepo.Create(ctx, command); err != nil {
		return nil, fmt.Errorf("fehler beim Erstellen des Advanced Commands: %w", err)
	}

	s.notifyBot(twitchUserID)
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

	s.notifyBot(twitchUserID)
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

	s.notifyBot(twitchUserID)
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

	s.notifyBot(twitchUserID)
	return command, nil
}

// HandleCommandExecution is called by the bot when a command is triggered
// Routes to simple handler or n8n based on command type
func (s *ChatCommandService) HandleCommandExecution(ctx context.Context, channelID, trigger string, payload N8NCommandPayload) (*N8NCommandResponse, error) {
	command, err := s.commandRepo.GetByTrigger(ctx, channelID, trigger)
	if err != nil {
		return nil, err
	}

	// Track usage
	_ = s.commandRepo.TrackUsage(ctx, command.ID)

	if command.IsSimple() {
		// Simple: format and return response directly
		response := formatSimpleResponse(command.Response, payload.Username)
		return &N8NCommandResponse{Response: response}, nil
	}

	if command.IsAdvanced() {
		// Advanced: forward to n8n
		// channelID here is the twitch_user_id of the channel owner
		return s.n8nSvc.TriggerCommandWebhook(ctx, channelID, payload)
	}

	return nil, fmt.Errorf("unbekannter command type: %s", command.CommandType)
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

func (s *ChatCommandService) notifyBot(twitchUserID string) {
	if err := s.redisService.SendRefreshCommandsSignal(twitchUserID); err != nil {
		fmt.Printf("⚠️ Fehler beim Senden des Bot-Signals: %v\n", err)
	}
}

func formatSimpleResponse(template, username string) string {
	// Replace {user} placeholder
	result := template
	for i := 0; i < len(result)-5; i++ {
		if result[i:i+6] == "{user}" {
			result = result[:i] + username + result[i+6:]
			break
		}
	}
	return result
}