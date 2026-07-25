package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/infrastructure/redis"
	"demo/backend-go/internal/service"
)

type GiveawayHandler struct {
	giveawayService *service.GiveawayService
	sessionStore    *sessions.CookieStore
	sessionName     string
	internalSecret  string
	redisService    *redis.RedisService
}

func NewGiveawayHandler(
	giveawayService *service.GiveawayService,
	sessionStore *sessions.CookieStore,
	sessionName string,
	internalSecret string,
	redisService *redis.RedisService,
) *GiveawayHandler {
	return &GiveawayHandler{
		giveawayService: giveawayService,
		sessionStore:    sessionStore,
		sessionName:     sessionName,
		internalSecret:  internalSecret,
		redisService:    redisService,
	}
}

func (h *GiveawayHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/dashboard/giveaways/status", h.GetStatus).Methods("GET")
	router.HandleFunc("/api/dashboard/giveaways/history", h.GetHistory).Methods("GET")
	router.HandleFunc("/api/dashboard/giveaways/start", h.StartForDashboard).Methods("POST")
	router.HandleFunc("/api/dashboard/giveaways/draw", h.DrawForDashboard).Methods("POST")
	router.HandleFunc("/api/dashboard/giveaways/cancel", h.CancelForDashboard).Methods("POST")

	router.HandleFunc("/api/bot/giveaways/start", h.StartForBot).Methods("POST")
	router.HandleFunc("/api/bot/giveaways/enter", h.EnterForBot).Methods("POST")
	router.HandleFunc("/api/bot/giveaways/draw", h.DrawForBot).Methods("POST")
	router.HandleFunc("/api/bot/giveaways/cancel", h.CancelForBot).Methods("POST")
	router.HandleFunc("/api/bot/giveaways/status", h.GetStatusForBot).Methods("GET")
	router.HandleFunc("/api/bot/giveaways/open", h.GetOpenForBot).Methods("GET")
}

// ============================================================
// DASHBOARD (Session-Auth, rein lesend - siehe Plan)
// ============================================================

func (h *GiveawayHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	user := h.requireUser(w, r)
	if user == nil {
		return
	}

	giveaway, entryCount, err := h.giveawayService.GetStatus(r.Context(), user.ID)
	if err != nil {
		log.Printf("Giveaway Fehler (Status): %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"giveaway":    giveaway,
		"entry_count": entryCount,
	})
}

