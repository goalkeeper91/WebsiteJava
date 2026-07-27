package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"demo/backend-go/internal/repository"
)

const (
	activitySubscriptionCheckInterval = 5 * time.Minute
)

// activitySubscriptionSpec describes one EventSub subscription type this
// manager keeps active per broadcaster. Unlike Subathon's manager (which
// hardcodes version "1" and a single condition shape, since all 3 of its
// event types share that shape), these 5 types genuinely differ: raid's
// condition key is broadcaster-as-recipient, not broadcaster-as-subject,
// and follow needs the newer v2 shape.
type activitySubscriptionSpec struct {
	Type      string
	Version   string
	Condition func(userID string) map[string]string
}

var activitySubscriptionSpecs = []activitySubscriptionSpec{
	{
		Type:    "channel.follow",
		Version: "2",
		Condition: func(userID string) map[string]string {
			return map[string]string{"broadcaster_user_id": userID, "moderator_user_id": userID}
		},
	},
	{
		Type:    "channel.subscribe",
		Version: "1",
		Condition: func(userID string) map[string]string {
			return map[string]string{"broadcaster_user_id": userID}
		},
	},
	{
		Type:    "channel.subscription.gift",
		Version: "1",
		Condition: func(userID string) map[string]string {
			return map[string]string{"broadcaster_user_id": userID}
		},
	},
	{
		Type:    "channel.cheer",
		Version: "1",
		Condition: func(userID string) map[string]string {
			return map[string]string{"broadcaster_user_id": userID}
		},
	},
	{
		// Condition key is to_broadcaster_user_id (the raid's recipient),
		// not broadcaster_user_id like every other type here - this is the
		// one Twitch raids into the streamer, which is what the activity
		// feed cares about.
		Type:    "channel.raid",
		Version: "1",
		Condition: func(userID string) map[string]string {
			return map[string]string{"to_broadcaster_user_id": userID}
		},
	},
}

// ActivityEventSubManager keeps follow/sub/gift-sub/cheer/raid webhook
// subscriptions active for every broadcaster with a stored token, feeding
// the Live-Dashboard's activity feed directly (no Redis/bot round-trip -
// the webhook notification lands on this Go process and calls
// ActivityService.CreateActivity directly, see activity_webhook_handler.go).
// Built parallel to SubathonEventSubManager rather than extending it, since
// that one is coupled to Subathon's timer-crediting semantics.
type ActivityEventSubManager struct {
	userRepo      repository.UserRepository
	tokenRepo     repository.AuthTokenRepository
	clientID      string
	webhookSecret string
	callbackURL   string
	httpClient    *http.Client
}

func NewActivityEventSubManager(
	userRepo repository.UserRepository,
	tokenRepo repository.AuthTokenRepository,
	clientID, webhookSecret, callbackURL string,
) *ActivityEventSubManager {
	return &ActivityEventSubManager{
		userRepo:      userRepo,
		tokenRepo:     tokenRepo,
		clientID:      clientID,
		webhookSecret: webhookSecret,
		callbackURL:   callbackURL,
		httpClient:    &http.Client{Timeout: 10 * time.Second},
	}
}

func (m *ActivityEventSubManager) Start(ctx context.Context) {
	log.Println("🚀 Activity EventSub Manager (Webhook) gestartet...")

	m.ensureAllSubscriptions(ctx)

	ticker := time.NewTicker(activitySubscriptionCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.ensureAllSubscriptions(ctx)
		}
	}
}

func (m *ActivityEventSubManager) ensureAllSubscriptions(ctx context.Context) {
	users, err := m.userRepo.GetAllNonBotUsers(ctx)
	if err != nil {
		log.Printf("Activity EventSub: failed to list users: %v", err)
		return
	}

	for _, user := range users {
		if err := m.ensureUserSubscriptions(ctx, user.TwitchID); err != nil {
			log.Printf("Activity EventSub: failed to ensure subscriptions for user %s: %v", user.TwitchID, err)
		}
	}
}

type activityEventSubSubscription struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	Type      string `json:"type"`
	Transport struct {
		Method   string `json:"method"`
		Callback string `json:"callback"`
	} `json:"transport"`
}

func (m *ActivityEventSubManager) ensureUserSubscriptions(ctx context.Context, userID string) error {
	token, err := m.tokenRepo.GetByTwitchUserID(ctx, userID)
	if err != nil || token == nil {
		return fmt.Errorf("no Twitch token available: %w", err)
	}

	existing, err := m.listSubscriptions(ctx, token.AccessToken, userID)
	if err != nil {
		return fmt.Errorf("failed to list existing subscriptions: %w", err)
	}

	haveEnabled := make(map[string]bool)
	for _, sub := range existing {
		if sub.Transport.Method == "webhook" && sub.Transport.Callback == m.callbackURL && sub.Status == "enabled" {
			haveEnabled[sub.Type] = true
		}
	}

	for _, spec := range activitySubscriptionSpecs {
		if haveEnabled[spec.Type] {
			continue
		}
		if err := m.createSubscription(ctx, token.AccessToken, userID, spec); err != nil {
			log.Printf("Activity EventSub: failed to create %s subscription for user %s: %v", spec.Type, userID, err)
		} else {
			log.Printf("✅ Activity EventSub: registered %s webhook subscription for user %s", spec.Type, userID)
		}
	}

	return nil
}

func (m *ActivityEventSubManager) listSubscriptions(ctx context.Context, accessToken, userID string) ([]activityEventSubSubscription, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		"https://api.twitch.tv/helix/eventsub/subscriptions?user_id="+url.QueryEscape(userID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Client-Id", m.clientID)

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	var result struct {
		Data []activityEventSubSubscription `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return result.Data, nil
}

func (m *ActivityEventSubManager) createSubscription(ctx context.Context, accessToken, userID string, spec activitySubscriptionSpec) error {
	body := map[string]interface{}{
		"type":      spec.Type,
		"version":   spec.Version,
		"condition": spec.Condition(userID),
		"transport": map[string]string{
			"method":   "webhook",
			"callback": m.callbackURL,
			"secret":   m.webhookSecret,
		},
	}
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.twitch.tv/helix/eventsub/subscriptions", strings.NewReader(string(data)))
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Client-Id", m.clientID)
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status %d", resp.StatusCode)
	}
	return nil
}
