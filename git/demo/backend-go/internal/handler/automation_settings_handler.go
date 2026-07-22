package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/service"
)

type AutomationSettingsHandler struct {
	settingsService *service.AutomationSettingsService
	sessionStore    *sessions.CookieStore
	sessionName     string
}

func NewAutomationSettingsHandler(
	settingsService *service.AutomationSettingsService,
	sessionStore *sessions.CookieStore,
	sessionName string,
) *AutomationSettingsHandler {
	return &AutomationSettingsHandler{
		settingsService: settingsService,
		sessionStore:    sessionStore,
		sessionName:     sessionName,
	}
}

func (h *AutomationSettingsHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/automation/settings", h.GetSettings).Methods("GET")
	router.HandleFunc("/api/automation/settings", h.UpdateSettings).Methods("PUT")
}

func (h *AutomationSettingsHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	user := h.requireUser(w, r)
	if user == nil {
		return
	}

	settings, err := h.settingsService.GetSettings(r.Context(), user.ID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.respondJSON(w, http.StatusOK, settings)
}

func (h *AutomationSettingsHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	user := h.requireUser(w, r)
	if user == nil {
		return
	}

	var input domain.AutomationSettingsUpdateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
		return
	}

	settings, err := h.settingsService.UpdateSettings(r.Context(), user.ID, input)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.respondJSON(w, http.StatusOK, settings)
}

// ============================================================
// HELPERS
// ============================================================

func (h *AutomationSettingsHandler) requireUser(w http.ResponseWriter, r *http.Request) *UserSession {
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

func (h *AutomationSettingsHandler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrFeatureNotAvailable):
		http.Error(w, "Clip-Automatisierung erfordert ein Pro oder Premium Abo", http.StatusForbidden)
	case errors.Is(err, domain.ErrSubscriptionExpired):
		http.Error(w, "Abo abgelaufen", http.StatusForbidden)
	case errors.Is(err, domain.ErrInvalidClipDurationRange),
		errors.Is(err, domain.ErrClipDurationTooShort),
		errors.Is(err, domain.ErrClipDurationTooLong),
		errors.Is(err, domain.ErrInvalidHashtagCount):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		log.Printf("Automation Settings Fehler: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
	}
}

func (h *AutomationSettingsHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
