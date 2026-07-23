package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

    "github.com/gorilla/mux"
	"demo/backend-go/internal/service"
	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/infrastructure/redis"
)

type DiscordBotStatusHandler struct {
	guildService *service.DiscordGuildService
	redisService *redis.RedisService
}

func NewDiscordBotStatusHandler(
	guildService *service.DiscordGuildService,
	redisService *redis.RedisService,
) *DiscordBotStatusHandler {
	return &DiscordBotStatusHandler{
		guildService: guildService,
		redisService: redisService,
	}
}

func (h *DiscordBotStatusHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/discord/bot/status", h.GetBotStatus)
}

// GetBotStatus returns Discord bot status
func (h *DiscordBotStatusHandler) GetBotStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get all guilds, keeping only ones the bot is currently in
	allGuilds, err := h.guildService.GetAllGuilds(r.Context())
	if err != nil {
		allGuilds = []domain.DiscordGuild{}
	}

	type guildSummary struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	totalMembers := 0
	guildList := make([]guildSummary, 0, len(allGuilds))
	for _, guild := range allGuilds {
		if !guild.IsActive {
			continue
		}
		if guild.MemberCount != nil {
			totalMembers += *guild.MemberCount
		}
		guildList = append(guildList, guildSummary{
			ID:   strconv.FormatInt(guild.ID, 10),
			Name: guild.Name,
		})
	}

	response := map[string]interface{}{
		"running":      false,
		"guildCount":   len(guildList),
		"totalMembers": totalMembers,
		"guilds":       guildList,
	}

	if h.redisService != nil {
		if val, err := h.redisService.GetDiscordBotStatusData(r.Context()); err == nil && val != "" {
			var status struct {
				Online        bool   `json:"online"`
				GuildCount    int    `json:"guild_count"`
				UptimeSeconds int    `json:"uptime_seconds"`
				LastHeartbeat string `json:"last_heartbeat"`
			}
			if err := json.Unmarshal([]byte(val), &status); err == nil {
				// The bot writes this key with a 90s TTL on every 30s heartbeat,
				// so its mere presence already means the bot is alive and recent.
				response["running"] = status.Online
				response["uptimeSeconds"] = status.UptimeSeconds
			} else {
				log.Printf("Failed to parse discord bot status: %v", err)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}