package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"demo/backend-go/internal/service"
)

type SubathonHandler struct {
	subathonService *service.SubathonService
	sessionStore    sessions.Store
	sessionName     string
}

func NewSubathonHandler(subathonService *service.SubathonService, sessionStore sessions.Store, sessionName string) *SubathonHandler {
	return &SubathonHandler{
		subathonService: subathonService,
		sessionStore:    sessionStore,
		sessionName:     sessionName,
	}
}

func (h *SubathonHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/subathon/state", h.GetState).Methods(http.MethodGet)
	router.HandleFunc("/api/subathon/start", h.Start).Methods(http.MethodPost)
	router.HandleFunc("/api/subathon/pause", h.Pause).Methods(http.MethodPost)
	router.HandleFunc("/api/subathon/reset", h.Reset).Methods(http.MethodPost)
	router.HandleFunc("/api/subathon/settings", h.UpdateSettings).Methods(http.MethodPut)
	// Public (no session) - read-only, for the OBS Browser Source overlay,
	// which can't carry a login cookie. Only exposes what's already visible
	// live on stream anyway (timer, sub/bit counts, log).
	router.HandleFunc("/api/subathon/overlay/{userId}", h.GetOverlayState).Methods(http.MethodGet)
}

func (h *SubathonHandler) userID(r *http.Request) (string, bool) {
	session, err := h.sessionStore.Get(r, h.sessionName)
	if err != nil {
		return "", false
	}
	userID, ok := session.Values["user_id"].(string)
	return userID, ok && userID != ""
}

func (h *SubathonHandler) GetState(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.userID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	state, err := h.subathonService.GetState(r.Context(), userID)
	if err != nil {
		log.Printf("Failed to get subathon state: %v", err)
		http.Error(w, "Failed to get state", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}

func (h *SubathonHandler) Start(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.userID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	state, err := h.subathonService.StartTimer(r.Context(), userID)
	if err != nil {
		log.Printf("Failed to start subathon timer: %v", err)
		http.Error(w, "Failed to start timer", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}

func (h *SubathonHandler) Pause(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.userID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	state, err := h.subathonService.PauseTimer(r.Context(), userID)
	if err != nil {
		log.Printf("Failed to pause subathon timer: %v", err)
		http.Error(w, "Failed to pause timer", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}

func (h *SubathonHandler) Reset(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.userID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	state, err := h.subathonService.ResetTimer(r.Context(), userID)
	if err != nil {
		log.Printf("Failed to reset subathon timer: %v", err)
		http.Error(w, "Failed to reset timer", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}

func (h *SubathonHandler) GetOverlayState(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["userId"]
	if userID == "" {
		http.Error(w, "userId required", http.StatusBadRequest)
		return
	}

	state, err := h.subathonService.GetState(r.Context(), userID)
	if err != nil {
		log.Printf("Failed to get subathon overlay state: %v", err)
		http.Error(w, "Failed to get state", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}

func (h *SubathonHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.userID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var input struct {
		InitialTime *int `json:"initialTime"`
		SubTime     *int `json:"subTime"`
		BitsTime    *int `json:"bitsTime"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	state, err := h.subathonService.UpdateSettings(r.Context(), userID, input.InitialTime, input.SubTime, input.BitsTime)
	if err != nil {
		log.Printf("Failed to update subathon settings: %v", err)
		http.Error(w, "Failed to update settings", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}
