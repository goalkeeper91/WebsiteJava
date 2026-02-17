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

type WorkflowTemplateHandler struct {
	templateService *service.WorkflowTemplateService
	sessionStore    *sessions.CookieStore
	sessionName     string
}

func NewWorkflowTemplateHandler(
	templateService *service.WorkflowTemplateService,
	sessionStore *sessions.CookieStore,
	sessionName string,
) *WorkflowTemplateHandler {
	return &WorkflowTemplateHandler{
		templateService: templateService,
		sessionStore:    sessionStore,
		sessionName:     sessionName,
	}
}

func (h *WorkflowTemplateHandler) RegisterRoutes(router *mux.Router) {
	// Public: listing without workflow JSON
	router.HandleFunc("/api/workflows/templates", h.GetPublicTemplates).Methods("GET")

	// Authenticated: get templates accessible by user's tier
	router.HandleFunc("/api/workflows/templates/available", h.GetAvailableTemplates).Methods("GET")
	router.HandleFunc("/api/workflows/templates/{id}", h.GetTemplate).Methods("GET")
}

// GetPublicTemplates returns all public templates as summaries (no auth needed, no JSON)
func (h *WorkflowTemplateHandler) GetPublicTemplates(w http.ResponseWriter, r *http.Request) {
	summaries, err := h.templateService.GetPublicTemplates(r.Context())
	if err != nil {
		log.Printf("Fehler beim Laden der öffentlichen Templates: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, summaries)
}

// GetAvailableTemplates returns full templates accessible for the user's subscription tier
func (h *WorkflowTemplateHandler) GetAvailableTemplates(w http.ResponseWriter, r *http.Request) {
	user := h.requireUser(w, r)
	if user == nil {
		return
	}

	// Optional: filter by category
	category := domain.WorkflowCategory(r.URL.Query().Get("category"))

	var templates interface{}
	var err error

	if category != "" {
		templates, err = h.templateService.GetTemplatesByCategory(r.Context(), user.ID, category)
	} else {
		templates, err = h.templateService.GetAvailableTemplates(r.Context(), user.ID)
	}

	if err != nil {
		h.handleError(w, err)
		return
	}

	h.respondJSON(w, http.StatusOK, templates)
}

// GetTemplate returns a specific template with full workflow JSON (tier check included)
func (h *WorkflowTemplateHandler) GetTemplate(w http.ResponseWriter, r *http.Request) {
	user := h.requireUser(w, r)
	if user == nil {
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Ungültige Template ID", http.StatusBadRequest)
		return
	}

	template, err := h.templateService.GetTemplate(r.Context(), user.ID, id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.respondJSON(w, http.StatusOK, template)
}

func (h *WorkflowTemplateHandler) requireUser(w http.ResponseWriter, r *http.Request) *UserSession {
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

func (h *WorkflowTemplateHandler) handleError(w http.ResponseWriter, err error) {
	switch err {
	case domain.ErrTemplateNotFound:
		http.Error(w, "Template nicht gefunden", http.StatusNotFound)
	case domain.ErrTemplateNotAllowed:
		http.Error(w, "Dieses Template erfordert ein höheres Abo", http.StatusForbidden)
	case domain.ErrFeatureNotAvailable:
		http.Error(w, "Diese Funktion ist in deinem aktuellen Abo nicht verfügbar", http.StatusForbidden)
	default:
		log.Printf("Template Fehler: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
	}
}

func (h *WorkflowTemplateHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}