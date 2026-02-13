package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"log"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"demo/backend-go/internal/repository"
	"demo/backend-go/internal/infrastructure/redis"
)

type BotStatusHandler struct {
	userRepo     repository.UserRepository
	tokenRepo    repository.AuthTokenRepository
	sessionStore *sessions.CookieStore
	redisService *redis.RedisService
	sessionName  string
}

func NewBotStatusHandler(
	userRepo repository.UserRepository,
	tokenRepo repository.AuthTokenRepository,
	sessionStore *sessions.CookieStore,
	redisService *redis.RedisService,
	sessionName string,
) *BotStatusHandler {
	return &BotStatusHandler{
		userRepo:     userRepo,
		tokenRepo:    tokenRepo,
		sessionStore: sessionStore,
		redisService: redisService,
		sessionName:  sessionName,
	}
}

func (h *BotStatusHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/bot/status", h.GetBotStatus).Methods("GET")
	router.HandleFunc("/api/bot/auth/status", h.GetBotAuthStatus).Methods("GET")
	router.HandleFunc("/api/bot/token/refresh", h.RefreshBotToken).Methods("POST")  // NEU!

	router.HandleFunc("/api/bot/start", h.StartBot).Methods("POST")
	router.HandleFunc("/api/bot/stop", h.StopBot).Methods("POST")
	router.HandleFunc("/api/bot/restart", h.RestartBot).Methods("POST")
}

type BotStatusResponse struct {
	Running        bool     `json:"running"`
	TokenPresent   bool     `json:"tokenPresent"`
	ExpiresAt      *string  `json:"expiresAt"`
	ChannelName    *string  `json:"channelName"`
	ActiveChannels []string `json:"activeChannels,omitempty"`
	UptimeSeconds  *int     `json:"uptimeSeconds,omitempty"`
}

