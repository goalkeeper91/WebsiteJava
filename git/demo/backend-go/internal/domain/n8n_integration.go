package domain

import "time"

type N8NIntegration struct {
	ID              int64     `json:"id" db:"id"`
	UserID          string    `json:"userId" db:"user_id"`
	Enabled         bool      `json:"enabled" db:"enabled"`
	WebhookBaseURL  *string   `json:"webhookBaseUrl" db:"webhook_base_url"`
	APIKey          *string   `json:"-" db:"api_key"` // Never serialize to JSON
	WorkflowsUsed   int       `json:"workflowsUsed" db:"workflows_used"`
	VotesThisMonth  int       `json:"votesThisMonth" db:"votes_this_month"`
	LastResetAt     time.Time `json:"lastResetAt" db:"last_reset_at"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time `json:"updatedAt" db:"updated_at"`
}

type N8NIntegrationCreateInput struct {
	UserID         string
	WebhookBaseURL *string
	APIKey         *string
}

type N8NIntegrationUpdateInput struct {
	Enabled        *bool
	WebhookBaseURL *string
	APIKey         *string
}

func NewN8NIntegration(userID string) *N8NIntegration {
	now := time.Now()
	return &N8NIntegration{
		UserID:      userID,
		Enabled:     false,
		LastResetAt: now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (n *N8NIntegration) Enable(webhookBaseURL string) {
	n.Enabled = true
	n.WebhookBaseURL = &webhookBaseURL
	n.UpdatedAt = time.Now()
}

func (n *N8NIntegration) Disable() {
	n.Enabled = false
	n.UpdatedAt = time.Now()
}

func (n *N8NIntegration) IsReady() bool {
	return n.Enabled && n.WebhookBaseURL != nil && *n.WebhookBaseURL != ""
}

func (n *N8NIntegration) GetWebhookURL(path string) string {
	if n.WebhookBaseURL == nil {
		return ""
	}
	return *n.WebhookBaseURL + "/webhook/" + path
}

// ShouldResetMonthlyCounter checks if votes_this_month should reset
func (n *N8NIntegration) ShouldResetMonthlyCounter() bool {
	now := time.Now()
	return now.Month() != n.LastResetAt.Month() ||
		now.Year() != n.LastResetAt.Year()
}

func (n *N8NIntegration) ResetMonthlyCounter() {
	n.VotesThisMonth = 0
	n.LastResetAt = time.Now()
	n.UpdatedAt = time.Now()
}

func (n *N8NIntegration) IncrementVotes() {
	n.VotesThisMonth++
	n.UpdatedAt = time.Now()
}