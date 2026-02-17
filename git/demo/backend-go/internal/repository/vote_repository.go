package repository

import (
	"context"
	"demo/backend-go/internal/domain"
)

type VoteRepository interface {
	// Sessions
	CreateSession(ctx context.Context, input domain.VoteSessionCreateInput) (*domain.VoteSession, error)

	GetSessionByID(ctx context.Context, id int64) (*domain.VoteSession, error)

	GetSessionsByUser(ctx context.Context, userID string, limit, offset int) ([]*domain.VoteSession, int64, error)

	GetActiveSession(ctx context.Context, userID string) (*domain.VoteSession, error)

	UpdateSession(ctx context.Context, id int64, input domain.VoteSessionUpdateInput) error

	CloseSession(ctx context.Context, id int64, winner string, results []domain.VoteResult) error

	// Votes
	CastVote(ctx context.Context, input domain.VoteCastInput) (*domain.Vote, error)

	GetVotesBySession(ctx context.Context, sessionID int64) ([]*domain.Vote, error)

	GetVoteCounts(ctx context.Context, sessionID int64) (map[string]int, error)

	HasVoted(ctx context.Context, sessionID int64, userID string) (bool, error)
}