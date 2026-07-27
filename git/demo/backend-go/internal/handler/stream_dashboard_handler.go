package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/gorilla/mux"

	"demo/backend-go/internal/service"
)

// StreamDashboardHandler backs the Live-Dashboard home page - stream
// info/title/category, live status, aggregate stats, category search, and
// ad breaks. All routes resolve team-delegated access via
// resolveEffectiveTwitchID from the start (unlike activity_handler.go,
// which needed a follow-up fix - see LD9).
type StreamDashboardHandler struct {
	streamService *service.StreamDashboardService
	sessionStore  *sessions.CookieStore
	sessionName   string
	teamService   *service.TeamService
}

func NewStreamDashboardHandler(
	streamService *service.StreamDashboardService,
	sessionStore *sessions.CookieStore,
	sessionName string,
	teamService *service.TeamService,
) *StreamDashboardHandler {
	return &StreamDashboardHandler{
		streamService: streamService,
		sessionStore:  sessionStore,
		sessionName:   sessionName,
		teamService:   teamService,
	}
}

func (h *StreamDashboardHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/dashboard/stream/info", h.GetStreamInfo).Methods("GET")
	router.HandleFunc("/api/dashboard/stream/info", h.UpdateStreamInfo).Methods("PATCH")
	router.HandleFunc("/api/dashboard/stream/live", h.GetLiveStatus).Methods("GET")
	router.HandleFunc("/api/dashboard/stream/stats", h.GetDashboardStats).Methods("GET")
	router.HandleFunc("/api/dashboard/stream/categories/search", h.SearchCategories).Methods("GET")
	router.HandleFunc("/api/dashboard/stream/commercial", h.StartCommercial).Methods("POST")
}

func (h *StreamDashboardHandler) GetStreamInfo(w http.ResponseWriter, r *http.Request) {
	_, channelID, ok := resolveEffectiveTwitchID(w, r, h.sessionStore, h.sessionName, h.teamService)
	if !ok {
		return
	}

	info, err := h.streamService.GetStreamInfo(r.Context(), channelID)
	if err != nil {
		log.Printf("Stream Fehler (Info): %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, info)
}

type UpdateStreamInfoRequest struct {
	Title  *string `json:"title,omitempty"`
	GameID *string `json:"gameId,omitempty"`
}

func (h *StreamDashboardHandler) UpdateStreamInfo(w http.ResponseWriter, r *http.Request) {
	_, channelID, ok := resolveEffectiveTwitchID(w, r, h.sessionStore, h.sessionName, h.teamService)
	if !ok {
		return
	}

	var req UpdateStreamInfoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
		return
	}

	info, err := h.streamService.UpdateStreamInfo(r.Context(), channelID, req.Title, req.GameID)
	if err != nil {
		log.Printf("Stream Fehler (Update Info): %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, info)
}

func (h *StreamDashboardHandler) GetLiveStatus(w http.ResponseWriter, r *http.Request) {
	_, channelID, ok := resolveEffectiveTwitchID(w, r, h.sessionStore, h.sessionName, h.teamService)
	if !ok {
		return
	}

	live, err := h.streamService.GetLiveStatus(r.Context(), channelID)
	if err != nil {
		log.Printf("Stream Fehler (Live-Status): %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, live)
}

func (h *StreamDashboardHandler) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	_, channelID, ok := resolveEffectiveTwitchID(w, r, h.sessionStore, h.sessionName, h.teamService)
	if !ok {
		return
	}

	stats, err := h.streamService.GetDashboardStats(r.Context(), channelID)
	if err != nil {
		log.Printf("Stream Fehler (Stats): %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, stats)
}

func (h *StreamDashboardHandler) SearchCategories(w http.ResponseWriter, r *http.Request) {
	_, _, ok := resolveEffectiveTwitchID(w, r, h.sessionStore, h.sessionName, h.teamService)
	if !ok {
		return
	}

	query := r.URL.Query().Get("query")
	if query == "" {
		h.respondJSON(w, http.StatusOK, []interface{}{})
		return
	}

	categories, err := h.streamService.SearchCategories(r.Context(), query)
	if err != nil {
		log.Printf("Stream Fehler (Kategorie-Suche): %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, categories)
}

type StartCommercialRequest struct {
	Length int `json:"length"`
}

func (h *StreamDashboardHandler) StartCommercial(w http.ResponseWriter, r *http.Request) {
	_, channelID, ok := resolveEffectiveTwitchID(w, r, h.sessionStore, h.sessionName, h.teamService)
	if !ok {
		return
	}

	var req StartCommercialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Length <= 0 {
		http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
		return
	}

	result, err := h.streamService.StartCommercial(r.Context(), channelID, req.Length)
	if err != nil {
		log.Printf("Stream Fehler (Werbepause): %v", err)
		http.Error(w, "Konnte keine Werbepause starten", http.StatusBadGateway)
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

func (h *StreamDashboardHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
