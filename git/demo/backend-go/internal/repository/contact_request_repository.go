package repository

import (
	"context"
	"demo/backend-go/internal/domain"
	"time"
)

type ContactRequestRepository interface {
	Create(ctx context.Context, request *domain.ContactRequest) error

	GetAll(ctx context.Context) ([]*domain.ContactRequest, error)

	DeleteOlderThan(ctx context.Context, date time.Time) error
}