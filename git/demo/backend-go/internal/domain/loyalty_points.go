package domain

import "time"

// LoyaltyPoints is one viewer's running points balance for one channel.
type LoyaltyPoints struct {
	UserTwitchID   string    `json:"user_twitch_id"`
	ViewerTwitchID string    `json:"viewer_twitch_id"`
	ViewerLogin    string    `json:"viewer_login"`
	Points         int64     `json:"points"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ViewerCredit is one chatter to credit points to during an accrual tick.
type ViewerCredit struct {
	ViewerTwitchID string
	ViewerLogin    string
	Amount         int64
}