func (h *GiveawayHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	user := h.requireUser(w, r)
	if user == nil {
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	giveaways, total, err := h.giveawayService.GetHistory(r.Context(), user.ID, pageSize, offset)
	if err != nil {
		log.Printf("Giveaway Fehler (Historie): %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	h.respondJSON(w, http.StatusOK, PaginatedResponse{
		Data:       giveaways,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

// StartForDashboard lets the streamer start a giveaway from the dashboard
// instead of chat - see DrawForDashboard for the matching draw action.
// Both post a chat announcement via SendBotAnnounce, since this codepath
// never runs the bot's own ctx.send confirmation.
func (h *GiveawayHandler) StartForDashboard(w http.ResponseWriter, r *http.Request) {
	user := h.requireUser(w, r)
	if user == nil {
		return
	}

	var req StartGiveawayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
		return
	}

	giveaway, err := h.giveawayService.StartGiveaway(r.Context(), user.ID, req.Keyword, req.SubBonus)
	if err == domain.ErrGiveawayAlreadyOpen {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	if err == domain.ErrGiveawayKeywordInvalid {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		log.Printf("Giveaway Fehler (Dashboard-Start): %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	if h.redisService != nil {
		message := fmt.Sprintf("🎉 Giveaway gestartet! Tippt \"%s\" in den Chat zum Teilnehmen.", giveaway.Keyword)
		if err := h.redisService.SendBotAnnounce(message, user.ID); err != nil {
			log.Printf("Giveaway Fehler (Ankündigung Start): %v", err)
		}
	}

	h.respondJSON(w, http.StatusOK, giveaway)
}

// DrawForDashboard mirrors DrawForBot but is session-auth-gated for
// dashboard use.
func (h *GiveawayHandler) DrawForDashboard(w http.ResponseWriter, r *http.Request) {
	user := h.requireUser(w, r)
	if user == nil {
		return
	}

	giveaway, err := h.giveawayService.DrawWinner(r.Context(), user.ID)
	if err == domain.ErrNoOpenGiveaway || err == domain.ErrNoGiveawayEntries {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	if err != nil {
		log.Printf("Giveaway Fehler (Dashboard-Ziehung): %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	if h.redisService != nil {
		winnerLogin := ""
		if giveaway.WinnerLogin != nil {
			winnerLogin = *giveaway.WinnerLogin
		}
		message := fmt.Sprintf("🎉 Der Gewinner ist @%s! Herzlichen Glückwunsch!", winnerLogin)
		if err := h.redisService.SendBotAnnounce(message, user.ID); err != nil {
			log.Printf("Giveaway Fehler (Ankündigung Ziehung): %v", err)
		}
	}

	h.respondJSON(w, http.StatusOK, giveaway)
}

// CancelForDashboard closes an open giveaway without drawing a winner - lets
// a streamer end one with zero (or unwanted) entries instead of being stuck
// until someone enters or forever hiding the start form.
func (h *GiveawayHandler) CancelForDashboard(w http.ResponseWriter, r *http.Request) {
	user := h.requireUser(w, r)
	if user == nil {
		return
	}

	giveaway, err := h.giveawayService.CancelGiveaway(r.Context(), user.ID)
	if err == domain.ErrNoOpenGiveaway {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	if err != nil {
		log.Printf("Giveaway Fehler (Dashboard-Abbruch): %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	if h.redisService != nil {
		if err := h.redisService.SendBotAnnounce("🚫 Das Giveaway wurde abgebrochen.", user.ID); err != nil {
			log.Printf("Giveaway Fehler (Ankündigung Abbruch): %v", err)
		}
	}

	h.respondJSON(w, http.StatusOK, giveaway)
}

// ============================================================
// BOT-INTERN (Shared-Secret-Auth, kein Browser-Login)
// ============================================================

type StartGiveawayRequest struct {
	BroadcasterTwitchID string `json:"broadcaster_id"`
	Keyword             string `json:"keyword"`
	SubBonus            bool   `json:"sub_bonus"`
}

func (h *GiveawayHandler) StartForBot(w http.ResponseWriter, r *http.Request) {
	if !requireInternalSecret(w, r, h.internalSecret) {
		return
	}

	var req StartGiveawayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
		return
	}
	if req.BroadcasterTwitchID == "" {
		http.Error(w, "broadcaster_id ist erforderlich", http.StatusBadRequest)
		return
	}

	giveaway, err := h.giveawayService.StartGiveaway(r.Context(), req.BroadcasterTwitchID, req.Keyword, req.SubBonus)
	if err == domain.ErrGiveawayAlreadyOpen {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	if err == domain.ErrGiveawayKeywordInvalid {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		log.Printf("Giveaway Fehler (Start): %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, giveaway)
}

type EnterGiveawayRequest struct {
	BroadcasterTwitchID string `json:"broadcaster_id"`
	ViewerTwitchID      string `json:"viewer_twitch_id"`
	ViewerLogin         string `json:"viewer_login"`
	IsSubscriber        bool   `json:"is_subscriber"`
}

func (h *GiveawayHandler) EnterForBot(w http.ResponseWriter, r *http.Request) {
	if !requireInternalSecret(w, r, h.internalSecret) {
		return
	}

	var req EnterGiveawayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
		return
	}
	if req.BroadcasterTwitchID == "" || req.ViewerTwitchID == "" || req.ViewerLogin == "" {
		http.Error(w, "broadcaster_id, viewer_twitch_id und viewer_login sind erforderlich", http.StatusBadRequest)
		return
	}

	inserted, err := h.giveawayService.EnterGiveaway(r.Context(), req.BroadcasterTwitchID, req.ViewerTwitchID, req.ViewerLogin, req.IsSubscriber)
	if err == domain.ErrNoOpenGiveaway {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Giveaway Fehler (Eintragen): %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]bool{"entered": inserted})
}

type DrawGiveawayRequest struct {
	BroadcasterTwitchID string `json:"broadcaster_id"`
}

func (h *GiveawayHandler) DrawForBot(w http.ResponseWriter, r *http.Request) {
	if !requireInternalSecret(w, r, h.internalSecret) {
		return
	}

	var req DrawGiveawayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
		return
	}
	if req.BroadcasterTwitchID == "" {
		http.Error(w, "broadcaster_id ist erforderlich", http.StatusBadRequest)
		return
	}

	giveaway, err := h.giveawayService.DrawWinner(r.Context(), req.BroadcasterTwitchID)
	if err == domain.ErrNoOpenGiveaway || err == domain.ErrNoGiveawayEntries {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	if err != nil {
		log.Printf("Giveaway Fehler (Ziehung): %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, giveaway)
}

func (h *GiveawayHandler) CancelForBot(w http.ResponseWriter, r *http.Request) {
	if !requireInternalSecret(w, r, h.internalSecret) {
		return
	}

	var req DrawGiveawayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
		return
	}
	if req.BroadcasterTwitchID == "" {
		http.Error(w, "broadcaster_id ist erforderlich", http.StatusBadRequest)
		return
	}

	giveaway, err := h.giveawayService.CancelGiveaway(r.Context(), req.BroadcasterTwitchID)
	if err == domain.ErrNoOpenGiveaway {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	if err != nil {
		log.Printf("Giveaway Fehler (Abbruch): %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, giveaway)
}

func (h *GiveawayHandler) GetStatusForBot(w http.ResponseWriter, r *http.Request) {
	if !requireInternalSecret(w, r, h.internalSecret) {
		return
	}

	broadcasterID := r.URL.Query().Get("broadcaster_id")
	if broadcasterID == "" {
		http.Error(w, "broadcaster_id ist erforderlich", http.StatusBadRequest)
		return
	}

	giveaway, entryCount, err := h.giveawayService.GetStatus(r.Context(), broadcasterID)
	if err != nil {
		log.Printf("Giveaway Fehler (Bot-Status): %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"giveaway":    giveaway,
		"entry_count": entryCount,
	})
}

func (h *GiveawayHandler) GetOpenForBot(w http.ResponseWriter, r *http.Request) {
	if !requireInternalSecret(w, r, h.internalSecret) {
		return
	}

	giveaways, err := h.giveawayService.GetOpenGiveawaysForBot(r.Context())
	if err != nil {
		log.Printf("Giveaway Fehler (Bot-Cache-Reload): %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, giveaways)
}

// ============================================================
// HELPERS
// ============================================================

func (h *GiveawayHandler) requireUser(w http.ResponseWriter, r *http.Request) *UserSession {
	session, err := h.sessionStore.Get(r, h.sessionName)
	if err != nil {
		log.Printf("Fehler beim Laden der Session: %v", err)
		http.Error(w, "Nicht authentifiziert", http.StatusUnauthorized)
		return nil
	}

	userID, ok := session.Values["user_id"].(string)
	if !ok || userID == "" {
		http.Error(w, "Nicht authentifiziert", http.StatusUnauthorized)
		return nil
	}

	return &UserSession{ID: userID}
}

func (h *GiveawayHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
