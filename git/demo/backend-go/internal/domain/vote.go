package domain

import (
	"encoding/json"
	"time"
)

type VoteSessionStatus string

const (
	VoteSessionActive   VoteSessionStatus = "active"
	VoteSessionClosed   VoteSessionStatus = "closed"
	VoteSessionCanceled VoteSessionStatus = "canceled"
)

// VoteOption represents a voteable option
type VoteOption struct {
	Key   string `json:"key"`   // e.g. "warrior"
	Label string `json:"label"` // e.g. "Warrior - Tank Class"
}

// VoteResult represents the result per option
type VoteResult struct {
	Key        string  `json:"key"`
	Label      string  `json:"label"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

type VoteSession struct {
	ID          int64             `json:"id" db:"id"`
	UserID      string            `json:"userId" db:"user_id"`
	Category    string            `json:"category" db:"category"`
	Title       *string           `json:"title" db:"title"`
	Description *string           `json:"description" db:"description"`
	Status      VoteSessionStatus `json:"status" db:"status"`
	StartedAt   time.Time         `json:"startedAt" db:"started_at"`
	EndedAt     *time.Time        `json:"endedAt" db:"ended_at"`
	Winner      *string           `json:"winner" db:"winner"`
	TotalVotes  int               `json:"totalVotes" db:"total_votes"`
	Options     []VoteOption      `json:"options" db:"-"`
	Results     []VoteResult      `json:"results,omitempty" db:"-"`

	// Raw JSON from DB
	OptionsJSON json.RawMessage `json:"-" db:"options"`
	ResultsJSON json.RawMessage `json:"-" db:"results"`
}

type VoteSessionCreateInput struct {
	UserID      string
	Category    string
	Title       *string
	Description *string
	Options     []VoteOption
}

type VoteSessionUpdateInput struct {
	Title       *string
	Description *string
	Status      *VoteSessionStatus
	EndedAt     *time.Time
	Winner      *string
	TotalVotes  *int
	Results     []VoteResult
}

type Vote struct {
	ID        int64     `json:"id" db:"id"`
	SessionID int64     `json:"sessionId" db:"session_id"`
	UserID    string    `json:"userId" db:"user_id"`
	Option    string    `json:"option" db:"option"`
	Weight    int       `json:"weight" db:"weight"` // 1 = normal, 2 = subscriber, etc.
	VotedAt   time.Time `json:"votedAt" db:"voted_at"`
}

type VoteCastInput struct {
	SessionID int64
	UserID    string
	Option    string
	Weight    int // Default 1
}

// VoteSessionWithResults combines session + computed results
type VoteSessionWithResults struct {
	VoteSession
	Results    []VoteResult `json:"results"`
	TotalVotes int          `json:"totalVotes"`
}

func NewVoteSession(userID, category string, options []VoteOption) *VoteSession {
	now := time.Now()
	return &VoteSession{
		UserID:    userID,
		Category:  category,
		Status:    VoteSessionActive,
		StartedAt: now,
		Options:   options,
	}
}

func (v *VoteSession) IsActive() bool {
	return v.Status == VoteSessionActive
}

func (v *VoteSession) Close(winner string, results []VoteResult) {
	now := time.Now()
	v.Status = VoteSessionClosed
	v.EndedAt = &now
	v.Winner = &winner
	v.Results = results
}

func (v *VoteSession) Cancel() {
	now := time.Now()
	v.Status = VoteSessionCanceled
	v.EndedAt = &now
}

func (v *VoteSession) ParseOptions() error {
	if v.OptionsJSON == nil {
		return nil
	}
	return json.Unmarshal(v.OptionsJSON, &v.Options)
}

func (v *VoteSession) ParseResults() error {
	if v.ResultsJSON == nil {
		return nil
	}
	return json.Unmarshal(v.ResultsJSON, &v.Results)
}

// DetermineWinner finds the option with the most votes
func DetermineWinner(results []VoteResult) string {
	if len(results) == 0 {
		return ""
	}
	winner := results[0]
	for _, r := range results[1:] {
		if r.Count > winner.Count {
			winner = r
		}
	}
	return winner.Key
}

// CalculateResults computes percentages from raw vote counts
func CalculateResults(votes map[string]int, options []VoteOption) []VoteResult {
	total := 0
	for _, count := range votes {
		total += count
	}

	results := make([]VoteResult, 0, len(options))
	for _, opt := range options {
		count := votes[opt.Key]
		percentage := 0.0
		if total > 0 {
			percentage = float64(count) / float64(total) * 100
		}
		results = append(results, VoteResult{
			Key:        opt.Key,
			Label:      opt.Label,
			Count:      count,
			Percentage: percentage,
		})
	}
	return results
}