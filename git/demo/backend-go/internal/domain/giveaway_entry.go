package domain

import "time"

// GiveawayEntry is one viewer's participation in a Giveaway - Entries is
// their ticket count in the weighted draw (2 if SubBonus applied and they
// were a subscriber at entry time, 1 otherwise).
type GiveawayEntry struct {
	ID             int64     `json:"id"`
	GiveawayID     int64     `json:"giveaway_id"`
	ViewerTwitchID string    `json:"viewer_twitch_id"`
	ViewerLogin    string    `json:"viewer_login"`
	Entries        int       `json:"entries"`
	EnteredAt      time.Time `json:"entered_at"`
}
