package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"demo/backend-go/internal/repository"
)

const (
	subathonEventSubURL   = "wss://eventsub.wss.twitch.tv/ws"
	subathonPollInterval  = 60 * time.Second
	subathonReconnectWait = 15 * time.Second
)

// SubathonEventSubClient keeps a live Twitch EventSub WebSocket connection
// per user who has ever used the Subathon Timer, so subs/cheers are counted
// even if nobody has the dashboard's Subathon tab open. This mirrors
// clip-detector's poll-then-manage-per-user-workers pattern (cmd/clip-detector).
type SubathonEventSubClient struct {
	subathonService *SubathonService
	tokenRepo       repository.AuthTokenRepository
	clientID        string
	httpClient      *http.Client

	mu     sync.Mutex
	active map[string]context.CancelFunc
}

func NewSubathonEventSubClient(subathonService *SubathonService, tokenRepo repository.AuthTokenRepository, clientID string) *SubathonEventSubClient {
	return &SubathonEventSubClient{
		subathonService: subathonService,
		tokenRepo:       tokenRepo,
		clientID:        clientID,
		httpClient:      &http.Client{Timeout: 10 * time.Second},
		active:          make(map[string]context.CancelFunc),
	}
}

// Start polls for users with Subathon state and ensures each has a running
// connection worker, until ctx is cancelled.
func (c *SubathonEventSubClient) Start(ctx context.Context) {
	log.Println("🚀 Subathon EventSub Client gestartet...")
	c.ensureConnections(ctx)

	ticker := time.NewTicker(subathonPollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.ensureConnections(ctx)
		}
	}
}

func (c *SubathonEventSubClient) ensureConnections(ctx context.Context) {
	userIDs, err := c.subathonService.GetAllUserIDs(ctx)
	if err != nil {
		log.Printf("Subathon: failed to list users: %v", err)
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for _, userID := range userIDs {
		if _, ok := c.active[userID]; ok {
			continue // worker already running for this user
		}

		userCtx, cancel := context.WithCancel(ctx)
		c.active[userID] = cancel
		go c.runUserConnection(userCtx, userID)
	}
}

func (c *SubathonEventSubClient) forgetUser(userID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.active, userID)
}

// runUserConnection keeps reconnecting for a single user until its context
// is cancelled (which only happens if the whole client shuts down - a user
// never gets permanently dropped once discovered, since Subathon has no
// on/off toggle beyond "has ever used it").
func (c *SubathonEventSubClient) runUserConnection(ctx context.Context, userID string) {
	defer c.forgetUser(userID)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		token, err := c.tokenRepo.GetByTwitchUserID(ctx, userID)
		if err != nil || token == nil {
			log.Printf("Subathon: no Twitch token for user %s, retrying later: %v", userID, err)
			select {
			case <-ctx.Done():
				return
			case <-time.After(subathonReconnectWait):
				continue
			}
		}

		if err := c.connectOnce(ctx, userID, token.AccessToken, subathonEventSubURL); err != nil {
			log.Printf("Subathon EventSub connection for user %s ended: %v", userID, err)
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(subathonReconnectWait):
		}
	}
}

type eventSubMessage struct {
	Metadata struct {
		MessageType string `json:"message_type"`
	} `json:"metadata"`
	Payload struct {
		Session struct {
			ID                      string `json:"id"`
			ReconnectURL            string `json:"reconnect_url"`
			KeepaliveTimeoutSeconds int    `json:"keepalive_timeout_seconds"`
		} `json:"session"`
		Subscription struct {
			Type string `json:"type"`
		} `json:"subscription"`
		Event json.RawMessage `json:"event"`
	} `json:"payload"`
}

// connectOnce runs a single WebSocket session to completion (or failure).
// Returning nil never happens in practice - Twitch always eventually closes
// or asks for a reconnect; the caller redials after a short wait either way.
func (c *SubathonEventSubClient) connectOnce(ctx context.Context, userID, accessToken, url string) error {
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, url, nil)
	if err != nil {
		return fmt.Errorf("dial failed: %w", err)
	}
	defer conn.Close()

	keepaliveTimeout := 30 * time.Second

	for {
		_ = conn.SetReadDeadline(time.Now().Add(keepaliveTimeout + 10*time.Second))

		_, raw, err := conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("read failed: %w", err)
		}

		var msg eventSubMessage
		if err := json.Unmarshal(raw, &msg); err != nil {
			log.Printf("Subathon: failed to parse EventSub message: %v", err)
			continue
		}

		switch msg.Metadata.MessageType {
		case "session_welcome":
			if msg.Payload.Session.KeepaliveTimeoutSeconds > 0 {
				keepaliveTimeout = time.Duration(msg.Payload.Session.KeepaliveTimeoutSeconds) * time.Second
			}
			if err := c.registerSubscriptions(ctx, accessToken, userID, msg.Payload.Session.ID); err != nil {
				log.Printf("Subathon: failed to register EventSub subscriptions for %s: %v", userID, err)
			} else {
				log.Printf("✅ Subathon EventSub verbunden für User %s", userID)
			}

		case "session_reconnect":
			_ = conn.Close()
			return c.connectOnce(ctx, userID, accessToken, msg.Payload.Session.ReconnectURL)

		case "notification":
			c.handleNotification(ctx, userID, msg)

		case "session_keepalive":
			// nothing to do, read deadline already reset above

		default:
			// revocation, session_disconnected, etc. - ignore, outer loop reconnects
		}
	}
}

func (c *SubathonEventSubClient) handleNotification(ctx context.Context, userID string, msg eventSubMessage) {
	subType := msg.Payload.Subscription.Type

	switch subType {
	case "channel.subscribe", "channel.subscription.gift":
		var event struct {
			UserName string `json:"user_name"`
			Tier     string `json:"tier"`
		}
		if err := json.Unmarshal(msg.Payload.Event, &event); err != nil {
			log.Printf("Subathon: failed to parse sub event: %v", err)
			return
		}
		if err := c.subathonService.ProcessEvent(ctx, userID, "sub", event.UserName, 1, event.Tier); err != nil {
			log.Printf("Subathon: failed to process sub event: %v", err)
		}

	case "channel.cheer":
		var event struct {
			UserName string `json:"user_name"`
			Bits     int    `json:"bits"`
		}
		if err := json.Unmarshal(msg.Payload.Event, &event); err != nil {
			log.Printf("Subathon: failed to parse cheer event: %v", err)
			return
		}
		if err := c.subathonService.ProcessEvent(ctx, userID, "bits", event.UserName, event.Bits, ""); err != nil {
			log.Printf("Subathon: failed to process cheer event: %v", err)
		}
	}
}

func (c *SubathonEventSubClient) registerSubscriptions(ctx context.Context, accessToken, userID, sessionID string) error {
	types := []string{"channel.subscribe", "channel.subscription.gift", "channel.cheer"}

	for _, subType := range types {
		body := map[string]interface{}{
			"type":      subType,
			"version":   "1",
			"condition": map[string]string{"broadcaster_user_id": userID},
			"transport": map[string]string{"method": "websocket", "session_id": sessionID},
		}
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal subscription request: %w", err)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost,
			"https://api.twitch.tv/helix/eventsub/subscriptions", strings.NewReader(string(data)))
		if err != nil {
			return fmt.Errorf("failed to build request: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+accessToken)
		req.Header.Set("Client-Id", c.clientID)
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("failed to send subscription request for %s: %w", subType, err)
		}
		resp.Body.Close()

		if resp.StatusCode >= 300 {
			return fmt.Errorf("subscription request for %s returned %d", subType, resp.StatusCode)
		}
	}

	return nil
}
