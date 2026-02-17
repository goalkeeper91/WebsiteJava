package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/service"
)

type N8NHandler struct {
	n8nService  *service.N8NService
	sessionStore *sessions.CookieStore
	sessionName  string
}

func NewN8NHandler(
	n8nService *service.N8NService,
	sessionStore *sessions.CookieStore,
	sessionName string,
) *N8NHandler {
	return &N8NHandler{
		n8nService:  n8nService,
		sessionStore: sessionStore,
		sessionName:  sessionName,
	}
}

func (h *N8NHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/n8n/integration", h.GetIntegration).Methods("GET")
	router.HandleFunc("/api/n8n/integration/enable", h.Enable).Methods("POST")
	router.HandleFunc("/api/n8n/integration/disable", h.Disable).Methods("POST")
	router.HandleFunc("/api/n8n/integration/test", h.TestWebhook).Methods("POST")
}

// GetIntegration returns the current n8n integration status
func (h *N8NHandler) GetIntegration(w http.ResponseWriter, r *http.Request) {
	user := h.requireUser(w, r)
	if user == nil {
		return
	}

	integration, err := h.n8nService.GetIntegration(r.Context(), user.ID)
	if err != nil {
		log.Printf("Fehler beim Laden der n8n Integration: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	// Never expose the API key in the response
	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"id":              integration.ID,
		"enabled":         integration.Enabled,
		"webhookBaseUrl":  integration.WebhookBaseURL,
		"workflowsUsed":   integration.WorkflowsUsed,
		"votesThisMonth":  integration.VotesThisMonth,
		"lastResetAt":     integration.LastResetAt,
		"isReady":         integration.IsReady(),
	})
}

// Enable activates n8n integration for the user
func (h *N8NHandler) Enable(w http.ResponseWriter, r *http.Request) {
	user := h.requireUser(w, r)
	if user == nil {
		return
	}

	var req struct {
		WebhookBaseURL string `json:"webhook_base_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
		return
	}

	if req.WebhookBaseURL == "" {
		http.Error(w, "webhook_base_url ist erforderlich", http.StatusBadRequest)
		return
	}

	integration, err := h.n8nService.Enable(r.Context(), user.ID, req.WebhookBaseURL)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"enabled":        integration.Enabled,
		"webhookBaseUrl": integration.WebhookBaseURL,
		"isReady":        integration.IsReady(),
	})
}

// Disable deactivates n8n integration
func (h *N8NHandler) Disable(w http.ResponseWriter, r *http.Request) {
	user := h.requireUser(w, r)
	if user == nil {
		return
	}

	if err := h.n8nService.Disable(r.Context(), user.ID); err != nil {
		log.Printf("Fehler beim Deaktivieren der n8n Integration: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "n8n Integration deaktiviert"})
}

// TestWebhook tests if n8n is reachable
func (h *N8NHandler) TestWebhook(w http.ResponseWriter, r *http.Request) {
	user := h.requireUser(w, r)
	if user == nil {
		return
	}

	if err := h.n8nService.TestWebhook(r.Context(), user.ID); err != nil {
		h.handleError(w, err)
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (h *N8NHandler) requireUser(w http.ResponseWriter, r *http.Request) *UserSession {
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

func (h *N8NHandler) handleError(w http.ResponseWriter, err error) {
	switch err {
	case domain.ErrFeatureNotAvailable:
		http.Error(w, "n8n Integration erfordert ein Pro oder Premium Abo", http.StatusForbidden)
	case domain.ErrN8NIntegrationNotReady:
		http.Error(w, "n8n Integration ist nicht konfiguriert oder deaktiviert", http.StatusBadRequest)
	case domain.ErrN8NWebhookFailed:
		http.Error(w, "n8n Webhook nicht erreichbar", http.StatusBadGateway)
	case domain.ErrN8NWorkflowNotFound:
		http.Error(w, "n8n Workflow nicht gefunden", http.StatusNotFound)
	default:
		log.Printf("n8n Fehler: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
	}
}

func (h *N8NHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}