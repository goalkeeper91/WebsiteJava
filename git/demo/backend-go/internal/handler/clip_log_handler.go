package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"demo/backend-go/internal/service"
)

type ClipLogHandler struct {
	clipService  *service.ClipLogService
	sessionStore *sessions.CookieStore
	sessionName  string
}

func NewClipLogHandler(
	clipService *service.ClipLogService,
	sessionStore *sessions.CookieStore,
	sessionName string,
) *ClipLogHandler {
	return &ClipLogHandler{
		clipService:  clipService,
		sessionStore: sessionStore,
		sessionName:  sessionName,
	}
}

func (h *ClipLogHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/clips", h.GetClips).Methods("GET")
	router.HandleFunc("/api/clips/stats", h.GetStats).Methods("GET")
}

func (h *ClipLogHandler) GetClips(w http.ResponseWriter, r *http.Request) {
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

	ctx := r.Context()

	clips, err := h.clipService.GetClips(ctx, user.ID, pageSize, offset)
	if err != nil {
		log.Printf("Fehler beim Laden der Clips: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	stats, err := h.clipService.GetStats(ctx, user.ID)
	if err != nil {
		log.Printf("Fehler beim Laden der Clip-Statistiken: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	total := stats.TotalClips
	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	h.respondJSON(w, http.StatusOK, PaginatedResponse{
		Data:       clips,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

func (h *ClipLogHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	user := h.requireUser(w, r)
	if user == nil {
		return
	}

	stats, err := h.clipService.GetStats(r.Context(), user.ID)
	if err != nil {
		log.Printf("Fehler beim Laden der Clip-Statistiken: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, stats)
}

// ============================================================
// HELPERS
// ============================================================

func (h *ClipLogHandler) requireUser(w http.ResponseWriter, r *http.Request) *UserSession {
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

func (h *ClipLogHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
