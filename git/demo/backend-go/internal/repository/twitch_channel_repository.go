package repository

import (
	"context"
	"demo/backend-go/internal/domain"
)

type TwitchChannelRepository interface {
	Create(ctx context.Context, channel *domain.TwitchChannel) error

	GetByTwitchUserID(ctx context.Context, twitchUserID string) (*domain.TwitchChannel, error)

	// GetByID looks a channel up by its internal ID - needed wherever a
	// stored foreign key (e.g. scheduled_messages.channel_id) only carries
	// the internal ID and must be resolved back to a Twitch user ID.
	GetByID(ctx context.Context, id int64) (*domain.TwitchChannel, error)

	Update(ctx context.Context, channel *domain.TwitchChannel) error

	Upsert(ctx context.Context, channel *domain.TwitchChannel) error

	Delete(ctx context.Context, twitchUserID string) error

	Exists(ctx context.Context, twitchUserID string) (bool, error)
}