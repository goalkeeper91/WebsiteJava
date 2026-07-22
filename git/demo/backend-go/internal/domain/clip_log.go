package domain

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ClipStatus ist eine Enum für Clip Processing Status
type ClipStatus string

const (
	StatusPending    ClipStatus = "pending"
	StatusProcessing ClipStatus = "processing"
	StatusCompleted  ClipStatus = "completed"
	StatusFailed     ClipStatus = "failed"
	StatusCancelled  ClipStatus = "cancelled"
)

// PlatformPost repräsentiert einen Post auf einer Platform
type PlatformPost struct {
	Platform Platform  `json:"platform"`
	PostID   string    `json:"post_id"`
	URL      string    `json:"url"`
	PostedAt time.Time `json:"posted_at"`
}

// PlatformPosts ist ein Array von Platform Posts (JSONB)
type PlatformPosts []PlatformPost

func (p PlatformPosts) Value() (driver.Value, error) {
	if p == nil {
		return json.Marshal([]PlatformPost{})
	}
	return json.Marshal(p)
}

func (p *PlatformPosts) Scan(value interface{}) error {
	if value == nil {
		*p = []PlatformPost{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}
	return json.Unmarshal(bytes, p)
}

// Metadata ist flexible JSONB für zusätzliche Daten
type Metadata map[string]interface{}

func (m Metadata) Value() (driver.Value, error) {
	if m == nil {
		return json.Marshal(map[string]interface{}{})
	}
	return json.Marshal(m)
}

func (m *Metadata) Scan(value interface{}) error {
	if value == nil {
		*m = make(map[string]interface{})
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}
	return json.Unmarshal(bytes, m)
}

// ClipLog repräsentiert die Historie eines verarbeiteten Clips
type ClipLog struct {
	ID               uuid.UUID     `json:"id" db:"id"`
	UserTwitchID     string        `json:"user_twitch_id" db:"user_twitch_id"`
	ClipID           string        `json:"clip_id" db:"clip_id"`
	ClipURL          string        `json:"clip_url" db:"clip_url"`
	ClipTitle        string        `json:"clip_title,omitempty" db:"clip_title"`
	BroadcasterName  string        `json:"broadcaster_name,omitempty" db:"broadcaster_name"`
	GameName         string        `json:"game_name,omitempty" db:"game_name"`
	DurationSeconds  int           `json:"duration_seconds" db:"duration_seconds"`
	ViewCount        int           `json:"view_count" db:"view_count"`
	CreatedAtTwitch  *time.Time    `json:"created_at_twitch,omitempty" db:"created_at_twitch"`
	Status           ClipStatus    `json:"status" db:"status"`
	N8NWorkflowID    string        `json:"n8n_workflow_id,omitempty" db:"n8n_workflow_id"`
	N8NExecutionID   string        `json:"n8n_execution_id,omitempty" db:"n8n_execution_id"`
	N8NJobData       Metadata      `json:"n8n_job_data,omitempty" db:"n8n_job_data"`
	AICaption        string        `json:"ai_caption,omitempty" db:"ai_caption"`
	AIHashtags       StringArray   `json:"ai_hashtags,omitempty" db:"ai_hashtags"`
	PostedPlatforms  PlatformPosts `json:"posted_platforms" db:"posted_platforms"`
	ErrorMessage     string        `json:"error_message,omitempty" db:"error_message"`
	RetryCount       int           `json:"retry_count" db:"retry_count"`
	LastRetryAt      *time.Time    `json:"last_retry_at,omitempty" db:"last_retry_at"`
	Metadata         Metadata      `json:"metadata" db:"metadata"`
	CreatedAt        time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at" db:"updated_at"`
	CompletedAt      *time.Time    `json:"completed_at,omitempty" db:"completed_at"`
}

// NewClipLog erstellt einen neuen ClipLog
func NewClipLog(userTwitchID, clipID, clipURL string) *ClipLog {
	now := time.Now()
	return &ClipLog{
		ID:              uuid.New(),
		UserTwitchID:    userTwitchID,
		ClipID:          clipID,
		ClipURL:         clipURL,
		Status:          StatusPending,
		ViewCount:       0,
		RetryCount:      0,
		PostedPlatforms: PlatformPosts{},
		Metadata:        Metadata{},
		N8NJobData:      Metadata{},
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// MarkAsProcessing setzt Status auf Processing
func (cl *ClipLog) MarkAsProcessing(n8nWorkflowID, n8nExecutionID string) {
	cl.Status = StatusProcessing
	cl.N8NWorkflowID = n8nWorkflowID
	cl.N8NExecutionID = n8nExecutionID
	cl.UpdatedAt = time.Now()
}

// MarkAsCompleted setzt Status auf Completed
func (cl *ClipLog) MarkAsCompleted() {
	cl.Status = StatusCompleted
	now := time.Now()
	cl.CompletedAt = &now
	cl.UpdatedAt = now
}

// MarkAsFailed setzt Status auf Failed
func (cl *ClipLog) MarkAsFailed(errorMsg string) {
	cl.Status = StatusFailed
	cl.ErrorMessage = errorMsg
	cl.UpdatedAt = time.Now()
}

// IncrementRetry erhöht den Retry Counter
func (cl *ClipLog) IncrementRetry() {
	cl.RetryCount++
	now := time.Now()
	cl.LastRetryAt = &now
	cl.UpdatedAt = now
}

// AddPlatformPost fügt einen Platform Post hinzu
func (cl *ClipLog) AddPlatformPost(platform Platform, postID, url string) {
	post := PlatformPost{
		Platform: platform,
		PostID:   postID,
		URL:      url,
		PostedAt: time.Now(),
	}
	cl.PostedPlatforms = append(cl.PostedPlatforms, post)
	cl.UpdatedAt = time.Now()
}

// SetAIContent setzt AI-generierte Inhalte
func (cl *ClipLog) SetAIContent(caption string, hashtags []string) {
	cl.AICaption = caption
	cl.AIHashtags = hashtags
	cl.UpdatedAt = time.Now()
}

// CanRetry prüft ob Retry möglich ist (max 3 Versuche)
func (cl *ClipLog) CanRetry() bool {
	return cl.RetryCount < 3 && cl.Status == StatusFailed
}
