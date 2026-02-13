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
}

func NewChatCommandHandler(
	commandService *service.ChatCommandService,
	sessionStore *sessions.CookieStore,
	sessionName string,
) *ChatCommandHandler {
	return &ChatCommandHandler{
		commandService: commandService,
		sessionStore:   sessionStore,
		sessionName:    sessionName,
	}
}

// RegisterRoutes registriert alle Command-Routen
func (h *ChatCommandHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/dashboard/commands", h.GetCommands).Methods("GET")
	router.HandleFunc("/api/dashboard/commands/{id}", h.GetCommand).Methods("GET")
	router.HandleFunc("/api/dashboard/commands", h.CreateCommand).Methods("POST")
	router.HandleFunc("/api/dashboard/commands/{id}", h.UpdateCommand).Methods("PUT")
	router.HandleFunc("/api/dashboard/commands/{id}", h.DeleteCommand).Methods("DELETE")
	router.HandleFunc("/api/dashboard/commands/{id}/toggle", h.ToggleCommand).Methods("PATCH")
}

type CreateCommandRequest struct {
	Trigger  string `json:"trigger"`
	Response string `json:"response"`
	Cooldown int    `json:"cooldown"`
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

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

func (h *ChatCommandHandler) GetCommands(w http.ResponseWriter, r *http.Request) {
	user := h.requireUser(w, r)
	if user == nil {
		return
	}

	search := r.URL.Query().Get("search")
	enabledStr := r.URL.Query().Get("enabled")
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

	if search != "" {
		commands, total, err = h.commandService.SearchCommands(ctx, user.ID, search, pageSize, offset)
	} else if enabledStr != "" {
		enabled := enabledStr == "true"
		commands, total, err = h.commandService.GetCommandsByStatus(ctx, user.ID, enabled, pageSize, offset)
	} else {
		commands, total, err = h.commandService.GetCommands(ctx, user.ID, pageSize, offset)
	}

	if err != nil {
		h.handleError(w, err)
		return
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	response := PaginatedResponse{
		Data:       commands,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	h.respondJSON(w, http.StatusOK, response)
}

func (h *ChatCommandHandler) GetCommand(w http.ResponseWriter, r *http.Request) {
	user := h.requireUser(w, r)
	if user == nil {
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Ungültige ID", http.StatusBadRequest)
		return
	}

	command, err := h.commandService.GetCommandByID(r.Context(), user.ID, id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.respondJSON(w, http.StatusOK, command)
}

func (h *ChatCommandHandler) CreateCommand(w http.ResponseWriter, r *http.Request) {
	user := h.requireUser(w, r)
	if user == nil {
		return
	}

	var req CreateCommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
		return
	}

	command, err := h.commandService.CreateCommand(
		r.Context(),
		user.ID,
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

func (h *ChatCommandHandler) UpdateCommand(w http.ResponseWriter, r *http.Request) {
	user := h.requireUser(w, r)
	if user == nil {
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
		user.ID,
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
	user := h.requireUser(w, r)
	if user == nil {
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Ungültige ID", http.StatusBadRequest)
		return
	}

	if err := h.commandService.DeleteCommand(r.Context(), user.ID, id); err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ChatCommandHandler) ToggleCommand(w http.ResponseWriter, r *http.Request) {
	user := h.requireUser(w, r)
	if user == nil {
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

	command, err := h.commandService.ToggleCommand(r.Context(), user.ID, id, req.Enabled)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.respondJSON(w, http.StatusOK, command)
}

func (h *ChatCommandHandler) requireUser(w http.ResponseWriter, r *http.Request) *UserSession {
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
	default:
		log.Printf("Fehler: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
	}
}

func (h *ChatCommandHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}