package service

import (
	"context"
	"encoding/json"
	"log"

	"demo/backend-go/internal/infrastructure/redis"
)

// ClipTriggerListener subscribes to the clips:trigger Redis channel (published
// by the clip-detector worker, Phase B) and dispatches each trigger to
// ClipCreationService. Structurally mirrors DiscordBotEventListener.
type ClipTriggerListener struct {
	redisService *redis.RedisService
	creationSvc  *ClipCreationService
}

func NewClipTriggerListener(redisService *redis.RedisService, creationSvc *ClipCreationService) *ClipTriggerListener {
	return &ClipTriggerListener{
		redisService: redisService,
		creationSvc:  creationSvc,
	}
}

// Start listens for clip triggers until ctx is cancelled.
func (l *ClipTriggerListener) Start(ctx context.Context) error {
	client := l.redisService.GetClient()
	if client == nil {
		return nil // Redis nicht verfügbar, überspringen
	}

	pubsub := client.Subscribe(ctx, redis.ClipsTriggerChannel)
	defer pubsub.Close()

	log.Printf("✅ Started listening for clip triggers on %s", redis.ClipsTriggerChannel)

	ch := pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			log.Println("Clip trigger listener stopped")
			return ctx.Err()
		case msg := <-ch:
			if msg == nil {
				continue
			}

			var trigger redis.ClipTrigger
			if err := json.Unmarshal([]byte(msg.Payload), &trigger); err != nil {
				log.Printf("Failed to parse clip trigger: %v", err)
				continue
			}

			// Jeder Trigger läuft in eigener Goroutine, damit ein langsamer
			// Poll für einen User nicht Trigger für andere User blockiert.
			go l.creationSvc.ProcessTrigger(context.Background(), trigger.UserTwitchID, trigger.Reason)
		}
	}
}
