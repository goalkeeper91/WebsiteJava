package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"demo/backend-go/internal/infrastructure/redis"
)

type BotAnnounceHandler struct {
	redisService   *redis.RedisService
	internalSecret string
}

func NewBotAnnounceHandler(redisService *redis.RedisService, internalSecret string) *BotAnnounceHandler {
	return &BotAnnounceHandler{
		redisService:   redisService,
		internalSecret: internalSecret,
	}
}

func (h *BotAnnounceHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/bot/announce", h.AnnounceMessage).Methods("POST", "OPTIONS")
}

type AnnounceRequest struct {
	Message   string `json:"message"`
	ChannelID string `json:"channel_id,omitempty"` // Optional: specific Twitch channel
}

// AnnounceMessage sends a message to Twitch chat via the bot.
// Called server-to-server by n8n/the bot, authenticated via a shared secret header
// instead of a browser session.
func (h *BotAnnounceHandler) AnnounceMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if !requireInternalSecret(w, r, h.internalSecret) {
		return
	}

	var req AnnounceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("❌ Invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Message == "" {
		http.Error(w, "Message is required", http.StatusBadRequest)
		return
	}

	if h.redisService == nil {
		log.Printf("❌ Redis Service not available")
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	// Send via existing Redis infrastructure
	if err := h.redisService.SendBotAnnounce(req.Message, req.ChannelID); err != nil {
       log.Printf("❌ Failed to send announcement: %v", err)
       http.Error(w, "Failed to send announcement", http.StatusInternalServerError)
       return
    }

	log.Printf("✅ Bot announcement queued: %s", req.Message)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Announcement sent to bot",
	})
}