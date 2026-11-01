package service

import (
	"context"
	"log"
	"time"
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
	client          *eventSubWebhookClient
}

// NewSubathonEventSubManager registers the webhook subscriptions with Twitch's
// APP access token (appToken) - webhook transport does not accept the
// broadcaster's user token, see eventSubWebhookClient.
func NewSubathonEventSubManager(
	subathonService *SubathonService,
	appToken AppTokenProvider,
	clientID, webhookSecret, callbackURL string,
) *SubathonEventSubManager {
	return &SubathonEventSubManager{
		subathonService: subathonService,
		client:          newEventSubWebhookClient(appToken, clientID, webhookSecret, callbackURL),
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

// All 3 Subathon event types share the same version and condition shape.
var subathonSubscriptionSpecs = []eventSubSpec{
	{Type: "channel.subscribe", Version: "1", Condition: broadcasterCondition},
	{Type: "channel.subscription.gift", Version: "1", Condition: broadcasterCondition},
	{Type: "channel.cheer", Version: "1", Condition: broadcasterCondition},
}

func (m *SubathonEventSubManager) ensureUserSubscriptions(ctx context.Context, userID string) error {
	return m.client.ensureUser(ctx, userID, subathonSubscriptionSpecs, "Subathon")
}
