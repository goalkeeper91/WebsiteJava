package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/gorilla/mux"
	"demo/backend-go/internal/repository"
	"demo/backend-go/internal/service"
)

type DiscordGuildHandler struct {
	guildService *service.DiscordGuildService
	userRepo     repository.UserRepository
	sessionStore sessions.Store
	sessionName  string
	botClientID  string
	botPermissions string
}

func NewDiscordGuildHandler(
	guildService *service.DiscordGuildService,
	userRepo repository.UserRepository,
	sessionStore sessions.Store,
	sessionName string,
	botClientID string,
	botPermissions string,
) *DiscordGuildHandler {
	return &DiscordGuildHandler{
		guildService:   guildService,
		userRepo:       userRepo,
		sessionStore:   sessionStore,
		sessionName:    sessionName,
		botClientID:    botClientID,
		botPermissions: botPermissions,
	}
}

// isAdmin reports whether the session user is a site admin, allowed to
// manage any guild (not just ones they personally added the bot to) — used
// by the admin bot-management panel, mirroring bot_status_handler.go.
func (h *DiscordGuildHandler) isAdmin(r *http.Request) bool {
	session, err := h.sessionStore.Get(r, h.sessionName)
	if err != nil {
		return false
	}

	userID, ok := session.Values["user_id"].(string)
	if !ok || userID == "" {
		return false
	}

	user, err := h.userRepo.GetByTwitchID(r.Context(), userID)
	if err != nil || user == nil {
		return false
	}

	return user.IsAdmin
}

func (h *DiscordGuildHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/discord/guilds", h.handleGuilds)
	router.HandleFunc("/api/discord/guilds/", h.handleGuildByID)
	router.HandleFunc("/api/discord/bot/invite-url", h.GetInviteURL)
}

func (h *DiscordGuildHandler) GetInviteURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	inviteURL := "https://discord.com/oauth2/authorize?client_id=" + h.botClientID +
		"&permissions=" + h.botPermissions +
		"&scope=bot%20applications.commands"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"inviteUrl": inviteURL,
	})
}

func (h *DiscordGuildHandler) handleGuilds(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.ListGuilds(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *DiscordGuildHandler) ListGuilds(w http.ResponseWriter, r *http.Request) {
	session, err := h.sessionStore.Get(r, h.sessionName)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]interface{}{})
		return
	}

	userID, ok := session.Values["user_id"].(string)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]interface{}{})
		return
	}

	guilds, err := h.guildService.GetUserGuilds(r.Context(), userID)
	if err != nil {
		log.Printf("Failed to get guilds: %v", err)
		http.Error(w, "Failed to get guilds", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(guilds)
}

func (h *DiscordGuildHandler) handleGuildByID(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/discord/guilds/")
	parts := strings.Split(path, "/")

	if len(parts) == 0 || parts[0] == "" {
		http.Error(w, "Guild ID required", http.StatusBadRequest)
		return
	}

	guildID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		http.Error(w, "Invalid guild ID", http.StatusBadRequest)
		return
	}

	if len(parts) > 1 {
		switch parts[1] {
		case "sync":
			if r.Method == http.MethodPost {
				h.SyncGuild(w, r, guildID)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		case "details":
			if r.Method == http.MethodGet {
				h.GetGuildDetails(w, r, guildID)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		default:
			http.Error(w, "Not found", http.StatusNotFound)
		}
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetGuild(w, r, guildID)
	case http.MethodDelete:
		h.RemoveBot(w, r, guildID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *DiscordGuildHandler) GetGuild(w http.ResponseWriter, r *http.Request, guildID int64) {
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

	guild, err := h.guildService.GetGuild(r.Context(), guildID, userID, h.isAdmin(r))
	if err != nil {
		log.Printf("Failed to get guild: %v", err)
		http.Error(w, "Zugriff auf diesen Server nicht möglich", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(guild)
}

func (h *DiscordGuildHandler) GetGuildDetails(w http.ResponseWriter, r *http.Request, guildID int64) {
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

	details, err := h.guildService.GetGuildDetails(r.Context(), guildID, userID)
	if err != nil {
		log.Printf("Failed to get guild details: %v", err)
		http.Error(w, "Zugriff auf diesen Server nicht möglich", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(details)
}

func (h *DiscordGuildHandler) SyncGuild(w http.ResponseWriter, r *http.Request, guildID int64) {
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

	err = h.guildService.SyncGuild(r.Context(), guildID, userID, h.isAdmin(r))
	if err != nil {
		log.Printf("Failed to sync guild: %v", err)
		http.Error(w, "Synchronisation fehlgeschlagen", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{
		"success": true,
	})
}

func (h *DiscordGuildHandler) RemoveBot(w http.ResponseWriter, r *http.Request, guildID int64) {
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

	err = h.guildService.RemoveBot(r.Context(), guildID, userID, h.isAdmin(r))
	if err != nil {
		log.Printf("Failed to remove bot: %v", err)
		http.Error(w, "Bot konnte nicht entfernt werden", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{
		"success": true,
	})
}