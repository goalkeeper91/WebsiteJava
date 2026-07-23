package domain

type DiscordChannel struct {
	ID       int64  `json:"id,string"`
	Name     string `json:"name"`
	Type     string `json:"type"`     // text, voice, category, etc.
	Position int    `json:"position"`
}

type DiscordRole struct {
	ID          int64  `json:"id,string"`
	Name        string `json:"name"`
	Color       int    `json:"color"`
	Position    int    `json:"position"`
	Permissions int64  `json:"permissions"`
}

type DiscordGuildDetails struct {
	Guild    DiscordGuild      `json:"guild"`
	Channels []DiscordChannel  `json:"channels"`
	Roles    []DiscordRole     `json:"roles"`
	Settings *DiscordGuildSettings `json:"settings,omitempty"`
}