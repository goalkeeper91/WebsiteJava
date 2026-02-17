package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
)

type usageAnalyticsRepository struct {
	db *sql.DB
}

func NewUsageAnalyticsRepository(db *sql.DB) repository.UsageAnalyticsRepository {
	return &usageAnalyticsRepository{db: db}
}

func (r *usageAnalyticsRepository) Track(ctx context.Context, input domain.UsageAnalyticsCreateInput) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO usage_analytics (user_id, metric_type, metric_value, recorded_at)
		VALUES ($1, $2, $3, NOW())
	`, input.UserID, input.MetricType, input.MetricValue)
	if err != nil {
		return fmt.Errorf("fehler beim Tracken der Nutzung: %w", err)
	}
	return nil
}

func (r *usageAnalyticsRepository) GetSummary(ctx context.Context, userID string, from, to time.Time) (*domain.UsageSummary, error) {
	query := `
		SELECT metric_type, SUM(metric_value) as total
		FROM usage_analytics
		WHERE user_id = $1
		  AND recorded_at >= $2
		  AND recorded_at <= $3
		GROUP BY metric_type
	`

	rows, err := r.db.QueryContext(ctx, query, userID, from, to)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Analytics: %w", err)
	}
	defer rows.Close()

	summary := &domain.UsageSummary{
		UserID: userID,
	}

	for rows.Next() {
		var metricType domain.MetricType
		var total int

		if err := rows.Scan(&metricType, &total); err != nil {
			return nil, fmt.Errorf("fehler beim Scannen: %w", err)
		}

		switch metricType {
		case domain.MetricCommandUsage:
			summary.CommandsUsed = total
		case domain.MetricVoteCast:
			summary.VotesCast = total
		case domain.MetricVoteSession:
			summary.VoteSessions = total
		case domain.MetricWorkflowRun:
			summary.WorkflowsRun = total
		}
	}

	return summary, nil
}

func (r *usageAnalyticsRepository) GetCountByType(ctx context.Context, userID string, metricType domain.MetricType, from, to time.Time) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(metric_value), 0)
		FROM usage_analytics
		WHERE user_id = $1
		  AND metric_type = $2
		  AND recorded_at >= $3
		  AND recorded_at <= $4
	`, userID, metricType, from, to).Scan(&count)

	if err != nil && err != sql.ErrNoRows {
		return 0, fmt.Errorf("fehler beim Zählen der Metriken: %w", err)
	}

	return count, nil
}