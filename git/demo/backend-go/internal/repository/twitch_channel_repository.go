package repository

import (
	"context"
	"demo/backend-go/internal/domain"
)

type TwitchChannelRepository interface {
	Create(ctx context.Context, channel *domain.TwitchChannel) error

	GetByTwitchUserID(ctx context.Context, twitchUserID string) (*domain.TwitchChannel, error)

	Update(ctx context.Context, channel *domain.TwitchChannel) error

	Upsert(ctx context.Context, channel *domain.TwitchChannel) error

	Delete(ctx context.Context, twitchUserID string) error

	Exists(ctx context.Context, twitchUserID string) (bool, error)
}