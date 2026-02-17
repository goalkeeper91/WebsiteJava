package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
)

type WorkflowTemplateRepository struct {
	db *sql.DB
}

func NewWorkflowTemplateRepository(db *sql.DB) repository.WorkflowTemplateRepository {
	return &WorkflowTemplateRepository{db: db}
}

func (r *WorkflowTemplateRepository) Create(ctx context.Context, input domain.WorkflowTemplateCreateInput) (*domain.WorkflowTemplate, error) {
	query := `
		INSERT INTO workflow_templates
		(name, description, category, n8n_workflow_json, tier_required, is_public, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id, name, description, category, n8n_workflow_json,
		          tier_required, is_public, usage_count, created_by, created_at, updated_at
	`

	var t domain.WorkflowTemplate
	err := r.db.QueryRowContext(
		ctx, query,
		input.Name,
		input.Description,
		input.Category,
		input.N8NWorkflowJSON,
		input.TierRequired,
		input.IsPublic,
		input.CreatedBy,
	).Scan(
		&t.ID,
		&t.Name,
		&t.Description,
		&t.Category,
		&t.N8NWorkflowJSON,
		&t.TierRequired,
		&t.IsPublic,
		&t.UsageCount,
		&t.CreatedBy,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Erstellen des Workflow Templates: %w", err)
	}

	return &t, nil
}

func (r *WorkflowTemplateRepository) GetByID(ctx context.Context, id int64) (*domain.WorkflowTemplate, error) {
	query := `
		SELECT id, name, description, category, n8n_workflow_json,
		       tier_required, is_public, usage_count, created_by, created_at, updated_at
		FROM workflow_templates
		WHERE id = $1
	`

	var t domain.WorkflowTemplate
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID,
		&t.Name,
		&t.Description,
		&t.Category,
		&t.N8NWorkflowJSON,
		&t.TierRequired,
		&t.IsPublic,
		&t.UsageCount,
		&t.CreatedBy,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden des Templates: %w", err)
	}

	return &t, nil
}

func (r *WorkflowTemplateRepository) GetAll(ctx context.Context) ([]*domain.WorkflowTemplate, error) {
	return r.queryTemplates(ctx, `
		SELECT id, name, description, category, n8n_workflow_json,
		       tier_required, is_public, usage_count, created_by, created_at, updated_at
		FROM workflow_templates
		ORDER BY category, name
	`)
}

func (r *WorkflowTemplateRepository) GetByCategory(ctx context.Context, category domain.WorkflowCategory) ([]*domain.WorkflowTemplate, error) {
	return r.queryTemplates(ctx, `
		SELECT id, name, description, category, n8n_workflow_json,
		       tier_required, is_public, usage_count, created_by, created_at, updated_at
		FROM workflow_templates
		WHERE category = $1
		ORDER BY name
	`, category)
}

func (r *WorkflowTemplateRepository) GetByTier(ctx context.Context, tierID domain.TierID) ([]*domain.WorkflowTemplate, error) {
	return r.queryTemplates(ctx, `
		SELECT id, name, description, category, n8n_workflow_json,
		       tier_required, is_public, usage_count, created_by, created_at, updated_at
		FROM workflow_templates
		WHERE tier_required = $1
		ORDER BY name
	`, tierID)
}

func (r *WorkflowTemplateRepository) GetPublic(ctx context.Context) ([]*domain.WorkflowTemplate, error) {
	return r.queryTemplates(ctx, `
		SELECT id, name, description, category, n8n_workflow_json,
		       tier_required, is_public, usage_count, created_by, created_at, updated_at
		FROM workflow_templates
		WHERE is_public = true
		ORDER BY usage_count DESC, name
	`)
}

func (r *WorkflowTemplateRepository) GetAccessibleByTier(ctx context.Context, tierID domain.TierID) ([]*domain.WorkflowTemplate, error) {
	// Returns templates accessible for the given tier (e.g. premium sees pro + free templates too)
	var allowedTiers []domain.TierID
	switch tierID {
	case domain.TierPremium:
		allowedTiers = []domain.TierID{domain.TierFree, domain.TierPro, domain.TierPremium}
	case domain.TierPro:
		allowedTiers = []domain.TierID{domain.TierFree, domain.TierPro}
	default:
		allowedTiers = []domain.TierID{domain.TierFree}
	}

	query := `
		SELECT id, name, description, category, n8n_workflow_json,
		       tier_required, is_public, usage_count, created_by, created_at, updated_at
		FROM workflow_templates
		WHERE tier_required = ANY($1) AND is_public = true
		ORDER BY tier_required, usage_count DESC
	`

	return r.queryTemplates(ctx, query, allowedTiers)
}

func (r *WorkflowTemplateRepository) Update(ctx context.Context, id int64, input domain.WorkflowTemplateUpdateInput) error {
	query := `
		UPDATE workflow_templates
		SET name = COALESCE($2, name),
		    description = COALESCE($3, description),
		    n8n_workflow_json = COALESCE($4, n8n_workflow_json),
		    tier_required = COALESCE($5, tier_required),
		    is_public = COALESCE($6, is_public),
		    updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.ExecContext(
		ctx, query,
		id,
		input.Name,
		input.Description,
		input.N8NWorkflowJSON,
		input.TierRequired,
		input.IsPublic,
	)
	if err != nil {
		return fmt.Errorf("fehler beim Update des Templates: %w", err)
	}

	return nil
}

func (r *WorkflowTemplateRepository) IncrementUsage(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE workflow_templates SET usage_count = usage_count + 1 WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("fehler beim Inkrementieren des Usage Counters: %w", err)
	}
	return nil
}

func (r *WorkflowTemplateRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM workflow_templates WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("fehler beim Löschen des Templates: %w", err)
	}
	return nil
}

// queryTemplates is a helper for common SELECT queries
func (r *WorkflowTemplateRepository) queryTemplates(ctx context.Context, query string, args ...interface{}) ([]*domain.WorkflowTemplate, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Templates: %w", err)
	}
	defer rows.Close()

	var templates []*domain.WorkflowTemplate
	for rows.Next() {
		var t domain.WorkflowTemplate
		err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.Description,
			&t.Category,
			&t.N8NWorkflowJSON,
			&t.TierRequired,
			&t.IsPublic,
			&t.UsageCount,
			&t.CreatedBy,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("fehler beim Scannen des Templates: %w", err)
		}
		templates = append(templates, &t)
	}

	return templates, nil
}