package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const twitchEventSubURL = "https://api.twitch.tv/helix/eventsub/subscriptions"

// AppTokenProvider hands out a valid Twitch APP access token (client
// credentials). *twitch.TwitchAppTokenClient implements it and caches/renews
// the token itself.
type AppTokenProvider interface {
	AppToken(ctx context.Context) (string, error)
}

// eventSubSpec describes one EventSub subscription type to keep active per
// broadcaster. Condition differs by type (raid is keyed on the recipient,
// follow v2 needs a moderator), so each spec builds its own.
type eventSubSpec struct {
	Type      string
	Version   string
	Condition func(userID string) map[string]string
}

func broadcasterCondition(userID string) map[string]string {
	return map[string]string{"broadcaster_user_id": userID}
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

// eventSubWebhookClient is the shared "make sure these webhook subscriptions
// exist" logic behind the Subathon and Activity managers. Both used to carry
// their own copy of it, which is how the same bug got into both: they
// authenticated with the broadcaster's stored USER token, but Twitch only
// accepts an APP access token for webhook-transport subscriptions (a user
// token is for WebSocket transport) - creation failed with 400, and the
// stored user tokens also expire (401) since only bot accounts get refreshed.
type eventSubWebhookClient struct {
	httpClient  *http.Client
	apiURL      string
	clientID    string
	appToken    AppTokenProvider
	secret      string
	callbackURL string
}

func newEventSubWebhookClient(appToken AppTokenProvider, clientID, secret, callbackURL string) *eventSubWebhookClient {
	return &eventSubWebhookClient{
		httpClient:  &http.Client{Timeout: 10 * time.Second},
		apiURL:      twitchEventSubURL,
		clientID:    clientID,
		appToken:    appToken,
		secret:      secret,
		callbackURL: callbackURL,
	}
}

// Subscription states in which Twitch will never deliver again on its own
// and in which re-creating is blocked by a 409 "already exists" until the
// dead one is deleted. (webhook_callback_verification_pending is NOT in this
// list - it is still in flight and must be left alone.)
var deadSubscriptionStatuses = map[string]bool{
	"webhook_callback_verification_failed": true,
	"notification_failures_exceeded":       true,
}

// ensureUser makes sure userID has an enabled webhook subscription for every
// spec, deleting dead ones first. logPrefix keeps each manager's existing log
// wording ("Subathon" / "Activity EventSub"). Returns an error only if the
// subscriptions couldn't even be listed; a failure to create one type is
// logged and the rest still get attempted.
func (c *eventSubWebhookClient) ensureUser(ctx context.Context, userID string, specs []eventSubSpec, logPrefix string) error {
	token, err := c.appToken.AppToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get Twitch app access token: %w", err)
	}

	existing, err := c.list(ctx, token, userID)
	if err != nil {
		return fmt.Errorf("failed to list existing subscriptions: %w", err)
	}

	present := make(map[string]bool)
	for _, sub := range existing {
		if sub.Transport.Method != "webhook" || sub.Transport.Callback != c.callbackURL {
			continue
		}
		if deadSubscriptionStatuses[sub.Status] {
			if err := c.delete(ctx, token, sub.ID); err != nil {
				log.Printf("%s: failed to delete dead %s subscription (%s) for user %s: %v", logPrefix, sub.Type, sub.Status, userID, err)
			} else {
				log.Printf("%s: deleted dead %s subscription (%s) for user %s - re-registering", logPrefix, sub.Type, sub.Status, userID)
			}
			continue
		}
		if sub.Status == "enabled" || sub.Status == "webhook_callback_verification_pending" {
			present[sub.Type] = true
		}
	}

	for _, spec := range specs {
		if present[spec.Type] {
			continue
		}
		if err := c.create(ctx, token, userID, spec); err != nil {
			log.Printf("%s: failed to create %s subscription for user %s: %v", logPrefix, spec.Type, userID, err)
		} else {
			log.Printf("✅ %s: registered %s webhook subscription for user %s", logPrefix, spec.Type, userID)
		}
	}
	return nil
}

func (c *eventSubWebhookClient) newRequest(ctx context.Context, method, rawURL, token string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, rawURL, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Client-Id", c.clientID)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

func (c *eventSubWebhookClient) list(ctx context.Context, token, userID string) ([]eventSubSubscription, error) {
	req, err := c.newRequest(ctx, http.MethodGet, c.apiURL+"?user_id="+url.QueryEscape(userID), token, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, twitchAPIError(resp)
	}

	var result struct {
		Data []eventSubSubscription `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return result.Data, nil
}

func (c *eventSubWebhookClient) create(ctx context.Context, token, userID string, spec eventSubSpec) error {
	data, err := json.Marshal(map[string]interface{}{
		"type":      spec.Type,
		"version":   spec.Version,
		"condition": spec.Condition(userID),
		"transport": map[string]string{
			"method":   "webhook",
			"callback": c.callbackURL,
			"secret":   c.secret,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, c.apiURL, token, strings.NewReader(string(data)))
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return twitchAPIError(resp)
	}
	return nil
}

func (c *eventSubWebhookClient) delete(ctx context.Context, token, subscriptionID string) error {
	req, err := c.newRequest(ctx, http.MethodDelete, c.apiURL+"?id="+url.QueryEscape(subscriptionID), token, nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return twitchAPIError(resp)
	}
	return nil
}

// twitchAPIError turns a non-2xx Twitch response into an error that includes
// Twitch's own explanation. A bare "unexpected status 400" gave no way to tell
// a wrong token type from a missing scope or a bad secret.
func twitchAPIError(resp *http.Response) error {
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
	if msg := strings.TrimSpace(string(body)); msg != "" {
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, msg)
	}
	return fmt.Errorf("unexpected status %d", resp.StatusCode)
}
