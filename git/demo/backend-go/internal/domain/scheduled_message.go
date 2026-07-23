package domain

import "time"

// MinScheduledMessageIntervalSeconds is the smallest interval a streamer can
// configure - prevents an accidental near-zero interval from spamming chat.
const MinScheduledMessageIntervalSeconds = 60

type ScheduledMessage struct {
	ID              int64      `json:"id"`
	ChannelID       string     `json:"channel_id"`
	Message         string     `json:"message"`
	IntervalSeconds int        `json:"interval_seconds"`
	Enabled         bool       `json:"enabled"`
	OnlyWhenLive    bool       `json:"only_when_live"`
	NextSendAt      time.Time  `json:"next_send_at"`
	LastSentAt      *time.Time `json:"last_sent_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func NewScheduledMessage(channelID, message string, intervalSeconds int) *ScheduledMessage {
	now := time.Now()
	return &ScheduledMessage{
		ChannelID:       channelID,
		Message:         message,
		IntervalSeconds: intervalSeconds,
		Enabled:         true,
		OnlyWhenLive:    true,
		NextSendAt:      now.Add(time.Duration(intervalSeconds) * time.Second),
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

func (m *ScheduledMessage) Update(message *string, intervalSeconds *int, enabled, onlyWhenLive *bool) {
	if message != nil && *message != "" {
		m.Message = *message
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

func ValidateScheduledMessage(message string, intervalSeconds int) error {
	if message == "" {
		return ErrEmptyMessage
	}
	if intervalSeconds < MinScheduledMessageIntervalSeconds {
		return ErrIntervalTooShort
	}
	return nil
}
