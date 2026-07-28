package handler

import (
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"demo/backend-go/internal/service"
)

// PaddleWebhookHandler receives Paddle Billing webhook deliveries (checkout
// completions, subscription lifecycle changes). Public route, no session
// auth - security comes entirely from the Paddle-Signature HMAC check, same
// posture as ActivityWebhookHandler/SubathonWebhookHandler for Twitch
// EventSub.
type PaddleWebhookHandler struct {
	paddleService *service.PaddleService
}

func NewPaddleWebhookHandler(paddleService *service.PaddleService) *PaddleWebhookHandler {
	return &PaddleWebhookHandler{paddleService: paddleService}
}

func (h *PaddleWebhookHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/billing/paddle/webhook", h.HandleWebhook).Methods(http.MethodPost)
}

func (h *PaddleWebhookHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	signature := r.Header.Get("Paddle-Signature")
	if err := h.paddleService.VerifyWebhookSignature(body, signature); err != nil {
		log.Printf("Paddle webhook: invalid signature, rejecting: %v", err)
		http.Error(w, "invalid signature", http.StatusForbidden)
		return
	}

	if err := h.paddleService.HandleWebhookEvent(r.Context(), body); err != nil {
		log.Printf("Paddle webhook: fehler bei der Verarbeitung: %v", err)
		http.Error(w, "processing failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
