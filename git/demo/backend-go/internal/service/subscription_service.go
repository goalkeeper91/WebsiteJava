package service

import (
	"context"
	"fmt"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
)

type SubscriptionService struct {
	subscriptionRepo repository.UserSubscriptionRepository
	tierRepo         repository.SubscriptionTierRepository
	analyticsRepo    repository.UsageAnalyticsRepository
	userRepo         repository.UserRepository // NEU: für Admin Check
}

func NewSubscriptionService(
	subscriptionRepo repository.UserSubscriptionRepository,
	tierRepo repository.SubscriptionTierRepository,
	analyticsRepo repository.UsageAnalyticsRepository,
	userRepo repository.UserRepository, // NEU
) *SubscriptionService {
	return &SubscriptionService{
		subscriptionRepo: subscriptionRepo,
		tierRepo:         tierRepo,
		analyticsRepo:    analyticsRepo,
		userRepo:         userRepo,
	}
}

// isAdmin checks if a user is an admin (bypasses all subscription checks)
func (s *SubscriptionService) isAdmin(ctx context.Context, userID string) bool {
	user, err := s.userRepo.GetByTwitchID(ctx, userID) // ← Fixed: GetByTwitchID statt GetByID
	if err != nil || user == nil {
		return false
	}
	return user.IsAdmin
}

// GetSubscription returns a user's subscription (or creates free tier)
func (s *SubscriptionService) GetSubscription(ctx context.Context, userID string) (*domain.UserSubscription, error) {
	// Check if subscription exists
	sub, err := s.subscriptionRepo.GetByUserIDWithTier(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Subscription: %w", err)
	}

	// If no subscription exists, create free tier
	if sub == nil {
		freeTier, err := s.tierRepo.GetByID(ctx, domain.TierFree)
		if err != nil {
			return nil, fmt.Errorf("fehler beim Laden des Free Tiers: %w", err)
		}

		sub, err = s.createFreeTier(ctx, userID)
		if err != nil {
			return nil, err
		}
		sub.Tier = freeTier
	}

	return sub, nil
}

// GetTiers returns all available subscription tiers
func (s *SubscriptionService) GetTiers(ctx context.Context) ([]*domain.SubscriptionTier, error) {
	tiers, err := s.tierRepo.GetActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Tiers: %w", err)
	}
	return tiers, nil
}

// CanUseFeature checks if a user can use a specific feature
// ADMIN BYPASS: Admins can use all features
func (s *SubscriptionService) CanUseFeature(ctx context.Context, userID string, feature string) (bool, error) {
	// ADMIN CHECK FIRST
	if s.isAdmin(ctx, userID) {
		return true, nil // Admins can use everything
	}

	sub, err := s.GetSubscription(ctx, userID)
	if err != nil {
		return false, err
	}

	if !sub.IsActive() {
		return false, domain.ErrSubscriptionExpired
	}

	return sub.Tier.HasFeature(feature), nil
}

// CanUseN8N checks if user can use n8n integration
// ADMIN BYPASS: Admins can always use n8n
func (s *SubscriptionService) CanUseN8N(ctx context.Context, userID string) (bool, error) {
	// ADMIN CHECK FIRST
	if s.isAdmin(ctx, userID) {
		return true, nil
	}

	sub, err := s.GetSubscription(ctx, userID)
	if err != nil {
		return false, err
	}

	if !sub.IsActive() {
		return false, domain.ErrSubscriptionExpired
	}

	// n8n requires Pro or Premium
	return sub.IsPro() || sub.IsPremium(), nil
}

// CheckCommandLimit checks if user has reached command limit
// ADMIN BYPASS: Admins have unlimited commands
func (s *SubscriptionService) CheckCommandLimit(ctx context.Context, userID string, currentCount int) error {
	// ADMIN CHECK FIRST
	if s.isAdmin(ctx, userID) {
		return nil // No limits for admins
	}

	sub, err := s.GetSubscription(ctx, userID)
	if err != nil {
		return err
	}

	if !sub.IsActive() {
		return domain.ErrSubscriptionExpired
	}

	if sub.Tier.IsCommandLimitReached(currentCount) {
		return domain.ErrCommandLimitReached
	}

	return nil
}

// CheckWorkflowLimit checks if user has reached workflow limit
// ADMIN BYPASS: Admins have unlimited workflows
func (s *SubscriptionService) CheckWorkflowLimit(ctx context.Context, userID string, currentCount int) error {
	// ADMIN CHECK FIRST
	if s.isAdmin(ctx, userID) {
		return nil // No limits for admins
	}

	sub, err := s.GetSubscription(ctx, userID)
	if err != nil {
		return err
	}

	if !sub.IsActive() {
		return domain.ErrSubscriptionExpired
	}

	if sub.Tier.IsWorkflowLimitReached(currentCount) {
		return domain.ErrWorkflowLimitReached
	}

	return nil
}

