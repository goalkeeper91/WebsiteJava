package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
)

type VoteRepository struct {
	db *sql.DB
}

func NewVoteRepository(db *sql.DB) repository.VoteRepository {
	return &VoteRepository{db: db}
}

// =====================================
// SESSION METHODS
// =====================================

func (r *VoteRepository) CreateSession(ctx context.Context, input domain.VoteSessionCreateInput) (*domain.VoteSession, error) {
	optionsJSON, err := json.Marshal(input.Options)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Serialisieren der Optionen: %w", err)
	}

	query := `
		INSERT INTO vote_sessions
		(user_id, category, title, description, status, options, total_votes, started_at)
		VALUES ($1, $2, $3, $4, 'active', $5, 0, NOW())
		RETURNING id, user_id, category, title, description,
		          status, started_at, ended_at, winner, total_votes,
		          options, results
	`

	var session domain.VoteSession
	err = r.db.QueryRowContext(
		ctx, query,
		input.UserID,
		input.Category,
		input.Title,
		input.Description,
		optionsJSON,
	).Scan(
		&session.ID,
		&session.UserID,
		&session.Category,
		&session.Title,
		&session.Description,
		&session.Status,
		&session.StartedAt,
		&session.EndedAt,
		&session.Winner,
		&session.TotalVotes,
		&session.OptionsJSON,
		&session.ResultsJSON,
	)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Erstellen der Vote Session: %w", err)
	}

	if err := session.ParseOptions(); err != nil {
		return nil, fmt.Errorf("fehler beim Parsen der Optionen: %w", err)
	}

	return &session, nil
}

func (r *VoteRepository) GetSessionByID(ctx context.Context, id int64) (*domain.VoteSession, error) {
	query := `
		SELECT id, user_id, category, title, description,
		       status, started_at, ended_at, winner, total_votes,
		       options, results
		FROM vote_sessions
		WHERE id = $1
	`

	var session domain.VoteSession
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&session.ID,
		&session.UserID,
		&session.Category,
		&session.Title,
		&session.Description,
		&session.Status,
		&session.StartedAt,
		&session.EndedAt,
		&session.Winner,
		&session.TotalVotes,
		&session.OptionsJSON,
		&session.ResultsJSON,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Vote Session: %w", err)
	}

	if err := session.ParseOptions(); err != nil {
		return nil, err
	}
	if err := session.ParseResults(); err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *VoteRepository) GetSessionsByUser(ctx context.Context, userID string, limit, offset int) ([]*domain.VoteSession, int64, error) {
	var total int64
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM vote_sessions WHERE user_id = $1`, userID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("fehler beim Zählen der Sessions: %w", err)
	}

	query := `
		SELECT id, user_id, category, title, description,
		       status, started_at, ended_at, winner, total_votes,
		       options, results
		FROM vote_sessions
		WHERE user_id = $1
		ORDER BY started_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("fehler beim Laden der Sessions: %w", err)
	}
	defer rows.Close()

	var sessions []*domain.VoteSession
	for rows.Next() {
		var session domain.VoteSession
		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.Category,
			&session.Title,
			&session.Description,
			&session.Status,
			&session.StartedAt,
			&session.EndedAt,
			&session.Winner,
			&session.TotalVotes,
			&session.OptionsJSON,
			&session.ResultsJSON,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("fehler beim Scannen der Session: %w", err)
		}
		session.ParseOptions()
		session.ParseResults()
		sessions = append(sessions, &session)
	}

	return sessions, total, nil
}

func (r *VoteRepository) GetActiveSession(ctx context.Context, userID string) (*domain.VoteSession, error) {
	query := `
		SELECT id, user_id, category, title, description,
		       status, started_at, ended_at, winner, total_votes,
		       options, results
		FROM vote_sessions
		WHERE user_id = $1 AND status = 'active'
		ORDER BY started_at DESC
		LIMIT 1
	`

	var session domain.VoteSession
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&session.ID,
		&session.UserID,
		&session.Category,
		&session.Title,
		&session.Description,
		&session.Status,
		&session.StartedAt,
		&session.EndedAt,
		&session.Winner,
		&session.TotalVotes,
		&session.OptionsJSON,
		&session.ResultsJSON,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der aktiven Session: %w", err)
	}

	session.ParseOptions()
	return &session, nil
}

