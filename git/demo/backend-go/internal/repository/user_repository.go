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

    GetBotUsers(ctx context.Context) ([]*domain.User, error)

	// GetAllNonBotUsers backs the activity-eventsub reconciliation loop -
	// every broadcaster (not the bot account itself) potentially needs
	// follow/sub/gift/cheer/raid subscriptions.
	GetAllNonBotUsers(ctx context.Context) ([]*domain.User, error)

	GetFirstBotUser(ctx context.Context) (*domain.User, error)
}