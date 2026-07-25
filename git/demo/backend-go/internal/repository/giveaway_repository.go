package repository

import (
	"context"

	"demo/backend-go/internal/domain"
)

type GiveawayRepository interface {
	CreateGiveaway(ctx context.Context, g *domain.Giveaway) error

	// GetOpenGiveaway returns the channel's currently open giveaway, or nil
	// if none is open.
	GetOpenGiveaway(ctx context.Context, userTwitchID string) (*domain.Giveaway, error)

	// GetAllOpenGiveaways returns every currently-open giveaway across all
	// channels - backs the bot's startup keyword-cache warmup.
	GetAllOpenGiveaways(ctx context.Context) ([]*domain.Giveaway, error)

	// AddEntry records a viewer's participation - a repeat attempt for the
	// same giveaway+viewer is a no-op (inserted=false), so typing the entry
	// command again never farms extra tickets.
	AddEntry(ctx context.Context, giveawayID int64, viewerTwitchID, viewerLogin string, entries int) (inserted bool, err error)

	GetEntries(ctx context.Context, giveawayID int64) ([]*domain.GiveawayEntry, error)

	GetEntryCount(ctx context.Context, giveawayID int64) (int, error)

	CloseGiveaway(ctx context.Context, giveawayID int64, winnerTwitchID, winnerLogin string) error

	// GetHistory returns a channel's past giveaways, newest first, paginated.
	GetHistory(ctx context.Context, userTwitchID string, limit, offset int) ([]*domain.Giveaway, int64, error)
}
