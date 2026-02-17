package repository

import (
	"context"
	"demo/backend-go/internal/domain"
)

type SubscriptionTierRepository interface {
	GetAll(ctx context.Context) ([]*domain.SubscriptionTier, error)

	GetByID(ctx context.Context, tierID domain.TierID) (*domain.SubscriptionTier, error)

	GetActive(ctx context.Context) ([]*domain.SubscriptionTier, error)
}