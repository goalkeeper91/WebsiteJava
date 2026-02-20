package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
)

// N8NCommandPayload is what we send to n8n when an advanced command is triggered
type N8NCommandPayload struct {
	Command   string   `json:"command"`
	Args      []string `json:"args"`
	UserID    string   `json:"userId"`
	Username  string   `json:"username"`
	ChannelID string   `json:"channelId"`
	IsSub     bool     `json:"isSub"`
	IsMod     bool     `json:"isMod"`
	IsVIP     bool     `json:"isVip"`
	Timestamp string   `json:"timestamp"`
}

// N8NCommandResponse is what n8n returns after handling a command
type N8NCommandResponse struct {
	Response string                 `json:"response"`
	Actions  []map[string]interface{} `json:"actions,omitempty"`
}

type N8NService struct {
	n8nRepo         repository.N8NIntegrationRepository
	subscriptionSvc *SubscriptionService
	httpClient      *http.Client
	webhookBaseURL  string // Centralized n8n instance URL
}

func NewN8NService(
	n8nRepo repository.N8NIntegrationRepository,
	subscriptionSvc *SubscriptionService,
) *N8NService {
	// Get centralized n8n webhook URL from environment
	webhookBaseURL := os.Getenv("N8N_WEBHOOK_BASE_URL")
	if webhookBaseURL == "" {
		log.Println("⚠️ Warning: N8N_WEBHOOK_BASE_URL not set - n8n integration will not work")
		webhookBaseURL = "http://localhost:5678/webhook"
	}

	return &N8NService{
		n8nRepo:         n8nRepo,
		subscriptionSvc: subscriptionSvc,
		webhookBaseURL:  webhookBaseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetIntegration returns the n8n integration for a user
func (s *N8NService) GetIntegration(ctx context.Context, userID string) (*domain.N8NIntegration, error) {
	integration, err := s.n8nRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der n8n Integration: %w", err)
	}

	// Auto-create entry if none exists
	if integration == nil {
		integration, err = s.n8nRepo.Create(ctx, domain.N8NIntegrationCreateInput{
			UserID: userID,
		})
		if err != nil {
			return nil, fmt.Errorf("fehler beim Erstellen der n8n Integration: %w", err)
		}
	}

	return integration, nil
}

// Enable enables n8n for a user (requires pro/premium subscription)
// For Option A: automatically sets the centralized webhook URL
func (s *N8NService) Enable(ctx context.Context, userID string) (*domain.N8NIntegration, error) {
	// Check subscription
	canUse, err := s.subscriptionSvc.CanUseN8N(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !canUse {
		return nil, domain.ErrFeatureNotAvailable
	}

	integration, err := s.GetIntegration(ctx, userID)
	if err != nil {
		return nil, err
	}

	enabled := true
	// Set centralized webhook URL automatically
	err = s.n8nRepo.Update(ctx, userID, domain.N8NIntegrationUpdateInput{
		Enabled:        &enabled,
		WebhookBaseURL: &s.webhookBaseURL,
	})
	if err != nil {
		return nil, fmt.Errorf("fehler beim Aktivieren der n8n Integration: %w", err)
	}

	log.Printf("✅ n8n Integration aktiviert für User: %s (using centralized n8n at %s)", userID, s.webhookBaseURL)

	// Reload to get updated data
	integration, err = s.n8nRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return integration, nil
}

// Disable disables n8n for a user
func (s *N8NService) Disable(ctx context.Context, userID string) error {
	enabled := false
	err := s.n8nRepo.Update(ctx, userID, domain.N8NIntegrationUpdateInput{
		Enabled: &enabled,
	})
	if err != nil {
		return fmt.Errorf("fehler beim Deaktivieren der n8n Integration: %w", err)
	}

	log.Printf("⚠️ n8n Integration deaktiviert für User: %s", userID)
	return nil
}

// IsReady checks if a user's n8n integration is ready to use
func (s *N8NService) IsReady(ctx context.Context, userID string) (bool, error) {
	integration, err := s.n8nRepo.GetByUserID(ctx, userID)
	if err != nil {
		return false, err
	}
	if integration == nil {
		return false, nil
	}
	return integration.IsReady(), nil
}

// TriggerCommandWebhook sends an advanced command event to n8n
func (s *N8NService) TriggerCommandWebhook(ctx context.Context, userID string, payload N8NCommandPayload) (*N8NCommandResponse, error) {
	integration, err := s.n8nRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if integration == nil || !integration.IsReady() {
		return nil, domain.ErrN8NIntegrationNotReady
	}

	// Check monthly reset
	if integration.ShouldResetMonthlyCounter() {
		if err := s.n8nRepo.ResetMonthlyCounter(ctx, userID); err != nil {
			log.Printf("⚠️ Fehler beim Zurücksetzen des Counters: %v", err)
		}
	}

	// Use centralized webhook URL with user-specific path
	webhookURL := fmt.Sprintf("%s/%s/command-handler", s.webhookBaseURL, userID)
	return s.callWebhook(ctx, webhookURL, payload)
}

// callWebhook is the generic HTTP call to any n8n webhook
func (s *N8NService) callWebhook(ctx context.Context, webhookURL string, payload interface{}) (*N8NCommandResponse, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Serialisieren des Payloads: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhookURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("fehler beim Erstellen des Requests: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrN8NWebhookFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, domain.ErrN8NWorkflowNotFound
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%w: status=%d body=%s", domain.ErrN8NWebhookFailed, resp.StatusCode, string(body))
	}

	// Parse response (n8n might return empty body or JSON)
	var result N8NCommandResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		// Empty response is OK
		return &result, nil
	}

	return &result, nil
}

// TestWebhook tests if n8n is reachable (centralized health check)
func (s *N8NService) TestWebhook(ctx context.Context, userID string) error {
	integration, err := s.n8nRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if integration == nil || !integration.Enabled {
		return domain.ErrN8NIntegrationNotReady
	}

	// Test centralized n8n health endpoint
	testURL := fmt.Sprintf("%s/healthz", s.webhookBaseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, testURL, nil)
	if err != nil {
		return fmt.Errorf("fehler beim Erstellen des Test-Requests: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %v", domain.ErrN8NWebhookFailed, err)
	}
	defer resp.Body.Close()

	log.Printf("✅ n8n Webhook Test für User %s: status=%d", userID, resp.StatusCode)
	return nil
}