// CheckVoteLimit checks if user has reached vote limit
// ADMIN BYPASS: Admins have unlimited votes
func (s *SubscriptionService) CheckVoteLimit(ctx context.Context, userID string, currentCount int) error {
	// ADMIN CHECK FIRST
	if s.isAdmin(ctx, userID) {
		return nil // No limits for admins
	}

	sub, err := s.GetSubscription(ctx, userID)
	if err != nil {
		return err
	}

	if !sub.IsActive() {
		return domain.ErrSubscriptionExpired
	}

	if sub.Tier.IsVoteLimitReached(currentCount) {
		return domain.ErrVoteLimitReached
	}

	return nil
}

// Upgrade upgrades a user to a new tier
func (s *SubscriptionService) Upgrade(ctx context.Context, userID string, newTierID domain.TierID) (*domain.UserSubscription, error) {
	// Validate new tier exists
	_, err := s.tierRepo.GetByID(ctx, newTierID)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden des neuen Tiers: %w", err)
	}

	// Update tier
	err = s.subscriptionRepo.Update(ctx, userID, domain.UserSubscriptionUpdateInput{
		TierID: &newTierID,
	})
	if err != nil {
		return nil, fmt.Errorf("fehler beim Upgrade: %w", err)
	}

	// Reload with tier details
	sub, err := s.subscriptionRepo.GetByUserIDWithTier(ctx, userID)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

// Cancel cancels a user's subscription
func (s *SubscriptionService) Cancel(ctx context.Context, userID string) error {
	sub, err := s.GetSubscription(ctx, userID)
	if err != nil {
		return err
	}

	if !sub.IsActive() {
		return fmt.Errorf("subscription ist bereits inaktiv")
	}

	// Call repository cancel method
	err = s.subscriptionRepo.Cancel(ctx, userID)
	if err != nil {
		return fmt.Errorf("fehler beim Kündigen: %w", err)
	}

	return nil
}

// AdminListCustomers lädt eine paginierte Kundenliste (alle Nicht-Bot-Nutzer,
// inkl. Nutzern ohne bisherige Subscription-Zeile) für das Admin-Panel.
func (s *SubscriptionService) AdminListCustomers(ctx context.Context, page, pageSize int) ([]*domain.AdminCustomerRow, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	return s.subscriptionRepo.ListAllWithUserAndTier(ctx, pageSize, offset)
}

// AdminSetTier überschreibt manuell den Tarif eines Nutzers (Support-Fall,
// z.B. Paddle-Webhook kam nicht an) - unabhängig vom Zahlungs-Fluss.
func (s *SubscriptionService) AdminSetTier(ctx context.Context, userID string, tierID domain.TierID) error {
	if tierID != domain.TierFree && tierID != domain.TierPro && tierID != domain.TierPremium {
		return fmt.Errorf("ungültiger Tarif: %s", tierID)
	}

	if err := s.subscriptionRepo.AdminOverrideTier(ctx, userID, tierID); err != nil {
		return fmt.Errorf("fehler beim manuellen Tarif-Override: %w", err)
	}

	return nil
}

// AdminGetStats berechnet eine grobe MRR-Schätzung + Anzahl aktiver Kunden
// aus allen existierenden Subscription-Zeilen. IsActive()/BillingCycle
// werden wiederverwendet statt in SQL dupliziert (gleiches Prinzip wie bei
// der laufenden Tier-Durchsetzung der Clip-Pipeline).
func (s *SubscriptionService) AdminGetStats(ctx context.Context) (mrr float64, activeCustomers int, err error) {
	subs, err := s.subscriptionRepo.GetAllActiveForMRR(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("fehler beim Laden der Subscriptions für MRR: %w", err)
	}

	for _, sub := range subs {
		if !sub.IsActive() || sub.Tier == nil {
			continue
		}
		activeCustomers++

		if sub.BillingCycle != nil && *sub.BillingCycle == domain.BillingYearly {
			mrr += sub.Tier.PriceYearly / 12
		} else {
			mrr += sub.Tier.PriceMonthly
		}
	}

	return mrr, activeCustomers, nil
}

// createFreeTier creates a free tier subscription for a new user
func (s *SubscriptionService) createFreeTier(ctx context.Context, userID string) (*domain.UserSubscription, error) {
	// Create subscription
	sub, err := s.subscriptionRepo.Create(ctx, domain.UserSubscriptionCreateInput{
		UserID: userID,
		TierID: domain.TierFree,
	})
	if err != nil {
		return nil, fmt.Errorf("fehler beim Erstellen der Free Tier Subscription: %w", err)
	}

	return sub, nil
}