package domain

import "time"

// AutomodSettings is a 1:1-per-channel configuration (like AutomationSettings)
// for the Twitch chat automod: a blocked-words list, an optional link-filter
// allowlist, generic spam filters (caps/symbols/repetition use fixed
// thresholds; the emote-spam filter's threshold is per-channel configurable,
// since emote spam in the right moment - e.g. an esports highlight - can be
// a legitimate reaction rather than spam), and exemptions beyond
// mods/broadcaster (VIPs, a manually-curated trusted-user list, and - since
// ExemptRegulars - viewers whose Loyalty points cross the channel's Regulars
// threshold, merged into ExemptUsers at read time by
// AutomodService.GetAllEnabledSettingsForBot, never persisted here). Off by
// default - a streamer has to opt in rather than have messages start
// disappearing unannounced.
type AutomodSettings struct {
	UserTwitchID             string      `json:"user_twitch_id"`
	Enabled                  bool        `json:"enabled"`
	BlockedWords             StringArray `json:"blocked_words"`
	LinkFilterEnabled        bool        `json:"link_filter_enabled"`
	AllowedDomains           StringArray `json:"allowed_domains"`
	CapsFilterEnabled        bool        `json:"caps_filter_enabled"`
	SymbolFilterEnabled      bool        `json:"symbol_filter_enabled"`
	EmoteFilterEnabled       bool        `json:"emote_filter_enabled"`
	EmoteThreshold           int         `json:"emote_threshold"`
	RepetitionFilterEnabled  bool        `json:"repetition_filter_enabled"`
	ExemptVips               bool        `json:"exempt_vips"`
	ExemptRegulars           bool        `json:"exempt_regulars"`
	ExemptUsers              StringArray `json:"exempt_users"`
	CreatedAt                time.Time   `json:"created_at"`
	UpdatedAt                time.Time   `json:"updated_at"`
}

func NewAutomodSettings(userTwitchID string) *AutomodSettings {
	now := time.Now()
	return &AutomodSettings{
		UserTwitchID:      userTwitchID,
		Enabled:           false,
		BlockedWords:      StringArray{},
		LinkFilterEnabled: false,
		AllowedDomains:    StringArray{},
		EmoteThreshold:    6,
		ExemptVips:        true,
		ExemptRegulars:    true,
		ExemptUsers:       StringArray{},
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

// AutomodSettingsUpdateInput holds the selectively-updatable fields (nil =
// leave unchanged).
type AutomodSettingsUpdateInput struct {
	Enabled                 *bool        `json:"enabled,omitempty"`
	BlockedWords            *StringArray `json:"blocked_words,omitempty"`
	LinkFilterEnabled       *bool        `json:"link_filter_enabled,omitempty"`
	AllowedDomains          *StringArray `json:"allowed_domains,omitempty"`
	CapsFilterEnabled       *bool        `json:"caps_filter_enabled,omitempty"`
	SymbolFilterEnabled     *bool        `json:"symbol_filter_enabled,omitempty"`
	EmoteFilterEnabled      *bool        `json:"emote_filter_enabled,omitempty"`
	EmoteThreshold          *int         `json:"emote_threshold,omitempty"`
	RepetitionFilterEnabled *bool        `json:"repetition_filter_enabled,omitempty"`
	ExemptVips              *bool        `json:"exempt_vips,omitempty"`
	ExemptRegulars          *bool        `json:"exempt_regulars,omitempty"`
	ExemptUsers             *StringArray `json:"exempt_users,omitempty"`
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
	if input.CapsFilterEnabled != nil {
		s.CapsFilterEnabled = *input.CapsFilterEnabled
	}
	if input.SymbolFilterEnabled != nil {
		s.SymbolFilterEnabled = *input.SymbolFilterEnabled
	}
	if input.EmoteFilterEnabled != nil {
		s.EmoteFilterEnabled = *input.EmoteFilterEnabled
	}
	if input.EmoteThreshold != nil && *input.EmoteThreshold >= 1 {
		s.EmoteThreshold = *input.EmoteThreshold
	}
	if input.RepetitionFilterEnabled != nil {
		s.RepetitionFilterEnabled = *input.RepetitionFilterEnabled
	}
	if input.ExemptVips != nil {
		s.ExemptVips = *input.ExemptVips
	}
	if input.ExemptRegulars != nil {
		s.ExemptRegulars = *input.ExemptRegulars
	}
	if input.ExemptUsers != nil {
		s.ExemptUsers = *input.ExemptUsers
	}
	s.UpdatedAt = time.Now()
}
