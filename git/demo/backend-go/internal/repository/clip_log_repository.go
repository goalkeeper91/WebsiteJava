package repository

import (
	"context"
	"time"

	"demo/backend-go/internal/domain"
	"github.com/google/uuid"
)

// ClipLogRepository definiert alle Operationen für Clip Logs
type ClipLogRepository interface {
	// Create erstellt einen neuen Clip Log
	Create(ctx context.Context, log *domain.ClipLog) error

	// GetByID lädt einen Clip Log by UUID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.ClipLog, error)

	// GetByClipID lädt einen Clip Log by Twitch Clip ID
	GetByClipID(ctx context.Context, clipID string) (*domain.ClipLog, error)

	// GetByUser lädt alle Clip Logs eines Users (mit Pagination)
	GetByUser(ctx context.Context, userTwitchID string, limit, offset int) ([]*domain.ClipLog, error)

	// GetByUserAndStatus lädt Clips eines Users mit bestimmtem Status
	GetByUserAndStatus(ctx context.Context, userTwitchID string, status domain.ClipStatus, limit, offset int) ([]*domain.ClipLog, error)

	// GetByStatus lädt alle Clips mit bestimmtem Status (für Processing)
	GetByStatus(ctx context.Context, status domain.ClipStatus, limit int) ([]*domain.ClipLog, error)

	// GetPending lädt alle Pending Clips (für Job Queue)
	GetPending(ctx context.Context, limit int) ([]*domain.ClipLog, error)

	// GetFailed lädt alle Failed Clips die retried werden können
	GetFailedRetryable(ctx context.Context, limit int) ([]*domain.ClipLog, error)

	// Update aktualisiert einen Clip Log
	Update(ctx context.Context, log *domain.ClipLog) error

	// UpdateStatus aktualisiert nur den Status
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ClipStatus) error

	// MarkAsProcessing setzt Status auf Processing
	MarkAsProcessing(ctx context.Context, id uuid.UUID, n8nWorkflowID, n8nExecutionID string) error

	// MarkAsCompleted setzt Status auf Completed
	MarkAsCompleted(ctx context.Context, id uuid.UUID) error

	// MarkAsFailed setzt Status auf Failed
	MarkAsFailed(ctx context.Context, id uuid.UUID, errorMsg string) error

	// IncrementRetry erhöht Retry Counter
	IncrementRetry(ctx context.Context, id uuid.UUID) error

	// AddPlatformPost fügt Platform Post hinzu
	AddPlatformPost(ctx context.Context, id uuid.UUID, post domain.PlatformPost) error

	// SetAIContent setzt AI-generierte Inhalte
	SetAIContent(ctx context.Context, id uuid.UUID, caption string, hashtags []string) error

	// Delete löscht einen Clip Log
	Delete(ctx context.Context, id uuid.UUID) error

	// GetStats lädt Statistiken für einen User
	GetStats(ctx context.Context, userTwitchID string) (*ClipLogStats, error)

	// GetByDateRange lädt Clips in einem Zeitraum
	GetByDateRange(ctx context.Context, userTwitchID string, from, to time.Time) ([]*domain.ClipLog, error)

	// ClipExists prüft ob ein Clip bereits existiert
	ClipExists(ctx context.Context, clipID string) (bool, error)

	// GetByN8NExecutionID lädt Clip by n8n Execution ID
	GetByN8NExecutionID(ctx context.Context, executionID string) (*domain.ClipLog, error)

	// BulkUpdateStatus aktualisiert Status für mehrere Clips
	BulkUpdateStatus(ctx context.Context, ids []uuid.UUID, status domain.ClipStatus) error
}

// ClipLogStats enthält Statistiken über Clip Processing
type ClipLogStats struct {
	TotalClips      int64 `json:"total_clips"`
	PendingClips    int64 `json:"pending_clips"`
	ProcessingClips int64 `json:"processing_clips"`
	CompletedClips  int64 `json:"completed_clips"`
	FailedClips     int64 `json:"failed_clips"`
	TotalPosts      int64 `json:"total_posts"`
	SuccessRate     float64 `json:"success_rate"`
}
