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

type ActivityHandler struct {
	activityService *service.ActivityService
	sessionStore    *sessions.CookieStore
	sessionName     string
}

func NewActivityHandler(
	activityService *service.ActivityService,
	sessionStore *sessions.CookieStore,
	sessionName string,
) *ActivityHandler {
	return &ActivityHandler{
		activityService: activityService,
		sessionStore:    sessionStore,
		sessionName:     sessionName,
	}
}

func (h *ActivityHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/dashboard/stream/activities", h.GetRecentActivities).Methods("GET")
}

func (h *ActivityHandler) GetRecentActivities(w http.ResponseWriter, r *http.Request) {
	user := h.requireUser(w, r)
	if user == nil {
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	activities, err := h.activityService.GetRecentActivities(r.Context(), user.ID, limit)
	if err != nil {
		log.Printf("Fehler beim Laden der Activities: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, activities)
}

func (h *ActivityHandler) requireUser(w http.ResponseWriter, r *http.Request) *UserSession {
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

func (h *ActivityHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}