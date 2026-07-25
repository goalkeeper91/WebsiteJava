package repository

import (
	"context"

	"demo/backend-go/internal/domain"
)

type AutomodRepository interface {
	GetSettings(ctx context.Context, userTwitchID string) (*domain.AutomodSettings, error)

	// UpsertSettings creates or updates the per-channel settings row.
	UpsertSettings(ctx context.Context, settings *domain.AutomodSettings) error

	// GetAllEnabledSettings is used by the bot to build its full in-memory
	// cache on startup/reload - only channels with automod switched on.
	GetAllEnabledSettings(ctx context.Context) ([]*domain.AutomodSettings, error)

	// RecordViolation increments (or, if the offender's last violation is
	// older than the streak TTL, resets-then-increments) the violation
	// counter for this channel+offender pair, atomically, and returns the
	// new count.
	RecordViolation(ctx context.Context, userTwitchID, offenderTwitchID string) (violationCount int, err error)

	// LogEvent appends one row to the automod_events log (separate from the
	// violation counter above).
	LogEvent(ctx context.Context, event *domain.AutomodEvent) error

	// GetEvents returns the most recent automod_events rows for a channel,
	// newest first.
	GetEvents(ctx context.Context, userTwitchID string, limit, offset int) ([]*domain.AutomodEvent, int64, error)
}
