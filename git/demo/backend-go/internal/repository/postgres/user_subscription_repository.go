package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
)

type UserSubscriptionRepository struct {
	db *sql.DB
}

func NewUserSubscriptionRepository(db *sql.DB) repository.UserSubscriptionRepository {
	return &UserSubscriptionRepository{db: db}
}

func (r *UserSubscriptionRepository) Create(ctx context.Context, input domain.UserSubscriptionCreateInput) (*domain.UserSubscription, error) {
	query := `
		INSERT INTO user_subscriptions
		(user_id, tier_id, status, billing_cycle, started_at, expires_at,
		 stripe_customer_id, stripe_subscription_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), $5, $6, $7, NOW(), NOW())
		RETURNING id, user_id, tier_id, status, billing_cycle,
		          started_at, expires_at, canceled_at,
		          stripe_customer_id, stripe_subscription_id,
		          created_at, updated_at
	`

	var sub domain.UserSubscription
	err := r.db.QueryRowContext(
		ctx, query,
		input.UserID,
		input.TierID,
		domain.SubscriptionActive,
		input.BillingCycle,
		input.ExpiresAt,
		input.StripeCustomerID,
		input.StripeSubscriptionID,
	).Scan(
		&sub.ID,
		&sub.UserID,
		&sub.TierID,
		&sub.Status,
		&sub.BillingCycle,
		&sub.StartedAt,
		&sub.ExpiresAt,
		&sub.CanceledAt,
		&sub.StripeCustomerID,
		&sub.StripeSubscriptionID,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Erstellen der Subscription: %w", err)
	}

	return &sub, nil
}

func (r *UserSubscriptionRepository) GetByUserID(ctx context.Context, userID string) (*domain.UserSubscription, error) {
	query := `
		SELECT id, user_id, tier_id, status, billing_cycle,
		       started_at, expires_at, canceled_at,
		       stripe_customer_id, stripe_subscription_id,
		       created_at, updated_at
		FROM user_subscriptions
		WHERE user_id = $1
	`

	var sub domain.UserSubscription
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&sub.ID,
		&sub.UserID,
		&sub.TierID,
		&sub.Status,
		&sub.BillingCycle,
		&sub.StartedAt,
		&sub.ExpiresAt,
		&sub.CanceledAt,
		&sub.StripeCustomerID,
		&sub.StripeSubscriptionID,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Subscription: %w", err)
	}

	return &sub, nil
}

func (r *UserSubscriptionRepository) GetByUserIDWithTier(ctx context.Context, userID string) (*domain.UserSubscription, error) {
	query := `
		SELECT s.id, s.user_id, s.tier_id, s.status, s.billing_cycle,
		       s.started_at, s.expires_at, s.canceled_at,
		       s.stripe_customer_id, s.stripe_subscription_id,
		       s.created_at, s.updated_at,
		       t.id, t.name, t.price_monthly, t.price_yearly,
		       t.max_commands, t.max_workflows, t.max_votes_per_month,
		       t.features, t.is_active, t.created_at
		FROM user_subscriptions s
		INNER JOIN subscription_tiers t ON s.tier_id = t.id
		WHERE s.user_id = $1
	`

	var sub domain.UserSubscription
	var tier domain.SubscriptionTier
	var featuresJSON []byte

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&sub.ID,
		&sub.UserID,
		&sub.TierID,
		&sub.Status,
		&sub.BillingCycle,
		&sub.StartedAt,
		&sub.ExpiresAt,
		&sub.CanceledAt,
		&sub.StripeCustomerID,
		&sub.StripeSubscriptionID,
		&sub.CreatedAt,
		&sub.UpdatedAt,
		&tier.ID,
		&tier.Name,
		&tier.PriceMonthly,
		&tier.PriceYearly,
		&tier.MaxCommands,
		&tier.MaxWorkflows,
		&tier.MaxVotesPerMonth,
		&featuresJSON,
		&tier.IsActive,
		&tier.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Subscription: %w", err)
	}

	if err := json.Unmarshal(featuresJSON, &tier.Features); err != nil {
		return nil, fmt.Errorf("fehler beim Parsen der Features: %w", err)
	}

	sub.Tier = &tier
	return &sub, nil
}

func (r *UserSubscriptionRepository) Update(ctx context.Context, userID string, input domain.UserSubscriptionUpdateInput) error {
	query := `
		UPDATE user_subscriptions
		SET tier_id = COALESCE($2, tier_id),
		    status = COALESCE($3, status),
		    billing_cycle = COALESCE($4, billing_cycle),
		    expires_at = COALESCE($5, expires_at),
		    canceled_at = COALESCE($6, canceled_at),
		    stripe_customer_id = COALESCE($7, stripe_customer_id),
		    stripe_subscription_id = COALESCE($8, stripe_subscription_id),
		    updated_at = NOW()
		WHERE user_id = $1
	`

	_, err := r.db.ExecContext(
		ctx, query,
		userID,
		input.TierID,
		input.Status,
		input.BillingCycle,
		input.ExpiresAt,
		input.CanceledAt,
		input.StripeCustomerID,
		input.StripeSubscriptionID,
	)
	if err != nil {
		return fmt.Errorf("fehler beim Update der Subscription: %w", err)
	}

	return nil
}

func (r *UserSubscriptionRepository) Cancel(ctx context.Context, userID string) error {
	query := `
		UPDATE user_subscriptions
		SET status = $2, canceled_at = NOW(), updated_at = NOW()
		WHERE user_id = $1
	`

	_, err := r.db.ExecContext(ctx, query, userID, domain.SubscriptionCanceled)
	if err != nil {
		return fmt.Errorf("fehler beim Kündigen der Subscription: %w", err)
	}

	return nil
}

func (r *UserSubscriptionRepository) Exists(ctx context.Context, userID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM user_subscriptions WHERE user_id = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("fehler beim Prüfen der Subscription: %w", err)
	}

	return exists, nil
}