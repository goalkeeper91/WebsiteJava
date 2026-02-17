package repository

import (
	"context"
	"demo/backend-go/internal/domain"
)

type N8NIntegrationRepository interface {
	Create(ctx context.Context, input domain.N8NIntegrationCreateInput) (*domain.N8NIntegration, error)

	GetByUserID(ctx context.Context, userID string) (*domain.N8NIntegration, error)

	Update(ctx context.Context, userID string, input domain.N8NIntegrationUpdateInput) error

	IncrementVotes(ctx context.Context, userID string) error

	ResetMonthlyCounter(ctx context.Context, userID string) error

	Exists(ctx context.Context, userID string) (bool, error)
}