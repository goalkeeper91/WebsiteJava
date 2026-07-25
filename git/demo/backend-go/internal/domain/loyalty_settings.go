package domain

import "time"

// MinLoyaltyIntervalMinutes is the smallest accrual interval a streamer can
// configure - mirrors MinScheduledMessageIntervalSeconds's role of preventing
// an accidental near-zero interval from hammering the Twitch API every tick.
const MinLoyaltyIntervalMinutes = 1

// LoyaltySettings is a 1:1-per-channel configuration for the Loyalty/
// Watchtime points system. PointsPerInterval viewers currently present in
// chat (via Helix "Get Chatters", not just those who type - watchtime should
// reward lurkers too) are credited every IntervalMinutes while the channel is
// live. NextAccrualAt drives the scheduler's due-query, same role as
// ScheduledMessage.NextSendAt. Off by default.
type LoyaltySettings struct {
	UserTwitchID      string    `json:"user_twitch_id"`
	Enabled           bool      `json:"enabled"`
	PointsName        string    `json:"points_name"`
	PointsPerInterval int       `json:"points_per_interval"`
	IntervalMinutes   int       `json:"interval_minutes"`
	NextAccrualAt     time.Time `json:"next_accrual_at"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func NewLoyaltySettings(userTwitchID string) *LoyaltySettings {
	now := time.Now()
	return &LoyaltySettings{
		UserTwitchID:      userTwitchID,
		Enabled:           false,
		PointsName:        "Punkte",
		PointsPerInterval: 1,
		IntervalMinutes:   5,
		NextAccrualAt:     now,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

// LoyaltySettingsUpdateInput holds the selectively-updatable fields (nil =
// leave unchanged).
type LoyaltySettingsUpdateInput struct {
	Enabled           *bool   `json:"enabled,omitempty"`
	PointsName        *string `json:"points_name,omitempty"`
	PointsPerInterval *int    `json:"points_per_interval,omitempty"`
	IntervalMinutes   *int    `json:"interval_minutes,omitempty"`
}

func (s *LoyaltySettings) ApplyUpdate(input LoyaltySettingsUpdateInput) {
	if input.Enabled != nil {
		s.Enabled = *input.Enabled
	}
	if input.PointsName != nil && *input.PointsName != "" {
		s.PointsName = *input.PointsName
	}
	if input.PointsPerInterval != nil && *input.PointsPerInterval >= 1 {
		s.PointsPerInterval = *input.PointsPerInterval
	}
	if input.IntervalMinutes != nil && *input.IntervalMinutes >= MinLoyaltyIntervalMinutes {
		s.IntervalMinutes = *input.IntervalMinutes
	}
	s.UpdatedAt = time.Now()
}

// MarkAccrued advances NextAccrualAt by the current interval - called after
// every tick regardless of whether the channel was live, so an offline
// channel doesn't "catch up" on points once it goes live again.
func (s *LoyaltySettings) MarkAccrued() {
	s.NextAccrualAt = time.Now().Add(time.Duration(s.IntervalMinutes) * time.Minute)
}
