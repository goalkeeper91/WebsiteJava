package domain

import "time"

type DiscordConnection struct {
	ID                   int64     `json:"id" db:"id"`
	UserID               int64     `json:"userId" db:"user_id"`
	DiscordUserID        int64     `json:"discordUserId" db:"discord_user_id"`
	DiscordUsername      string    `json:"discordUsername" db:"discord_username"`
	DiscordDiscriminator string    `json:"discordDiscriminator" db:"discord_discriminator"`
	AccessToken          string    `json:"-" db:"access_token"`           // Encrypted, not in JSON
	RefreshToken         string    `json:"-" db:"refresh_token"`          // Encrypted, not in JSON
	TokenExpiresAt       time.Time `json:"tokenExpiresAt" db:"token_expires_at"`
	CreatedAt            time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt            time.Time `json:"updatedAt" db:"updated_at"`
}

type DiscordConnectionCreateInput struct {
	UserID               int64
	DiscordUserID        int64
	DiscordUsername      string
	DiscordDiscriminator string
	AccessToken          string
	RefreshToken         string
	TokenExpiresAt       time.Time
}

type DiscordConnectionUpdateInput struct {
	DiscordUsername      *string
	DiscordDiscriminator *string
	AccessToken          *string
	RefreshToken         *string
	TokenExpiresAt       *time.Time
}