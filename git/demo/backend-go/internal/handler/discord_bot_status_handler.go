package handler

import (
	"encoding/json"
	"net/http"

    "github.com/gorilla/mux"
	"demo/backend-go/internal/service"
	"demo/backend-go/internal/domain"
)

type DiscordBotStatusHandler struct {
	guildService *service.DiscordGuildService
}

func NewDiscordBotStatusHandler(
	guildService *service.DiscordGuildService,
) *DiscordBotStatusHandler {
	return &DiscordBotStatusHandler{
		guildService: guildService,
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

	// Get guild count
	guildCount, err := h.guildService.GetGuildCount(r.Context())
	if err != nil {
		guildCount = 0
	}

	// Get all guilds (for detailed info)
	guilds, err := h.guildService.GetAllGuilds(r.Context())
	if err != nil {
		guilds = []domain.DiscordGuild{}
	}

	// Calculate total member count
	totalMembers := 0
	for _, guild := range guilds {
		if guild.MemberCount != nil {
			totalMembers += *guild.MemberCount
		}
	}

	response := map[string]interface{}{
		"running":      true, // We'll enhance this with actual bot status from Redis later
		"guilds":       guildCount,
		"totalMembers": totalMembers,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}