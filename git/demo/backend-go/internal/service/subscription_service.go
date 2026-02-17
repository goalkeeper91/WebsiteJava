package service

import (
	"context"
	"fmt"
	"log"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
)

type SubscriptionService struct {
	subscriptionRepo repository.UserSubscriptionRepository
	tierRepo         repository.SubscriptionTierRepository
	analyticsRepo    repository.UsageAnalyticsRepository
}

func NewSubscriptionService(
	subscriptionRepo repository.UserSubscriptionRepository,
	tierRepo repository.SubscriptionTierRepository,
	analyticsRepo repository.UsageAnalyticsRepository,
) *SubscriptionService {
	return &SubscriptionService{
		subscriptionRepo: subscriptionRepo,
		tierRepo:         tierRepo,
		analyticsRepo:    analyticsRepo,
	}
}

// GetSubscription returns a user's subscription with tier details
func (s *SubscriptionService) GetSubscription(ctx context.Context, userID string) (*domain.UserSubscription, error) {
	sub, err := s.subscriptionRepo.GetByUserIDWithTier(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Subscription: %w", err)
	}

	// Auto-create free tier if none exists
	if sub == nil {
		sub, err = s.createFreeTier(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("fehler beim Erstellen des Free Tiers: %w", err)
		}
	}

	return sub, nil
}

// GetTiers returns all active subscription tiers (for pricing page)
func (s *SubscriptionService) GetTiers(ctx context.Context) ([]*domain.SubscriptionTier, error) {
	return s.tierRepo.GetActive(ctx)
}

// CanUseFeature checks if a user's subscription includes a specific feature
func (s *SubscriptionService) CanUseFeature(ctx context.Context, userID string, feature string) (bool, error) {
	sub, err := s.GetSubscription(ctx, userID)
	if err != nil {
		return false, err
	}

	if !sub.IsActive() {
		return false, nil
	}

	if sub.Tier == nil {
		tier, err := s.tierRepo.GetByID(ctx, sub.TierID)
		if err != nil || tier == nil {
			return false, nil
		}
		sub.Tier = tier
	}

	return sub.Tier.HasFeature(feature), nil
}

// CanUseN8N checks if a user can use n8n features
func (s *SubscriptionService) CanUseN8N(ctx context.Context, userID string) (bool, error) {
	return s.CanUseFeature(ctx, userID, "advanced_commands")
}

// CheckCommandLimit checks if user has reached their command limit
func (s *SubscriptionService) CheckCommandLimit(ctx context.Context, userID string, currentCount int) error {
	sub, err := s.GetSubscription(ctx, userID)
	if err != nil {
		return err
	}

	if sub.Tier == nil {
		return nil
	}

	if sub.Tier.IsCommandLimitReached(currentCount) {
		return domain.ErrCommandLimitReached
	}

	return nil
}

// CheckWorkflowLimit checks if user has reached their workflow limit
func (s *SubscriptionService) CheckWorkflowLimit(ctx context.Context, userID string, currentCount int) error {
	sub, err := s.GetSubscription(ctx, userID)
	if err != nil {
		return err
	}

	if sub.Tier == nil {
		return nil
	}

	if sub.Tier.IsWorkflowLimitReached(currentCount) {
		return domain.ErrWorkflowLimitReached
	}

	return nil
}

// CheckVoteLimit checks if user has reached their monthly vote limit
func (s *SubscriptionService) CheckVoteLimit(ctx context.Context, userID string, currentCount int) error {
	sub, err := s.GetSubscription(ctx, userID)
	if err != nil {
		return err
	}

	if sub.Tier == nil {
		return nil
	}

	if sub.Tier.IsVoteLimitReached(currentCount) {
		return domain.ErrVoteLimitReached
	}

	return nil
}

// Upgrade upgrades a user to a new tier
func (s *SubscriptionService) Upgrade(ctx context.Context, userID string, tierID domain.TierID) (*domain.UserSubscription, error) {
	// Verify tier exists
	tier, err := s.tierRepo.GetByID(ctx, tierID)
	if err != nil || tier == nil {
		return nil, fmt.Errorf("tier nicht gefunden: %s", tierID)
	}

	exists, err := s.subscriptionRepo.Exists(ctx, userID)
	if err != nil {
		return nil, err
	}

	if exists {
		err = s.subscriptionRepo.Update(ctx, userID, domain.UserSubscriptionUpdateInput{
			TierID: &tierID,
			Status: statusPtr(domain.SubscriptionActive),
		})
	} else {
		_, err = s.subscriptionRepo.Create(ctx, domain.UserSubscriptionCreateInput{
			UserID: userID,
			TierID: tierID,
		})
	}

	if err != nil {
		return nil, fmt.Errorf("fehler beim Upgrade der Subscription: %w", err)
	}

	log.Printf("✅ User %s upgraded to tier: %s", userID, tierID)
	return s.subscriptionRepo.GetByUserIDWithTier(ctx, userID)
}

// Cancel cancels a user's subscription (downgrades to free at period end)
func (s *SubscriptionService) Cancel(ctx context.Context, userID string) error {
	if err := s.subscriptionRepo.Cancel(ctx, userID); err != nil {
		return fmt.Errorf("fehler beim Kündigen der Subscription: %w", err)
	}

	log.Printf("⚠️ Subscription gekündigt für User: %s", userID)
	return nil
}

// createFreeTier auto-creates a free subscription for new users
func (s *SubscriptionService) createFreeTier(ctx context.Context, userID string) (*domain.UserSubscription, error) {
	_, err := s.subscriptionRepo.Create(ctx, domain.UserSubscriptionCreateInput{
		UserID: userID,
		TierID: domain.TierFree,
	})
	if err != nil {
		return nil, err
	}

	return s.subscriptionRepo.GetByUserIDWithTier(ctx, userID)
}

func statusPtr(s domain.SubscriptionStatus) *domain.SubscriptionStatus {
	return &s
}