package domain

import "time"

type DiscordGuildSettings struct {
	ID                         int64     `json:"id" db:"id"`
	GuildID                    int64     `json:"guildId" db:"guild_id"`
	UserID                     string    `json:"userId" db:"user_id"`

	NotificationChannelID      *int64    `json:"notificationChannelId" db:"notification_channel_id"`
	CommandChannelID           *int64    `json:"commandChannelId" db:"command_channel_id"`
	ActivityChannelID          *int64    `json:"activityChannelId" db:"activity_channel_id"`

	TwitchNotificationsEnabled bool      `json:"twitchNotificationsEnabled" db:"twitch_notifications_enabled"`
	JoinToCreateEnabled        bool      `json:"joinToCreateEnabled" db:"join_to_create_enabled"`

	AdminRoleID                *int64    `json:"adminRoleId" db:"admin_role_id"`

	CreatedAt                  time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt                  time.Time `json:"updatedAt" db:"updated_at"`
}

type DiscordGuildSettingsCreateInput struct {
	GuildID                    int64
	UserID                     string
	NotificationChannelID      *int64
	CommandChannelID           *int64
	ActivityChannelID          *int64
	TwitchNotificationsEnabled bool
	JoinToCreateEnabled        bool
	AdminRoleID                *int64
}

type DiscordGuildSettingsUpdateInput struct {
	NotificationChannelID      *int64
	CommandChannelID           *int64
	ActivityChannelID          *int64
	TwitchNotificationsEnabled *bool
	JoinToCreateEnabled        *bool
	AdminRoleID                *int64
}

type DiscordGuildWithSettings struct {
	Guild    DiscordGuild         `json:"guild"`
	Settings *DiscordGuildSettings `json:"settings,omitempty"`
}