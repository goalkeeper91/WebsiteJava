package domain

import "time"

// AutomodEscalationTiers are the fixed timeout durations (seconds) applied
// on the 1st, 2nd, 3rd, ... violation within a "streak" (see
// AutomodViolationStreakTTL). Not configurable in v1.
var AutomodEscalationTiers = []int{10, 60, 600, 86400}

// AutomodViolationStreakTTL: if a viewer's last violation is older than
// this, their streak has "cooled off" and the next violation restarts at
// tier 1 instead of continuing to escalate.
const AutomodViolationStreakTTL = 24 * time.Hour

// NextTimeoutSeconds returns the timeout duration for the Nth violation
// (1-indexed) in an active streak, clamped to the last configured tier.
func NextTimeoutSeconds(violationCount int) int {
	if violationCount < 1 {
		violationCount = 1
	}
	index := violationCount - 1
	if index >= len(AutomodEscalationTiers) {
		index = len(AutomodEscalationTiers) - 1
	}
	return AutomodEscalationTiers[index]
}

type AutomodViolation struct {
	UserTwitchID     string     `json:"user_twitch_id"`
	OffenderTwitchID string     `json:"offender_twitch_id"`
	ViolationCount   int        `json:"violation_count"`
	LastViolationAt  *time.Time `json:"last_violation_at,omitempty"`
}
