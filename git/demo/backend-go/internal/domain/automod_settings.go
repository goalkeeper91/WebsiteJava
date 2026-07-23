package domain

import "time"

// AutomodSettings is a 1:1-per-channel configuration (like AutomationSettings)
// for the Twitch chat automod: a blocked-words list and an optional
// link-filter allowlist. Off by default - a streamer has to opt in rather
// than have messages start disappearing unannounced.
type AutomodSettings struct {
	UserTwitchID       string      `json:"user_twitch_id"`
	Enabled            bool        `json:"enabled"`
	BlockedWords       StringArray `json:"blocked_words"`
	LinkFilterEnabled  bool        `json:"link_filter_enabled"`
	AllowedDomains     StringArray `json:"allowed_domains"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`
}

func NewAutomodSettings(userTwitchID string) *AutomodSettings {
	now := time.Now()
	return &AutomodSettings{
		UserTwitchID:      userTwitchID,
		Enabled:           false,
		BlockedWords:      StringArray{},
		LinkFilterEnabled: false,
		AllowedDomains:    StringArray{},
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

// AutomodSettingsUpdateInput holds the selectively-updatable fields (nil =
// leave unchanged).
type AutomodSettingsUpdateInput struct {
	Enabled           *bool        `json:"enabled,omitempty"`
	BlockedWords      *StringArray `json:"blocked_words,omitempty"`
	LinkFilterEnabled *bool        `json:"link_filter_enabled,omitempty"`
	AllowedDomains    *StringArray `json:"allowed_domains,omitempty"`
}

func (s *AutomodSettings) ApplyUpdate(input AutomodSettingsUpdateInput) {
	if input.Enabled != nil {
		s.Enabled = *input.Enabled
	}
	if input.BlockedWords != nil {
		s.BlockedWords = *input.BlockedWords
	}
	if input.LinkFilterEnabled != nil {
		s.LinkFilterEnabled = *input.LinkFilterEnabled
	}
	if input.AllowedDomains != nil {
		s.AllowedDomains = *input.AllowedDomains
	}
	s.UpdatedAt = time.Now()
}
