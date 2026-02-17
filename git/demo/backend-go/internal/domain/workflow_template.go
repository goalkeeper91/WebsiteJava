package domain

import (
	"encoding/json"
	"time"
)

type WorkflowCategory string

const (
	WorkflowCategoryVoting     WorkflowCategory = "voting"
	WorkflowCategoryGiveaway   WorkflowCategory = "giveaway"
	WorkflowCategoryModeration WorkflowCategory = "moderation"
	WorkflowCategoryAnalytics  WorkflowCategory = "analytics"
	WorkflowCategoryCustom     WorkflowCategory = "custom"
)

type WorkflowTemplate struct {
	ID               int64            `json:"id" db:"id"`
	Name             string           `json:"name" db:"name"`
	Description      *string          `json:"description" db:"description"`
	Category         WorkflowCategory `json:"category" db:"category"`
	N8NWorkflowJSON  json.RawMessage  `json:"n8nWorkflowJson,omitempty" db:"n8n_workflow_json"`
	TierRequired     TierID           `json:"tierRequired" db:"tier_required"`
	IsPublic         bool             `json:"isPublic" db:"is_public"`
	UsageCount       int              `json:"usageCount" db:"usage_count"`
	CreatedBy        *string          `json:"createdBy" db:"created_by"`
	CreatedAt        time.Time        `json:"createdAt" db:"created_at"`
	UpdatedAt        time.Time        `json:"updatedAt" db:"updated_at"`
}

type WorkflowTemplateCreateInput struct {
	Name            string
	Description     *string
	Category        WorkflowCategory
	N8NWorkflowJSON json.RawMessage
	TierRequired    TierID
	IsPublic        bool
	CreatedBy       *string
}

type WorkflowTemplateUpdateInput struct {
	Name            *string
	Description     *string
	N8NWorkflowJSON json.RawMessage
	TierRequired    *TierID
	IsPublic        *bool
}

// WorkflowTemplateSummary is used in listings (without workflow JSON)
type WorkflowTemplateSummary struct {
	ID           int64            `json:"id"`
	Name         string           `json:"name"`
	Description  *string          `json:"description"`
	Category     WorkflowCategory `json:"category"`
	TierRequired TierID           `json:"tierRequired"`
	IsPublic     bool             `json:"isPublic"`
	UsageCount   int              `json:"usageCount"`
}

func (t *WorkflowTemplate) ToSummary() WorkflowTemplateSummary {
	return WorkflowTemplateSummary{
		ID:           t.ID,
		Name:         t.Name,
		Description:  t.Description,
		Category:     t.Category,
		TierRequired: t.TierRequired,
		IsPublic:     t.IsPublic,
		UsageCount:   t.UsageCount,
	}
}

func (t *WorkflowTemplate) IsAccessibleBy(tierID TierID) bool {
	switch t.TierRequired {
	case TierFree:
		return true
	case TierPro:
		return tierID == TierPro || tierID == TierPremium
	case TierPremium:
		return tierID == TierPremium
	}
	return false
}