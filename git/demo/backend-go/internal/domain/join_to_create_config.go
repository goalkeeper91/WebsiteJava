package domain

import "time"

type JoinToCreateConfig struct {
	ID                int64     `json:"id" db:"id"`
	GuildID           int64     `json:"guildId" db:"guild_id"`
	UserID            string    `json:"userId" db:"user_id"`
	JoinChannelID     int64     `json:"joinChannelId" db:"join_channel_id"`
	CategoryID        int64     `json:"categoryId" db:"category_id"`
	ChannelNamePrefix string    `json:"channelNamePrefix" db:"channel_name_prefix"`
	UserLimit         int       `json:"userLimit" db:"user_limit"`
	PrivateChannel    bool      `json:"privateChannel" db:"private_channel"`
	Enabled           bool      `json:"enabled" db:"enabled"`
	CreatedAt         time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt         time.Time `json:"updatedAt" db:"updated_at"`
}

type JoinToCreateConfigCreateInput struct {
	GuildID           int64
	UserID            string
	JoinChannelID     int64
	CategoryID        int64
	ChannelNamePrefix string
	UserLimit         int
	PrivateChannel    bool
	Enabled           bool
}

type JoinToCreateConfigUpdateInput struct {
	JoinChannelID     *int64
	CategoryID        *int64
	ChannelNamePrefix *string
	UserLimit         *int
	PrivateChannel    *bool
	Enabled           *bool
}