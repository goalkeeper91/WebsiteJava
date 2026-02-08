package domain

import (
	"time"
)

type Session struct {
	ID        string    `json:"id"`
	TwitchID  string    `json:"twitch_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

func NewSession(id, twitchID string, duration time.Duration) *Session {
	now := time.Now()
	return &Session{
		ID:        id,
		TwitchID:  twitchID,
		ExpiresAt: now.Add(duration),
		CreatedAt: now,
	}
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

func (s *Session) IsValid() bool {
	return !s.IsExpired()
}