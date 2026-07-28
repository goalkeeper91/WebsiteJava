package repository

import (
	"context"
	"time"

	"demo/backend-go/internal/domain"
)

type StreamSessionRepository interface {
	// GetOpenSessionsByUsers returns every channel's currently open session
	// (ended_at IS NULL), keyed by twitch_user_id - one batch query for the
	// sampler tick instead of a query per channel.
	GetOpenSessionsByUsers(ctx context.Context, twitchUserIDs []string) (map[string]*domain.StreamSession, error)

	CreateSession(ctx context.Context, twitchUserID string, startedAt time.Time) (*domain.StreamSession, error)

	// RecordSample adds one viewer-count reading to a session's running
	// sum/count and bumps peak_viewers if this sample is higher.
	RecordSample(ctx context.Context, sessionID int64, viewerCount int) error

	CloseSession(ctx context.Context, sessionID int64) error

	// GetLastClosedSession returns a channel's most recently ended session,
	// or nil if it has never had one - backs the offline "letzter Stream"
	// average.
	GetLastClosedSession(ctx context.Context, twitchUserID string) (*domain.StreamSession, error)
}
