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

type ChatCommandHandler struct {
	commandService *service.ChatCommandService
	sessionStore   *sessions.CookieStore
	sessionName    string
	teamService    *service.TeamService
}

func NewChatCommandHandler(
	commandService *service.ChatCommandService,
	sessionStore *sessions.CookieStore,
	sessionName string,
	teamService *service.TeamService,
) *ChatCommandHandler {
	return &ChatCommandHandler{
		commandService: commandService,
		sessionStore:   sessionStore,
		sessionName:    sessionName,
		teamService:    teamService,
	}
}

func (h *ChatCommandHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/dashboard/commands", h.GetCommands).Methods("GET")
	router.HandleFunc("/api/dashboard/commands/{id}", h.GetCommand).Methods("GET")
	router.HandleFunc("/api/dashboard/commands", h.CreateCommand).Methods("POST")
	router.HandleFunc("/api/dashboard/commands/advanced", h.CreateAdvancedCommand).Methods("POST")
	router.HandleFunc("/api/dashboard/commands/{id}", h.UpdateCommand).Methods("PUT")
	router.HandleFunc("/api/dashboard/commands/{id}", h.DeleteCommand).Methods("DELETE")
	router.HandleFunc("/api/dashboard/commands/{id}/toggle", h.ToggleCommand).Methods("PATCH")
}

// ============================================================
// REQUEST TYPES
// ============================================================

type CreateCommandRequest struct {
	Trigger  string `json:"trigger"`
	Response string `json:"response"`
	Cooldown int    `json:"cooldown"`
}

type CreateAdvancedCommandRequest struct {
	Trigger    string `json:"trigger"`
	WorkflowID string `json:"workflow_id"`
	Cooldown   int    `json:"cooldown"`
}

type UpdateCommandRequest struct {
	Trigger  *string `json:"trigger,omitempty"`
	Response *string `json:"response,omitempty"`
	Cooldown *int    `json:"cooldown,omitempty"`
	Enabled  *bool   `json:"enabled,omitempty"`
}

type ToggleCommandRequest struct {
	Enabled bool `json:"enabled"`
}

// ============================================================
// HANDLERS
// ============================================================

func (h *ChatCommandHandler) GetCommands(w http.ResponseWriter, r *http.Request) {
	_, channelID, ok := resolveEffectiveTwitchID(w, r, h.sessionStore, h.sessionName, h.teamService)
	if !ok {
		return
	}

	search := r.URL.Query().Get("search")
	enabledStr := r.URL.Query().Get("enabled")
	commandType := r.URL.Query().Get("type") // "simple" or "advanced"
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

	var commands []*domain.ChatCommand
	var total int64
	var err error

	switch {
	case search != "":
		commands, total, err = h.commandService.SearchCommands(ctx, channelID, search, pageSize, offset)
	case enabledStr != "":
		enabled := enabledStr == "true"
		commands, total, err = h.commandService.GetCommandsByStatus(ctx, channelID, enabled, pageSize, offset)
	case commandType == "simple" || commandType == "advanced":
		commands, total, err = h.commandService.GetCommandsByType(ctx, channelID, domain.CommandType(commandType), pageSize, offset)
	default:
		commands, total, err = h.commandService.GetCommands(ctx, channelID, pageSize, offset)
	}

	if err != nil {
		h.handleError(w, err)
		return
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	h.respondJSON(w, http.StatusOK, PaginatedResponse{
		Data:       commands,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

func (h *ChatCommandHandler) GetCommand(w http.ResponseWriter, r *http.Request) {
	_, channelID, ok := resolveEffectiveTwitchID(w, r, h.sessionStore, h.sessionName, h.teamService)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Ungültige ID", http.StatusBadRequest)
		return
	}

	command, err := h.commandService.GetCommandByID(r.Context(), channelID, id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.respondJSON(w, http.StatusOK, command)
}

// CreateCommand creates a simple text command (free tier)
func (h *ChatCommandHandler) CreateCommand(w http.ResponseWriter, r *http.Request) {
	_, channelID, ok := resolveEffectiveTwitchID(w, r, h.sessionStore, h.sessionName, h.teamService)
	if !ok {
		return
	}

	var req CreateCommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
		return
	}

	command, err := h.commandService.CreateCommand(
		r.Context(),
		channelID,
		req.Trigger,
		req.Response,
		req.Cooldown,
	)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.respondJSON(w, http.StatusCreated, command)
}

// CreateAdvancedCommand creates an n8n-powered command (pro/premium)
func (h *ChatCommandHandler) CreateAdvancedCommand(w http.ResponseWriter, r *http.Request) {
	_, channelID, ok := resolveEffectiveTwitchID(w, r, h.sessionStore, h.sessionName, h.teamService)
	if !ok {
		return
	}

	var req CreateAdvancedCommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
		return
	}

	if req.WorkflowID == "" {
		http.Error(w, "workflow_id ist erforderlich", http.StatusBadRequest)
		return
	}

	command, err := h.commandService.CreateAdvancedCommand(
		r.Context(),
		channelID,
		req.Trigger,
		req.WorkflowID,
		req.Cooldown,
	)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.respondJSON(w, http.StatusCreated, command)
}

func (h *ChatCommandHandler) UpdateCommand(w http.ResponseWriter, r *http.Request) {
	_, channelID, ok := resolveEffectiveTwitchID(w, r, h.sessionStore, h.sessionName, h.teamService)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Ungültige ID", http.StatusBadRequest)
		return
	}

	var req UpdateCommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
		return
	}

	command, err := h.commandService.UpdateCommand(
		r.Context(),
		channelID,
		id,
		req.Trigger,
		req.Response,
		req.Cooldown,
		req.Enabled,
	)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.respondJSON(w, http.StatusOK, command)
}

