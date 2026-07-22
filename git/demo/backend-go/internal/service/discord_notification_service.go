package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"demo/backend-go/internal/infrastructure/redis"
	"demo/backend-go/internal/repository/postgres"
)

type DiscordNotificationService struct {
	settingsRepo *postgres.DiscordGuildSettingsRepository
	guildRepo    *postgres.DiscordGuildRepository
	redisService *redis.RedisService
}

func NewDiscordNotificationService(
	settingsRepo *postgres.DiscordGuildSettingsRepository,
	guildRepo *postgres.DiscordGuildRepository,
	redisService *redis.RedisService,
) *DiscordNotificationService {
	return &DiscordNotificationService{
		settingsRepo: settingsRepo,
		guildRepo:    guildRepo,
		redisService: redisService,
	}
}

func (s *DiscordNotificationService) SendTwitchCommandNotification(ctx context.Context, userID string, commandData map[string]interface{}) error {
	settings, err := s.settingsRepo.GetByUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user settings: %w", err)
	}

	for _, setting := range settings {
		if !setting.TwitchNotificationsEnabled || setting.CommandChannelID == nil {
			continue
		}

		guild, err := s.guildRepo.GetByID(ctx, setting.GuildID)
		if err != nil || guild == nil || !guild.IsActive {
			continue
		}

		embed := s.buildCommandEmbed(commandData)

		message := map[string]interface{}{
			"type":              "SEND_NOTIFICATION",
			"user_id":           userID,
			"guild_id":          setting.GuildID,
			"channel_id":        *setting.CommandChannelID,
			"notification_type": "twitch_command",
			"embed":             embed,
		}

		err = s.publishToDiscordBot(message)
		if err != nil {
			fmt.Printf("Warning: failed to send command notification to guild %d: %v\n", setting.GuildID, err)
		}
	}

	return nil
}

func (s *DiscordNotificationService) SendActivityNotification(ctx context.Context, userID string, activityType string, activityData map[string]interface{}) error {
	settings, err := s.settingsRepo.GetByUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user settings: %w", err)
	}

	for _, setting := range settings {
		if !setting.TwitchNotificationsEnabled || setting.ActivityChannelID == nil {
			continue
		}

		guild, err := s.guildRepo.GetByID(ctx, setting.GuildID)
		if err != nil || guild == nil || !guild.IsActive {
			continue
		}

		embed := s.buildActivityEmbed(activityType, activityData)

		message := map[string]interface{}{
			"type":              "SEND_NOTIFICATION",
			"user_id":           userID,
			"guild_id":          setting.GuildID,
			"channel_id":        *setting.ActivityChannelID,
			"notification_type": "stream_activity",
			"activity_type":     activityType,
			"embed":             embed,
		}

		err = s.publishToDiscordBot(message)
		if err != nil {
			fmt.Printf("Warning: failed to send activity notification to guild %d: %v\n", setting.GuildID, err)
		}
	}

	return nil
}

// SendClipNotification announces a newly created clip in every Discord guild
// the user has connected with a configured notification channel. Returns
// whether at least one channel was actually configured, so the caller can
// tell "nothing to do" apart from a real failure.
func (s *DiscordNotificationService) SendClipNotification(ctx context.Context, userID, clipTitle, clipURL, caption string, hashtags []string) (bool, error) {
	settings, err := s.settingsRepo.GetByUser(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user settings: %w", err)
	}

	posted := false
	for _, setting := range settings {
		if setting.NotificationChannelID == nil {
			continue
		}

		guild, err := s.guildRepo.GetByID(ctx, setting.GuildID)
		if err != nil || guild == nil || !guild.IsActive {
			continue
		}

		embed := s.buildClipEmbed(clipTitle, clipURL, caption, hashtags)

		message := map[string]interface{}{
			"type":              "SEND_NOTIFICATION",
			"user_id":           userID,
			"guild_id":          setting.GuildID,
			"channel_id":        *setting.NotificationChannelID,
			"notification_type": "new_clip",
			"embed":             embed,
		}

		if err := s.publishToDiscordBot(message); err != nil {
			fmt.Printf("Warning: failed to send clip notification to guild %d: %v\n", setting.GuildID, err)
			continue
		}
		posted = true
	}

	return posted, nil
}

func (s *DiscordNotificationService) buildClipEmbed(clipTitle, clipURL, caption string, hashtags []string) map[string]interface{} {
	description := caption
	if len(hashtags) > 0 {
		description += "\n\n" + strings.Join(hashtags, " ")
	}

	return map[string]interface{}{
		"title":       "🎬 Neuer Clip: " + clipTitle,
		"description": description,
		"url":         clipURL,
		"color":       9520895, // Twitch-Lila
	}
}

func (s *DiscordNotificationService) SendAdminContactNotification(ctx context.Context, contactData map[string]interface{}) error {
	adminChannelID := os.Getenv("DISCORD_ADMIN_CONTACT_CHANNEL")
	if adminChannelID == "" {
		return fmt.Errorf("admin contact channel not configured")
	}

	embed := s.buildContactEmbed(contactData)

	message := map[string]interface{}{
		"type":              "SEND_ADMIN_NOTIFICATION",
		"channel_id":        adminChannelID,
		"notification_type": "contact_form",
		"embed":             embed,
	}

	return s.publishToDiscordBot(message)
}

