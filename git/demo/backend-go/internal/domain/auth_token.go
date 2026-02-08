package domain

import (
	"time"
)

type AuthToken struct {
	ID            int64     `json:"id"`
	TwitchUserID  string    `json:"twitch_user_id"`
	Username      string    `json:"username"`
	AccessToken   string    `json:"-"` // Nie in JSON serialisieren
	RefreshToken  string    `json:"-"` // Nie in JSON serialisieren
	TokenType     string    `json:"token_type"`
	ExpiresIn     int64     `json:"expires_in"`
	Scope         string    `json:"scope"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func NewAuthToken(twitchUserID, username, accessToken, refreshToken, tokenType, scope string, expiresIn int64) *AuthToken {
	now := time.Now()
	return &AuthToken{
		TwitchUserID: twitchUserID,
		Username:     username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    tokenType,
		ExpiresIn:    expiresIn,
		Scope:        scope,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func (t *AuthToken) Update(accessToken, refreshToken string, expiresIn int64) {
	t.AccessToken = accessToken
	t.RefreshToken = refreshToken
	t.ExpiresIn = expiresIn
	t.UpdatedAt = time.Now()
}

func (t *AuthToken) IsExpired() bool {
	expirationTime := t.UpdatedAt.Add(time.Duration(t.ExpiresIn) * time.Second)
	return time.Now().After(expirationTime.Add(-5 * time.Minute))
}