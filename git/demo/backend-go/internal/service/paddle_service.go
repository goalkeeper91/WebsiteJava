package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
	"demo/backend-go/pkg/config"
)

// paddleSignatureTolerance guards against replayed webhook deliveries by
// rejecting a Paddle-Signature whose ts is too far from now. Paddle's docs
// recommend a short tolerance window here; 5 minutes is a conservative,
// commonly-used default (matches other providers' webhook tolerances) -
// worth re-checking against Paddle's current documentation before relying
// on it in production, since the exact recommended value wasn't pinned down
// with full certainty during research.
const paddleSignatureTolerance = 5 * time.Minute

type PaddleService struct {
	subscriptionRepo repository.UserSubscriptionRepository
	eventRepo        repository.PaddleEventRepository
	cfg              config.PaddleConfig
}

func NewPaddleService(
	subscriptionRepo repository.UserSubscriptionRepository,
	eventRepo repository.PaddleEventRepository,
	cfg config.PaddleConfig,
) *PaddleService {
	return &PaddleService{
		subscriptionRepo: subscriptionRepo,
		eventRepo:        eventRepo,
		cfg:              cfg,
	}
}

// VerifyWebhookSignature checks the Paddle-Signature header
// ("ts=<unix>;h1=<hex>") against HMAC-SHA256("{ts}:{rawBody}", secret),
// per Paddle's webhook signing scheme.
func (s *PaddleService) VerifyWebhookSignature(rawBody []byte, signatureHeader string) error {
	if s.cfg.WebhookSecret == "" {
		return fmt.Errorf("PADDLE_WEBHOOK_SECRET ist nicht konfiguriert")
	}

	ts, h1, err := parsePaddleSignatureHeader(signatureHeader)
	if err != nil {
		return err
	}

	tsSeconds, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return fmt.Errorf("ungueltiger ts in Paddle-Signature: %w", err)
	}
	age := time.Since(time.Unix(tsSeconds, 0))
	if age > paddleSignatureTolerance || age < -paddleSignatureTolerance {
		return fmt.Errorf("paddle-signature timestamp ausserhalb der Toleranz")
	}

	mac := hmac.New(sha256.New, []byte(s.cfg.WebhookSecret))
	mac.Write([]byte(ts))
	mac.Write([]byte(":"))
	mac.Write(rawBody)
	expected := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(expected), []byte(h1)) {
		return fmt.Errorf("paddle-signature stimmt nicht ueberein")
	}
	return nil
}

func parsePaddleSignatureHeader(header string) (ts, h1 string, err error) {
	for _, part := range strings.Split(header, ";") {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}
		switch kv[0] {
		case "ts":
			ts = kv[1]
		case "h1":
			h1 = kv[1]
		}
	}
	if ts == "" || h1 == "" {
		return "", "", fmt.Errorf("paddle-signature header fehlt ts/h1")
	}
	return ts, h1, nil
}

type paddleWebhookEnvelope struct {
	EventID    string          `json:"event_id"`
	EventType  string          `json:"event_type"`
	OccurredAt time.Time       `json:"occurred_at"`
	Data       json.RawMessage `json:"data"`
}

type paddleCustomData struct {
	TwitchUserID string `json:"twitch_user_id"`
}

type paddlePriceRef struct {
	ID string `json:"id"`
}

type paddleItem struct {
	Price paddlePriceRef `json:"price"`
}

type paddleBillingPeriod struct {
	StartsAt time.Time `json:"starts_at"`
	EndsAt   time.Time `json:"ends_at"`
}

type paddleTransactionData struct {
	CustomerID     string               `json:"customer_id"`
	SubscriptionID string               `json:"subscription_id"`
	CustomData     *paddleCustomData    `json:"custom_data"`
	Items          []paddleItem         `json:"items"`
	BillingPeriod  *paddleBillingPeriod `json:"billing_period"`
}

type paddleSubscriptionData struct {
	ID                   string               `json:"id"`
	CustomerID           string               `json:"customer_id"`
	Status               string               `json:"status"`
	CustomData           *paddleCustomData    `json:"custom_data"`
	CurrentBillingPeriod *paddleBillingPeriod `json:"current_billing_period"`
	Items                []paddleItem         `json:"items"`
}

// HandleWebhookEvent dispatches an already-signature-verified Paddle
// webhook body. Deduplicated via PaddleEventRepository first - Paddle does
// not guarantee exactly-once or in-order delivery.
func (s *PaddleService) HandleWebhookEvent(ctx context.Context, rawBody []byte) error {
	var envelope paddleWebhookEnvelope
	if err := json.Unmarshal(rawBody, &envelope); err != nil {
		return fmt.Errorf("fehler beim Parsen des Paddle-Webhooks: %w", err)
	}

	alreadyProcessed, err := s.eventRepo.HasProcessed(ctx, envelope.EventID)
	if err != nil {
		return err
	}
	if alreadyProcessed {
		log.Printf("Paddle-Webhook %s (%s) bereits verarbeitet, ueberspringe", envelope.EventID, envelope.EventType)
		return nil
	}

	switch envelope.EventType {
	case "transaction.completed":
		err = s.handleTransactionCompleted(ctx, envelope)
	case "subscription.updated", "subscription.canceled":
		err = s.handleSubscriptionUpdated(ctx, envelope)
	case "transaction.payment_failed":
		// Paddle runs its own dunning/retry flow; the eventual status change
		// (past_due, then possibly canceled) arrives via a later
		// subscription.updated event, so no bespoke handling here.
		log.Printf("Paddle: Zahlung fehlgeschlagen (event=%s) - Paddles Dunning uebernimmt Retries", envelope.EventID)
	default:
		log.Printf("Paddle-Webhook: unbehandelter Event-Typ %q, ignoriert", envelope.EventType)
	}
	if err != nil {
		return err
	}

	return s.eventRepo.MarkProcessed(ctx, envelope.EventID, envelope.EventType)
}

