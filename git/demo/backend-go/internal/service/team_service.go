package service

import (
	"context"
	"fmt"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
	"demo/backend-go/internal/twitch"
)

type TeamService struct {
	teamRepo repository.TeamMemberRepository
	userRepo repository.UserRepository
	appToken *twitch.TwitchAppTokenClient
}

func NewTeamService(teamRepo repository.TeamMemberRepository, userRepo repository.UserRepository, appToken *twitch.TwitchAppTokenClient) *TeamService {
	return &TeamService{teamRepo: teamRepo, userRepo: userRepo, appToken: appToken}
}

// InviteMember resolves memberLogin to a Twitch ID via Helix immediately -
// the invited person doesn't need to have ever logged into the site before.
func (s *TeamService) InviteMember(ctx context.Context, ownerTwitchID, memberLogin string) (*domain.TeamMember, error) {
	twitchUser, err := s.appToken.GetUserByLogin(ctx, memberLogin)
	if err != nil {
		return nil, fmt.Errorf("fehler bei der Twitch-Nutzerauflösung: %w", err)
	}
	if twitchUser == nil {
		return nil, domain.ErrTwitchUserNotFound
	}
	if twitchUser.ID == ownerTwitchID {
		return nil, domain.ErrCannotInviteSelf
	}

	inserted, err := s.teamRepo.Create(ctx, ownerTwitchID, twitchUser.ID, twitchUser.Login)
	if err != nil {
		return nil, err
	}
	if !inserted {
		return nil, domain.ErrAlreadyTeamMember
	}

	return &domain.TeamMember{
		OwnerTwitchID:  ownerTwitchID,
		MemberTwitchID: twitchUser.ID,
		MemberLogin:    twitchUser.Login,
	}, nil
}

// RemoveMember also covers a member "leaving" a team they don't own - the
// handler authorizes either the owner or the member themselves to call this.
func (s *TeamService) RemoveMember(ctx context.Context, ownerTwitchID, memberTwitchID string) error {
	return s.teamRepo.Delete(ctx, ownerTwitchID, memberTwitchID)
}

// ListMembers is for the owner's "Team" tab - who has access to my channel.
func (s *TeamService) ListMembers(ctx context.Context, ownerTwitchID string) ([]*domain.TeamMember, error) {
	return s.teamRepo.ListByOwner(ctx, ownerTwitchID)
}

// ManagedChannel is a hydrated view of a TeamMember row from the member's
// side - which channel they can help manage, with the owner's display name
// for the dashboard's channel switcher.
type ManagedChannel struct {
	OwnerTwitchID string `json:"owner_twitch_id"`
	OwnerLogin    string `json:"owner_login"`
}

// ListManagedChannels backs the channel switcher - which channels can the
// logged-in user (memberTwitchID) help manage.
func (s *TeamService) ListManagedChannels(ctx context.Context, memberTwitchID string) ([]*ManagedChannel, error) {
	members, err := s.teamRepo.ListByMember(ctx, memberTwitchID)
	if err != nil {
		return nil, err
	}

	channels := make([]*ManagedChannel, 0, len(members))
	for _, m := range members {
		owner, err := s.userRepo.GetByTwitchID(ctx, m.OwnerTwitchID)
		if err != nil {
			return nil, fmt.Errorf("fehler beim Laden des Kanalinhabers: %w", err)
		}
		login := m.OwnerTwitchID
		if owner != nil {
			login = owner.Username
		}
		channels = append(channels, &ManagedChannel{OwnerTwitchID: m.OwnerTwitchID, OwnerLogin: login})
	}

	return channels, nil
}

// HasAccess is the actual authorization check used by every request that
// might act on behalf of another channel - re-checked per request (never
// cached), so revoking access takes effect immediately, not just at the
// requester's next login.
func (s *TeamService) HasAccess(ctx context.Context, ownerTwitchID, requestingTwitchID string) (bool, error) {
	if requestingTwitchID == ownerTwitchID {
		return true, nil
	}
	return s.teamRepo.Exists(ctx, ownerTwitchID, requestingTwitchID)
}
