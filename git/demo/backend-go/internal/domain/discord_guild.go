package domain

import "time"

type DiscordGuild struct {
	ID            int64      `json:"id" db:"id"`                       // Discord Guild ID (Snowflake)
	OwnerUserID   *int64     `json:"ownerUserId" db:"owner_user_id"`   // User who added the bot
	Name          string     `json:"name" db:"name"`
	IconURL       *string    `json:"iconUrl" db:"icon_url"`
	MemberCount   *int       `json:"memberCount" db:"member_count"`
	BotAddedAt    *time.Time `json:"botAddedAt" db:"bot_added_at"`
	BotRemovedAt  *time.Time `json:"botRemovedAt" db:"bot_removed_at"`
	IsActive      bool       `json:"isActive" db:"is_active"`
	CreatedAt     time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time  `json:"updatedAt" db:"updated_at"`
}

type DiscordGuildCreateInput struct {
	ID          int64
	OwnerUserID *int64
	Name        string
	IconURL     *string
	MemberCount *int
}

type DiscordGuildUpdateInput struct {
	Name         *string
	IconURL      *string
	MemberCount  *int
	BotAddedAt   *time.Time
	BotRemovedAt *time.Time
	IsActive     *bool
}