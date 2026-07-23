package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/gorilla/mux"

	"demo/backend-go/internal/infrastructure/redis"
	"demo/backend-go/internal/service"
)

type DiscordAuthHandler struct {
	authService  *service.DiscordAuthService
	guildService *service.DiscordGuildService
	sessionStore sessions.Store
	sessionName  string
	redisService *redis.RedisService
}

func NewDiscordAuthHandler(
	authService *service.DiscordAuthService,
	guildService *service.DiscordGuildService,
	sessionStore sessions.Store,
	sessionName string,
	redisService *redis.RedisService,
) *DiscordAuthHandler {
	return &DiscordAuthHandler{
		authService:  authService,
		guildService: guildService,
		sessionStore: sessionStore,
		sessionName:  sessionName,
		redisService: redisService,
	}
}

func (h *DiscordAuthHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/discord/auth/url", h.GetAuthURL)
	router.HandleFunc("/api/discord/auth/callback", h.HandleCallback)
	router.HandleFunc("/api/discord/auth/disconnect", h.Disconnect)
	router.HandleFunc("/api/discord/auth/status", h.GetStatus)
}

// GetAuthURL generates Discord OAuth URL
func (h *DiscordAuthHandler) GetAuthURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from session
	session, err := h.sessionStore.Get(r, h.sessionName)
	if err != nil {
		log.Printf("Failed to get session: %v", err)
		http.Error(w, "Unauthorized - Please login first", http.StatusUnauthorized)
		return
	}

	userID, ok := session.Values["user_id"].(string)
	if !ok || userID == "" {
		log.Printf("No user_id in session")
		http.Error(w, "Unauthorized - Please login first", http.StatusUnauthorized)
		return
	}

	log.Printf("Processing auth URL for user %s", userID)

	authURL, state, err := h.authService.GetAuthURL()
	if err != nil {
		log.Printf("Failed to generate auth URL: %v", err)
		http.Error(w, "Failed to generate auth URL", http.StatusInternalServerError)
		return
	}

	// Store state -> userID mapping in Redis (expires in 10 minutes)
	err = h.storeState(state, userID, 10*time.Minute)
	if err != nil {
		log.Printf("Failed to store state: %v", err)
		http.Error(w, "Failed to store state", http.StatusInternalServerError)
		return
	}

	log.Printf("Stored state %s for user %s", state, userID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"authUrl": authURL,
		"state":   state,
	})
}

// HandleCallback processes OAuth callback
func (h *DiscordAuthHandler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get state from URL
	state := r.URL.Query().Get("state")
	if state == "" {
		log.Printf("Missing state parameter")
		http.Error(w, "Missing state parameter", http.StatusBadRequest)
		return
	}

	// Get authorization code
	code := r.URL.Query().Get("code")
	if code == "" {
		log.Printf("Missing code parameter")
		http.Error(w, "Missing code parameter", http.StatusBadRequest)
		return
	}

	// Get userID from state
	userIDFromState, err := h.getUserFromState(state)
	if err != nil {
		log.Printf("Invalid or expired state: %s, error: %v", state, err)
		http.Error(w, "Invalid or expired state parameter", http.StatusBadRequest)
		return
	}

	log.Printf("Processing Discord callback for user %s", userIDFromState)

	// Handle callback
	discordUserID, err := h.authService.HandleCallback(r.Context(), code, userIDFromState)
	if err != nil {
		log.Printf("Failed to handle callback: %v", err)
		http.Error(w, "Discord-Verbindung fehlgeschlagen", http.StatusInternalServerError)
		return
	}

	// Clean up state
	h.deleteState(state)
	log.Printf("Successfully connected Discord for user %s", userIDFromState)

	// Backfill ownership for any guild the bot already joined before this
	// user connected their Discord account (GUILD_JOINED only knows the raw
	// Discord owner ID until this link exists).
	if h.guildService != nil {
		if err := h.guildService.LinkOwnedGuilds(r.Context(), discordUserID, userIDFromState); err != nil {
			log.Printf("Failed to link owned guilds for user %s: %v", userIDFromState, err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{
		"success": true,
	})
}

// Disconnect removes Discord connection
func (h *DiscordAuthHandler) Disconnect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from session
	session, err := h.sessionStore.Get(r, h.sessionName)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := session.Values["user_id"].(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Disconnect
	err = h.authService.Disconnect(r.Context(), userID)
	if err != nil {
		log.Printf("Failed to disconnect: %v", err)
		http.Error(w, "Failed to disconnect", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{
		"success": true,
	})
}

// GetStatus checks if user has Discord connected
func (h *DiscordAuthHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from session
	session, err := h.sessionStore.Get(r, h.sessionName)
	if err != nil {
		// Not logged in - return not connected
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"connected": false,
		})
		return
	}

	userID, ok := session.Values["user_id"].(string)
	if !ok || userID == "" {
		// Not logged in - return not connected
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"connected": false,
		})
		return
	}

	// Check connection
	connection, err := h.authService.GetConnection(r.Context(), userID)
	if err != nil {
		log.Printf("Failed to get connection: %v", err)
		http.Error(w, "Failed to get connection status", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"connected": connection != nil,
	}

	if connection != nil {
		response["discordUsername"] = connection.DiscordUsername
		response["discordDiscriminator"] = connection.DiscordDiscriminator
		response["discordUserId"] = connection.DiscordUserID
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Helper functions for state management

func (h *DiscordAuthHandler) storeState(state string, userID string, ttl time.Duration) error {
	if h.redisService == nil {
		return fmt.Errorf("redis service not available")
	}

	client := h.redisService.GetClient()
	if client == nil {
		return fmt.Errorf("redis client not available")
	}

	key := fmt.Sprintf("discord:oauth:state:%s", state)

	err := client.Set(context.Background(), key, userID, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to store state in redis: %w", err)
	}

	return nil
}

func (h *DiscordAuthHandler) getUserFromState(state string) (string, error) {
	if h.redisService == nil {
		return "", fmt.Errorf("redis service not available")
	}

	client := h.redisService.GetClient()
	if client == nil {
		return "", fmt.Errorf("redis client not available")
	}

	key := fmt.Sprintf("discord:oauth:state:%s", state)

	value, err := client.Get(context.Background(), key).Result()
	if err != nil {
		return "", fmt.Errorf("state not found or expired: %w", err)
	}

	if value == "" {
		return "", fmt.Errorf("invalid user ID in state")
	}

	return value, nil
}

func (h *DiscordAuthHandler) deleteState(state string) {
	if h.redisService == nil {
		return
	}

	client := h.redisService.GetClient()
	if client == nil {
		return
	}

	key := fmt.Sprintf("discord:oauth:state:%s", state)
	client.Del(context.Background(), key)
}