func (h *ChatCommandHandler) DeleteCommand(w http.ResponseWriter, r *http.Request) {
	_, channelID, ok := resolveEffectiveTwitchID(w, r, h.sessionStore, h.sessionName, h.teamService)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Ungültige ID", http.StatusBadRequest)
		return
	}

	if err := h.commandService.DeleteCommand(r.Context(), channelID, id); err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ChatCommandHandler) ToggleCommand(w http.ResponseWriter, r *http.Request) {
	_, channelID, ok := resolveEffectiveTwitchID(w, r, h.sessionStore, h.sessionName, h.teamService)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Ungültige ID", http.StatusBadRequest)
		return
	}

	var req ToggleCommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
		return
	}

	command, err := h.commandService.ToggleCommand(r.Context(), channelID, id, req.Enabled)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.respondJSON(w, http.StatusOK, command)
}

// ============================================================
// HELPERS
// ============================================================

func (h *ChatCommandHandler) handleError(w http.ResponseWriter, err error) {
	switch err {
	case domain.ErrCommandNotFound:
		http.Error(w, "Command nicht gefunden", http.StatusNotFound)
	case domain.ErrCommandAlreadyExists:
		http.Error(w, "Command existiert bereits", http.StatusConflict)
	case domain.ErrEmptyTrigger:
		http.Error(w, "Trigger darf nicht leer sein", http.StatusBadRequest)
	case domain.ErrTriggerTooLong:
		http.Error(w, "Trigger zu lang (max 100 Zeichen)", http.StatusBadRequest)
	case domain.ErrTriggerTaken:
		http.Error(w, "Trigger bereits vergeben", http.StatusConflict)
	case domain.ErrChannelNotRegistered:
		http.Error(w, "Channel nicht registriert. Bitte erneut einloggen.", http.StatusBadRequest)
	case domain.ErrCommandLimitReached:
		http.Error(w, "Command-Limit deines Abos erreicht", http.StatusForbidden)
	case domain.ErrFeatureNotAvailable:
		http.Error(w, "Advanced Commands erfordern ein Pro oder Premium Abo", http.StatusForbidden)
	case domain.ErrN8NIntegrationNotReady:
		http.Error(w, "n8n Integration ist nicht konfiguriert", http.StatusBadRequest)
	default:
		log.Printf("Command Fehler: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
	}
}

func (h *ChatCommandHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}