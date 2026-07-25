package repository

import (
	"context"
	"time"

	"demo/backend-go/internal/domain"
)

type LoyaltyRepository interface {
	GetSettings(ctx context.Context, userTwitchID string) (*domain.LoyaltySettings, error)

	// UpsertSettings creates or updates the per-channel settings row.
	UpsertSettings(ctx context.Context, settings *domain.LoyaltySettings) error

	// GetDueForAccrual returns enabled channels whose NextAccrualAt has
	// passed - the scheduler's due-query, mirroring ScheduledMessage's
	// GetDue.
	GetDueForAccrual(ctx context.Context, limit int) ([]*domain.LoyaltySettings, error)

	// MarkAccrued advances a channel's next_accrual_at after a tick,
	// regardless of whether it was live - an offline channel's timer keeps
	// running, it doesn't "catch up" once live again.
	MarkAccrued(ctx context.Context, userTwitchID string, nextAccrualAt time.Time) error

	// CreditPoints adds each credit's Amount to that viewer's running
	// balance for userTwitchID, in one transaction.
	CreditPoints(ctx context.Context, userTwitchID string, credits []domain.ViewerCredit) error

	// GetLeaderboard returns a channel's viewers ordered by points
	// descending, paginated.
	GetLeaderboard(ctx context.Context, userTwitchID string, limit, offset int) ([]*domain.LoyaltyPoints, int64, error)

	// GetViewerPoints looks up one viewer's balance by login
	// (case-insensitive). Returns nil, nil if the viewer has no rows yet.
	GetViewerPoints(ctx context.Context, userTwitchID, viewerLogin string) (*domain.LoyaltyPoints, error)

	// GetViewerLoginsAboveThreshold returns the logins of every viewer
	// whose points meet or exceed threshold - backs the Regulars merge in
	// AutomodService and the !regulars chat command.
	GetViewerLoginsAboveThreshold(ctx context.Context, userTwitchID string, threshold int) ([]string, error)
}
