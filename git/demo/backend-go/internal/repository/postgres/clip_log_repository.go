package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
	"github.com/google/uuid"
)

var (
	ErrClipLogNotFound = errors.New("clip log not found")
)

type clipLogRepository struct {
	db *sql.DB
}

func NewClipLogRepository(db *sql.DB) repository.ClipLogRepository {
	return &clipLogRepository{db: db}
}

func (r *clipLogRepository) Create(ctx context.Context, log *domain.ClipLog) error {
	query := `
		INSERT INTO clip_logs
		(user_twitch_id, clip_id, clip_url, clip_title, broadcaster_name,
		 game_name, duration_seconds, view_count, created_at_twitch,
		 status, n8n_job_data, metadata, posted_platforms, ai_hashtags,
		 retry_count, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		log.UserTwitchID,
		log.ClipID,
		log.ClipURL,
		log.ClipTitle,
		log.BroadcasterName,
		log.GameName,
		log.DurationSeconds,
		log.ViewCount,
		log.CreatedAtTwitch,
		log.Status,
		log.N8NJobData,
		log.Metadata,
		log.PostedPlatforms,
		log.AIHashtags,
		log.RetryCount,
		log.CreatedAt,
		log.UpdatedAt,
	).Scan(&log.ID)

	if err != nil {
		return fmt.Errorf("fehler beim Erstellen des Clip Logs: %w", err)
	}

	return nil
}

func (r *clipLogRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.ClipLog, error) {
	query := `
		SELECT id, user_twitch_id, clip_id, clip_url, clip_title, broadcaster_name,
		       game_name, duration_seconds, view_count, created_at_twitch,
		       status, n8n_workflow_id, n8n_execution_id, n8n_job_data,
		       ai_caption, ai_hashtags, posted_platforms, error_message,
		       retry_count, last_retry_at, metadata, created_at, updated_at, completed_at
		FROM clip_logs
		WHERE id = $1
	`

	var log domain.ClipLog
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&log.ID,
		&log.UserTwitchID,
		&log.ClipID,
		&log.ClipURL,
		&log.ClipTitle,
		&log.BroadcasterName,
		&log.GameName,
		&log.DurationSeconds,
		&log.ViewCount,
		&log.CreatedAtTwitch,
		&log.Status,
		&log.N8NWorkflowID,
		&log.N8NExecutionID,
		&log.N8NJobData,
		&log.AICaption,
		&log.AIHashtags,
		&log.PostedPlatforms,
		&log.ErrorMessage,
		&log.RetryCount,
		&log.LastRetryAt,
		&log.Metadata,
		&log.CreatedAt,
		&log.UpdatedAt,
		&log.CompletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrClipLogNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden des Clip Logs: %w", err)
	}

	return &log, nil
}

func (r *clipLogRepository) GetByClipID(ctx context.Context, clipID string) (*domain.ClipLog, error) {
	query := `
		SELECT id, user_twitch_id, clip_id, clip_url, clip_title, broadcaster_name,
		       game_name, duration_seconds, view_count, created_at_twitch,
		       status, n8n_workflow_id, n8n_execution_id, n8n_job_data,
		       ai_caption, ai_hashtags, posted_platforms, error_message,
		       retry_count, last_retry_at, metadata, created_at, updated_at, completed_at
		FROM clip_logs
		WHERE clip_id = $1
	`

	var log domain.ClipLog
	err := r.db.QueryRowContext(ctx, query, clipID).Scan(
		&log.ID, &log.UserTwitchID, &log.ClipID, &log.ClipURL, &log.ClipTitle,
		&log.BroadcasterName, &log.GameName, &log.DurationSeconds, &log.ViewCount,
		&log.CreatedAtTwitch, &log.Status, &log.N8NWorkflowID, &log.N8NExecutionID,
		&log.N8NJobData, &log.AICaption, &log.AIHashtags, &log.PostedPlatforms,
		&log.ErrorMessage, &log.RetryCount, &log.LastRetryAt, &log.Metadata,
		&log.CreatedAt, &log.UpdatedAt, &log.CompletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrClipLogNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden: %w", err)
	}

	return &log, nil
}

func (r *clipLogRepository) GetByUser(ctx context.Context, userTwitchID string, limit, offset int) ([]*domain.ClipLog, error) {
	query := `
		SELECT id, user_twitch_id, clip_id, clip_url, clip_title, broadcaster_name,
		       game_name, duration_seconds, view_count, created_at_twitch,
		       status, n8n_workflow_id, n8n_execution_id, n8n_job_data,
		       ai_caption, ai_hashtags, posted_platforms, error_message,
		       retry_count, last_retry_at, metadata, created_at, updated_at, completed_at
		FROM clip_logs
		WHERE user_twitch_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	return r.queryMultiple(ctx, query, userTwitchID, limit, offset)
}

