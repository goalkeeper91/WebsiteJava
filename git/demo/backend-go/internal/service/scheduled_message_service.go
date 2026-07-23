package service

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/infrastructure/redis"
	"demo/backend-go/internal/repository"
	"demo/backend-go/internal/twitch"
)

// dueMessagesBatchSize caps how many due messages one scheduler tick
// processes - keeps a single slow tick bounded, remaining due messages just
// get picked up on the next tick a few seconds later.
const dueMessagesBatchSize = 100

type ScheduledMessageService struct {
	messageRepo  repository.ScheduledMessageRepository
	channelRepo  repository.TwitchChannelRepository
	redisService *redis.RedisService
	appToken     *twitch.TwitchAppTokenClient
}

func NewScheduledMessageService(
	messageRepo repository.ScheduledMessageRepository,
	channelRepo repository.TwitchChannelRepository,
	redisService *redis.RedisService,
	appToken *twitch.TwitchAppTokenClient,
) *ScheduledMessageService {
	return &ScheduledMessageService{
		messageRepo:  messageRepo,
		channelRepo:  channelRepo,
		redisService: redisService,
		appToken:     appToken,
	}
}

func (s *ScheduledMessageService) GetMessages(ctx context.Context, twitchUserID string, limit, offset int) ([]*domain.ScheduledMessage, int64, error) {
	channelID, err := s.getChannelID(ctx, twitchUserID)
	if err != nil {
		return nil, 0, err
	}
	return s.messageRepo.GetAll(ctx, channelID, limit, offset)
}

func (s *ScheduledMessageService) GetMessageByID(ctx context.Context, twitchUserID string, id int64) (*domain.ScheduledMessage, error) {
	channelID, err := s.getChannelID(ctx, twitchUserID)
	if err != nil {
		return nil, err
	}
	return s.messageRepo.GetByID(ctx, id, channelID)
}

func (s *ScheduledMessageService) CreateMessage(ctx context.Context, twitchUserID, message string, intervalSeconds int) (*domain.ScheduledMessage, error) {
	if err := domain.ValidateScheduledMessage(message, intervalSeconds); err != nil {
		return nil, err
	}

	channelID, err := s.getChannelID(ctx, twitchUserID)
	if err != nil {
		return nil, err
	}

	scheduled := domain.NewScheduledMessage(channelID, message, intervalSeconds)
	if err := s.messageRepo.Create(ctx, scheduled); err != nil {
		return nil, fmt.Errorf("fehler beim Erstellen der automatisierten Nachricht: %w", err)
	}

	return scheduled, nil
}

func (s *ScheduledMessageService) UpdateMessage(
	ctx context.Context,
	twitchUserID string,
	id int64,
	message *string,
	intervalSeconds *int,
	enabled, onlyWhenLive *bool,
) (*domain.ScheduledMessage, error) {
	channelID, err := s.getChannelID(ctx, twitchUserID)
	if err != nil {
		return nil, err
	}

	scheduled, err := s.messageRepo.GetByID(ctx, id, channelID)
	if err != nil {
		return nil, err
	}

	newMessage := scheduled.Message
	if message != nil {
		newMessage = *message
	}
	newInterval := scheduled.IntervalSeconds
	if intervalSeconds != nil {
		newInterval = *intervalSeconds
	}
	if err := domain.ValidateScheduledMessage(newMessage, newInterval); err != nil {
		return nil, err
	}

	scheduled.Update(message, intervalSeconds, enabled, onlyWhenLive)

	if err := s.messageRepo.Update(ctx, scheduled); err != nil {
		return nil, fmt.Errorf("fehler beim Aktualisieren der automatisierten Nachricht: %w", err)
	}

	return scheduled, nil
}

func (s *ScheduledMessageService) DeleteMessage(ctx context.Context, twitchUserID string, id int64) error {
	channelID, err := s.getChannelID(ctx, twitchUserID)
	if err != nil {
		return err
	}
	return s.messageRepo.Delete(ctx, id, channelID)
}

func (s *ScheduledMessageService) ToggleMessage(ctx context.Context, twitchUserID string, id int64, enabled bool) (*domain.ScheduledMessage, error) {
	channelID, err := s.getChannelID(ctx, twitchUserID)
	if err != nil {
		return nil, err
	}

	scheduled, err := s.messageRepo.GetByID(ctx, id, channelID)
	if err != nil {
		return nil, err
	}

	scheduled.Toggle(enabled)

	if err := s.messageRepo.Update(ctx, scheduled); err != nil {
		return nil, err
	}

	return scheduled, nil
}

// RunDueMessages is called on every scheduler tick. It resolves each due
// message's internal channel_id to a Twitch user ID, batch-checks live
// status for messages that require it, posts via the existing bot-announce
// path for live (or gating-exempt) channels, and always advances the
// schedule - a skipped offline channel isn't "caught up" later, it just
// waits for its next regular slot.
func (s *ScheduledMessageService) RunDueMessages(ctx context.Context) {
	due, err := s.messageRepo.GetDue(ctx, dueMessagesBatchSize)
	if err != nil {
		log.Printf("⚠️ Fehler beim Laden fälliger automatisierter Nachrichten: %v", err)
		return
	}
	if len(due) == 0 {
		return
	}

	twitchUserIDs := make(map[int64]string, len(due))
	liveCheckIDs := make([]string, 0, len(due))

	for _, m := range due {
		internalID, err := strconv.ParseInt(m.ChannelID, 10, 64)
		if err != nil {
			log.Printf("⚠️ Ungültige channel_id bei automatisierter Nachricht %d: %v", m.ID, err)
			continue
		}
		channel, err := s.channelRepo.GetByID(ctx, internalID)
		if err != nil {
			log.Printf("⚠️ Channel für automatisierte Nachricht %d nicht gefunden: %v", m.ID, err)
			continue
		}
		twitchUserIDs[m.ID] = channel.TwitchUserID
		if m.OnlyWhenLive {
			liveCheckIDs = append(liveCheckIDs, channel.TwitchUserID)
		}
	}

	var liveStreams map[string]twitch.LiveStream
	if len(liveCheckIDs) > 0 {
		liveStreams, err = s.appToken.GetLiveStreams(ctx, liveCheckIDs)
		if err != nil {
			log.Printf("⚠️ Fehler beim Live-Status-Check: %v", err)
			liveStreams = map[string]twitch.LiveStream{}
		}
	}

	now := time.Now()
	for _, m := range due {
		twitchUserID, ok := twitchUserIDs[m.ID]
		if ok {
			_, isLive := liveStreams[twitchUserID]
			if !m.OnlyWhenLive || isLive {
				if err := s.redisService.SendBotAnnounce(m.Message, twitchUserID); err != nil {
					log.Printf("⚠️ Fehler beim Posten der automatisierten Nachricht %d: %v", m.ID, err)
				}
			}
		}

		nextSendAt := now.Add(time.Duration(m.IntervalSeconds) * time.Second)
		if err := s.messageRepo.MarkSent(ctx, m.ID, now, nextSendAt); err != nil {
			log.Printf("⚠️ Fehler beim Fortschreiben von Nachricht %d: %v", m.ID, err)
		}
	}
}

func (s *ScheduledMessageService) getChannelID(ctx context.Context, twitchUserID string) (string, error) {
	channel, err := s.channelRepo.GetByTwitchUserID(ctx, twitchUserID)
	if err != nil {
		if err == domain.ErrChannelNotFound {
			return "", domain.ErrChannelNotRegistered
		}
		return "", fmt.Errorf("fehler beim Laden des Channels: %w", err)
	}
	return strconv.FormatInt(channel.ID, 10), nil
}