func (r *VoteRepository) UpdateSession(ctx context.Context, id int64, input domain.VoteSessionUpdateInput) error {
	query := `
		UPDATE vote_sessions
		SET title = COALESCE($2, title),
		    description = COALESCE($3, description),
		    status = COALESCE($4, status),
		    ended_at = COALESCE($5, ended_at),
		    winner = COALESCE($6, winner),
		    total_votes = COALESCE($7, total_votes)
		WHERE id = $1
	`

	_, err := r.db.ExecContext(
		ctx, query,
		id,
		input.Title,
		input.Description,
		input.Status,
		input.EndedAt,
		input.Winner,
		input.TotalVotes,
	)
	if err != nil {
		return fmt.Errorf("fehler beim Update der Session: %w", err)
	}

	return nil
}

func (r *VoteRepository) CloseSession(ctx context.Context, id int64, winner string, results []domain.VoteResult) error {
	resultsJSON, err := json.Marshal(results)
	if err != nil {
		return fmt.Errorf("fehler beim Serialisieren der Ergebnisse: %w", err)
	}

	query := `
		UPDATE vote_sessions
		SET status = 'closed',
		    ended_at = NOW(),
		    winner = $2,
		    results = $3,
		    total_votes = (SELECT COUNT(*) FROM votes WHERE session_id = $1)
		WHERE id = $1
	`

	_, err = r.db.ExecContext(ctx, query, id, winner, resultsJSON)
	if err != nil {
		return fmt.Errorf("fehler beim Schließen der Session: %w", err)
	}

	return nil
}

// =====================================
// VOTE METHODS
// =====================================

func (r *VoteRepository) CastVote(ctx context.Context, input domain.VoteCastInput) (*domain.Vote, error) {
	weight := input.Weight
	if weight == 0 {
		weight = 1
	}

	query := `
		INSERT INTO votes (session_id, user_id, option, weight, voted_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (session_id, user_id) DO NOTHING
		RETURNING id, session_id, user_id, option, weight, voted_at
	`

	var vote domain.Vote
	err := r.db.QueryRowContext(
		ctx, query,
		input.SessionID,
		input.UserID,
		input.Option,
		weight,
	).Scan(
		&vote.ID,
		&vote.SessionID,
		&vote.UserID,
		&vote.Option,
		&vote.Weight,
		&vote.VotedAt,
	)
	if err == sql.ErrNoRows {
		// User already voted (ON CONFLICT DO NOTHING)
		return nil, domain.ErrAlreadyVoted
	}
	if err != nil {
		return nil, fmt.Errorf("fehler beim Abstimmen: %w", err)
	}

	return &vote, nil
}

func (r *VoteRepository) GetVotesBySession(ctx context.Context, sessionID int64) ([]*domain.Vote, error) {
	query := `
		SELECT id, session_id, user_id, option, weight, voted_at
		FROM votes
		WHERE session_id = $1
		ORDER BY voted_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Votes: %w", err)
	}
	defer rows.Close()

	var votes []*domain.Vote
	for rows.Next() {
		var vote domain.Vote
		err := rows.Scan(
			&vote.ID,
			&vote.SessionID,
			&vote.UserID,
			&vote.Option,
			&vote.Weight,
			&vote.VotedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("fehler beim Scannen des Votes: %w", err)
		}
		votes = append(votes, &vote)
	}

	return votes, nil
}

func (r *VoteRepository) GetVoteCounts(ctx context.Context, sessionID int64) (map[string]int, error) {
	query := `
		SELECT option, SUM(weight) as total
		FROM votes
		WHERE session_id = $1
		GROUP BY option
		ORDER BY total DESC
	`

	rows, err := r.db.QueryContext(ctx, query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Zählen der Votes: %w", err)
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var option string
		var total int
		if err := rows.Scan(&option, &total); err != nil {
			return nil, fmt.Errorf("fehler beim Scannen: %w", err)
		}
		counts[option] = total
	}

	return counts, nil
}

func (r *VoteRepository) HasVoted(ctx context.Context, sessionID int64, userID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM votes WHERE session_id = $1 AND user_id = $2)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, sessionID, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("fehler beim Prüfen des Votes: %w", err)
	}

	return exists, nil
}