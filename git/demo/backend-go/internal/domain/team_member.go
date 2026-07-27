package domain

import (
	"errors"
	"time"
)

// TeamMember grants member_twitch_id full parity access to owner_twitch_id's
// Twitch-chatbot dashboard (Automod, Loyalty, Chat Commands, Scheduled
// Messages, Giveaways) - no roles, no accept/decline step. Removal is a hard
// delete; there is no pending/accepted status.
type TeamMember struct {
	ID             int64     `json:"id"`
	OwnerTwitchID  string    `json:"owner_twitch_id"`
	MemberTwitchID string    `json:"member_twitch_id"`
	MemberLogin    string    `json:"member_login"`
	CreatedAt      time.Time `json:"created_at"`
}

var (
	ErrTwitchUserNotFound = errors.New("dieser Twitch-Nutzername existiert nicht")
	ErrCannotInviteSelf   = errors.New("du kannst dich nicht selbst einladen")
	ErrAlreadyTeamMember  = errors.New("diese Person hat bereits Zugriff")
	ErrNotTeamMember      = errors.New("kein Team-Zugriff gefunden")
)
