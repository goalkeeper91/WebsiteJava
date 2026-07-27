package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/infrastructure/redis"
	"demo/backend-go/internal/repository"
)

type ActivityService struct {
	activityRepo repository.StreamActivityRepository
}

func NewActivityService(activityRepo repository.StreamActivityRepository) *ActivityService {
	return &ActivityService{
		activityRepo: activityRepo,
	}
}

// eventMessageID is only set for EventSub-sourced activities (see
// ActivityWebhookHandler) - guards against Twitch's at-least-once delivery
// via the unique index on stream_activities.event_message_id. Redis-sourced
// activities (HandleActivityEvent below) pass nil, since that path has no
// Twitch message ID and no duplicate-delivery risk to guard against.
func (s *ActivityService) CreateActivity(
	ctx context.Context,
	twitchUserID string,
	activityType domain.ActivityType,
	username, displayName string,
	viewers, bits *int,
	tier, message *string,
	eventMessageID *string,
) (*domain.StreamActivity, error) {
	activity := domain.NewStreamActivity(twitchUserID, activityType, username, displayName)

	if viewers != nil {
		activity.SetViewers(*viewers)
	}
	if bits != nil {
		activity.SetBits(*bits)
	}
	if tier != nil {
		activity.SetTier(*tier)
	}
	if message != nil {
		activity.SetMessage(*message)
	}
	activity.EventMessageID = eventMessageID

	if err := s.activityRepo.Create(ctx, activity); err != nil {
		return nil, fmt.Errorf("fehler beim Erstellen der Activity: %w", err)
	}

	log.Printf("✅ Activity erstellt: type=%s, user=%s", activityType, username)
	return activity, nil
}

func (s *ActivityService) GetRecentActivities(ctx context.Context, twitchUserID string, limit int) ([]*domain.StreamActivity, error) {
	return s.activityRepo.GetRecent(ctx, twitchUserID, limit)
}

func (s *ActivityService) GetActivitiesByType(ctx context.Context, twitchUserID string, activityType domain.ActivityType, limit int) ([]*domain.StreamActivity, error) {
	return s.activityRepo.GetByType(ctx, twitchUserID, activityType, limit)
}

func (s *ActivityService) CleanupOldActivities(ctx context.Context) error {
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30).Unix()
	return s.activityRepo.DeleteOlderThan(ctx, thirtyDaysAgo)
}

type ActivityEventHandler struct {
	activityService *ActivityService
	// Später: WebSocket Handler für Broadcasting
}

func NewActivityEventHandler(activityService *ActivityService) *ActivityEventHandler {
	return &ActivityEventHandler{
		activityService: activityService,
	}
}

func (h *ActivityEventHandler) HandleActivityEvent(event *redis.BotEvent) error {
	ctx := context.Background()

	activityType := domain.ActivityType(event.ActivityType)

	_, err := h.activityService.CreateActivity(
		ctx,
		event.TwitchUserID,
		activityType,
		event.Username,
		event.DisplayName,
		event.Viewers,
		event.Bits,
		event.Tier,
		event.Message,
		nil,
	)

	if err != nil {
		return fmt.Errorf("fehler beim Verarbeiten der Activity: %w", err)
	}

	// TODO: WebSocket Broadcasting an Frontend
	log.Printf("📡 Activity Event verarbeitet: type=%s, user=%s", activityType, event.Username)

	return nil
}

func (h *ActivityEventHandler) HandleChatMessage(event *redis.BotEvent) error {
	log.Printf("💬 Chat Message empfangen von %s", event.Username)
	// TODO: WebSocket an Frontend weiterleiten
	return nil
}

func (h *ActivityEventHandler) HandleBotStatus(event *redis.BotEvent) error {
	log.Printf("🤖 Bot Status: %s", event.Status)
	// TODO: Status in DB speichern oder WebSocket broadcast
	return nil
}