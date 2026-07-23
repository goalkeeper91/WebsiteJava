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

type ScheduledMessageHandler struct {
	messageService *service.ScheduledMessageService
	sessionStore   *sessions.CookieStore
	sessionName    string
}

func NewScheduledMessageHandler(
	messageService *service.ScheduledMessageService,
	sessionStore *sessions.CookieStore,
	sessionName string,
) *ScheduledMessageHandler {
	return &ScheduledMessageHandler{
		messageService: messageService,
		sessionStore:   sessionStore,
		sessionName:    sessionName,
	}
}

func (h *ScheduledMessageHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/dashboard/scheduled-messages", h.GetMessages).Methods("GET")
	router.HandleFunc("/api/dashboard/scheduled-messages/{id}", h.GetMessage).Methods("GET")
	router.HandleFunc("/api/dashboard/scheduled-messages", h.CreateMessage).Methods("POST")
	router.HandleFunc("/api/dashboard/scheduled-messages/{id}", h.UpdateMessage).Methods("PUT")
	router.HandleFunc("/api/dashboard/scheduled-messages/{id}", h.DeleteMessage).Methods("DELETE")
	router.HandleFunc("/api/dashboard/scheduled-messages/{id}/toggle", h.ToggleMessage).Methods("PATCH")
}

// ============================================================
// REQUEST TYPES
// ============================================================

type CreateScheduledMessageRequest struct {
	Message         string `json:"message"`
	CommandID       *int64 `json:"command_id,omitempty"`
	IntervalSeconds int    `json:"interval_seconds"`
}

type UpdateScheduledMessageRequest struct {
	Message         *string `json:"message,omitempty"`
	IntervalSeconds *int    `json:"interval_seconds,omitempty"`
	Enabled         *bool   `json:"enabled,omitempty"`
	OnlyWhenLive    *bool   `json:"only_when_live,omitempty"`
}

type ToggleScheduledMessageRequest struct {
	Enabled bool `json:"enabled"`
}

// ============================================================
// HANDLERS
// ============================================================

func (h *ScheduledMessageHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
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

	messages, total, err := h.messageService.GetMessages(r.Context(), user.ID, pageSize, offset)
	if err != nil {
		h.handleError(w, err)
		return
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	h.respondJSON(w, http.StatusOK, PaginatedResponse{
		Data:       messages,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

func (h *ScheduledMessageHandler) GetMessage(w http.ResponseWriter, r *http.Request) {
	user := h.requireUser(w, r)
	if user == nil {
		return
	}

	id, err := h.parseID(r)
	if err != nil {
		http.Error(w, "Ungültige ID", http.StatusBadRequest)
		return
	}

	message, err := h.messageService.GetMessageByID(r.Context(), user.ID, id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.respondJSON(w, http.StatusOK, message)
}

func (h *ScheduledMessageHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	user := h.requireUser(w, r)
	if user == nil {
		return
	}

	var req CreateScheduledMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
		return
	}

	var message *domain.ScheduledMessage
	var err error
	if req.CommandID != nil {
		message, err = h.messageService.CreateCommandSchedule(r.Context(), user.ID, *req.CommandID, req.IntervalSeconds)
	} else {
		message, err = h.messageService.CreateMessage(r.Context(), user.ID, req.Message, req.IntervalSeconds)
	}
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.respondJSON(w, http.StatusCreated, message)
}

func (h *ScheduledMessageHandler) UpdateMessage(w http.ResponseWriter, r *http.Request) {
	user := h.requireUser(w, r)
	if user == nil {
		return
	}

	id, err := h.parseID(r)
	if err != nil {
		http.Error(w, "Ungültige ID", http.StatusBadRequest)
		return
	}

	var req UpdateScheduledMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
		return
	}

	message, err := h.messageService.UpdateMessage(
		r.Context(), user.ID, id, req.Message, req.IntervalSeconds, req.Enabled, req.OnlyWhenLive,
	)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.respondJSON(w, http.StatusOK, message)
}

func (h *ScheduledMessageHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	user := h.requireUser(w, r)
	if user == nil {
		return
	}

	id, err := h.parseID(r)
	if err != nil {
		http.Error(w, "Ungültige ID", http.StatusBadRequest)
		return
	}

	if err := h.messageService.DeleteMessage(r.Context(), user.ID, id); err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ScheduledMessageHandler) ToggleMessage(w http.ResponseWriter, r *http.Request) {
	user := h.requireUser(w, r)
	if user == nil {
		return
	}

	id, err := h.parseID(r)
	if err != nil {
		http.Error(w, "Ungültige ID", http.StatusBadRequest)
		return
	}

	var req ToggleScheduledMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
		return
	}

	message, err := h.messageService.ToggleMessage(r.Context(), user.ID, id, req.Enabled)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.respondJSON(w, http.StatusOK, message)
}

// ============================================================
// HELPERS
// ============================================================

func (h *ScheduledMessageHandler) parseID(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	return strconv.ParseInt(vars["id"], 10, 64)
}

func (h *ScheduledMessageHandler) requireUser(w http.ResponseWriter, r *http.Request) *UserSession {
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

func (h *ScheduledMessageHandler) handleError(w http.ResponseWriter, err error) {
	switch err {
	case domain.ErrScheduledMessageNotFound:
		http.Error(w, "Automatisierte Nachricht nicht gefunden", http.StatusNotFound)
	case domain.ErrEmptyMessage:
		http.Error(w, "Nachricht darf nicht leer sein", http.StatusBadRequest)
	case domain.ErrIntervalTooShort:
		http.Error(w, "Intervall muss mindestens 60 Sekunden betragen", http.StatusBadRequest)
	case domain.ErrChannelNotRegistered:
		http.Error(w, "Channel nicht registriert. Bitte erneut einloggen.", http.StatusBadRequest)
	case domain.ErrCommandNotFound:
		http.Error(w, "Command nicht gefunden", http.StatusNotFound)
	case domain.ErrFeatureNotAvailable:
		http.Error(w, "Nur einfache Commands (ohne n8n) können mit einem Timer verknüpft werden", http.StatusBadRequest)
	default:
		log.Printf("Scheduled-Message Fehler: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
	}
}

func (h *ScheduledMessageHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
