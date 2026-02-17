package repository

import (
	"context"
	"demo/backend-go/internal/domain"
)

type WorkflowTemplateRepository interface {
	Create(ctx context.Context, input domain.WorkflowTemplateCreateInput) (*domain.WorkflowTemplate, error)

	GetByID(ctx context.Context, id int64) (*domain.WorkflowTemplate, error)

	GetAll(ctx context.Context) ([]*domain.WorkflowTemplate, error)

	GetByCategory(ctx context.Context, category domain.WorkflowCategory) ([]*domain.WorkflowTemplate, error)

	GetByTier(ctx context.Context, tierID domain.TierID) ([]*domain.WorkflowTemplate, error)

	GetPublic(ctx context.Context) ([]*domain.WorkflowTemplate, error)

	// GetAccessibleByTier returns all templates accessible by a given tier (including lower tiers)
	GetAccessibleByTier(ctx context.Context, tierID domain.TierID) ([]*domain.WorkflowTemplate, error)

	Update(ctx context.Context, id int64, input domain.WorkflowTemplateUpdateInput) error

	IncrementUsage(ctx context.Context, id int64) error

	Delete(ctx context.Context, id int64) error
}