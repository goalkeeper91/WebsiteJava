package domain

import (
	"time"
)

type TwitchChannel struct {
	ID           int64     `json:"id"`
	TwitchUserID string    `json:"twitch_user_id"`
	Username     string    `json:"username"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func NewTwitchChannel(twitchUserID, username string) *TwitchChannel {
	now := time.Now()
	return &TwitchChannel{
		TwitchUserID: twitchUserID,
		Username:     username,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func (c *TwitchChannel) Update(username string) {
	c.Username = username
	c.UpdatedAt = time.Now()
}

func (c *TwitchChannel) Activate() {
	c.IsActive = true
	c.UpdatedAt = time.Now()
}

func (c *TwitchChannel) Deactivate() {
	c.IsActive = false
	c.UpdatedAt = time.Now()
}