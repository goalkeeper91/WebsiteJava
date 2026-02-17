package service

import (
	"context"
	"fmt"
	"log"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
)

type VoteService struct {
	voteRepo        repository.VoteRepository
	n8nSvc          *N8NService
	subscriptionSvc *SubscriptionService
	analyticsRepo   repository.UsageAnalyticsRepository
}

func NewVoteService(
	voteRepo repository.VoteRepository,
	n8nSvc *N8NService,
	subscriptionSvc *SubscriptionService,
	analyticsRepo repository.UsageAnalyticsRepository,
) *VoteService {
	return &VoteService{
		voteRepo:        voteRepo,
		n8nSvc:          n8nSvc,
		subscriptionSvc: subscriptionSvc,
		analyticsRepo:   analyticsRepo,
	}
}

// StartSession creates a new vote session (requires n8n integration)
func (s *VoteService) StartSession(ctx context.Context, userID string, input domain.VoteSessionCreateInput) (*domain.VoteSession, error) {
	// Check if n8n is ready
	ready, err := s.n8nSvc.IsReady(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !ready {
		return nil, domain.ErrN8NIntegrationNotReady
	}

	// Check if there's already an active session
	existing, err := s.voteRepo.GetActiveSession(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Prüfen aktiver Sessions: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("es gibt bereits eine aktive Vote Session (ID: %d)", existing.ID)
	}

	// Check vote limit
	integration, err := s.n8nSvc.GetIntegration(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err := s.subscriptionSvc.CheckVoteLimit(ctx, userID, integration.VotesThisMonth); err != nil {
		return nil, err
	}

	// Ensure options are set
	input.UserID = userID
	if len(input.Options) == 0 {
		return nil, fmt.Errorf("mindestens eine Option ist erforderlich")
	}

	session, err := s.voteRepo.CreateSession(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Erstellen der Vote Session: %w", err)
	}

	// Track analytics
	_ = s.analyticsRepo.Track(ctx, domain.UsageAnalyticsCreateInput{
		UserID:      userID,
		MetricType:  domain.MetricVoteSession,
		MetricValue: 1,
	})

	log.Printf("✅ Vote Session gestartet: ID=%d, Category=%s, User=%s", session.ID, session.Category, userID)
	return session, nil
}

// GetSession returns a vote session by ID
func (s *VoteService) GetSession(ctx context.Context, sessionID int64) (*domain.VoteSession, error) {
	session, err := s.voteRepo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Session: %w", err)
	}
	if session == nil {
		return nil, domain.ErrVoteSessionNotFound
	}
	return session, nil
}

// GetActiveSession returns the currently active session for a user
func (s *VoteService) GetActiveSession(ctx context.Context, userID string) (*domain.VoteSession, error) {
	return s.voteRepo.GetActiveSession(ctx, userID)
}

// GetSessions returns all sessions for a user with pagination
func (s *VoteService) GetSessions(ctx context.Context, userID string, limit, offset int) ([]*domain.VoteSession, int64, error) {
	return s.voteRepo.GetSessionsByUser(ctx, userID, limit, offset)
}

// CastVote casts a vote in an active session
func (s *VoteService) CastVote(ctx context.Context, sessionID int64, voterID, option string, weight int) (*domain.Vote, error) {
	session, err := s.voteRepo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Session: %w", err)
	}
	if session == nil {
		return nil, domain.ErrVoteSessionNotFound
	}
	if !session.IsActive() {
		return nil, domain.ErrVoteSessionNotActive
	}

	// Validate option is valid
	if !s.isValidOption(option, session.Options) {
		return nil, domain.ErrInvalidVoteOption
	}

	// Default weight
	if weight == 0 {
		weight = 1
	}

	vote, err := s.voteRepo.CastVote(ctx, domain.VoteCastInput{
		SessionID: sessionID,
		UserID:    voterID,
		Option:    option,
		Weight:    weight,
	})
	if err != nil {
		return nil, err
	}

	// Track analytics for session owner
	_ = s.analyticsRepo.Track(ctx, domain.UsageAnalyticsCreateInput{
		UserID:      session.UserID,
		MetricType:  domain.MetricVoteCast,
		MetricValue: 1,
	})

	// Increment vote counter on n8n integration
	_ = s.n8nSvc.n8nRepo.IncrementVotes(ctx, session.UserID)

	return vote, nil
}

// GetResults returns live vote counts for a session
func (s *VoteService) GetResults(ctx context.Context, sessionID int64) ([]domain.VoteResult, int, error) {
	session, err := s.voteRepo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, 0, fmt.Errorf("fehler beim Laden der Session: %w", err)
	}
	if session == nil {
		return nil, 0, domain.ErrVoteSessionNotFound
	}

	// If session is closed, return stored results
	if !session.IsActive() && len(session.Results) > 0 {
		total := 0
		for _, r := range session.Results {
			total += r.Count
		}
		return session.Results, total, nil
	}

	// Otherwise calculate live
	counts, err := s.voteRepo.GetVoteCounts(ctx, sessionID)
	if err != nil {
		return nil, 0, fmt.Errorf("fehler beim Zählen der Votes: %w", err)
	}

	results := domain.CalculateResults(counts, session.Options)
	total := 0
	for _, count := range counts {
		total += count
	}

	return results, total, nil
}

// CloseSession closes an active session and determines the winner
func (s *VoteService) CloseSession(ctx context.Context, sessionID int64, ownerID string) (*domain.VoteSession, error) {
	session, err := s.voteRepo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Session: %w", err)
	}
	if session == nil {
		return nil, domain.ErrVoteSessionNotFound
	}

	// Only the owner can close
	if session.UserID != ownerID {
		return nil, domain.ErrForbidden
	}

	if !session.IsActive() {
		return nil, domain.ErrVoteSessionClosed
	}

	// Get final counts
	counts, err := s.voteRepo.GetVoteCounts(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Zählen: %w", err)
	}

	results := domain.CalculateResults(counts, session.Options)
	winner := domain.DetermineWinner(results)

	if err := s.voteRepo.CloseSession(ctx, sessionID, winner, results); err != nil {
		return nil, fmt.Errorf("fehler beim Schließen der Session: %w", err)
	}

	log.Printf("🏆 Vote Session %d geschlossen. Gewinner: %s", sessionID, winner)

	// Return updated session
	return s.voteRepo.GetSessionByID(ctx, sessionID)
}

// HasVoted checks if a user has already voted in a session
func (s *VoteService) HasVoted(ctx context.Context, sessionID int64, userID string) (bool, error) {
	return s.voteRepo.HasVoted(ctx, sessionID, userID)
}

func (s *VoteService) isValidOption(option string, options []domain.VoteOption) bool {
	if len(options) == 0 {
		return true // No restriction if no options defined
	}
	for _, opt := range options {
		if opt.Key == option {
			return true
		}
	}
	return false
}