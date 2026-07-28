package repository

import "context"

type PaddleEventRepository interface {
	// HasProcessed reports whether this webhook event_id has already been
	// handled - Paddle can redeliver the same event.
	HasProcessed(ctx context.Context, eventID string) (bool, error)

	MarkProcessed(ctx context.Context, eventID, eventType string) error
}
