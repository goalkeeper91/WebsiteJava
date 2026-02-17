package domain

import (
	"strings"
	"time"
)

type CommandType string

const (
	CommandTypeSimple   CommandType = "simple"
	CommandTypeAdvanced CommandType = "advanced"
)

type ChatCommand struct {
	ID             int64       `json:"id"`
	ChannelID      string      `json:"channel_id"`
	Trigger        string      `json:"trigger"`
	Response       string      `json:"response"`
	Cooldown       int         `json:"cooldown"` // in Sekunden
	Enabled        bool        `json:"enabled"`
	CommandType    CommandType `json:"command_type"`
	N8NWorkflowID  *string     `json:"n8n_workflow_id,omitempty"`
	UsageCount     int64       `json:"usage_count"`
	LastUsedAt     *time.Time  `json:"last_used_at,omitempty"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

func NewChatCommand(channelID, trigger, response string, cooldown int) *ChatCommand {
	now := time.Now()
	return &ChatCommand{
		ChannelID:   channelID,
		Trigger:     NormalizeTrigger(trigger),
		Response:    response,
		Cooldown:    cooldown,
		Enabled:     true,
		CommandType: CommandTypeSimple, // Default: simple
		UsageCount:  0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func NewAdvancedChatCommand(channelID, trigger, workflowID string, cooldown int) *ChatCommand {
	now := time.Now()
	return &ChatCommand{
		ChannelID:     channelID,
		Trigger:       NormalizeTrigger(trigger),
		Response:      "", // No static response for advanced
		Cooldown:      cooldown,
		Enabled:       true,
		CommandType:   CommandTypeAdvanced,
		N8NWorkflowID: &workflowID,
		UsageCount:    0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func (c *ChatCommand) Update(trigger, response *string, cooldown *int, enabled *bool) {
	if trigger != nil && *trigger != "" {
		c.Trigger = NormalizeTrigger(*trigger)
	}
	if response != nil {
		c.Response = *response
	}
	if cooldown != nil {
		c.Cooldown = *cooldown
	}
	if enabled != nil {
		c.Enabled = *enabled
	}
	c.UpdatedAt = time.Now()
}

func (c *ChatCommand) Toggle(enabled bool) {
	c.Enabled = enabled
	c.UpdatedAt = time.Now()
}

func (c *ChatCommand) TrackUsage() {
	c.UsageCount++
	now := time.Now()
	c.LastUsedAt = &now
}

func (c *ChatCommand) IsSimple() bool {
	return c.CommandType == CommandTypeSimple
}

func (c *ChatCommand) IsAdvanced() bool {
	return c.CommandType == CommandTypeAdvanced
}

func (c *ChatCommand) IsReadyForN8N() bool {
	return c.IsAdvanced() && c.N8NWorkflowID != nil && *c.N8NWorkflowID != ""
}

func NormalizeTrigger(trigger string) string {
	return strings.ToLower(strings.TrimSpace(trigger))
}

func ValidateTrigger(trigger string) error {
	if trigger == "" {
		return ErrEmptyTrigger
	}
	if len(trigger) > 100 {
		return ErrTriggerTooLong
	}
	return nil
}