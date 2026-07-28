package domain

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// AITone ist eine Enum für AI Tonalität
type AITone string

const (
	ToneNeutral      AITone = "neutral"
	ToneCasual       AITone = "casual"
	ToneProfessional AITone = "professional"
	ToneEnergetic    AITone = "energetic"
	ToneFunny        AITone = "funny"
)

// AIStyle ist eine Enum für AI Stil
type AIStyle string

const (
	StyleStandard     AIStyle = "standard"
	StyleShort        AIStyle = "short"
	StyleDetailed     AIStyle = "detailed"
	StyleViral        AIStyle = "viral"
	StyleStorytelling AIStyle = "storytelling"
)

// PostingFrequency ist eine Enum für Posting Häufigkeit
type PostingFrequency string

const (
	FrequencyRealtime PostingFrequency = "realtime"
	FrequencyHourly   PostingFrequency = "hourly"
	FrequencyDaily    PostingFrequency = "daily"
	FrequencyWeekly   PostingFrequency = "weekly"
	FrequencyManual   PostingFrequency = "manual"
)

// StringArray ist ein Helper für JSONB String Arrays
type StringArray []string

// Value implementiert driver.Valuer für database/sql
func (a StringArray) Value() (driver.Value, error) {
	if a == nil {
		return json.Marshal([]string{})
	}
	return json.Marshal(a)
}

// Scan implementiert sql.Scanner für database/sql
func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = []string{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	return json.Unmarshal(bytes, a)
}

// PlatformArray ist ein Helper für JSONB Platform Arrays
type PlatformArray []Platform

// Value implementiert driver.Valuer
func (a PlatformArray) Value() (driver.Value, error) {
	if a == nil {
		return json.Marshal([]Platform{})
	}
	return json.Marshal(a)
}

