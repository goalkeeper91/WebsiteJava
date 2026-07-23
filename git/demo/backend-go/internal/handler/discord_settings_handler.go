package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
    "github.com/gorilla/mux"
	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/service"
)

type DiscordSettingsHandler struct {
	settingsService     *service.DiscordSettingsService
	notificationService *service.DiscordNotificationService
	sessionStore        sessions.Store
	sessionName         string
}

func NewDiscordSettingsHandler(
	settingsService *service.DiscordSettingsService,
	notificationService *service.DiscordNotificationService,
	sessionStore sessions.Store,
	sessionName string,
) *DiscordSettingsHandler {
	return &DiscordSettingsHandler{
		settingsService:     settingsService,
		notificationService: notificationService,
		sessionStore:        sessionStore,
		sessionName:         sessionName,
	}
}

func (h *DiscordSettingsHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/discord/guilds/{id}/settings", h.handleSettingsRoute)
	router.HandleFunc("/api/discord/guilds/{id}/test-notification", h.handleTestNotificationRoute).Methods(http.MethodPost)
}

// handleSettingsRoute and handleTestNotificationRoute replace a single
// catch-all that was registered as the plain string "/api/discord/guilds/"
// with manual strings.Split parsing - gorilla/mux only matches that
// literally, so every nested path 404'd before ever reaching the handler.
// Path variables ({id}) fix that.
func (h *DiscordSettingsHandler) handleSettingsRoute(w http.ResponseWriter, r *http.Request) {
	guildID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid guild ID", http.StatusBadRequest)
		return
	}
	h.handleSettings(w, r, guildID)
}

func (h *DiscordSettingsHandler) handleTestNotificationRoute(w http.ResponseWriter, r *http.Request) {
	guildID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid guild ID", http.StatusBadRequest)
		return
	}
	h.SendTestNotification(w, r, guildID)
}

// handleSettings handles settings CRUD
func (h *DiscordSettingsHandler) handleSettings(w http.ResponseWriter, r *http.Request, guildID int64) {
	switch r.Method {
	case http.MethodGet:
		h.GetSettings(w, r, guildID)
	case http.MethodPut:
		h.UpdateSettings(w, r, guildID)
	case http.MethodDelete:
		h.DeleteSettings(w, r, guildID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetSettings gets settings for a guild
func (h *DiscordSettingsHandler) GetSettings(w http.ResponseWriter, r *http.Request, guildID int64) {
	// Get user ID from session
	session, err := h.sessionStore.Get(r, h.sessionName)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := session.Values["user_id"].(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get settings
	settings, err := h.settingsService.GetSettings(r.Context(), guildID, userID)
	if err != nil {
		log.Printf("Failed to get settings: %v", err)
		http.Error(w, "Anfrage fehlgeschlagen", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}

// UpdateSettings updates settings for a guild
func (h *DiscordSettingsHandler) UpdateSettings(w http.ResponseWriter, r *http.Request, guildID int64) {
	// Get user ID from session
	session, err := h.sessionStore.Get(r, h.sessionName)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := session.Values["user_id"].(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse request body
	var input domain.DiscordGuildSettingsUpdateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update settings
	settings, err := h.settingsService.UpdateSettings(r.Context(), guildID, userID, input)
	if err != nil {
		log.Printf("Failed to update settings: %v", err)
		http.Error(w, "Anfrage fehlgeschlagen", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}

// DeleteSettings deletes settings for a guild
func (h *DiscordSettingsHandler) DeleteSettings(w http.ResponseWriter, r *http.Request, guildID int64) {
	// Get user ID from session
	session, err := h.sessionStore.Get(r, h.sessionName)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := session.Values["user_id"].(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Delete settings
	err = h.settingsService.DeleteSettings(r.Context(), guildID, userID)
	if err != nil {
		log.Printf("Failed to delete settings: %v", err)
		http.Error(w, "Anfrage fehlgeschlagen", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{
		"success": true,
	})
}

// SendTestNotification sends a test notification to verify channel configuration
func (h *DiscordSettingsHandler) SendTestNotification(w http.ResponseWriter, r *http.Request, guildID int64) {
	// Get user ID from session
	session, err := h.sessionStore.Get(r, h.sessionName)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := session.Values["user_id"].(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse request body
	var input struct {
		ChannelID int64 `json:"channelId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if input.ChannelID == 0 {
		http.Error(w, "Channel ID required", http.StatusBadRequest)
		return
	}

	// Send test notification
	err = h.notificationService.SendTestNotification(r.Context(), userID, guildID, input.ChannelID)
	if err != nil {
		log.Printf("Failed to send test notification: %v", err)
		http.Error(w, "Anfrage fehlgeschlagen", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{
		"success": true,
	})
}