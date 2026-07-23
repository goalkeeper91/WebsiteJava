package repository

import (
	"context"
	"time"

	"demo/backend-go/internal/domain"
)

type ScheduledMessageRepository interface {
	Create(ctx context.Context, message *domain.ScheduledMessage) error

	GetByID(ctx context.Context, id int64, channelID string) (*domain.ScheduledMessage, error)

	GetAll(ctx context.Context, channelID string, limit, offset int) ([]*domain.ScheduledMessage, int64, error)

	Update(ctx context.Context, message *domain.ScheduledMessage) error

	Delete(ctx context.Context, id int64, channelID string) error

	// GetDue returns enabled messages whose next_send_at has passed, across
	// all channels - the scheduler's periodic tick uses this instead of a
	// per-channel query since it has to check everyone at once.
	GetDue(ctx context.Context, limit int) ([]*domain.ScheduledMessage, error)

	// MarkSent advances next_send_at and records last_sent_at after a tick
	// has processed a message (whether or not it actually posted - e.g. a
	// skipped offline channel still needs its timer to move forward).
	MarkSent(ctx context.Context, id int64, sentAt, nextSendAt time.Time) error
}
