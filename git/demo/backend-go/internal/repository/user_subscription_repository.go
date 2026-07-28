package repository

import (
	"context"
	"demo/backend-go/internal/domain"
)

type UserSubscriptionRepository interface {
	Create(ctx context.Context, input domain.UserSubscriptionCreateInput) (*domain.UserSubscription, error)

	GetByUserID(ctx context.Context, userID string) (*domain.UserSubscription, error)

	GetByUserIDWithTier(ctx context.Context, userID string) (*domain.UserSubscription, error)

	// GetByPaddleCustomerID resolves a webhook event back to a local user
	// when the event only carries the Paddle customer ID (not customData).
	GetByPaddleCustomerID(ctx context.Context, paddleCustomerID string) (*domain.UserSubscription, error)

	Update(ctx context.Context, userID string, input domain.UserSubscriptionUpdateInput) error

	Cancel(ctx context.Context, userID string) error

	Exists(ctx context.Context, userID string) (bool, error)
}