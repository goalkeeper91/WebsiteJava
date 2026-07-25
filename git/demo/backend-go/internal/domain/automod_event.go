package domain

import "time"

// AutomodEvent is one append-only log row recording an automod violation -
// separate from AutomodViolation, which just tracks the live escalation
// counter (one row per offender, overwritten). This is what the dashboard's
// "recent violations" view reads.
type AutomodEvent struct {
	ID                int64     `json:"id"`
	UserTwitchID      string    `json:"user_twitch_id"`
	OffenderTwitchID  string    `json:"offender_twitch_id"`
	OffenderName      string    `json:"offender_name"`
	Reason            string    `json:"reason"`
	MessageExcerpt    string    `json:"message_excerpt,omitempty"`
	TimeoutSeconds    int       `json:"timeout_seconds"`
	CreatedAt         time.Time `json:"created_at"`
}

func NewAutomodEvent(userTwitchID, offenderTwitchID, offenderName, reason, messageExcerpt string, timeoutSeconds int) *AutomodEvent {
	return &AutomodEvent{
		UserTwitchID:     userTwitchID,
		OffenderTwitchID: offenderTwitchID,
		OffenderName:     offenderName,
		Reason:           reason,
		MessageExcerpt:   messageExcerpt,
		TimeoutSeconds:   timeoutSeconds,
		CreatedAt:        time.Now(),
	}
}