func (r *clipLogRepository) GetByUserAndStatus(ctx context.Context, userTwitchID string, status domain.ClipStatus, limit, offset int) ([]*domain.ClipLog, error) {
	query := `
		SELECT id, user_twitch_id, clip_id, clip_url, clip_title, broadcaster_name,
		       game_name, duration_seconds, view_count, created_at_twitch,
		       status, n8n_workflow_id, n8n_execution_id, n8n_job_data,
		       ai_caption, ai_hashtags, posted_platforms, error_message,
		       retry_count, last_retry_at, metadata, created_at, updated_at, completed_at
		FROM clip_logs
		WHERE user_twitch_id = $1 AND status = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	return r.queryMultiple(ctx, query, userTwitchID, status, limit, offset)
}

func (r *clipLogRepository) GetByStatus(ctx context.Context, status domain.ClipStatus, limit int) ([]*domain.ClipLog, error) {
	query := `
		SELECT id, user_twitch_id, clip_id, clip_url, clip_title, broadcaster_name,
		       game_name, duration_seconds, view_count, created_at_twitch,
		       status, n8n_workflow_id, n8n_execution_id, n8n_job_data,
		       ai_caption, ai_hashtags, posted_platforms, error_message,
		       retry_count, last_retry_at, metadata, created_at, updated_at, completed_at
		FROM clip_logs
		WHERE status = $1
		ORDER BY created_at ASC
		LIMIT $2
	`

	return r.queryMultiple(ctx, query, status, limit)
}

func (r *clipLogRepository) GetPending(ctx context.Context, limit int) ([]*domain.ClipLog, error) {
	return r.GetByStatus(ctx, domain.StatusPending, limit)
}

func (r *clipLogRepository) GetFailedRetryable(ctx context.Context, limit int) ([]*domain.ClipLog, error) {
	query := `
		SELECT id, user_twitch_id, clip_id, clip_url, clip_title, broadcaster_name,
		       game_name, duration_seconds, view_count, created_at_twitch,
		       status, n8n_workflow_id, n8n_execution_id, n8n_job_data,
		       ai_caption, ai_hashtags, posted_platforms, error_message,
		       retry_count, last_retry_at, metadata, created_at, updated_at, completed_at
		FROM clip_logs
		WHERE status = 'failed' AND retry_count < 3
		ORDER BY created_at ASC
		LIMIT $1
	`

	return r.queryMultiple(ctx, query, limit)
}

func (r *clipLogRepository) queryMultiple(ctx context.Context, query string, args ...interface{}) ([]*domain.ClipLog, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Query: %w", err)
	}
	defer rows.Close()

	logs := make([]*domain.ClipLog, 0)
	for rows.Next() {
		var log domain.ClipLog
		err := rows.Scan(
			&log.ID, &log.UserTwitchID, &log.ClipID, &log.ClipURL, &log.ClipTitle,
			&log.BroadcasterName, &log.GameName, &log.DurationSeconds, &log.ViewCount,
			&log.CreatedAtTwitch, &log.Status, &log.N8NWorkflowID, &log.N8NExecutionID,
			&log.N8NJobData, &log.AICaption, &log.AIHashtags, &log.PostedPlatforms,
			&log.ErrorMessage, &log.RetryCount, &log.LastRetryAt, &log.Metadata,
			&log.CreatedAt, &log.UpdatedAt, &log.CompletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("fehler beim Scan: %w", err)
		}
		logs = append(logs, &log)
	}

	return logs, rows.Err()
}

func (r *clipLogRepository) Update(ctx context.Context, log *domain.ClipLog) error {
	query := `
		UPDATE clip_logs
		SET clip_title = $2,
		    broadcaster_name = $3,
		    game_name = $4,
		    duration_seconds = $5,
		    view_count = $6,
		    status = $7,
		    n8n_workflow_id = $8,
		    n8n_execution_id = $9,
		    n8n_job_data = $10,
		    ai_caption = $11,
		    ai_hashtags = $12,
		    posted_platforms = $13,
		    error_message = $14,
		    retry_count = $15,
		    last_retry_at = $16,
		    metadata = $17,
		    updated_at = $18,
		    completed_at = $19
		WHERE id = $1
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		log.ID,
		log.ClipTitle,
		log.BroadcasterName,
		log.GameName,
		log.DurationSeconds,
		log.ViewCount,
		log.Status,
		log.N8NWorkflowID,
		log.N8NExecutionID,
		log.N8NJobData,
		log.AICaption,
		log.AIHashtags,
		log.PostedPlatforms,
		log.ErrorMessage,
		log.RetryCount,
		log.LastRetryAt,
		log.Metadata,
		time.Now(),
		log.CompletedAt,
	)

	if err != nil {
		return fmt.Errorf("fehler beim Update: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrClipLogNotFound
	}

	return nil
}

