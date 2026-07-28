package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"demo/backend-go/internal/repository"
)

type paddleEventRepository struct {
	db *sql.DB
}

func NewPaddleEventRepository(db *sql.DB) repository.PaddleEventRepository {
	return &paddleEventRepository{db: db}
}

func (r *paddleEventRepository) HasProcessed(ctx context.Context, eventID string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM paddle_events WHERE event_id = $1)", eventID,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("fehler beim Pruefen des Paddle-Events: %w", err)
	}
	return exists, nil
}

func (r *paddleEventRepository) MarkProcessed(ctx context.Context, eventID, eventType string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO paddle_events (event_id, event_type)
		VALUES ($1, $2)
		ON CONFLICT (event_id) DO NOTHING
	`, eventID, eventType)
	if err != nil {
		return fmt.Errorf("fehler beim Speichern des Paddle-Events: %w", err)
	}
	return nil
}
