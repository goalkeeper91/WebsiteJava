package service

import (
	"context"
	"fmt"
	"log"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
)

type WorkflowTemplateService struct {
	templateRepo    repository.WorkflowTemplateRepository
	subscriptionSvc *SubscriptionService
}

func NewWorkflowTemplateService(
	templateRepo repository.WorkflowTemplateRepository,
	subscriptionSvc *SubscriptionService,
) *WorkflowTemplateService {
	return &WorkflowTemplateService{
		templateRepo:    templateRepo,
		subscriptionSvc: subscriptionSvc,
	}
}

// GetAvailableTemplates returns templates accessible by the user's tier
func (s *WorkflowTemplateService) GetAvailableTemplates(ctx context.Context, userID string) ([]*domain.WorkflowTemplate, error) {
	sub, err := s.subscriptionSvc.GetSubscription(ctx, userID)
	if err != nil {
		return nil, err
	}

	templates, err := s.templateRepo.GetAccessibleByTier(ctx, sub.TierID)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Templates: %w", err)
	}

	return templates, nil
}

// GetTemplate returns a specific template if the user's tier allows it
func (s *WorkflowTemplateService) GetTemplate(ctx context.Context, userID string, templateID int64) (*domain.WorkflowTemplate, error) {
	template, err := s.templateRepo.GetByID(ctx, templateID)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden des Templates: %w", err)
	}
	if template == nil {
		return nil, domain.ErrTemplateNotFound
	}

	// Check tier access
	sub, err := s.subscriptionSvc.GetSubscription(ctx, userID)
	if err != nil {
		return nil, err
	}

	if !template.IsAccessibleBy(sub.TierID) {
		return nil, domain.ErrTemplateNotAllowed
	}

	// Track usage
	_ = s.templateRepo.IncrementUsage(ctx, templateID)

	return template, nil
}

// GetTemplatesByCategory returns templates by category
func (s *WorkflowTemplateService) GetTemplatesByCategory(ctx context.Context, userID string, category domain.WorkflowCategory) ([]*domain.WorkflowTemplate, error) {
	sub, err := s.subscriptionSvc.GetSubscription(ctx, userID)
	if err != nil {
		return nil, err
	}

	all, err := s.templateRepo.GetByCategory(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Templates: %w", err)
	}

	// Filter by accessible tiers
	accessible := make([]*domain.WorkflowTemplate, 0)
	for _, t := range all {
		if t.IsAccessibleBy(sub.TierID) {
			accessible = append(accessible, t)
		}
	}

	return accessible, nil
}

// GetPublicTemplates returns all public templates with summaries (no workflow JSON)
func (s *WorkflowTemplateService) GetPublicTemplates(ctx context.Context) ([]*domain.WorkflowTemplateSummary, error) {
	templates, err := s.templateRepo.GetPublic(ctx)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der öffentlichen Templates: %w", err)
	}

	summaries := make([]*domain.WorkflowTemplateSummary, 0, len(templates))
	for _, t := range templates {
		summary := t.ToSummary()
		summaries = append(summaries, &summary)
	}

	return summaries, nil
}

// CreateTemplate creates a new workflow template (admin only)
func (s *WorkflowTemplateService) CreateTemplate(ctx context.Context, input domain.WorkflowTemplateCreateInput) (*domain.WorkflowTemplate, error) {
	template, err := s.templateRepo.Create(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Erstellen des Templates: %w", err)
	}

	log.Printf("✅ Workflow Template erstellt: %s (Tier: %s)", template.Name, template.TierRequired)
	return template, nil
}