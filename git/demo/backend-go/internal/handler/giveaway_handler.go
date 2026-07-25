package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/service"
)

type GiveawayHandler struct {
	giveawayService *service.GiveawayService
	sessionStore    *sessions.CookieStore
	sessionName     string
	internalSecret  string
}

func NewGiveawayHandler(
	giveawayService *service.GiveawayService,
	sessionStore *sessions.CookieStore,
	sessionName string,
	internalSecret string,
) *GiveawayHandler {
	return &GiveawayHandler{
		giveawayService: giveawayService,
		sessionStore:    sessionStore,
		sessionName:     sessionName,
		internalSecret:  internalSecret,
	}
}

func (h *GiveawayHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/dashboard/giveaways/status", h.GetStatus).Methods("GET")
	router.HandleFunc("/api/dashboard/giveaways/history", h.GetHistory).Methods("GET")

	router.HandleFunc("/api/bot/giveaways/start", h.StartForBot).Methods("POST")
	router.HandleFunc("/api/bot/giveaways/enter", h.EnterForBot).Methods("POST")
	router.HandleFunc("/api/bot/giveaways/draw", h.DrawForBot).Methods("POST")
	router.HandleFunc("/api/bot/giveaways/status", h.GetStatusForBot).Methods("GET")
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

// ============================================================
// BOT-INTERN (Shared-Secret-Auth, kein Browser-Login)
// ============================================================

type StartGiveawayRequest struct {
	BroadcasterTwitchID string `json:"broadcaster_id"`
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

	giveaway, err := h.giveawayService.StartGiveaway(r.Context(), req.BroadcasterTwitchID, req.SubBonus)
	if err == domain.ErrGiveawayAlreadyOpen {
		http.Error(w, err.Error(), http.StatusConflict)
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