func (s *PaddleService) resolveTwitchUserID(ctx context.Context, customData *paddleCustomData, paddleCustomerID string) (string, error) {
	if customData != nil && customData.TwitchUserID != "" {
		return customData.TwitchUserID, nil
	}
	sub, err := s.subscriptionRepo.GetByPaddleCustomerID(ctx, paddleCustomerID)
	if err != nil {
		return "", err
	}
	if sub == nil {
		return "", fmt.Errorf("konnte Paddle-Event keinem Nutzer zuordnen (customer_id=%s)", paddleCustomerID)
	}
	return sub.UserID, nil
}

func (s *PaddleService) handleTransactionCompleted(ctx context.Context, envelope paddleWebhookEnvelope) error {
	var data paddleTransactionData
	if err := json.Unmarshal(envelope.Data, &data); err != nil {
		return fmt.Errorf("fehler beim Parsen der Paddle-Transaction: %w", err)
	}

	twitchUserID, err := s.resolveTwitchUserID(ctx, data.CustomData, data.CustomerID)
	if err != nil {
		return err
	}
	if len(data.Items) == 0 {
		return fmt.Errorf("paddle transaction ohne items (event=%s)", envelope.EventID)
	}

	tierIDStr, billingCycleStr, ok := s.cfg.TierForPriceID(data.Items[0].Price.ID)
	if !ok {
		return fmt.Errorf("unbekannte paddle price_id %q (event=%s)", data.Items[0].Price.ID, envelope.EventID)
	}
	tier := domain.TierID(tierIDStr)
	cycle := domain.BillingCycle(billingCycleStr)
	status := domain.SubscriptionActive

	var expiresAt *time.Time
	if data.BillingPeriod != nil {
		expiresAt = &data.BillingPeriod.EndsAt
	}
	occurredAt := envelope.OccurredAt

	return s.subscriptionRepo.Update(ctx, twitchUserID, domain.UserSubscriptionUpdateInput{
		TierID:               &tier,
		Status:               &status,
		BillingCycle:         &cycle,
		ExpiresAt:            expiresAt,
		PaddleCustomerID:     &data.CustomerID,
		PaddleSubscriptionID: &data.SubscriptionID,
		LastPaddleEventAt:    &occurredAt,
	})
}

func (s *PaddleService) handleSubscriptionUpdated(ctx context.Context, envelope paddleWebhookEnvelope) error {
	var data paddleSubscriptionData
	if err := json.Unmarshal(envelope.Data, &data); err != nil {
		return fmt.Errorf("fehler beim Parsen der Paddle-Subscription: %w", err)
	}

	twitchUserID, err := s.resolveTwitchUserID(ctx, data.CustomData, data.CustomerID)
	if err != nil {
		return err
	}

	existing, err := s.subscriptionRepo.GetByUserID(ctx, twitchUserID)
	if err != nil {
		return err
	}
	if existing != nil && existing.LastPaddleEventAt != nil && envelope.OccurredAt.Before(*existing.LastPaddleEventAt) {
		// Out-of-order delivery (Paddle does not guarantee event ordering) -
		// this event is older than the state we already applied, discard it.
		log.Printf("Paddle-Webhook %s ist aelter als der zuletzt verarbeitete Stand fuer %s, ignoriert", envelope.EventID, twitchUserID)
		return nil
	}

	occurredAt := envelope.OccurredAt
	update := domain.UserSubscriptionUpdateInput{
		PaddleCustomerID:     &data.CustomerID,
		PaddleSubscriptionID: &data.ID,
		LastPaddleEventAt:    &occurredAt,
	}

	if data.CurrentBillingPeriod != nil {
		update.ExpiresAt = &data.CurrentBillingPeriod.EndsAt
	}

	if status, ok := mapPaddleStatus(data.Status); ok {
		update.Status = &status
		if status == domain.SubscriptionCanceled {
			update.CanceledAt = &occurredAt
		}
	} else {
		log.Printf("Paddle: unbekannter subscription-status %q fuer %s, Status bleibt unveraendert", data.Status, twitchUserID)
	}

	if len(data.Items) > 0 {
		if tierIDStr, billingCycleStr, ok := s.cfg.TierForPriceID(data.Items[0].Price.ID); ok {
			tier := domain.TierID(tierIDStr)
			cycle := domain.BillingCycle(billingCycleStr)
			update.TierID = &tier
			update.BillingCycle = &cycle
		}
	}

	return s.subscriptionRepo.Update(ctx, twitchUserID, update)
}

// mapPaddleStatus maps Paddle's subscription status onto the local domain
// enum. "paused" has no dedicated local status (pausing isn't a supported
// customer action surfaced by this app in this phase) and is treated like
// "canceled" - access continues until ExpiresAt via IsActive()'s canceled
// case, which is the closest existing semantic.
func mapPaddleStatus(paddleStatus string) (domain.SubscriptionStatus, bool) {
	switch paddleStatus {
	case "active":
		return domain.SubscriptionActive, true
	case "trialing":
		return domain.SubscriptionTrialing, true
	case "past_due":
		return domain.SubscriptionPastDue, true
	case "paused", "canceled":
		return domain.SubscriptionCanceled, true
	default:
		return "", false
	}
}
