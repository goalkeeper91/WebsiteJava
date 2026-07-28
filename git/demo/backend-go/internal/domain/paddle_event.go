package domain

import "time"

// PaddleEvent is a dedup log entry for processed Paddle webhook events -
// Paddle does not guarantee exactly-once delivery, so each event_id is
// recorded once before its effects are applied (same shape as
// stream_activities.event_message_id's unique-index dedup pattern).
type PaddleEvent struct {
	ID         int64     `json:"id" db:"id"`
	EventID    string    `json:"eventId" db:"event_id"`
	EventType  string    `json:"eventType" db:"event_type"`
	ReceivedAt time.Time `json:"receivedAt" db:"received_at"`
}
