package repository

import (
	"context"
	"demo/backend-go/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error

	GetByTwitchID(ctx context.Context, twitchID string) (*domain.User, error)

	GetByDiscordID(ctx context.Context, discordID string) (*domain.User, error)

	Update(ctx context.Context, user *domain.User) error

	Delete(ctx context.Context, twitchID string) error

	Exists(ctx context.Context, twitchID string) (bool, error)
}