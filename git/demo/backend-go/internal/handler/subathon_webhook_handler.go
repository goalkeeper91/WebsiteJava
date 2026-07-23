package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"demo/backend-go/internal/service"
)

// SubathonWebhookHandler receives Twitch EventSub webhook deliveries for
// sub/gift-sub/cheer events. Public route - Twitch calls this directly, no
// session cookie exists to check, so authenticity is verified via HMAC
// signature instead (see verifySignature).
type SubathonWebhookHandler struct {
	subathonService *service.SubathonService
	webhookSecret   string
}

func NewSubathonWebhookHandler(subathonService *service.SubathonService, webhookSecret string) *SubathonWebhookHandler {
	return &SubathonWebhookHandler{
		subathonService: subathonService,
		webhookSecret:   webhookSecret,
	}
}

func (h *SubathonWebhookHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/subathon/eventsub/callback", h.HandleCallback).Methods(http.MethodPost)
}

type eventSubEnvelope struct {
	Challenge    string `json:"challenge"`
	Subscription struct {
		Type      string `json:"type"`
		Condition struct {
			BroadcasterUserID string `json:"broadcaster_user_id"`
		} `json:"condition"`
	} `json:"subscription"`
	Event json.RawMessage `json:"event"`
}

func (h *SubathonWebhookHandler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	messageID := r.Header.Get("Twitch-Eventsub-Message-Id")
	timestamp := r.Header.Get("Twitch-Eventsub-Message-Timestamp")
	signature := r.Header.Get("Twitch-Eventsub-Message-Signature")

	if !h.verifySignature(messageID, timestamp, body, signature) {
		log.Printf("Subathon webhook: invalid signature, rejecting")
		http.Error(w, "invalid signature", http.StatusForbidden)
		return
	}

	var envelope eventSubEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	switch r.Header.Get("Twitch-Eventsub-Message-Type") {
	case "webhook_callback_verification":
		// One-time handshake when a subscription is first created - Twitch
		// expects the raw challenge string echoed back, no JSON wrapper.
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(envelope.Challenge))
		return

	case "notification":
		h.handleNotification(r, messageID, envelope)
		w.WriteHeader(http.StatusOK)
		return

	case "revocation":
		log.Printf("Subathon webhook: subscription revoked (type=%s) - EventSub manager will recreate it on its next check", envelope.Subscription.Type)
		w.WriteHeader(http.StatusOK)
		return

	default:
		w.WriteHeader(http.StatusOK)
	}
}

func (h *SubathonWebhookHandler) handleNotification(r *http.Request, messageID string, envelope eventSubEnvelope) {
	userID := envelope.Subscription.Condition.BroadcasterUserID
	if userID == "" {
		log.Printf("Subathon webhook: notification missing broadcaster_user_id")
		return
	}

	var eventType, userName, subTier string
	var amount int

	switch envelope.Subscription.Type {
	case "channel.subscribe", "channel.subscription.gift":
		var event struct {
			UserName string `json:"user_name"`
			Tier     string `json:"tier"`
		}
		if err := json.Unmarshal(envelope.Event, &event); err != nil {
			log.Printf("Subathon webhook: failed to parse sub event: %v", err)
			return
		}
		eventType, userName, subTier, amount = "sub", event.UserName, event.Tier, 1

	case "channel.cheer":
		var event struct {
			UserName string `json:"user_name"`
			Bits     int    `json:"bits"`
		}
		if err := json.Unmarshal(envelope.Event, &event); err != nil {
			log.Printf("Subathon webhook: failed to parse cheer event: %v", err)
			return
		}
		eventType, userName, amount = "bits", event.UserName, event.Bits

	default:
		return
	}

	// Stored alongside the failed-event record (if processing fails) so
	// RetryFailedEvents can reconstruct this exact call later.
	retryPayload, _ := json.Marshal(map[string]interface{}{
		"userName": userName,
		"amount":   amount,
		"subTier":  subTier,
	})

	if err := h.subathonService.ProcessEventDeduped(
		r.Context(), messageID, userID, eventType, userName, amount, subTier, retryPayload,
	); err != nil {
		log.Printf("Subathon webhook: %v", err)
	}
}

// verifySignature checks Twitch-Eventsub-Message-Signature so a forged
// request can't inject fake subs/bits into someone's timer. Per Twitch's
// docs: sha256=hex(HMAC-SHA256(secret, message_id + timestamp + body)).
func (h *SubathonWebhookHandler) verifySignature(messageID, timestamp string, body []byte, signatureHeader string) bool {
	if h.webhookSecret == "" || signatureHeader == "" {
		return false
	}

	mac := hmac.New(sha256.New, []byte(h.webhookSecret))
	mac.Write([]byte(messageID))
	mac.Write([]byte(timestamp))
	mac.Write(body)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expected), []byte(signatureHeader))
}
