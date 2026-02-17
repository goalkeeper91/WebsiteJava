package domain

import "time"

type MetricType string

const (
	MetricCommandUsage   MetricType = "command_usage"
	MetricVoteCast       MetricType = "vote_cast"
	MetricVoteSession    MetricType = "vote_session"
	MetricWorkflowRun    MetricType = "workflow_run"
	MetricN8NWebhook     MetricType = "n8n_webhook"
)

type UsageAnalytics struct {
	ID           int64      `json:"id" db:"id"`
	UserID       string     `json:"userId" db:"user_id"`
	MetricType   MetricType `json:"metricType" db:"metric_type"`
	MetricValue  int        `json:"metricValue" db:"metric_value"`
	RecordedAt   time.Time  `json:"recordedAt" db:"recorded_at"`
}

type UsageAnalyticsCreateInput struct {
	UserID      string
	MetricType  MetricType
	MetricValue int
}

// UsageSummary aggregates metrics for a user
type UsageSummary struct {
	UserID          string         `json:"userId"`
	Period          string         `json:"period"` // "daily", "weekly", "monthly"
	CommandsUsed    int            `json:"commandsUsed"`
	VotesCast       int            `json:"votesCast"`
	VoteSessions    int            `json:"voteSessions"`
	WorkflowsRun    int            `json:"workflowsRun"`
	ByDay           []DailyMetric  `json:"byDay,omitempty"`
}

type DailyMetric struct {
	Date        time.Time `json:"date"`
	MetricType  MetricType `json:"metricType"`
	Total       int        `json:"total"`
}

func NewUsageAnalytics(userID string, metricType MetricType, value int) *UsageAnalytics {
	return &UsageAnalytics{
		UserID:      userID,
		MetricType:  metricType,
		MetricValue: value,
		RecordedAt:  time.Now(),
	}
}

func TrackCommandUsage(userID string) *UsageAnalytics {
	return NewUsageAnalytics(userID, MetricCommandUsage, 1)
}

func TrackVoteCast(userID string) *UsageAnalytics {
	return NewUsageAnalytics(userID, MetricVoteCast, 1)
}

func TrackVoteSession(userID string) *UsageAnalytics {
	return NewUsageAnalytics(userID, MetricVoteSession, 1)
}

func TrackWorkflowRun(userID string) *UsageAnalytics {
	return NewUsageAnalytics(userID, MetricWorkflowRun, 1)
}