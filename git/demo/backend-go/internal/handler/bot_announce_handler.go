package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"demo/backend-go/internal/infrastructure/redis"
)

type BotAnnounceHandler struct {
	redisService *redis.RedisService
}

func NewBotAnnounceHandler(redisService *redis.RedisService) *BotAnnounceHandler {
	return &BotAnnounceHandler{
		redisService: redisService,
	}
}

func (h *BotAnnounceHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/bot/announce", h.AnnounceMessage).Methods("POST", "OPTIONS")
}

type AnnounceRequest struct {
	Message   string `json:"message"`
	ChannelID string `json:"channel_id,omitempty"` // Optional: specific Twitch channel
}

// AnnounceMessage sends a message to Twitch chat via the bot
func (h *BotAnnounceHandler) AnnounceMessage(w http.ResponseWriter, r *http.Request) {
	// CORS for n8n
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
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
	if err := h.sendBotAnnounce(req.Message, req.ChannelID); err != nil {
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

func (h *BotAnnounceHandler) sendBotAnnounce(message string, channelID string) error {
	// Use existing Redis publish infrastructure
	signal := redis.BotSignal{
		Type:         "BOT_ANNOUNCE",
		TwitchUserID: channelID, // Reuse this field for channel_id
	}

	// Since BotSignal doesn't have Message field, we need to publish raw JSON
	payload := map[string]interface{}{
		"type":    "BOT_ANNOUNCE",
		"message": message,
	}

	if channelID != "" {
		payload["channel_id"] = channelID
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := h.redisService.GetClient()
	if client == nil {
		return fmt.Errorf("redis client not available")
	}

	// Publish to bot:events channel (same as other bot signals)
	err = client.Publish(context.Background(), redis.BotEventsChannel, data).Err()
	if err != nil {
		return err
	}

	log.Printf("📤 Bot announce signal sent: %s", message)
	return nil
}