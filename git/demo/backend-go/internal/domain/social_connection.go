package domain

import "time"

// Platform ist eine Enum für Social Media Plattformen
type Platform string

const (
	PlatformTikTok    Platform = "tiktok"
	PlatformYouTube   Platform = "youtube"
	PlatformInstagram Platform = "instagram"

	// Link-sharing targets (v1 clip distribution): we post the Twitch clip URL,
	// no native video re-upload or OAuth needed. TikTok/YouTube/Instagram above
	// are native re-upload targets, planned for a later phase.
	PlatformDiscord Platform = "discord"
	PlatformTwitter Platform = "twitter"
)

// SocialConnection repräsentiert eine OAuth-Verbindung zu einer Social Media Plattform
type SocialConnection struct {
	ID                int64     `json:"id" db:"id"`
	UserTwitchID      string    `json:"user_twitch_id" db:"user_twitch_id"`
	Platform          Platform  `json:"platform" db:"platform"`
	PlatformUserID    string    `json:"platform_user_id" db:"platform_user_id"`
	PlatformUsername  string    `json:"platform_username,omitempty" db:"platform_username"`
	AccessToken       string    `json:"-" db:"access_token"` // Nie im JSON ausgeben!
	RefreshToken      string    `json:"-" db:"refresh_token"`
	TokenType         string    `json:"token_type" db:"token_type"`
	ExpiresAt         *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	Scope             string    `json:"scope,omitempty" db:"scope"`
	IsActive          bool      `json:"is_active" db:"is_active"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
	LastUsedAt        *time.Time `json:"last_used_at,omitempty" db:"last_used_at"`
}

// NewSocialConnection erstellt eine neue Social Connection
func NewSocialConnection(
	userTwitchID string,
	platform Platform,
	platformUserID string,
	platformUsername string,
	accessToken string,
	refreshToken string,
	expiresAt *time.Time,
) *SocialConnection {
	now := time.Now()
	return &SocialConnection{
		UserTwitchID:     userTwitchID,
		Platform:         platform,
		PlatformUserID:   platformUserID,
		PlatformUsername: platformUsername,
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		TokenType:        "Bearer",
		ExpiresAt:        expiresAt,
		IsActive:         true,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

// UpdateTokens aktualisiert die OAuth Tokens
func (sc *SocialConnection) UpdateTokens(accessToken, refreshToken string, expiresAt *time.Time) {
	sc.AccessToken = accessToken
	sc.RefreshToken = refreshToken
	sc.ExpiresAt = expiresAt
	sc.UpdatedAt = time.Now()
}

// MarkAsUsed aktualisiert LastUsedAt
func (sc *SocialConnection) MarkAsUsed() {
	now := time.Now()
	sc.LastUsedAt = &now
	sc.UpdatedAt = now
}

// Deactivate deaktiviert die Connection
func (sc *SocialConnection) Deactivate() {
	sc.IsActive = false
	sc.UpdatedAt = time.Now()
}

// Activate aktiviert die Connection
func (sc *SocialConnection) Activate() {
	sc.IsActive = true
	sc.UpdatedAt = time.Now()
}

// IsExpired prüft ob das Token abgelaufen ist
func (sc *SocialConnection) IsExpired() bool {
	if sc.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*sc.ExpiresAt)
}

// NeedsRefresh prüft ob Token bald abläuft (5 Minuten Buffer)
func (sc *SocialConnection) NeedsRefresh() bool {
	if sc.ExpiresAt == nil {
		return false
	}
	return time.Now().Add(5 * time.Minute).After(*sc.ExpiresAt)
}
