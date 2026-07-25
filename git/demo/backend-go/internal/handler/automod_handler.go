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

type AutomodHandler struct {
	automodService *service.AutomodService
	sessionStore   *sessions.CookieStore
	sessionName    string
	internalSecret string
}

func NewAutomodHandler(
	automodService *service.AutomodService,
	sessionStore *sessions.CookieStore,
	sessionName string,
	internalSecret string,
) *AutomodHandler {
	return &AutomodHandler{
		automodService: automodService,
		sessionStore:   sessionStore,
		sessionName:    sessionName,
		internalSecret: internalSecret,
	}
}

func (h *AutomodHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/dashboard/automod", h.GetSettings).Methods("GET")
	router.HandleFunc("/api/dashboard/automod", h.UpdateSettings).Methods("PUT")
	router.HandleFunc("/api/dashboard/automod/events", h.GetEvents).Methods("GET")

	router.HandleFunc("/api/bot/automod-settings", h.GetSettingsForBot).Methods("GET")
	router.HandleFunc("/api/bot/automod-violation", h.RecordViolation).Methods("POST")
}

// ============================================================
// DASHBOARD (Session-Auth)
// ============================================================

func (h *AutomodHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	user := h.requireUser(w, r)
	if user == nil {
		return
	}

	settings, err := h.automodService.GetSettings(r.Context(), user.ID)
	if err != nil {
		log.Printf("Automod Fehler: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, settings)
}

func (h *AutomodHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	user := h.requireUser(w, r)
	if user == nil {
		return
	}

	var input domain.AutomodSettingsUpdateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
		return
	}

	settings, err := h.automodService.UpdateSettings(r.Context(), user.ID, input)
	if err != nil {
		log.Printf("Automod Fehler: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, settings)
}

func (h *AutomodHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
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

	events, total, err := h.automodService.GetEvents(r.Context(), user.ID, pageSize, offset)
	if err != nil {
		log.Printf("Automod Fehler (Events laden): %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	h.respondJSON(w, http.StatusOK, PaginatedResponse{
		Data:       events,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

// ============================================================
// BOT-INTERN (Shared-Secret-Auth, kein Browser-Login)
// ============================================================

func (h *AutomodHandler) GetSettingsForBot(w http.ResponseWriter, r *http.Request) {
	if !requireInternalSecret(w, r, h.internalSecret) {
		return
	}

	settings, err := h.automodService.GetAllEnabledSettingsForBot(r.Context())
	if err != nil {
		log.Printf("Automod Fehler (Bot-Cache-Reload): %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, settings)
}

type RecordViolationRequest struct {
	BroadcasterTwitchID string `json:"broadcaster_twitch_id"`
	OffenderTwitchID    string `json:"offender_twitch_id"`
	OffenderName        string `json:"offender_name"`
	Reason              string `json:"reason"`
	MessageExcerpt      string `json:"message_excerpt"`
}

func (h *AutomodHandler) RecordViolation(w http.ResponseWriter, r *http.Request) {
	if !requireInternalSecret(w, r, h.internalSecret) {
		return
	}

	var req RecordViolationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
		return
	}
	if req.BroadcasterTwitchID == "" || req.OffenderTwitchID == "" {
		http.Error(w, "broadcaster_twitch_id und offender_twitch_id sind erforderlich", http.StatusBadRequest)
		return
	}

	timeoutSeconds, err := h.automodService.RecordViolation(
		r.Context(), req.BroadcasterTwitchID, req.OffenderTwitchID, req.OffenderName, req.Reason, req.MessageExcerpt,
	)
	if err != nil {
		log.Printf("Automod Fehler (Verstoß erfassen): %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]int{"timeout_seconds": timeoutSeconds})
}

// ============================================================
// HELPERS
// ============================================================

func (h *AutomodHandler) requireUser(w http.ResponseWriter, r *http.Request) *UserSession {
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

func (h *AutomodHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
