package domain

import "time"

type DiscordGuildSettings struct {
	ID                         int64     `json:"id" db:"id"`
	GuildID                    int64     `json:"guildId,string" db:"guild_id"`
	UserID                     string    `json:"userId" db:"user_id"`

	NotificationChannelID      *int64    `json:"notificationChannelId,string" db:"notification_channel_id"`
	CommandChannelID           *int64    `json:"commandChannelId,string" db:"command_channel_id"`
	ActivityChannelID          *int64    `json:"activityChannelId,string" db:"activity_channel_id"`

	TwitchNotificationsEnabled bool      `json:"twitchNotificationsEnabled" db:"twitch_notifications_enabled"`
	JoinToCreateEnabled        bool      `json:"joinToCreateEnabled" db:"join_to_create_enabled"`

	AdminRoleID                *int64    `json:"adminRoleId,string" db:"admin_role_id"`

	CreatedAt                  time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt                  time.Time `json:"updatedAt" db:"updated_at"`
}

type DiscordGuildSettingsCreateInput struct {
	GuildID                    int64    `json:"guildId,string"`
	UserID                     string   `json:"userId"`
	NotificationChannelID      *int64   `json:"notificationChannelId,string"`
	CommandChannelID           *int64   `json:"commandChannelId,string"`
	ActivityChannelID          *int64   `json:"activityChannelId,string"`
	TwitchNotificationsEnabled bool     `json:"twitchNotificationsEnabled"`
	JoinToCreateEnabled        bool     `json:"joinToCreateEnabled"`
	AdminRoleID                *int64   `json:"adminRoleId,string"`
}

type DiscordGuildSettingsUpdateInput struct {
	NotificationChannelID      *int64 `json:"notificationChannelId,string"`
	CommandChannelID           *int64 `json:"commandChannelId,string"`
	ActivityChannelID          *int64 `json:"activityChannelId,string"`
	TwitchNotificationsEnabled *bool  `json:"twitchNotificationsEnabled"`
	JoinToCreateEnabled        *bool  `json:"joinToCreateEnabled"`
	AdminRoleID                *int64 `json:"adminRoleId,string"`
}

type DiscordGuildWithSettings struct {
	Guild    DiscordGuild         `json:"guild"`
	Settings *DiscordGuildSettings `json:"settings,omitempty"`
}