package domain

import "time"

// Giveaway is one raffle round for a channel. Only one "open" giveaway per
// channel is allowed at a time - enforced by GiveawayService, not a DB
// constraint, since this is an episodic record table (like AutomodEvent),
// not a singleton settings row.
type Giveaway struct {
	ID             int64      `json:"id"`
	UserTwitchID   string     `json:"user_twitch_id"`
	Status         string     `json:"status"`
	SubBonus       bool       `json:"sub_bonus"`
	WinnerTwitchID *string    `json:"winner_twitch_id,omitempty"`
	WinnerLogin    *string    `json:"winner_login,omitempty"`
	StartedAt      time.Time  `json:"started_at"`
	EndedAt        *time.Time `json:"ended_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	// EntryCount is only populated by GiveawayRepository.GetHistory (a
	// display-only aggregate, not a stored column) - zero elsewhere.
	EntryCount int `json:"entry_count"`
}

const (
	GiveawayStatusOpen   = "open"
	GiveawayStatusClosed = "closed"
)

func NewGiveaway(userTwitchID string, subBonus bool) *Giveaway {
	now := time.Now()
	return &Giveaway{
		UserTwitchID: userTwitchID,
		Status:       GiveawayStatusOpen,
		SubBonus:     subBonus,
		StartedAt:    now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
