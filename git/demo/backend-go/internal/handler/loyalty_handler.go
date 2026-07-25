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

type LoyaltyHandler struct {
	loyaltyService *service.LoyaltyService
	sessionStore   *sessions.CookieStore
	sessionName    string
	internalSecret string
}

func NewLoyaltyHandler(
	loyaltyService *service.LoyaltyService,
	sessionStore *sessions.CookieStore,
	sessionName string,
	internalSecret string,
) *LoyaltyHandler {
	return &LoyaltyHandler{
		loyaltyService: loyaltyService,
		sessionStore:   sessionStore,
		sessionName:    sessionName,
		internalSecret: internalSecret,
	}
}

func (h *LoyaltyHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/dashboard/loyalty", h.GetSettings).Methods("GET")
	router.HandleFunc("/api/dashboard/loyalty", h.UpdateSettings).Methods("PUT")
	router.HandleFunc("/api/dashboard/loyalty/leaderboard", h.GetLeaderboard).Methods("GET")

	router.HandleFunc("/api/bot/loyalty/points", h.GetPointsForBot).Methods("GET")
	router.HandleFunc("/api/bot/loyalty/leaderboard", h.GetLeaderboardForBot).Methods("GET")
}

// ============================================================
// DASHBOARD (Session-Auth)
// ============================================================

func (h *LoyaltyHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	user := h.requireUser(w, r)
	if user == nil {
		return
	}

	settings, err := h.loyaltyService.GetSettings(r.Context(), user.ID)
	if err != nil {
		log.Printf("Loyalty Fehler: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, settings)
}

func (h *LoyaltyHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	user := h.requireUser(w, r)
	if user == nil {
		return
	}

	var input domain.LoyaltySettingsUpdateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
		return
	}

	settings, err := h.loyaltyService.UpdateSettings(r.Context(), user.ID, input)
	if err != nil {
		log.Printf("Loyalty Fehler: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, settings)
}

func (h *LoyaltyHandler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
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

	entries, total, err := h.loyaltyService.GetLeaderboard(r.Context(), user.ID, pageSize, offset)
	if err != nil {
		log.Printf("Loyalty Fehler (Bestenliste laden): %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	h.respondJSON(w, http.StatusOK, PaginatedResponse{
		Data:       entries,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

// ============================================================
// BOT-INTERN (Shared-Secret-Auth, kein Browser-Login)
// ============================================================

func (h *LoyaltyHandler) GetPointsForBot(w http.ResponseWriter, r *http.Request) {
	if !requireInternalSecret(w, r, h.internalSecret) {
		return
	}

	broadcasterID := r.URL.Query().Get("broadcaster_id")
	viewerLogin := r.URL.Query().Get("viewer_login")
	if broadcasterID == "" || viewerLogin == "" {
		http.Error(w, "broadcaster_id und viewer_login sind erforderlich", http.StatusBadRequest)
		return
	}

	points, pointsName, err := h.loyaltyService.GetViewerPointsForBot(r.Context(), broadcasterID, viewerLogin)
	if err != nil {
		log.Printf("Loyalty Fehler (Bot Punkte-Abfrage): %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"points":      points,
		"points_name": pointsName,
	})
}

func (h *LoyaltyHandler) GetLeaderboardForBot(w http.ResponseWriter, r *http.Request) {
	if !requireInternalSecret(w, r, h.internalSecret) {
		return
	}

	broadcasterID := r.URL.Query().Get("broadcaster_id")
	if broadcasterID == "" {
		http.Error(w, "broadcaster_id ist erforderlich", http.StatusBadRequest)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 25 {
		limit = 5
	}

	entries, pointsName, err := h.loyaltyService.GetLeaderboardForBot(r.Context(), broadcasterID, limit)
	if err != nil {
		log.Printf("Loyalty Fehler (Bot Bestenliste): %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"entries":     entries,
		"points_name": pointsName,
	})
}

// ============================================================
// HELPERS
// ============================================================

func (h *LoyaltyHandler) requireUser(w http.ResponseWriter, r *http.Request) *UserSession {
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

func (h *LoyaltyHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
