package repository

import (
	"context"
	"demo/backend-go/internal/domain"
	"time"
)

type UsageAnalyticsRepository interface {
	Track(ctx context.Context, input domain.UsageAnalyticsCreateInput) error

	GetSummary(ctx context.Context, userID string, from, to time.Time) (*domain.UsageSummary, error)

	GetCountByType(ctx context.Context, userID string, metricType domain.MetricType, from, to time.Time) (int, error)
}