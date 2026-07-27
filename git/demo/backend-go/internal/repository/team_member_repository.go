package repository

import (
	"context"

	"demo/backend-go/internal/domain"
)

type TeamMemberRepository interface {
	// Create is a no-op (inserted=false) if the pair already exists, so the
	// service can translate that into ErrAlreadyTeamMember without a
	// separate existence check first.
	Create(ctx context.Context, ownerTwitchID, memberTwitchID, memberLogin string) (inserted bool, err error)

	Delete(ctx context.Context, ownerTwitchID, memberTwitchID string) error

	// ListByOwner returns everyone with access to ownerTwitchID's channel.
	ListByOwner(ctx context.Context, ownerTwitchID string) ([]*domain.TeamMember, error)

	// ListByMember returns every channel memberTwitchID has been granted
	// access to - backs the dashboard's channel switcher.
	ListByMember(ctx context.Context, memberTwitchID string) ([]*domain.TeamMember, error)

	Exists(ctx context.Context, ownerTwitchID, memberTwitchID string) (bool, error)
}
