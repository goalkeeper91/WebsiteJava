package domain

import (
	"strings"
	"time"
)

// Giveaway is one raffle round for a channel. Only one "open" giveaway per
// channel is allowed at a time - enforced by GiveawayService, not a DB
// constraint, since this is an episodic record table (like AutomodEvent),
// not a singleton settings row.
//
// Entry is via a streamer-chosen Keyword typed plainly in chat (no ! prefix -
// ! commands are reserved for direct bot addressing: start/draw/status).
// Keyword is always stored lowercased so matching a chat message is a plain
// string comparison, case-insensitive by construction.
type Giveaway struct {
	ID             int64      `json:"id"`
	UserTwitchID   string     `json:"user_twitch_id"`
	Status         string     `json:"status"`
	Keyword        string     `json:"keyword"`
	SubBonus       bool       `json:"sub_bonus"`
	WinnerTwitchID *string    `json:"winner_twitch_id,omitempty"`
	WinnerLogin    *string    `json:"winner_login,omitempty"`
	StartedAt      time.Time  `json:"started_at"`
	EndedAt        *time.Time `json:"ended_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	// EntryCount is only populated by GiveawayRepository.GetHistory (a
	// display-only aggregate, not a stored column) - zero elsewhere.
	EntryCount int `json:"entry_count"`
}

const (
	GiveawayStatusOpen   = "open"
	GiveawayStatusClosed = "closed"

	MinGiveawayKeywordLength = 2
	MaxGiveawayKeywordLength = 50
)

// NormalizeGiveawayKeyword trims and lowercases a keyword so storage and
// chat-message matching always compare like-for-like.
func NormalizeGiveawayKeyword(keyword string) string {
	return strings.ToLower(strings.TrimSpace(keyword))
}

// ValidateGiveawayKeyword checks an already-normalized keyword.
func ValidateGiveawayKeyword(keyword string) error {
	if len(keyword) < MinGiveawayKeywordLength || len(keyword) > MaxGiveawayKeywordLength {
		return ErrGiveawayKeywordInvalid
	}
	return nil
}

func NewGiveaway(userTwitchID, keyword string, subBonus bool) *Giveaway {
	now := time.Now()
	return &Giveaway{
		UserTwitchID: userTwitchID,
		Status:       GiveawayStatusOpen,
		Keyword:      keyword,
		SubBonus:     subBonus,
		StartedAt:    now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
