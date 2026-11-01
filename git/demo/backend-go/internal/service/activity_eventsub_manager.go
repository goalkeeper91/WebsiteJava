package service

import (
	"context"
	"log"
	"time"

	"demo/backend-go/internal/repository"
)

const (
	activitySubscriptionCheckInterval = 5 * time.Minute
)

// Unlike Subathon's manager (all 3 of its event types share version "1" and
// one condition shape), these 5 types genuinely differ: raid's condition key
// is broadcaster-as-recipient, not broadcaster-as-subject, and follow needs
// the newer v2 shape.
var activitySubscriptionSpecs = []eventSubSpec{
	{
		Type:    "channel.follow",
		Version: "2",
		Condition: func(userID string) map[string]string {
			return map[string]string{"broadcaster_user_id": userID, "moderator_user_id": userID}
		},
	},
	{Type: "channel.subscribe", Version: "1", Condition: broadcasterCondition},
	{Type: "channel.subscription.gift", Version: "1", Condition: broadcasterCondition},
	{Type: "channel.cheer", Version: "1", Condition: broadcasterCondition},
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
	userRepo repository.UserRepository
	client   *eventSubWebhookClient
}

// NewActivityEventSubManager registers the webhook subscriptions with
// Twitch's APP access token (appToken) - webhook transport does not accept
// the broadcaster's user token, see eventSubWebhookClient.
func NewActivityEventSubManager(
	userRepo repository.UserRepository,
	appToken AppTokenProvider,
	clientID, webhookSecret, callbackURL string,
) *ActivityEventSubManager {
	return &ActivityEventSubManager{
		userRepo: userRepo,
		client:   newEventSubWebhookClient(appToken, clientID, webhookSecret, callbackURL),
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

func (m *ActivityEventSubManager) ensureUserSubscriptions(ctx context.Context, userID string) error {
	return m.client.ensureUser(ctx, userID, activitySubscriptionSpecs, "Activity EventSub")
}
