package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/sessions"
    "github.com/gorilla/mux"
	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/service"
)

type JoinToCreateHandler struct {
	jtcService   *service.JoinToCreateService
	sessionStore sessions.Store
	sessionName  string
}

func NewJoinToCreateHandler(
	jtcService *service.JoinToCreateService,
	sessionStore sessions.Store,
	sessionName string,
) *JoinToCreateHandler {
	return &JoinToCreateHandler{
		jtcService:   jtcService,
		sessionStore: sessionStore,
		sessionName:  sessionName,
	}
}

func (h *JoinToCreateHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/discord/join-to-create", h.handleConfigs)
	router.HandleFunc("/api/discord/join-to-create/", h.handleConfigByID)
}

// handleConfigs handles listing and creating configs
func (h *JoinToCreateHandler) handleConfigs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.ListConfigs(w, r)
	case http.MethodPost:
		h.CreateConfig(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleConfigByID handles operations on specific config
func (h *JoinToCreateHandler) handleConfigByID(w http.ResponseWriter, r *http.Request) {
	// Extract config ID from path
	path := strings.TrimPrefix(r.URL.Path, "/api/discord/join-to-create/")

	if path == "" {
		http.Error(w, "Config ID required", http.StatusBadRequest)
		return
	}

	configID, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		http.Error(w, "Invalid config ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetConfig(w, r, configID)
	case http.MethodPut:
		h.UpdateConfig(w, r, configID)
	case http.MethodDelete:
		h.DeleteConfig(w, r, configID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// ListConfigs lists user's Join-to-Create configs
func (h *JoinToCreateHandler) ListConfigs(w http.ResponseWriter, r *http.Request) {
	// Get user ID from session
	session, err := h.sessionStore.Get(r, h.sessionName)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := session.Values["user_id"].(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get configs
	configs, err := h.jtcService.ListConfigs(r.Context(), userID)
	if err != nil {
		log.Printf("Failed to list configs: %v", err)
		http.Error(w, "Failed to list configs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(configs)
}

// GetConfig gets a specific config
func (h *JoinToCreateHandler) GetConfig(w http.ResponseWriter, r *http.Request, configID int64) {
	// Get user ID from session
	session, err := h.sessionStore.Get(r, h.sessionName)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := session.Values["user_id"].(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get config
	config, err := h.jtcService.GetConfig(r.Context(), configID, userID)
	if err != nil {
		log.Printf("Failed to get config: %v", err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}

// CreateConfig creates a new Join-to-Create config
func (h *JoinToCreateHandler) CreateConfig(w http.ResponseWriter, r *http.Request) {
	// Get user ID from session
	session, err := h.sessionStore.Get(r, h.sessionName)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := session.Values["user_id"].(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse request body
	var input domain.JoinToCreateConfigCreateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create config
	config, err := h.jtcService.CreateConfig(r.Context(), input, userID)
	if err != nil {
		log.Printf("Failed to create config: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(config)
}

// UpdateConfig updates a Join-to-Create config
func (h *JoinToCreateHandler) UpdateConfig(w http.ResponseWriter, r *http.Request, configID int64) {
	// Get user ID from session
	session, err := h.sessionStore.Get(r, h.sessionName)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := session.Values["user_id"].(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse request body
	var input domain.JoinToCreateConfigUpdateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update config
	config, err := h.jtcService.UpdateConfig(r.Context(), configID, userID, input)
	if err != nil {
		log.Printf("Failed to update config: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}

// DeleteConfig deletes a Join-to-Create config
func (h *JoinToCreateHandler) DeleteConfig(w http.ResponseWriter, r *http.Request, configID int64) {
	// Get user ID from session
	session, err := h.sessionStore.Get(r, h.sessionName)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := session.Values["user_id"].(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Delete config
	err = h.jtcService.DeleteConfig(r.Context(), configID, userID)
	if err != nil {
		log.Printf("Failed to delete config: %v", err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{
		"success": true,
	})
}