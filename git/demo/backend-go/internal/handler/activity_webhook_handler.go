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

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/service"
)

// ActivityWebhookHandler receives Twitch EventSub webhook deliveries for
// follow/subscribe/gift-sub/cheer/raid events, feeding the Live-Dashboard's
// activity feed directly - no Redis/bot round-trip, this Go process is the
// EventSub subscriber itself (see ActivityEventSubManager). Public route,
// same HMAC-signature verification as SubathonWebhookHandler (copied, not
// shared, since the two features are otherwise independent).
type ActivityWebhookHandler struct {
	activityService *service.ActivityService
	webhookSecret   string
}

func NewActivityWebhookHandler(activityService *service.ActivityService, webhookSecret string) *ActivityWebhookHandler {
	return &ActivityWebhookHandler{
		activityService: activityService,
		webhookSecret:   webhookSecret,
	}
}

func (h *ActivityWebhookHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/activity/eventsub/callback", h.HandleCallback).Methods(http.MethodPost)
}

type activityEventSubEnvelope struct {
	Challenge    string `json:"challenge"`
	Subscription struct {
		Type      string `json:"type"`
		Condition struct {
			BroadcasterUserID   string `json:"broadcaster_user_id"`
			ToBroadcasterUserID string `json:"to_broadcaster_user_id"`
		} `json:"condition"`
	} `json:"subscription"`
	Event json.RawMessage `json:"event"`
}

func (h *ActivityWebhookHandler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	messageID := r.Header.Get("Twitch-Eventsub-Message-Id")
	timestamp := r.Header.Get("Twitch-Eventsub-Message-Timestamp")
	signature := r.Header.Get("Twitch-Eventsub-Message-Signature")

	if !h.verifySignature(messageID, timestamp, body, signature) {
		log.Printf("Activity webhook: invalid signature, rejecting")
		http.Error(w, "invalid signature", http.StatusForbidden)
		return
	}

	var envelope activityEventSubEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	switch r.Header.Get("Twitch-Eventsub-Message-Type") {
	case "webhook_callback_verification":
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(envelope.Challenge))
		return

	case "notification":
		h.handleNotification(r, messageID, envelope)
		w.WriteHeader(http.StatusOK)
		return

	case "revocation":
		log.Printf("Activity webhook: subscription revoked (type=%s) - EventSub manager will recreate it on its next check", envelope.Subscription.Type)
		w.WriteHeader(http.StatusOK)
		return

	default:
		w.WriteHeader(http.StatusOK)
	}
}

func (h *ActivityWebhookHandler) handleNotification(r *http.Request, messageID string, envelope activityEventSubEnvelope) {
	// channel.raid's broadcaster is the recipient (to_broadcaster_user_id),
	// every other type here uses broadcaster_user_id - see
	// activity_eventsub_manager.go's activitySubscriptionSpecs comment.
	broadcasterID := envelope.Subscription.Condition.BroadcasterUserID
	if broadcasterID == "" {
		broadcasterID = envelope.Subscription.Condition.ToBroadcasterUserID
	}
	if broadcasterID == "" {
		log.Printf("Activity webhook: notification missing broadcaster id")
		return
	}

	var activityType domain.ActivityType
	var username, displayName string
	var viewers, bits *int
	var tier, message *string

	switch envelope.Subscription.Type {
	case "channel.follow":
		var event struct {
			UserName  string `json:"user_name"`
			UserLogin string `json:"user_login"`
		}
		if err := json.Unmarshal(envelope.Event, &event); err != nil {
			log.Printf("Activity webhook: failed to parse follow event: %v", err)
			return
		}
		activityType, username, displayName = domain.ActivityFollow, event.UserLogin, event.UserName

	case "channel.subscribe":
		var event struct {
			UserName  string `json:"user_name"`
			UserLogin string `json:"user_login"`
			Tier      string `json:"tier"`
		}
		if err := json.Unmarshal(envelope.Event, &event); err != nil {
			log.Printf("Activity webhook: failed to parse subscribe event: %v", err)
			return
		}
		activityType, username, displayName, tier = domain.ActivitySubscribe, event.UserLogin, event.UserName, &event.Tier

	case "channel.subscription.gift":
		var event struct {
			UserName  string `json:"user_name"`
			UserLogin string `json:"user_login"`
			Tier      string `json:"tier"`
		}
		if err := json.Unmarshal(envelope.Event, &event); err != nil {
			log.Printf("Activity webhook: failed to parse gift-sub event: %v", err)
			return
		}
		activityType, username, displayName, tier = domain.ActivityGiftSub, event.UserLogin, event.UserName, &event.Tier

	case "channel.cheer":
		var event struct {
			UserName  string `json:"user_name"`
			UserLogin string `json:"user_login"`
			Bits      int    `json:"bits"`
			Message   string `json:"message"`
		}
		if err := json.Unmarshal(envelope.Event, &event); err != nil {
			log.Printf("Activity webhook: failed to parse cheer event: %v", err)
			return
		}
		activityType, username, displayName, bits, message = domain.ActivityCheer, event.UserLogin, event.UserName, &event.Bits, &event.Message

	case "channel.raid":
		var event struct {
			FromBroadcasterUserName  string `json:"from_broadcaster_user_name"`
			FromBroadcasterUserLogin string `json:"from_broadcaster_user_login"`
			Viewers                  int    `json:"viewers"`
		}
		if err := json.Unmarshal(envelope.Event, &event); err != nil {
			log.Printf("Activity webhook: failed to parse raid event: %v", err)
			return
		}
		activityType, username, displayName, viewers = domain.ActivityRaid, event.FromBroadcasterUserLogin, event.FromBroadcasterUserName, &event.Viewers

	default:
		return
	}

	if _, err := h.activityService.CreateActivity(
		r.Context(), broadcasterID, activityType, username, displayName, viewers, bits, tier, message, &messageID,
	); err != nil {
		log.Printf("Activity webhook: %v", err)
	}
}

// verifySignature - identical scheme to SubathonWebhookHandler's, per
// Twitch's docs: sha256=hex(HMAC-SHA256(secret, message_id+timestamp+body)).
func (h *ActivityWebhookHandler) verifySignature(messageID, timestamp string, body []byte, signatureHeader string) bool {
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