func (s *DiscordNotificationService) SendTestNotification(ctx context.Context, userID string, guildID int64, channelID int64) error {
	owns, err := s.guildRepo.UserOwnsGuild(ctx, userID, guildID)
	if err != nil {
		return fmt.Errorf("failed to check guild ownership: %w", err)
	}

	if !owns {
		return fmt.Errorf("user does not own this guild")
	}

	embed := map[string]interface{}{
		"title":       "🧪 Test Notification",
		"description": "This is a test notification from your Twitch Bot Dashboard!",
		"color":       5814783, // Blue color
		"fields": []map[string]interface{}{
			{
				"name":   "Status",
				"value":  "✅ Channel configured successfully!",
				"inline": false,
			},
			{
				"name":   "Guild ID",
				"value":  fmt.Sprintf("%d", guildID),
				"inline": true,
			},
			{
				"name":   "Channel ID",
				"value":  fmt.Sprintf("%d", channelID),
				"inline": true,
			},
		},
	}

	message := map[string]interface{}{
		"type":              "SEND_NOTIFICATION",
		"user_id":           userID,
		"guild_id":          guildID,
		"channel_id":        channelID,
		"notification_type": "test",
		"embed":             embed,
	}

	return s.publishToDiscordBot(message)
}

func (s *DiscordNotificationService) buildCommandEmbed(data map[string]interface{}) map[string]interface{} {
	trigger, _ := data["trigger"].(string)
	response, _ := data["response"].(string)

	return map[string]interface{}{
		"title":       "🎮 New Twitch Command",
		"description": "A new command has been created",
		"color":       9442302, // Purple
		"fields": []map[string]interface{}{
			{
				"name":   "Trigger",
				"value":  trigger,
				"inline": true,
			},
			{
				"name":   "Response",
				"value":  response,
				"inline": false,
			},
		},
	}
}

func (s *DiscordNotificationService) buildActivityEmbed(activityType string, data map[string]interface{}) map[string]interface{} {
	var title, color string
	var emoji string

	switch activityType {
	case "follow":
		emoji = "👥"
		title = "New Follower"
		color = "3066993" // Green
	case "subscription":
		emoji = "⭐"
		title = "New Subscriber"
		color = "15105570" // Gold
	case "bits":
		emoji = "💎"
		title = "Bits Cheered"
		color = "10181046" // Purple
	default:
		emoji = "📢"
		title = "Stream Activity"
		color = "5814783" // Blue
	}

	username, _ := data["username"].(string)

	embed := map[string]interface{}{
		"title": fmt.Sprintf("%s %s", emoji, title),
		"color": color,
		"fields": []map[string]interface{}{
			{
				"name":   "User",
				"value":  username,
				"inline": true,
			},
		},
	}

	if activityType == "subscription" {
		if tier, ok := data["tier"].(string); ok {
			embed["fields"] = append(embed["fields"].([]map[string]interface{}), map[string]interface{}{
				"name":   "Tier",
				"value":  tier,
				"inline": true,
			})
		}
	}

	if activityType == "bits" {
		if amount, ok := data["bits"].(int); ok {
			embed["fields"] = append(embed["fields"].([]map[string]interface{}), map[string]interface{}{
				"name":   "Amount",
				"value":  fmt.Sprintf("%d bits", amount),
				"inline": true,
			})
		}
	}

	if message, ok := data["message"].(string); ok && message != "" {
		embed["fields"] = append(embed["fields"].([]map[string]interface{}), map[string]interface{}{
			"name":   "Message",
			"value":  message,
			"inline": false,
		})
	}

	return embed
}

func (s *DiscordNotificationService) buildContactEmbed(data map[string]interface{}) map[string]interface{} {
	name, _ := data["name"].(string)
	email, _ := data["email"].(string)
	subject, _ := data["subject"].(string)
	message, _ := data["message"].(string)

	fields := []map[string]interface{}{
		{
			"name":   "👤 Name",
			"value":  name,
			"inline": false,
		},
		{
			"name":   "📧 Email",
			"value":  email,
			"inline": false,
		},
	}

	if phone, ok := data["phone"].(string); ok && phone != "" {
		fields = append(fields, map[string]interface{}{
			"name":   "📞 Phone",
			"value":  phone,
			"inline": false,
		})
	}

	fields = append(fields, map[string]interface{}{
		"name":   "📝 Subject",
		"value":  subject,
		"inline": false,
	})

	fields = append(fields, map[string]interface{}{
		"name":   "💬 Message",
		"value":  message,
		"inline": false,
	})

	return map[string]interface{}{
		"title":  "📩 New Contact Request",
		"color":  3447003, // Blue
		"fields": fields,
	}
}

func (s *DiscordNotificationService) publishToDiscordBot(message map[string]interface{}) error {
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