// Fortsetzung von clip_log_repository_part1.go
// Diese Methoden gehören zur clipLogRepository struct

func (r *clipLogRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ClipStatus) error {
	query := `UPDATE clip_logs SET status = $2, updated_at = $3 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id, status, time.Now())
	return err
}

func (r *clipLogRepository) MarkAsProcessing(ctx context.Context, id uuid.UUID, n8nWorkflowID, n8nExecutionID string) error {
	query := `
		UPDATE clip_logs
		SET status = 'processing',
		    n8n_workflow_id = $2,
		    n8n_execution_id = $3,
		    updated_at = $4
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id, n8nWorkflowID, n8nExecutionID, time.Now())
	return err
}

func (r *clipLogRepository) MarkAsCompleted(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE clip_logs
		SET status = 'completed',
		    completed_at = $2,
		    updated_at = $2
		WHERE id = $1
	`

	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, id, now)
	return err
}

func (r *clipLogRepository) MarkAsFailed(ctx context.Context, id uuid.UUID, errorMsg string) error {
	query := `
		UPDATE clip_logs
		SET status = 'failed',
		    error_message = $2,
		    updated_at = $3
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id, errorMsg, time.Now())
	return err
}

func (r *clipLogRepository) IncrementRetry(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE clip_logs
		SET retry_count = retry_count + 1,
		    last_retry_at = $2,
		    updated_at = $2
		WHERE id = $1
	`

	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, id, now)
	return err
}

func (r *clipLogRepository) AddPlatformPost(ctx context.Context, id uuid.UUID, post domain.PlatformPost) error {
	// JSONB Array Append
	query := `
		UPDATE clip_logs
		SET posted_platforms = posted_platforms || $2::jsonb,
		    updated_at = $3
		WHERE id = $1
	`

	postJSON, err := json.Marshal([]domain.PlatformPost{post})
	if err != nil {
		return fmt.Errorf("fehler beim Marshal: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, id, postJSON, time.Now())
	return err
}

func (r *clipLogRepository) SetAIContent(ctx context.Context, id uuid.UUID, caption string, hashtags []string) error {
	query := `
		UPDATE clip_logs
		SET ai_caption = $2,
		    ai_hashtags = $3,
		    updated_at = $4
		WHERE id = $1
	`

	hashtagArray := domain.StringArray(hashtags)
	_, err := r.db.ExecContext(ctx, query, id, caption, hashtagArray, time.Now())
	return err
}

func (r *clipLogRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM clip_logs WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrClipLogNotFound
	}

	return nil
}

func (r *clipLogRepository) GetStats(ctx context.Context, userTwitchID string) (*repository.ClipLogStats, error) {
	query := `
		SELECT
		    COUNT(*) as total_clips,
		    COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_clips,
		    COUNT(CASE WHEN status = 'processing' THEN 1 END) as processing_clips,
		    COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_clips,
		    COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_clips,
		    COALESCE(SUM(jsonb_array_length(posted_platforms)), 0) as total_posts
		FROM clip_logs
		WHERE user_twitch_id = $1
	`

	var stats repository.ClipLogStats
	err := r.db.QueryRowContext(ctx, query, userTwitchID).Scan(
		&stats.TotalClips,
		&stats.PendingClips,
		&stats.ProcessingClips,
		&stats.CompletedClips,
		&stats.FailedClips,
		&stats.TotalPosts,
	)

	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Stats: %w", err)
	}

	// Success Rate berechnen
	if stats.TotalClips > 0 {
		stats.SuccessRate = float64(stats.CompletedClips) / float64(stats.TotalClips) * 100
	}

	return &stats, nil
}

func (r *clipLogRepository) GetByDateRange(ctx context.Context, userTwitchID string, from, to time.Time) ([]*domain.ClipLog, error) {
	query := `
		SELECT id, user_twitch_id, clip_id, clip_url, clip_title, broadcaster_name,
		       game_name, duration_seconds, view_count, created_at_twitch,
		       status, n8n_workflow_id, n8n_execution_id, n8n_job_data,
		       ai_caption, ai_hashtags, posted_platforms, error_message,
		       retry_count, last_retry_at, metadata, created_at, updated_at, completed_at
		FROM clip_logs
		WHERE user_twitch_id = $1
		  AND created_at >= $2
		  AND created_at <= $3
		ORDER BY created_at DESC
	`

	return r.queryMultiple(ctx, query, userTwitchID, from, to)
}

func (r *clipLogRepository) ClipExists(ctx context.Context, clipID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM clip_logs WHERE clip_id = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, clipID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("fehler beim Check: %w", err)
	}

	return exists, nil
}

func (r *clipLogRepository) GetByN8NExecutionID(ctx context.Context, executionID string) (*domain.ClipLog, error) {
	query := `
		SELECT id, user_twitch_id, clip_id, clip_url, clip_title, broadcaster_name,
		       game_name, duration_seconds, view_count, created_at_twitch,
		       status, n8n_workflow_id, n8n_execution_id, n8n_job_data,
		       ai_caption, ai_hashtags, posted_platforms, error_message,
		       retry_count, last_retry_at, metadata, created_at, updated_at, completed_at
		FROM clip_logs
		WHERE n8n_execution_id = $1
	`

	var log domain.ClipLog
	err := r.db.QueryRowContext(ctx, query, executionID).Scan(
		&log.ID, &log.UserTwitchID, &log.ClipID, &log.ClipURL, &log.ClipTitle,
		&log.BroadcasterName, &log.GameName, &log.DurationSeconds, &log.ViewCount,
		&log.CreatedAtTwitch, &log.Status, &log.N8NWorkflowID, &log.N8NExecutionID,
		&log.N8NJobData, &log.AICaption, &log.AIHashtags, &log.PostedPlatforms,
		&log.ErrorMessage, &log.RetryCount, &log.LastRetryAt, &log.Metadata,
		&log.CreatedAt, &log.UpdatedAt, &log.CompletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrClipLogNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden: %w", err)
	}

	return &log, nil
}

func (r *clipLogRepository) BulkUpdateStatus(ctx context.Context, ids []uuid.UUID, status domain.ClipStatus) error {
	if len(ids) == 0 {
		return nil
	}

	// Build placeholders for IN clause: $1, $2, $3...
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids)+2)

	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	args[len(ids)] = status
	args[len(ids)+1] = time.Now()

	query := fmt.Sprintf(`
		UPDATE clip_logs
		SET status = $%d, updated_at = $%d
		WHERE id IN (%s)
	`, len(ids)+1, len(ids)+2, strings.Join(placeholders, ", "))

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}