// Scan implementiert sql.Scanner
func (a *PlatformArray) Scan(value interface{}) error {
	if value == nil {
		*a = []Platform{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	return json.Unmarshal(bytes, a)
}

// AutomationSettings repräsentiert User-spezifische Automation Konfiguration
type AutomationSettings struct {
	UserTwitchID           string           `json:"user_twitch_id" db:"user_twitch_id"`
	IsEnabled              bool             `json:"is_enabled" db:"is_enabled"`
	AITone                 AITone           `json:"ai_tone" db:"ai_tone"`
	AIStyle                AIStyle          `json:"ai_style" db:"ai_style"`
	UseHashtags            bool             `json:"use_hashtags" db:"use_hashtags"`
	MaxHashtags            int              `json:"max_hashtags" db:"max_hashtags"`
	MinClipDuration        int              `json:"min_clip_duration" db:"min_clip_duration"`
	MaxClipDuration        int              `json:"max_clip_duration" db:"max_clip_duration"`
	MinViewCount           int              `json:"min_view_count" db:"min_view_count"`
	OnlyVerifiedClips      bool             `json:"only_verified_clips" db:"only_verified_clips"`
	TargetPlatforms        PlatformArray    `json:"target_platforms" db:"target_platforms"`
	AutoPostEnabled        bool             `json:"auto_post_enabled" db:"auto_post_enabled"`
	PostingFrequency       PostingFrequency `json:"posting_frequency" db:"posting_frequency"`
	PreferredPostingTimes  StringArray      `json:"preferred_posting_times,omitempty" db:"preferred_posting_times"`
	BlockedWords           StringArray      `json:"blocked_words,omitempty" db:"blocked_words"`
	AllowedGames           StringArray      `json:"allowed_games,omitempty" db:"allowed_games"`
	NotifyOnPost           bool             `json:"notify_on_post" db:"notify_on_post"`
	NotifyOnError          bool             `json:"notify_on_error" db:"notify_on_error"`
	CreatedAt              time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time        `json:"updated_at" db:"updated_at"`
}

// NewAutomationSettings erstellt Default Settings für einen User
func NewAutomationSettings(userTwitchID string) *AutomationSettings {
	now := time.Now()
	return &AutomationSettings{
		UserTwitchID:          userTwitchID,
		IsEnabled:             false,
		AITone:                ToneNeutral,
		AIStyle:               StyleStandard,
		UseHashtags:           true,
		MaxHashtags:           5,
		MinClipDuration:       10,
		MaxClipDuration:       60,
		MinViewCount:          0,
		OnlyVerifiedClips:     false,
		TargetPlatforms:       PlatformArray{},
		AutoPostEnabled:       false,
		PostingFrequency:      FrequencyDaily,
		PreferredPostingTimes: StringArray{},
		BlockedWords:          StringArray{},
		AllowedGames:          StringArray{},
		NotifyOnPost:          true,
		NotifyOnError:         true,
		CreatedAt:             now,
		UpdatedAt:             now,
	}
}

// Update aktualisiert die Settings
func (as *AutomationSettings) Update(updates map[string]interface{}) {
	// Hier könntest du selektive Updates implementieren
	as.UpdatedAt = time.Now()
}

// AutomationSettingsUpdateInput enthält die selektiv änderbaren Felder der
// Automation Settings (nil = unverändert lassen).
type AutomationSettingsUpdateInput struct {
	IsEnabled         *bool             `json:"is_enabled,omitempty"`
	AITone            *AITone           `json:"ai_tone,omitempty"`
	AIStyle           *AIStyle          `json:"ai_style,omitempty"`
	UseHashtags       *bool             `json:"use_hashtags,omitempty"`
	MaxHashtags       *int              `json:"max_hashtags,omitempty"`
	MinClipDuration   *int              `json:"min_clip_duration,omitempty"`
	MaxClipDuration   *int              `json:"max_clip_duration,omitempty"`
	MinViewCount      *int              `json:"min_view_count,omitempty"`
	OnlyVerifiedClips *bool             `json:"only_verified_clips,omitempty"`
	TargetPlatforms   *PlatformArray    `json:"target_platforms,omitempty"`
	AutoPostEnabled   *bool             `json:"auto_post_enabled,omitempty"`
	PostingFrequency  *PostingFrequency `json:"posting_frequency,omitempty"`
	BlockedWords      *StringArray      `json:"blocked_words,omitempty"`
	AllowedGames      *StringArray      `json:"allowed_games,omitempty"`
	NotifyOnPost      *bool             `json:"notify_on_post,omitempty"`
	NotifyOnError     *bool             `json:"notify_on_error,omitempty"`
}

// ApplyUpdate übernimmt alle gesetzten Felder aus dem Input in die Settings.
func (as *AutomationSettings) ApplyUpdate(input AutomationSettingsUpdateInput) {
	if input.IsEnabled != nil {
		as.IsEnabled = *input.IsEnabled
	}
	if input.AITone != nil {
		as.AITone = *input.AITone
	}
	if input.AIStyle != nil {
		as.AIStyle = *input.AIStyle
	}
	if input.UseHashtags != nil {
		as.UseHashtags = *input.UseHashtags
	}
	if input.MaxHashtags != nil {
		as.MaxHashtags = *input.MaxHashtags
	}
	if input.MinClipDuration != nil {
		as.MinClipDuration = *input.MinClipDuration
	}
	if input.MaxClipDuration != nil {
		as.MaxClipDuration = *input.MaxClipDuration
	}
	if input.MinViewCount != nil {
		as.MinViewCount = *input.MinViewCount
	}
	if input.OnlyVerifiedClips != nil {
		as.OnlyVerifiedClips = *input.OnlyVerifiedClips
	}
	if input.TargetPlatforms != nil {
		as.TargetPlatforms = *input.TargetPlatforms
	}
	if input.AutoPostEnabled != nil {
		as.AutoPostEnabled = *input.AutoPostEnabled
	}
	if input.PostingFrequency != nil {
		as.PostingFrequency = *input.PostingFrequency
	}
	if input.BlockedWords != nil {
		as.BlockedWords = *input.BlockedWords
	}
	if input.AllowedGames != nil {
		as.AllowedGames = *input.AllowedGames
	}
	if input.NotifyOnPost != nil {
		as.NotifyOnPost = *input.NotifyOnPost
	}
	if input.NotifyOnError != nil {
		as.NotifyOnError = *input.NotifyOnError
	}
	as.UpdatedAt = time.Now()
}

// Enable aktiviert die Automation
func (as *AutomationSettings) Enable() {
	as.IsEnabled = true
	as.UpdatedAt = time.Now()
}

// Disable deaktiviert die Automation
func (as *AutomationSettings) Disable() {
	as.IsEnabled = false
	as.UpdatedAt = time.Now()
}

// Validate validiert die Settings
func (as *AutomationSettings) Validate() error {
	if as.MinClipDuration > as.MaxClipDuration {
		return ErrInvalidClipDurationRange
	}

	if as.MaxHashtags < 0 || as.MaxHashtags > 30 {
		return ErrInvalidHashtagCount
	}

	if as.MinClipDuration < 5 {
		return ErrClipDurationTooShort
	}

	if as.MaxClipDuration > 300 {
		return ErrClipDurationTooLong
	}

	return nil
}

// HasPlatform prüft ob eine Platform in den Targets ist
func (as *AutomationSettings) HasPlatform(platform Platform) bool {
	for _, p := range as.TargetPlatforms {
		if p == platform {
			return true
		}
	}
	return false
}

// AutomationSettingsWithSubscription bündelt Settings + Subscription (inkl.
// Tier) + Admin-Status für einen Kandidaten der Clip-Detector-Poll-Schleife -
// reine Datenbündelung, die eigentliche Berechtigungsprüfung passiert beim
// Aufrufer über die bereits vorhandenen UserSubscription.IsActive() /
// SubscriptionTier.HasFeature()-Methoden, nicht hier.
type AutomationSettingsWithSubscription struct {
	Settings     *AutomationSettings
	Subscription *UserSubscription
	IsAdmin      bool
}
