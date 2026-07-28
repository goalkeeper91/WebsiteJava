package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
)

type streamSessionRepository struct {
	db *sql.DB
}

func NewStreamSessionRepository(db *sql.DB) repository.StreamSessionRepository {
	return &streamSessionRepository{db: db}
}

const streamSessionColumns = `
	id, twitch_user_id, started_at, ended_at, viewer_sum, sample_count, peak_viewers, created_at, updated_at
`

func scanStreamSession(row interface{ Scan(dest ...interface{}) error }) (*domain.StreamSession, error) {
	s := &domain.StreamSession{}
	var endedAt sql.NullTime

	err := row.Scan(
		&s.ID, &s.TwitchUserID, &s.StartedAt, &endedAt, &s.ViewerSum, &s.SampleCount, &s.PeakViewers,
		&s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if endedAt.Valid {
		s.EndedAt = &endedAt.Time
	}

	return s, nil
}

func (r *streamSessionRepository) GetOpenSessionsByUsers(ctx context.Context, twitchUserIDs []string) (map[string]*domain.StreamSession, error) {
	result := make(map[string]*domain.StreamSession)
	if len(twitchUserIDs) == 0 {
		return result, nil
	}

	query := `
		SELECT ` + streamSessionColumns + `
		FROM stream_sessions
		WHERE ended_at IS NULL AND twitch_user_id = ANY($1)
	`

	rows, err := r.db.QueryContext(ctx, query, pq.Array(twitchUserIDs))
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden offener Stream-Sessions: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		s, err := scanStreamSession(rows)
		if err != nil {
			return nil, fmt.Errorf("fehler beim Scannen einer Stream-Session: %w", err)
		}
		result[s.TwitchUserID] = s
	}

	return result, nil
}

func (r *streamSessionRepository) CreateSession(ctx context.Context, twitchUserID string, startedAt time.Time) (*domain.StreamSession, error) {
	now := time.Now()
	s := &domain.StreamSession{
		TwitchUserID: twitchUserID,
		StartedAt:    startedAt,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	query := `
		INSERT INTO stream_sessions (twitch_user_id, started_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query, s.TwitchUserID, s.StartedAt, s.CreatedAt, s.UpdatedAt).Scan(&s.ID)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Anlegen der Stream-Session: %w", err)
	}

	return s, nil
}

func (r *streamSessionRepository) RecordSample(ctx context.Context, sessionID int64, viewerCount int) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE stream_sessions
		SET viewer_sum = viewer_sum + $2,
		    sample_count = sample_count + 1,
		    peak_viewers = GREATEST(peak_viewers, $2),
		    updated_at = NOW()
		WHERE id = $1
	`, sessionID, viewerCount)
	if err != nil {
		return fmt.Errorf("fehler beim Speichern des Viewer-Samples: %w", err)
	}
	return nil
}

func (r *streamSessionRepository) CloseSession(ctx context.Context, sessionID int64) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE stream_sessions
		SET ended_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`, sessionID)
	if err != nil {
		return fmt.Errorf("fehler beim Schließen der Stream-Session: %w", err)
	}
	return nil
}

func (r *streamSessionRepository) GetLastClosedSession(ctx context.Context, twitchUserID string) (*domain.StreamSession, error) {
	query := `
		SELECT ` + streamSessionColumns + `
		FROM stream_sessions
		WHERE twitch_user_id = $1 AND ended_at IS NOT NULL
		ORDER BY ended_at DESC
		LIMIT 1
	`

	s, err := scanStreamSession(r.db.QueryRowContext(ctx, query, twitchUserID))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der letzten Stream-Session: %w", err)
	}

	return s, nil
}
