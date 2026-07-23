package domain

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// SubathonEventLogEntry is one line in the timer's rolling event log
// (a sub, a bit cheer) shown in the dashboard's "System Log" panel.
type SubathonEventLogEntry struct {
	Time string `json:"time"`
	Text string `json:"text"`
}

// SubathonEventLog is a JSONB-backed array, capped to the most recent 50
// entries by the service layer (matches the original Next.js implementation).
type SubathonEventLog []SubathonEventLogEntry

func (l SubathonEventLog) Value() (driver.Value, error) {
	if l == nil {
		return json.Marshal([]SubathonEventLogEntry{})
	}
	return json.Marshal(l)
}

func (l *SubathonEventLog) Scan(value interface{}) error {
	if value == nil {
		*l = []SubathonEventLogEntry{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}
	return json.Unmarshal(bytes, l)
}

// SubathonState is the per-user, DB-backed state for the Subathon Timer —
// replaces the original Next.js app's single in-memory global object, since
// this dashboard is multi-tenant.
type SubathonState struct {
	UserID             string           `json:"userId" db:"user_id"`
	TimeRemaining      int              `json:"timeRemaining" db:"time_remaining"`
	IsRunning          bool             `json:"isRunning" db:"is_running"`
	TargetTimestamp    *int64           `json:"targetTimestamp,omitempty" db:"target_timestamp"`
	TotalSubs          int              `json:"totalSubs" db:"total_subs"`
	TotalBits          int              `json:"totalBits" db:"total_bits"`
	TotalEvents        int              `json:"totalEvents" db:"total_events"`
	InitialTimeMinutes int              `json:"initialTime" db:"initial_time_minutes"`
	SubTimeSeconds     int              `json:"subTime" db:"sub_time_seconds"`
	BitsTimeSeconds    int              `json:"bitsTime" db:"bits_time_seconds"`
	EventLog           SubathonEventLog `json:"eventLog" db:"event_log"`
	CreatedAt          time.Time        `json:"-" db:"created_at"`
	UpdatedAt          time.Time        `json:"-" db:"updated_at"`
}

// SubathonFailedEvent is a sub/cheer that failed to persist (e.g. a
// transient DB error) - kept visible and retried automatically instead of
// silently vanishing, since a lost event during a subathon means real
// support the streamer never gets credit for.
type SubathonFailedEvent struct {
	ID           int64     `json:"id" db:"id"`
	UserID       string    `json:"userId" db:"user_id"`
	MessageID    string    `json:"messageId" db:"message_id"`
	EventType    string    `json:"eventType" db:"event_type"`
	RawPayload   []byte    `json:"-" db:"raw_payload"`
	ErrorMessage string    `json:"errorMessage" db:"error_message"`
	RetryCount   int       `json:"retryCount" db:"retry_count"`
	Resolved     bool      `json:"resolved" db:"resolved"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}

// SubathonStateUpdateInput carries a partial update from the dashboard
// (start/pause toggles isRunning, settings form updates the *TimeSeconds
// fields, reset sets everything back to defaults).
type SubathonStateUpdateInput struct {
	IsRunning *bool
	// TargetTimestamp is always written as-is (nil -> SQL NULL), never
	// COALESCE'd - the service layer recomputes it fully on every
	// start/pause/reset, since it's the one field that must be explicitly
	// clearable (paused/reset state has no target timestamp at all).
	TargetTimestamp    *int64
	TimeRemaining      *int
	TotalSubs          *int
	TotalBits          *int
	TotalEvents        *int
	InitialTimeMinutes *int
	SubTimeSeconds     *int
	BitsTimeSeconds    *int
	EventLog           SubathonEventLog
}
