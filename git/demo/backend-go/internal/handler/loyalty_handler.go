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
	teamService    *service.TeamService
}

func NewLoyaltyHandler(
	loyaltyService *service.LoyaltyService,
	sessionStore *sessions.CookieStore,
	sessionName string,
	internalSecret string,
	teamService *service.TeamService,
) *LoyaltyHandler {
	return &LoyaltyHandler{
		loyaltyService: loyaltyService,
		sessionStore:   sessionStore,
		sessionName:    sessionName,
		internalSecret: internalSecret,
		teamService:    teamService,
	}
}

func (h *LoyaltyHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/dashboard/loyalty", h.GetSettings).Methods("GET")
	router.HandleFunc("/api/dashboard/loyalty", h.UpdateSettings).Methods("PUT")
	router.HandleFunc("/api/dashboard/loyalty/leaderboard", h.GetLeaderboard).Methods("GET")

	router.HandleFunc("/api/bot/loyalty/points", h.GetPointsForBot).Methods("GET")
	router.HandleFunc("/api/bot/loyalty/leaderboard", h.GetLeaderboardForBot).Methods("GET")
	router.HandleFunc("/api/bot/loyalty/regulars", h.GetRegularsForBot).Methods("GET")
}

// ============================================================
// DASHBOARD (Session-Auth)
// ============================================================

func (h *LoyaltyHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	_, channelID, ok := resolveEffectiveTwitchID(w, r, h.sessionStore, h.sessionName, h.teamService)
	if !ok {
		return
	}

	settings, err := h.loyaltyService.GetSettings(r.Context(), channelID)
	if err != nil {
		log.Printf("Loyalty Fehler: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, settings)
}

func (h *LoyaltyHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	_, channelID, ok := resolveEffectiveTwitchID(w, r, h.sessionStore, h.sessionName, h.teamService)
	if !ok {
		return
	}

	var input domain.LoyaltySettingsUpdateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
		return
	}

	settings, err := h.loyaltyService.UpdateSettings(r.Context(), channelID, input)
	if err != nil {
		log.Printf("Loyalty Fehler: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, settings)
}

func (h *LoyaltyHandler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	_, channelID, ok := resolveEffectiveTwitchID(w, r, h.sessionStore, h.sessionName, h.teamService)
	if !ok {
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

	entries, total, err := h.loyaltyService.GetLeaderboard(r.Context(), channelID, pageSize, offset)
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

func (h *LoyaltyHandler) GetRegularsForBot(w http.ResponseWriter, r *http.Request) {
	if !requireInternalSecret(w, r, h.internalSecret) {
		return
	}

	broadcasterID := r.URL.Query().Get("broadcaster_id")
	if broadcasterID == "" {
		http.Error(w, "broadcaster_id ist erforderlich", http.StatusBadRequest)
		return
	}

	logins, err := h.loyaltyService.GetRegularsForBot(r.Context(), broadcasterID)
	if err != nil {
		log.Printf("Loyalty Fehler (Bot Regulars): %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"logins": logins,
	})
}

// ============================================================
// HELPERS
// ============================================================

func (h *LoyaltyHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
