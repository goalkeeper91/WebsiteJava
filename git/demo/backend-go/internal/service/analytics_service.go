package service

import (
	"context"
	"fmt"
	"time"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
)

type AnalyticsService struct {
	analyticsRepo   repository.UsageAnalyticsRepository
	subscriptionSvc *SubscriptionService
}

func NewAnalyticsService(
	analyticsRepo repository.UsageAnalyticsRepository,
	subscriptionSvc *SubscriptionService,
) *AnalyticsService {
	return &AnalyticsService{
		analyticsRepo:   analyticsRepo,
		subscriptionSvc: subscriptionSvc,
	}
}

// GetDashboardStats returns an overview of stats for the current month
func (s *AnalyticsService) GetDashboardStats(ctx context.Context, userID string) (*domain.UsageSummary, error) {
	// Check if user has analytics feature
	canUse, err := s.subscriptionSvc.CanUseFeature(ctx, userID, "analytics")
	if err != nil {
		return nil, err
	}
	if !canUse {
		return nil, domain.ErrFeatureNotAvailable
	}

	now := time.Now()
	from := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	summary, err := s.analyticsRepo.GetSummary(ctx, userID, from, now)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Analytics: %w", err)
	}

	summary.Period = "monthly"
	return summary, nil
}

// GetWeeklyStats returns stats for the last 7 days
func (s *AnalyticsService) GetWeeklyStats(ctx context.Context, userID string) (*domain.UsageSummary, error) {
	canUse, err := s.subscriptionSvc.CanUseFeature(ctx, userID, "analytics")
	if err != nil {
		return nil, err
	}
	if !canUse {
		return nil, domain.ErrFeatureNotAvailable
	}

	now := time.Now()
	from := now.AddDate(0, 0, -7)

	summary, err := s.analyticsRepo.GetSummary(ctx, userID, from, now)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Weekly Stats: %w", err)
	}

	summary.Period = "weekly"
	return summary, nil
}

// Track is a convenience method for tracking individual metric events
func (s *AnalyticsService) Track(ctx context.Context, userID string, metricType domain.MetricType) {
	_ = s.analyticsRepo.Track(ctx, domain.UsageAnalyticsCreateInput{
		UserID:      userID,
		MetricType:  metricType,
		MetricValue: 1,
	})
}