package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/paddle"
	"demo/backend-go/internal/service"
)

type SubscriptionHandler struct {
	subscriptionService *service.SubscriptionService
	paddleClient        *paddle.Client
	sessionStore        *sessions.CookieStore
	sessionName         string
}

func NewSubscriptionHandler(
	subscriptionService *service.SubscriptionService,
	paddleClient *paddle.Client,
	sessionStore *sessions.CookieStore,
	sessionName string,
) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: subscriptionService,
		paddleClient:        paddleClient,
		sessionStore:        sessionStore,
		sessionName:         sessionName,
	}
}

func (h *SubscriptionHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/subscription", h.GetSubscription).Methods("GET")
	router.HandleFunc("/api/subscription/tiers", h.GetTiers).Methods("GET")
	router.HandleFunc("/api/subscription/upgrade", h.Upgrade).Methods("POST")
	router.HandleFunc("/api/subscription/cancel", h.Cancel).Methods("POST")
	router.HandleFunc("/api/subscription/portal-link", h.GetPortalLink).Methods("GET")
}

// GetSubscription returns the current user's subscription with tier details
func (h *SubscriptionHandler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	user := h.requireUser(w, r)
	if user == nil {
		return
	}

	sub, err := h.subscriptionService.GetSubscription(r.Context(), user.ID)
	if err != nil {
		log.Printf("Fehler beim Laden der Subscription: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, sub)
}

// GetTiers returns all available subscription tiers (public, no auth needed)
func (h *SubscriptionHandler) GetTiers(w http.ResponseWriter, r *http.Request) {
	tiers, err := h.subscriptionService.GetTiers(r.Context())
	if err != nil {
		log.Printf("Fehler beim Laden der Tiers: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, tiers)
}

// Upgrade only allows a self-service downgrade to Free (no payment
// involved). Any paid tier now requires a real Paddle checkout - the
// resulting transaction.completed webhook is what actually sets a paid
// tier (see PaddleService.handleTransactionCompleted), never this endpoint.
func (h *SubscriptionHandler) Upgrade(w http.ResponseWriter, r *http.Request) {
	user := h.requireUser(w, r)
	if user == nil {
		return
	}

	var req struct {
		TierID domain.TierID `json:"tier_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
		return
	}

	if req.TierID == "" {
		http.Error(w, "tier_id ist erforderlich", http.StatusBadRequest)
		return
	}

	if req.TierID != domain.TierFree {
		http.Error(w, "Kostenpflichtige Tiers erfordern einen Paddle-Checkout, siehe /api/subscription/tiers", http.StatusBadRequest)
		return
	}

	sub, err := h.subscriptionService.Upgrade(r.Context(), user.ID, req.TierID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.respondJSON(w, http.StatusOK, sub)
}

// GetPortalLink mints a fresh Paddle Customer Portal session link so the
// user can manage/cancel their subscription or change payment method via
// Paddle's own self-service UI (including its Cancellation Flow).
func (h *SubscriptionHandler) GetPortalLink(w http.ResponseWriter, r *http.Request) {
	user := h.requireUser(w, r)
	if user == nil {
		return
	}

	sub, err := h.subscriptionService.GetSubscription(r.Context(), user.ID)
	if err != nil {
		log.Printf("Fehler beim Laden der Subscription: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	if sub.PaddleCustomerID == nil || *sub.PaddleCustomerID == "" {
		http.Error(w, "Kein Paddle-Kunde vorhanden - noch kein kostenpflichtiges Abo abgeschlossen", http.StatusNotFound)
		return
	}

	url, err := h.paddleClient.GetOrCreatePortalSession(r.Context(), *sub.PaddleCustomerID)
	if err != nil {
		log.Printf("Fehler beim Erstellen der Paddle-Portal-Session: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"url": url})
}

// Cancel cancels the user's current subscription
func (h *SubscriptionHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	user := h.requireUser(w, r)
	if user == nil {
		return
	}

	if err := h.subscriptionService.Cancel(r.Context(), user.ID); err != nil {
		h.handleError(w, err)
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Subscription erfolgreich gekündigt"})
}

func (h *SubscriptionHandler) requireUser(w http.ResponseWriter, r *http.Request) *UserSession {
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

func (h *SubscriptionHandler) handleError(w http.ResponseWriter, err error) {
	switch err {
	case domain.ErrFeatureNotAvailable:
		http.Error(w, "Diese Funktion ist in deinem aktuellen Abo nicht verfügbar", http.StatusForbidden)
	case domain.ErrSubscriptionExpired:
		http.Error(w, "Dein Abo ist abgelaufen", http.StatusPaymentRequired)
	case domain.ErrCommandLimitReached:
		http.Error(w, "Command-Limit deines Abos erreicht", http.StatusForbidden)
	case domain.ErrWorkflowLimitReached:
		http.Error(w, "Workflow-Limit deines Abos erreicht", http.StatusForbidden)
	default:
		log.Printf("Subscription Fehler: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
	}
}

func (h *SubscriptionHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}