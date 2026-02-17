package domain

import "time"

type SubscriptionStatus string
type BillingCycle string

const (
	SubscriptionActive   SubscriptionStatus = "active"
	SubscriptionExpired  SubscriptionStatus = "expired"
	SubscriptionCanceled SubscriptionStatus = "canceled"
	SubscriptionTrialing SubscriptionStatus = "trialing"

	BillingMonthly BillingCycle = "monthly"
	BillingYearly  BillingCycle = "yearly"
)

type UserSubscription struct {
	ID                   int64              `json:"id" db:"id"`
	UserID               string             `json:"userId" db:"user_id"`
	TierID               TierID             `json:"tierId" db:"tier_id"`
	Status               SubscriptionStatus `json:"status" db:"status"`
	BillingCycle         *BillingCycle      `json:"billingCycle" db:"billing_cycle"`
	StartedAt            time.Time          `json:"startedAt" db:"started_at"`
	ExpiresAt            *time.Time         `json:"expiresAt" db:"expires_at"`
	CanceledAt           *time.Time         `json:"canceledAt" db:"canceled_at"`
	StripeCustomerID     *string            `json:"-" db:"stripe_customer_id"`
	StripeSubscriptionID *string            `json:"-" db:"stripe_subscription_id"`
	CreatedAt            time.Time          `json:"createdAt" db:"created_at"`
	UpdatedAt            time.Time          `json:"updatedAt" db:"updated_at"`

	// Joined from subscription_tiers
	Tier *SubscriptionTier `json:"tier,omitempty" db:"-"`
}

type UserSubscriptionCreateInput struct {
	UserID               string
	TierID               TierID
	BillingCycle         *BillingCycle
	ExpiresAt            *time.Time
	StripeCustomerID     *string
	StripeSubscriptionID *string
}

type UserSubscriptionUpdateInput struct {
	TierID               *TierID
	Status               *SubscriptionStatus
	BillingCycle         *BillingCycle
	ExpiresAt            *time.Time
	CanceledAt           *time.Time
	StripeCustomerID     *string
	StripeSubscriptionID *string
}

func NewUserSubscription(userID string, tierID TierID) *UserSubscription {
	now := time.Now()
	return &UserSubscription{
		UserID:    userID,
		TierID:    tierID,
		Status:    SubscriptionActive,
		StartedAt: now,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (s *UserSubscription) IsActive() bool {
	if s.Status != SubscriptionActive && s.Status != SubscriptionTrialing {
		return false
	}
	if s.ExpiresAt != nil && time.Now().After(*s.ExpiresAt) {
		return false
	}
	return true
}

func (s *UserSubscription) IsFree() bool {
	return s.TierID == TierFree
}

func (s *UserSubscription) IsPro() bool {
	return s.TierID == TierPro
}

func (s *UserSubscription) IsPremium() bool {
	return s.TierID == TierPremium
}

func (s *UserSubscription) CanUseN8N() bool {
	return s.IsActive() && (s.IsPro() || s.IsPremium())
}

func (s *UserSubscription) Cancel() {
	now := time.Now()
	s.Status = SubscriptionCanceled
	s.CanceledAt = &now
	s.UpdatedAt = now
}

func (s *UserSubscription) Upgrade(tierID TierID) {
	s.TierID = tierID
	s.Status = SubscriptionActive
	s.UpdatedAt = time.Now()
}