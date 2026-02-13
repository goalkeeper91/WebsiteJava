package domain

import (
	"time"
)

type User struct {
	TwitchID  string     `json:"twitch_id"`
	Username  string     `json:"username"`
	Email     string     `json:"email,omitempty"`
	DiscordID *string    `json:"discord_id,omitempty"`
	IsAdmin   bool       `json:"is_admin"`
	IsBot     bool       `json:"is_bot"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func NewUser(twitchID, username, email string) *User {
	now := time.Now()
	return &User{
		TwitchID:  twitchID,
		Username:  username,
		Email:     email,
		IsAdmin:   false,
		IsBot:     false,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func NewBotUser(twitchID, username, email string) *User {
	now := time.Now()
	return &User{
		TwitchID:  twitchID,
		Username:  username,
		Email:     email,
		IsAdmin:   false,
		IsBot:     true,   // Bot Account!
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (u *User) Update(username, email string) {
	u.Username = username
	u.Email = email
	u.UpdatedAt = time.Now()
}

func (u *User) SetDiscordID(discordID string) {
	u.DiscordID = &discordID
	u.UpdatedAt = time.Now()
}

func (u *User) PromoteToAdmin() {
	u.IsAdmin = true
	u.UpdatedAt = time.Now()
}

func (u *User) MarkAsBot() {
	u.IsBot = true
	u.UpdatedAt = time.Now()
}

func (u *User) IsUserBot() bool {
	return u.IsBot
}