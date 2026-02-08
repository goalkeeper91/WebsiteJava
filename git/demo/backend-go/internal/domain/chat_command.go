package domain

import (
	"strings"
	"time"
)

type ChatCommand struct {
	ID        int64     `json:"id"`
	ChannelID string    `json:"channel_id"`
	Trigger   string    `json:"trigger"`
	Response  string    `json:"response"`
	Cooldown  int       `json:"cooldown"` // in Sekunden
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewChatCommand(channelID, trigger, response string, cooldown int) *ChatCommand {
	now := time.Now()
	return &ChatCommand{
		ChannelID: channelID,
		Trigger:   NormalizeTrigger(trigger),
		Response:  response,
		Cooldown:  cooldown,
		Enabled:   true,
		CreatedAt: now,
		UpdatedAt: now,
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