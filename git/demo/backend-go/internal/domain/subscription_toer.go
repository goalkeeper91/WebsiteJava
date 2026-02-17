package domain

import (
	"encoding/json"
	"time"
)

// TierID represents a subscription tier identifier
type TierID string

const (
	TierFree    TierID = "free"
	TierPro     TierID = "pro"
	TierPremium TierID = "premium"
)

// TierFeatures represents the features available in a tier
type TierFeatures struct {
	SimpleCommands    bool `json:"simple_commands"`
	AdvancedCommands  bool `json:"advanced_commands"`
	DiscordIntegration bool `json:"discord_integration"`
	Analytics         bool `json:"analytics"`
	WorkflowTemplates bool `json:"workflow_templates"`
	CustomWorkflows   bool `json:"custom_workflows"`
	APIAccess         bool `json:"api_access"`
	PrioritySupport   bool `json:"priority_support"`
}

type SubscriptionTier struct {
	ID                 TierID          `json:"id" db:"id"`
	Name               string          `json:"name" db:"name"`
	PriceMonthly       float64         `json:"priceMonthly" db:"price_monthly"`
	PriceYearly        float64         `json:"priceYearly" db:"price_yearly"`
	MaxCommands        *int            `json:"maxCommands" db:"max_commands"`         // nil = unlimited
	MaxWorkflows       *int            `json:"maxWorkflows" db:"max_workflows"`        // nil = unlimited
	MaxVotesPerMonth   *int            `json:"maxVotesPerMonth" db:"max_votes_per_month"` // nil = unlimited
	Features           TierFeatures    `json:"features" db:"features"`
	IsActive           bool            `json:"isActive" db:"is_active"`
	CreatedAt          time.Time       `json:"createdAt" db:"created_at"`
}

func (t *SubscriptionTier) ScanFeatures(data []byte) error {
	return json.Unmarshal(data, &t.Features)
}

func (t *SubscriptionTier) HasFeature(feature string) bool {
	switch feature {
	case "simple_commands":
		return t.Features.SimpleCommands
	case "advanced_commands":
		return t.Features.AdvancedCommands
	case "discord_integration":
		return t.Features.DiscordIntegration
	case "analytics":
		return t.Features.Analytics
	case "workflow_templates":
		return t.Features.WorkflowTemplates
	case "custom_workflows":
		return t.Features.CustomWorkflows
	case "api_access":
		return t.Features.APIAccess
	case "priority_support":
		return t.Features.PrioritySupport
	}
	return false
}

func (t *SubscriptionTier) IsCommandLimitReached(currentCount int) bool {
	if t.MaxCommands == nil {
		return false // unlimited
	}
	return currentCount >= *t.MaxCommands
}

func (t *SubscriptionTier) IsWorkflowLimitReached(currentCount int) bool {
	if t.MaxWorkflows == nil {
		return false // unlimited
	}
	return currentCount >= *t.MaxWorkflows
}

func (t *SubscriptionTier) IsVoteLimitReached(currentCount int) bool {
	if t.MaxVotesPerMonth == nil {
		return false // unlimited
	}
	return currentCount >= *t.MaxVotesPerMonth
}