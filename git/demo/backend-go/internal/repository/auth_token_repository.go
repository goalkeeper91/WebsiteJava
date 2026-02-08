package repository

import (
	"context"
	"demo/backend-go/internal/domain"
)

type AuthTokenRepository interface {
	Create(ctx context.Context, token *domain.AuthToken) error

	GetByTwitchUserID(ctx context.Context, twitchUserID string) (*domain.AuthToken, error)

	Update(ctx context.Context, token *domain.AuthToken) error

	Delete(ctx context.Context, twitchUserID string) error

	Upsert(ctx context.Context, token *domain.AuthToken) error
}