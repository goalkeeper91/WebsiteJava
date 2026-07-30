package domain

import (
	"fmt"
	"time"
)

// CS2NoteSubjectType constrains cs2_notes.subject_type to the two supported
// kinds - notes are matched live against GSI team/player names, no free-form
// subject types needed.
type CS2NoteSubjectType string

const (
	CS2NoteSubjectTeam   CS2NoteSubjectType = "team"
	CS2NoteSubjectPlayer CS2NoteSubjectType = "player"
)

func (t CS2NoteSubjectType) Validate() error {
	if t != CS2NoteSubjectTeam && t != CS2NoteSubjectPlayer {
		return fmt.Errorf("ungültiger subject_type: %s", t)
	}
	return nil
}

// CS2CasterSettings is a per-user singleton (like loyalty_settings) - the
// gsi_token identifies incoming GSI POSTs back to this user, since CS2 posts
// directly to our public URL rather than through any session-authenticated
// path.
type CS2CasterSettings struct {
	UserTwitchID             string    `json:"user_twitch_id" db:"user_twitch_id"`
	GSIToken                 string    `json:"gsi_token" db:"gsi_token"`
	PredictionsEnabled       bool      `json:"predictions_enabled" db:"predictions_enabled"`
	MultikillAnnounceEnabled bool      `json:"multikill_announce_enabled" db:"multikill_announce_enabled"`
	MapEndAnnounceEnabled    bool      `json:"map_end_announce_enabled" db:"map_end_announce_enabled"`
	TitleUpdateEnabled       bool      `json:"title_update_enabled" db:"title_update_enabled"`
	CreatedAt                time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                time.Time `json:"updated_at" db:"updated_at"`
}

// CS2CasterSettingsUpdateInput lets a user toggle each automation
// independently (nil = unverändert lassen), same ApplyUpdate-Stil as
// automation_settings.go.
type CS2CasterSettingsUpdateInput struct {
	PredictionsEnabled       *bool `json:"predictions_enabled,omitempty"`
	MultikillAnnounceEnabled *bool `json:"multikill_announce_enabled,omitempty"`
	MapEndAnnounceEnabled    *bool `json:"map_end_announce_enabled,omitempty"`
	TitleUpdateEnabled       *bool `json:"title_update_enabled,omitempty"`
}

func (s *CS2CasterSettings) ApplyUpdate(input CS2CasterSettingsUpdateInput) {
	if input.PredictionsEnabled != nil {
		s.PredictionsEnabled = *input.PredictionsEnabled
	}
	if input.MultikillAnnounceEnabled != nil {
		s.MultikillAnnounceEnabled = *input.MultikillAnnounceEnabled
	}
	if input.MapEndAnnounceEnabled != nil {
		s.MapEndAnnounceEnabled = *input.MapEndAnnounceEnabled
	}
	if input.TitleUpdateEnabled != nil {
		s.TitleUpdateEnabled = *input.TitleUpdateEnabled
	}
	s.UpdatedAt = time.Now()
}

// CS2Note is a free-text note tied to a team or player name - matched
// case-insensitively against live GSI data to surface the right note while
// that team/player is on screen (see CS2CasterService.GetLiveStatus).
type CS2Note struct {
	ID           int64              `json:"id" db:"id"`
	UserTwitchID string             `json:"user_twitch_id" db:"user_twitch_id"`
	SubjectType  CS2NoteSubjectType `json:"subject_type" db:"subject_type"`
	SubjectName  string             `json:"subject_name" db:"subject_name"`
	Content      string             `json:"content" db:"content"`
	CreatedAt    time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" db:"updated_at"`
}

type CS2NoteCreateInput struct {
	SubjectType CS2NoteSubjectType
	SubjectName string
	Content     string
}

type CS2NoteUpdateInput struct {
	SubjectName *string
	Content     *string
}
