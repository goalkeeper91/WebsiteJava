package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	BotEventsChannel     = "bot:events"
	BackendEventsChannel = "backend:events"
)

type BotSignal struct {
	Type         string `json:"type"`
	TwitchUserID string `json:"twitch_user_id"`
}

type BotEvent struct {
	Type         string  `json:"type"`
	TwitchUserID string  `json:"twitch_user_id"`
	ActivityType string  `json:"activity_type,omitempty"`
	Username     string  `json:"username,omitempty"`
	DisplayName  string  `json:"display_name,omitempty"`
	Viewers      *int    `json:"viewers,omitempty"`
	Bits         *int    `json:"bits,omitempty"`
	Tier         *string `json:"tier,omitempty"`
	Message      *string `json:"message,omitempty"`
	Status       string  `json:"status,omitempty"`
}

type RedisService struct {
	client  *redis.Client
	pubsub  *redis.PubSub
	handler EventHandler
}

type EventHandler interface {
	HandleActivityEvent(event *BotEvent) error
	HandleChatMessage(event *BotEvent) error
	HandleBotStatus(event *BotEvent) error
}

func NewRedisService(addr, password string, db int, handler EventHandler) (*RedisService, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis verbindungsfehler: %w", err)
	}

	pubsub := client.Subscribe(context.Background(), BackendEventsChannel)

	service := &RedisService{
		client:  client,
		pubsub:  pubsub,
		handler: handler,
	}

	go service.listenForMessages()

	log.Printf("✅ Redis Service verbunden und subscribed auf: %s", BackendEventsChannel)
	return service, nil
}

func (r *RedisService) GetClient() *redis.Client {
    return r.client
}

func (r *RedisService) SendBotControlSignal(command string) error {
    signal := BotSignal{
        Type: command,
    }
    return r.publish(signal)
}

func (r *RedisService) GetBotStatusData(ctx context.Context) (string, error) {
    return r.client.Get(ctx, "bot:status").Result()
}

func (r *RedisService) SendJoinChannelSignal(twitchUserID string) error {
	signal := BotSignal{
		Type:         "JOIN_CHANNEL",
		TwitchUserID: twitchUserID,
	}
	return r.publish(signal)
}

func (r *RedisService) SendRefreshCommandsSignal(twitchUserID string) error {
	signal := BotSignal{
		Type:         "REFRESH_COMMANDS",
		TwitchUserID: twitchUserID,
	}
	return r.publish(signal)
}

func (r *RedisService) publish(signal BotSignal) error {
	data, err := json.Marshal(signal)
	if err != nil {
		return fmt.Errorf("fehler beim serialisieren: %w", err)
	}

	ctx := context.Background()
	if err := r.client.Publish(ctx, BotEventsChannel, data).Err(); err != nil {
		return fmt.Errorf("fehler beim publishen: %w", err)
	}

	if signal.TwitchUserID != "" {
		log.Printf("📤 Signal gesendet: type=%s, user=%s", signal.Type, signal.TwitchUserID)
	} else {
		log.Printf("📤 Signal gesendet: type=%s", signal.Type)
	}
	return nil
}

func (r *RedisService) listenForMessages() {
	ch := r.pubsub.Channel()

	for msg := range ch {
		if err := r.handleMessage(msg.Payload); err != nil {
			log.Printf("❌ Fehler beim Verarbeiten der Nachricht: %v", err)
		}
	}
}

func (r *RedisService) handleMessage(payload string) error {
	var event BotEvent
	if err := json.Unmarshal([]byte(payload), &event); err != nil {
		return fmt.Errorf("fehler beim deserialisieren: %w", err)
	}

	log.Printf("📥 Event empfangen: type=%s", event.Type)

	switch event.Type {
	case "ACTIVITY":
		return r.handler.HandleActivityEvent(&event)
	case "CHAT_MESSAGE":
		return r.handler.HandleChatMessage(&event)
	case "BOT_STATUS":
		return r.handler.HandleBotStatus(&event)
	default:
		log.Printf("⚠️ Unbekannter Event-Typ: %s", event.Type)
		return nil
	}
}

func (r *RedisService) Close() error {
	if err := r.pubsub.Close(); err != nil {
		return err
	}
	return r.client.Close()
}