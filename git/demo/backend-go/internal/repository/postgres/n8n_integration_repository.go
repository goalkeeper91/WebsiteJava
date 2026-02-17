package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
)

type N8NIntegrationRepository struct {
	db *sql.DB
}

func NewN8NIntegrationRepository(db *sql.DB) repository.N8NIntegrationRepository {
	return &N8NIntegrationRepository{db: db}
}

func (r *N8NIntegrationRepository) Create(ctx context.Context, input domain.N8NIntegrationCreateInput) (*domain.N8NIntegration, error) {
	query := `
		INSERT INTO n8n_integrations
		(user_id, enabled, webhook_base_url, api_key, last_reset_at, created_at, updated_at)
		VALUES ($1, false, $2, $3, NOW(), NOW(), NOW())
		RETURNING id, user_id, enabled, webhook_base_url, api_key,
		          workflows_used, votes_this_month, last_reset_at,
		          created_at, updated_at
	`

	var integration domain.N8NIntegration
	err := r.db.QueryRowContext(
		ctx, query,
		input.UserID,
		input.WebhookBaseURL,
		input.APIKey,
	).Scan(
		&integration.ID,
		&integration.UserID,
		&integration.Enabled,
		&integration.WebhookBaseURL,
		&integration.APIKey,
		&integration.WorkflowsUsed,
		&integration.VotesThisMonth,
		&integration.LastResetAt,
		&integration.CreatedAt,
		&integration.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Erstellen der n8n Integration: %w", err)
	}

	return &integration, nil
}

func (r *N8NIntegrationRepository) GetByUserID(ctx context.Context, userID string) (*domain.N8NIntegration, error) {
	query := `
		SELECT id, user_id, enabled, webhook_base_url, api_key,
		       workflows_used, votes_this_month, last_reset_at,
		       created_at, updated_at
		FROM n8n_integrations
		WHERE user_id = $1
	`

	var integration domain.N8NIntegration
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&integration.ID,
		&integration.UserID,
		&integration.Enabled,
		&integration.WebhookBaseURL,
		&integration.APIKey,
		&integration.WorkflowsUsed,
		&integration.VotesThisMonth,
		&integration.LastResetAt,
		&integration.CreatedAt,
		&integration.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der n8n Integration: %w", err)
	}

	return &integration, nil
}

func (r *N8NIntegrationRepository) Update(ctx context.Context, userID string, input domain.N8NIntegrationUpdateInput) error {
	query := `
		UPDATE n8n_integrations
		SET enabled = COALESCE($2, enabled),
		    webhook_base_url = COALESCE($3, webhook_base_url),
		    api_key = COALESCE($4, api_key),
		    updated_at = NOW()
		WHERE user_id = $1
	`

	_, err := r.db.ExecContext(
		ctx, query,
		userID,
		input.Enabled,
		input.WebhookBaseURL,
		input.APIKey,
	)
	if err != nil {
		return fmt.Errorf("fehler beim Update der n8n Integration: %w", err)
	}

	return nil
}

func (r *N8NIntegrationRepository) IncrementVotes(ctx context.Context, userID string) error {
	query := `
		UPDATE n8n_integrations
		SET votes_this_month = votes_this_month + 1, updated_at = NOW()
		WHERE user_id = $1
	`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("fehler beim Inkrementieren der Votes: %w", err)
	}

	return nil
}

func (r *N8NIntegrationRepository) ResetMonthlyCounter(ctx context.Context, userID string) error {
	query := `
		UPDATE n8n_integrations
		SET votes_this_month = 0, last_reset_at = NOW(), updated_at = NOW()
		WHERE user_id = $1
	`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("fehler beim Zurücksetzen des Monats-Counters: %w", err)
	}

	return nil
}

func (r *N8NIntegrationRepository) Exists(ctx context.Context, userID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM n8n_integrations WHERE user_id = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("fehler beim Prüfen der n8n Integration: %w", err)
	}

	return exists, nil
}