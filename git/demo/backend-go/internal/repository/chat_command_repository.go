package repository

import (
	"context"
	"demo/backend-go/internal/domain"
)

type ChatCommandRepository interface {
	Create(ctx context.Context, command *domain.ChatCommand) error

	GetByID(ctx context.Context, id int64, channelID string) (*domain.ChatCommand, error)

	GetByTrigger(ctx context.Context, channelID, trigger string) (*domain.ChatCommand, error)

	GetAll(ctx context.Context, channelID string, limit, offset int) ([]*domain.ChatCommand, int64, error)

	Search(ctx context.Context, channelID, search string, limit, offset int) ([]*domain.ChatCommand, int64, error)

	GetByStatus(ctx context.Context, channelID string, enabled bool, limit, offset int) ([]*domain.ChatCommand, int64, error)

	// NEW: Filter by command type (simple vs advanced)
	GetByType(ctx context.Context, channelID string, commandType domain.CommandType, limit, offset int) ([]*domain.ChatCommand, int64, error)

	// NEW: Get all active advanced commands (for bot initialization)
	GetAdvancedCommands(ctx context.Context, channelID string) ([]*domain.ChatCommand, error)

	Update(ctx context.Context, command *domain.ChatCommand) error

	Delete(ctx context.Context, id int64, channelID string) error

	Exists(ctx context.Context, channelID, trigger string) (bool, error)

	// NEW: Track usage count
	TrackUsage(ctx context.Context, id int64) error
}