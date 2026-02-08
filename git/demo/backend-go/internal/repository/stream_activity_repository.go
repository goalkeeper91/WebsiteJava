package repository

import (
	"context"
	"demo/backend-go/internal/domain"
)

type StreamActivityRepository interface {
	Create(ctx context.Context, activity *domain.StreamActivity) error

	GetRecent(ctx context.Context, twitchUserID string, limit int) ([]*domain.StreamActivity, error)

	GetByType(ctx context.Context, twitchUserID string, activityType domain.ActivityType, limit int) ([]*domain.StreamActivity, error)

	DeleteOlderThan(ctx context.Context, timestamp int64) error
}