func (h *BotStatusHandler) GetBotStatus(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    response := BotStatusResponse{
        Running:        false,
        TokenPresent:   false,
        ActiveChannels: []string{},
    }

    botUser, err := h.userRepo.GetFirstBotUser(ctx)
    if err == nil && botUser != nil {
        if token, err := h.tokenRepo.GetByTwitchUserID(ctx, botUser.TwitchID); err == nil {
            response.TokenPresent = true
            response.ChannelName = &botUser.Username
            expiry := token.UpdatedAt.Add(time.Duration(token.ExpiresIn) * time.Second).Format(time.RFC3339)
            response.ExpiresAt = &expiry
        }
    }

    if h.redisService != nil {
        val, err := h.redisService.GetBotStatusData(ctx)
        if err == nil && val != "" {
            var redisData struct {
                IsRunning      bool     `json:"running"`
                Uptime         string   `json:"uptime"`
                UptimeSeconds  int      `json:"uptime_seconds"`
                ActiveChannels []string `json:"active_channels"`
                LastHeartbeat  string   `json:"last_heartbeat"`
            }

            if err := json.Unmarshal([]byte(val), &redisData); err == nil {
                if redisData.LastHeartbeat != "" {
                    lastBeat, parseErr := time.Parse(time.RFC3339, redisData.LastHeartbeat)
                    if parseErr == nil && time.Since(lastBeat) < 90*time.Second {
                        response.Running = redisData.IsRunning
                        if redisData.ActiveChannels != nil {
                            response.ActiveChannels = redisData.ActiveChannels
                        }
                        response.UptimeSeconds = &redisData.UptimeSeconds
                        log.Printf("✅ Bot Status aus Redis: running=%v, uptime=%ds, channels=%d",
                            redisData.IsRunning, redisData.UptimeSeconds, len(redisData.ActiveChannels))
                    } else {
                        log.Printf("⚠️ Heartbeat zu alt oder ungültig: %v (age: %v)", redisData.LastHeartbeat, time.Since(lastBeat))
                    }
                }
            } else {
                log.Printf("❌ Redis JSON parse error: %v", err)
            }
        } else if err != nil {
            log.Printf("⚠️ Redis GetBotStatusData error: %v", err)
        }
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (h *BotStatusHandler) GetBotAuthStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if !h.isAuthenticated(r) {
		http.Error(w, "Nicht authentifiziert", http.StatusUnauthorized)
		return
	}

	if !h.isAdmin(r, ctx) {
		http.Error(w, "Nicht authorisiert", http.StatusForbidden)
		return
	}

	botUser, err := h.userRepo.GetFirstBotUser(ctx)

	response := struct {
		Authenticated bool   `json:"authenticated"`
		BotUsername   string `json:"botUsername,omitempty"`
		BotTwitchId   string `json:"botTwitchId,omitempty"`
		TokenPresent  bool   `json:"tokenPresent"`
		ExpiresAt     string `json:"expiresAt,omitempty"`
	}{
		Authenticated: false,
		TokenPresent:  false,
	}

	if err == nil && botUser != nil {
		response.Authenticated = true
		response.BotUsername = botUser.Username
		response.BotTwitchId = botUser.TwitchID

		token, err := h.tokenRepo.GetByTwitchUserID(ctx, botUser.TwitchID)
		if err == nil {
			response.TokenPresent = true
			expiryTime := token.UpdatedAt.Add(time.Duration(token.ExpiresIn) * time.Second)
			response.ExpiresAt = expiryTime.Format(time.RFC3339)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *BotStatusHandler) isAuthenticated(r *http.Request) bool {
	session, err := h.sessionStore.Get(r, h.sessionName)
	if err != nil {
		return false
	}

	userID, ok := session.Values["user_id"].(string)
	return ok && userID != ""
}

func (h *BotStatusHandler) isAdmin(r *http.Request, ctx context.Context) bool {
	session, err := h.sessionStore.Get(r, h.sessionName)
	if err != nil {
		return false
	}

	userID, ok := session.Values["user_id"].(string)
	if !ok || userID == "" {
		return false
	}

	user, err := h.userRepo.GetByTwitchID(ctx, userID)
	if err != nil {
		return false
	}

	return user.IsAdmin
}

func (h *BotStatusHandler) StartBot(w http.ResponseWriter, r *http.Request) {
	h.publishBotCommand(w, r, "START_BOT")
}

func (h *BotStatusHandler) StopBot(w http.ResponseWriter, r *http.Request) {
	h.publishBotCommand(w, r, "STOP_BOT")
}

func (h *BotStatusHandler) RestartBot(w http.ResponseWriter, r *http.Request) {
	h.publishBotCommand(w, r, "RESTART_BOT")
}

func (h *BotStatusHandler) RefreshBotToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if !h.isAuthenticated(r) {
		http.Error(w, "Nicht authentifiziert", http.StatusUnauthorized)
		return
	}

	if !h.isAdmin(r, ctx) {
		http.Error(w, "Nicht authorisiert", http.StatusForbidden)
		return
	}

	botUser, err := h.userRepo.GetFirstBotUser(ctx)
	if err != nil || botUser == nil {
		http.Error(w, "Kein Bot User gefunden", http.StatusNotFound)
		return
	}

	token, err := h.tokenRepo.GetByTwitchUserID(ctx, botUser.TwitchID)
	if err != nil {
		http.Error(w, "Kein Token gefunden", http.StatusNotFound)
		return
	}

	// TODO: Token Refresh über AuthService
	// Für jetzt: Placeholder Response

	log.Printf("🔄 Manuelles Token Refresh angefordert für Bot: %s", botUser.Username)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Token refresh angefordert",
		"botUsername": botUser.Username,
		"expiresAt": token.UpdatedAt.Add(time.Duration(token.ExpiresIn) * time.Second).Format(time.RFC3339),
	})
}

func (h *BotStatusHandler) publishBotCommand(w http.ResponseWriter, r *http.Request, command string) {
    if h.redisService == nil {
        http.Error(w, "Redis Service nicht verfügbar", http.StatusInternalServerError)
        return
    }

    if err := h.redisService.SendBotControlSignal(command); err != nil {
        log.Printf("❌ Fehler beim Senden des Befehls %s: %v", command, err)
        http.Error(w, "Fehler beim Senden", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "ok", "command": command})
}