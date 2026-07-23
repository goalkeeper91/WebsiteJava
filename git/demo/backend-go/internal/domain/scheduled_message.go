package domain

import "time"

// MinScheduledMessageIntervalSeconds is the smallest interval a streamer can
// configure - prevents an accidental near-zero interval from spamming chat.
const MinScheduledMessageIntervalSeconds = 60

// ScheduledMessage is either free text (Message set, CommandID nil) or a
// link to an existing chat command (CommandID set, Message nil) - never
// both, never neither (enforced by a DB CHECK constraint too). The actual
// text for a command-linked entry is resolved at send time from the
// command's current response, not stored here.
type ScheduledMessage struct {
	ID              int64      `json:"id"`
	ChannelID       string     `json:"channel_id"`
	Message         *string    `json:"message"`
	CommandID       *int64     `json:"command_id,omitempty"`
	IntervalSeconds int        `json:"interval_seconds"`
	Enabled         bool       `json:"enabled"`
	OnlyWhenLive    bool       `json:"only_when_live"`
	NextSendAt      time.Time  `json:"next_send_at"`
	LastSentAt      *time.Time `json:"last_sent_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func NewFreeTextScheduledMessage(channelID, message string, intervalSeconds int) *ScheduledMessage {
	now := time.Now()
	return &ScheduledMessage{
		ChannelID:       channelID,
		Message:         &message,
		IntervalSeconds: intervalSeconds,
		Enabled:         true,
		OnlyWhenLive:    true,
		NextSendAt:      now.Add(time.Duration(intervalSeconds) * time.Second),
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

func NewCommandScheduledMessage(channelID string, commandID int64, intervalSeconds int) *ScheduledMessage {
	now := time.Now()
	return &ScheduledMessage{
		ChannelID:       channelID,
		CommandID:       &commandID,
		IntervalSeconds: intervalSeconds,
		Enabled:         true,
		OnlyWhenLive:    true,
		NextSendAt:      now.Add(time.Duration(intervalSeconds) * time.Second),
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// Update only applies to free-text entries - a command-linked entry's text
// always tracks the command itself, so message updates are ignored there.
func (m *ScheduledMessage) Update(message *string, intervalSeconds *int, enabled, onlyWhenLive *bool) {
	if message != nil && *message != "" && m.CommandID == nil {
		m.Message = message
	}
	if intervalSeconds != nil {
		m.IntervalSeconds = *intervalSeconds
	}
	if enabled != nil {
		m.Enabled = *enabled
	}
	if onlyWhenLive != nil {
		m.OnlyWhenLive = *onlyWhenLive
	}
	m.UpdatedAt = time.Now()
}

func (m *ScheduledMessage) Toggle(enabled bool) {
	m.Enabled = enabled
	m.UpdatedAt = time.Now()
}

func (m *ScheduledMessage) IsCommandLinked() bool {
	return m.CommandID != nil
}

func ValidateScheduledMessage(message string, intervalSeconds int) error {
	if message == "" {
		return ErrEmptyMessage
	}
	return ValidateScheduledInterval(intervalSeconds)
}

func ValidateScheduledInterval(intervalSeconds int) error {
	if intervalSeconds < MinScheduledMessageIntervalSeconds {
		return ErrIntervalTooShort
	}
	return nil
}
