package service

import (
	"context"
	"encoding/json"
	"log"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/infrastructure/redis"
	"demo/backend-go/internal/repository/postgres"
)

type DiscordBotEventListener struct {
	redisService *redis.RedisService
	guildRepo    *postgres.DiscordGuildRepository
}

func NewDiscordBotEventListener(
	redisService *redis.RedisService,
	guildRepo *postgres.DiscordGuildRepository,
) *DiscordBotEventListener {
	return &DiscordBotEventListener{
		redisService: redisService,
		guildRepo:    guildRepo,
	}
}

// Start starts listening for Discord bot events
func (l *DiscordBotEventListener) Start(ctx context.Context) error {
	client := l.redisService.GetClient()
	if client == nil {
		return nil // Redis not available, skip
	}

	pubsub := client.Subscribe(ctx, "discord:status")
	defer pubsub.Close()

	log.Println("✅ Started listening for Discord bot events on discord:status")

	// Listen for messages
	ch := pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			log.Println("Discord bot event listener stopped")
			return ctx.Err()
		case msg := <-ch:
			if msg == nil {
				continue
			}

			// Parse message
			var event map[string]interface{}
			if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
				log.Printf("Failed to parse bot event: %v", err)
				continue
			}

			// Handle event
			if err := l.handleEvent(ctx, event); err != nil {
				log.Printf("Failed to handle bot event: %v", err)
			}
		}
	}
}

func (l *DiscordBotEventListener) handleEvent(ctx context.Context, event map[string]interface{}) error {
	eventType, ok := event["type"].(string)
	if !ok {
		log.Printf("Missing event type in bot event")
		return nil
	}

	log.Printf("📨 Received bot event: %s", eventType)

	switch eventType {
	case "GUILD_JOINED":
		return l.handleGuildJoined(ctx, event)
	case "GUILD_LEFT":
		return l.handleGuildLeft(ctx, event)
	default:
		log.Printf("Unknown bot event type: %s", eventType)
	}

	return nil
}

func (l *DiscordBotEventListener) handleGuildJoined(ctx context.Context, event map[string]interface{}) error {
	// Extract guild data
	guildIDFloat, ok := event["guild_id"].(float64)
	if !ok {
		log.Printf("Invalid guild_id in GUILD_JOINED event")
		return nil
	}
	guildID := int64(guildIDFloat)

	guildName, _ := event["guild_name"].(string)
	iconURL, _ := event["icon_url"].(string)

	var memberCount *int
	if mc, ok := event["member_count"].(float64); ok {
		count := int(mc)
		memberCount = &count
	}

	log.Printf("🎉 Bot joined guild: %s (ID: %d)", guildName, guildID)

	// Try to get existing guild
	// Since GetByGuildID doesn't exist, we'll try GetByID with the guildID
	// Or we create it and handle duplicate key error
	input := domain.DiscordGuildCreateInput{
		ID:          guildID,
		OwnerUserID: nil, // Will be set when user connects
		Name:        guildName,
		IconURL:     &iconURL,
		MemberCount: memberCount,
	}

	guild, err := l.guildRepo.Create(ctx, input)
	if err != nil {
		// If it's a duplicate key error, the guild already exists
		// Try to update it instead
		if isUniqueViolation(err) {
			log.Printf("Guild %d already exists, updating...", guildID)

			// Update using the guild ID
			updateInput := domain.DiscordGuildUpdateInput{
				Name:        &guildName,
				IconURL:     &iconURL,
				MemberCount: memberCount,
				IsActive:    boolPtr(true),
			}

			err = l.guildRepo.Update(ctx, guildID, updateInput)
			if err != nil {
				log.Printf("Failed to update guild: %v", err)
				return err
			}

			log.Printf("✅ Updated existing guild: %s", guildName)
			return nil
		}

		log.Printf("Failed to create guild: %v", err)
		return err
	}

	log.Printf("✅ Created new guild: %s (ID: %d)", guildName, guild.ID)
	return nil
}

func (l *DiscordBotEventListener) handleGuildLeft(ctx context.Context, event map[string]interface{}) error {
	// Extract guild ID
	guildIDFloat, ok := event["guild_id"].(float64)
	if !ok {
		log.Printf("Invalid guild_id in GUILD_LEFT event")
		return nil
	}
	guildID := int64(guildIDFloat)

	log.Printf("👋 Bot left guild: %d", guildID)

	// Mark as inactive using MarkAsLeft
	err := l.guildRepo.MarkAsLeft(ctx, guildID)
	if err != nil {
		log.Printf("Failed to mark guild as left: %v", err)
		return err
	}

	log.Printf("✅ Marked guild as left: %d", guildID)
	return nil
}

// Helper functions
func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	// Check if error message contains "duplicate key" or "unique constraint"
	errMsg := err.Error()
	return contains(errMsg, "duplicate key") || contains(errMsg, "unique constraint")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
		containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func boolPtr(b bool) *bool {
	return &b
}