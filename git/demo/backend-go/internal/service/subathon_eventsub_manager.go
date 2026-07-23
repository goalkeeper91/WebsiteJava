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
	subathonSubscriptionCheckInterval = 5 * time.Minute
	subathonFailedEventRetryInterval  = 5 * time.Minute
)

// SubathonEventSubManager replaces an earlier WebSocket-based client.
// WebSocket EventSub subscriptions die the instant the connection drops
// (network blip, or the backend container restarting during a deploy -
// which happened more than once today), and Twitch does not replay
// anything missed during that gap - exactly the highest-risk window during
// a hype train. Webhook transport instead has Twitch call a stable HTTPS
// endpoint directly; the subscription lives on Twitch's side independent
// of our process, so a restart or brief outage doesn't require rebuilding a
// session, and Twitch's own retry logic covers momentary delivery failures.
type SubathonEventSubManager struct {
	subathonService *SubathonService
	tokenRepo       repository.AuthTokenRepository
	clientID        string
	webhookSecret   string
	callbackURL     string
	httpClient      *http.Client
}

func NewSubathonEventSubManager(
	subathonService *SubathonService,
	tokenRepo repository.AuthTokenRepository,
	clientID, webhookSecret, callbackURL string,
) *SubathonEventSubManager {
	return &SubathonEventSubManager{
		subathonService: subathonService,
		tokenRepo:       tokenRepo,
		clientID:        clientID,
		webhookSecret:   webhookSecret,
		callbackURL:     callbackURL,
		httpClient:      &http.Client{Timeout: 10 * time.Second},
	}
}

// Start periodically ensures every Subathon user has the 3 needed webhook
// subscriptions active, and retries any previously-failed events.
func (m *SubathonEventSubManager) Start(ctx context.Context) {
	log.Println("🚀 Subathon EventSub Manager (Webhook) gestartet...")

	m.ensureAllSubscriptions(ctx)
	m.subathonService.RetryFailedEvents(ctx)

	subTicker := time.NewTicker(subathonSubscriptionCheckInterval)
	defer subTicker.Stop()
	retryTicker := time.NewTicker(subathonFailedEventRetryInterval)
	defer retryTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-subTicker.C:
			m.ensureAllSubscriptions(ctx)
		case <-retryTicker.C:
			m.subathonService.RetryFailedEvents(ctx)
		}
	}
}

func (m *SubathonEventSubManager) ensureAllSubscriptions(ctx context.Context) {
	userIDs, err := m.subathonService.GetAllUserIDs(ctx)
	if err != nil {
		log.Printf("Subathon: failed to list users: %v", err)
		return
	}

	for _, userID := range userIDs {
		if err := m.ensureUserSubscriptions(ctx, userID); err != nil {
			log.Printf("Subathon: failed to ensure subscriptions for user %s: %v", userID, err)
		}
	}
}

type eventSubSubscription struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	Type      string `json:"type"`
	Transport struct {
		Method   string `json:"method"`
		Callback string `json:"callback"`
	} `json:"transport"`
}

type getSubscriptionsResponse struct {
	Data []eventSubSubscription `json:"data"`
}

var subathonSubscriptionTypes = []string{"channel.subscribe", "channel.subscription.gift", "channel.cheer"}

func (m *SubathonEventSubManager) ensureUserSubscriptions(ctx context.Context, userID string) error {
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

	for _, subType := range subathonSubscriptionTypes {
		if haveEnabled[subType] {
			continue
		}
		if err := m.createSubscription(ctx, token.AccessToken, userID, subType); err != nil {
			log.Printf("Subathon: failed to create %s subscription for user %s: %v", subType, userID, err)
		} else {
			log.Printf("✅ Subathon: registered %s webhook subscription for user %s", subType, userID)
		}
	}

	return nil
}

func (m *SubathonEventSubManager) listSubscriptions(ctx context.Context, accessToken, userID string) ([]eventSubSubscription, error) {
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

	var result getSubscriptionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return result.Data, nil
}

func (m *SubathonEventSubManager) createSubscription(ctx context.Context, accessToken, userID, subType string) error {
	body := map[string]interface{}{
		"type":      subType,
		"version":   "1",
		"condition": map[string]string{"broadcaster_user_id": userID},
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
