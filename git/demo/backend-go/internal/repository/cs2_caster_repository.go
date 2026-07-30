package repository

import (
	"context"

	"demo/backend-go/internal/domain"
)

type CS2CasterRepository interface {
	// GetOrCreateSettings lazily creates a settings row (incl. a fresh
	// gsi_token) on first access - same pattern as Free-Tier-Subscriptions.
	GetOrCreateSettings(ctx context.Context, userTwitchID string) (*domain.CS2CasterSettings, error)

	// GetByGSIToken resolves an incoming GSI POST back to the owning user -
	// the token in the URL path is the sole authentication for that endpoint.
	GetByGSIToken(ctx context.Context, token string) (*domain.CS2CasterSettings, error)

	UpdateSettings(ctx context.Context, userTwitchID string, input domain.CS2CasterSettingsUpdateInput) error

	ListNotes(ctx context.Context, userTwitchID string) ([]*domain.CS2Note, error)
	CreateNote(ctx context.Context, userTwitchID string, input domain.CS2NoteCreateInput) (*domain.CS2Note, error)
	UpdateNote(ctx context.Context, userTwitchID string, noteID int64, input domain.CS2NoteUpdateInput) error
	DeleteNote(ctx context.Context, userTwitchID string, noteID int64) error
}
