package domain

import "time"

type JoinToCreateConfig struct {
	ID                int64     `json:"id" db:"id"`
	GuildID           int64     `json:"guildId,string" db:"guild_id"`
	UserID            string    `json:"userId" db:"user_id"`
	JoinChannelID     int64     `json:"joinChannelId,string" db:"join_channel_id"`
	CategoryID        int64     `json:"categoryId,string" db:"category_id"`
	ChannelNamePrefix string    `json:"channelNamePrefix" db:"channel_name_prefix"`
	UserLimit         int       `json:"userLimit" db:"user_limit"`
	PrivateChannel    bool      `json:"privateChannel" db:"private_channel"`
	Enabled           bool      `json:"enabled" db:"enabled"`
	CreatedAt         time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt         time.Time `json:"updatedAt" db:"updated_at"`
}

type JoinToCreateConfigCreateInput struct {
	GuildID           int64  `json:"guildId,string"`
	UserID            string `json:"userId"`
	JoinChannelID     int64  `json:"joinChannelId,string"`
	CategoryID        int64  `json:"categoryId,string"`
	ChannelNamePrefix string `json:"channelNamePrefix"`
	UserLimit         int    `json:"userLimit"`
	PrivateChannel    bool   `json:"privateChannel"`
	Enabled           bool   `json:"enabled"`
}

type JoinToCreateConfigUpdateInput struct {
	JoinChannelID     *int64  `json:"joinChannelId,string"`
	CategoryID        *int64  `json:"categoryId,string"`
	ChannelNamePrefix *string `json:"channelNamePrefix"`
	UserLimit         *int    `json:"userLimit"`
	PrivateChannel    *bool   `json:"privateChannel"`
	Enabled           *bool   `json:"enabled"`
}