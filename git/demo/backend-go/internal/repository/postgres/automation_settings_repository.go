package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
)

var (
	ErrAutomationSettingsNotFound = errors.New("automation settings not found")
)

type automationSettingsRepository struct {
	db *sql.DB
}

func NewAutomationSettingsRepository(db *sql.DB) repository.AutomationSettingsRepository {
	return &automationSettingsRepository{db: db}
}

func (r *automationSettingsRepository) Create(ctx context.Context, settings *domain.AutomationSettings) error {
	query := `
		INSERT INTO automation_settings
		(user_twitch_id, is_enabled, ai_tone, ai_style, use_hashtags, 
		 max_hashtags, min_clip_duration, max_clip_duration, min_view_count,
		 only_verified_clips, target_platforms, auto_post_enabled, 
		 posting_frequency, preferred_posting_times, blocked_words,
		 allowed_games, notify_on_post, notify_on_error, 
		 created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		settings.UserTwitchID,
		settings.IsEnabled,
		settings.AITone,
		settings.AIStyle,
		settings.UseHashtags,
		settings.MaxHashtags,
		settings.MinClipDuration,
		settings.MaxClipDuration,
		settings.MinViewCount,
		settings.OnlyVerifiedClips,
		settings.TargetPlatforms,
		settings.AutoPostEnabled,
		settings.PostingFrequency,
		settings.PreferredPostingTimes,
		settings.BlockedWords,
		settings.AllowedGames,
		settings.NotifyOnPost,
		settings.NotifyOnError,
		settings.CreatedAt,
		settings.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("fehler beim Erstellen der Automation Settings: %w", err)
	}

	return nil
}

func (r *automationSettingsRepository) GetByUser(ctx context.Context, userTwitchID string) (*domain.AutomationSettings, error) {
	query := `
		SELECT user_twitch_id, is_enabled, ai_tone, ai_style, use_hashtags,
		       max_hashtags, min_clip_duration, max_clip_duration, min_view_count,
		       only_verified_clips, target_platforms, auto_post_enabled,
		       posting_frequency, preferred_posting_times, blocked_words,
		       allowed_games, notify_on_post, notify_on_error,
		       created_at, updated_at
		FROM automation_settings
		WHERE user_twitch_id = $1
	`

	var settings domain.AutomationSettings

	err := r.db.QueryRowContext(ctx, query, userTwitchID).Scan(
		&settings.UserTwitchID,
		&settings.IsEnabled,
		&settings.AITone,
		&settings.AIStyle,
		&settings.UseHashtags,
		&settings.MaxHashtags,
		&settings.MinClipDuration,
		&settings.MaxClipDuration,
		&settings.MinViewCount,
		&settings.OnlyVerifiedClips,
		&settings.TargetPlatforms,
		&settings.AutoPostEnabled,
		&settings.PostingFrequency,
		&settings.PreferredPostingTimes,
		&settings.BlockedWords,
		&settings.AllowedGames,
		&settings.NotifyOnPost,
		&settings.NotifyOnError,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrAutomationSettingsNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Automation Settings: %w", err)
	}

	return &settings, nil
}

func (r *automationSettingsRepository) Update(ctx context.Context, settings *domain.AutomationSettings) error {
	query := `
		UPDATE automation_settings
		SET is_enabled = $2,
		    ai_tone = $3,
		    ai_style = $4,
		    use_hashtags = $5,
		    max_hashtags = $6,
		    min_clip_duration = $7,
		    max_clip_duration = $8,
		    min_view_count = $9,
		    only_verified_clips = $10,
		    target_platforms = $11,
		    auto_post_enabled = $12,
		    posting_frequency = $13,
		    preferred_posting_times = $14,
		    blocked_words = $15,
		    allowed_games = $16,
		    notify_on_post = $17,
		    notify_on_error = $18,
		    updated_at = $19
		WHERE user_twitch_id = $1
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		settings.UserTwitchID,
		settings.IsEnabled,
		settings.AITone,
		settings.AIStyle,
		settings.UseHashtags,
		settings.MaxHashtags,
		settings.MinClipDuration,
		settings.MaxClipDuration,
		settings.MinViewCount,
		settings.OnlyVerifiedClips,
		settings.TargetPlatforms,
		settings.AutoPostEnabled,
		settings.PostingFrequency,
		settings.PreferredPostingTimes,
		settings.BlockedWords,
		settings.AllowedGames,
		settings.NotifyOnPost,
		settings.NotifyOnError,
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("fehler beim Update der Automation Settings: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrAutomationSettingsNotFound
	}

	return nil
}

func (r *automationSettingsRepository) Upsert(ctx context.Context, settings *domain.AutomationSettings) error {
	query := `
		INSERT INTO automation_settings
		(user_twitch_id, is_enabled, ai_tone, ai_style, use_hashtags, 
		 max_hashtags, min_clip_duration, max_clip_duration, min_view_count,
		 only_verified_clips, target_platforms, auto_post_enabled, 
		 posting_frequency, preferred_posting_times, blocked_words,
		 allowed_games, notify_on_post, notify_on_error, 
		 created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
		ON CONFLICT (user_twitch_id)
		DO UPDATE SET
		    is_enabled = EXCLUDED.is_enabled,
		    ai_tone = EXCLUDED.ai_tone,
		    ai_style = EXCLUDED.ai_style,
		    use_hashtags = EXCLUDED.use_hashtags,
		    max_hashtags = EXCLUDED.max_hashtags,
		    min_clip_duration = EXCLUDED.min_clip_duration,
		    max_clip_duration = EXCLUDED.max_clip_duration,
		    min_view_count = EXCLUDED.min_view_count,
		    only_verified_clips = EXCLUDED.only_verified_clips,
		    target_platforms = EXCLUDED.target_platforms,
		    auto_post_enabled = EXCLUDED.auto_post_enabled,
		    posting_frequency = EXCLUDED.posting_frequency,
		    preferred_posting_times = EXCLUDED.preferred_posting_times,
		    blocked_words = EXCLUDED.blocked_words,
		    allowed_games = EXCLUDED.allowed_games,
		    notify_on_post = EXCLUDED.notify_on_post,
		    notify_on_error = EXCLUDED.notify_on_error,
		    updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		settings.UserTwitchID,
		settings.IsEnabled,
		settings.AITone,
		settings.AIStyle,
		settings.UseHashtags,
		settings.MaxHashtags,
		settings.MinClipDuration,
		settings.MaxClipDuration,
		settings.MinViewCount,
		settings.OnlyVerifiedClips,
		settings.TargetPlatforms,
		settings.AutoPostEnabled,
		settings.PostingFrequency,
		settings.PreferredPostingTimes,
		settings.BlockedWords,
		settings.AllowedGames,
		settings.NotifyOnPost,
		settings.NotifyOnError,
		settings.CreatedAt,
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("fehler beim Upsert der Automation Settings: %w", err)
	}

	return nil
}

func (r *automationSettingsRepository) Delete(ctx context.Context, userTwitchID string) error {
	query := `DELETE FROM automation_settings WHERE user_twitch_id = $1`

	result, err := r.db.ExecContext(ctx, query, userTwitchID)
	if err != nil {
		return fmt.Errorf("fehler beim Löschen: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrAutomationSettingsNotFound
	}

	return nil
}

func (r *automationSettingsRepository) GetAllEnabled(ctx context.Context) ([]*domain.AutomationSettings, error) {
	query := `
		SELECT user_twitch_id, is_enabled, ai_tone, ai_style, use_hashtags,
		       max_hashtags, min_clip_duration, max_clip_duration, min_view_count,
		       only_verified_clips, target_platforms, auto_post_enabled,
		       posting_frequency, preferred_posting_times, blocked_words,
		       allowed_games, notify_on_post, notify_on_error,
		       created_at, updated_at
		FROM automation_settings
		WHERE is_enabled = true
		ORDER BY created_at DESC
	`

	return r.queryMultiple(ctx, query)
}

func (r *automationSettingsRepository) GetAllEnabledWithSubscription(ctx context.Context) ([]*domain.AutomationSettingsWithSubscription, error) {
	query := `
		SELECT a.user_twitch_id, a.is_enabled, a.ai_tone, a.ai_style, a.use_hashtags,
		       a.max_hashtags, a.min_clip_duration, a.max_clip_duration, a.min_view_count,
		       a.only_verified_clips, a.target_platforms, a.auto_post_enabled,
		       a.posting_frequency, a.preferred_posting_times, a.blocked_words,
		       a.allowed_games, a.notify_on_post, a.notify_on_error,
		       a.created_at, a.updated_at,
		       s.id, s.user_id, s.tier_id, s.status, s.billing_cycle,
		       s.started_at, s.expires_at, s.canceled_at,
		       s.paddle_customer_id, s.paddle_subscription_id,
		       s.created_at, s.updated_at,
		       t.id, t.name, t.price_monthly, t.price_yearly,
		       t.max_commands, t.max_workflows, t.max_votes_per_month,
		       t.features, t.is_active, t.created_at,
		       u.is_admin
		FROM automation_settings a
		INNER JOIN user_subscriptions s ON s.user_id = a.user_twitch_id
		INNER JOIN subscription_tiers t ON s.tier_id = t.id
		INNER JOIN users u ON u.twitch_id = a.user_twitch_id
		WHERE a.is_enabled = true
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Query: %w", err)
	}
	defer rows.Close()

	var result []*domain.AutomationSettingsWithSubscription
	for rows.Next() {
		var settings domain.AutomationSettings
		var sub domain.UserSubscription
		var tier domain.SubscriptionTier
		var featuresJSON []byte
		var isAdmin bool

		err := rows.Scan(
			&settings.UserTwitchID,
			&settings.IsEnabled,
			&settings.AITone,
			&settings.AIStyle,
			&settings.UseHashtags,
			&settings.MaxHashtags,
			&settings.MinClipDuration,
			&settings.MaxClipDuration,
			&settings.MinViewCount,
			&settings.OnlyVerifiedClips,
			&settings.TargetPlatforms,
			&settings.AutoPostEnabled,
			&settings.PostingFrequency,
			&settings.PreferredPostingTimes,
			&settings.BlockedWords,
			&settings.AllowedGames,
			&settings.NotifyOnPost,
			&settings.NotifyOnError,
			&settings.CreatedAt,
			&settings.UpdatedAt,
			&sub.ID,
			&sub.UserID,
			&sub.TierID,
			&sub.Status,
			&sub.BillingCycle,
			&sub.StartedAt,
			&sub.ExpiresAt,
			&sub.CanceledAt,
			&sub.PaddleCustomerID,
			&sub.PaddleSubscriptionID,
			&sub.CreatedAt,
			&sub.UpdatedAt,
			&tier.ID,
			&tier.Name,
			&tier.PriceMonthly,
			&tier.PriceYearly,
			&tier.MaxCommands,
			&tier.MaxWorkflows,
			&tier.MaxVotesPerMonth,
			&featuresJSON,
			&tier.IsActive,
			&tier.CreatedAt,
			&isAdmin,
		)
		if err != nil {
			return nil, fmt.Errorf("fehler beim Scan: %w", err)
		}

		if err := json.Unmarshal(featuresJSON, &tier.Features); err != nil {
			return nil, fmt.Errorf("fehler beim Parsen der Features: %w", err)
		}

		sub.Tier = &tier
		result = append(result, &domain.AutomationSettingsWithSubscription{
			Settings:     &settings,
			Subscription: &sub,
			IsAdmin:      isAdmin,
		})
	}

	return result, rows.Err()
}

func (r *automationSettingsRepository) GetAllEnabledWithAutoPost(ctx context.Context) ([]*domain.AutomationSettings, error) {
	query := `
		SELECT user_twitch_id, is_enabled, ai_tone, ai_style, use_hashtags,
		       max_hashtags, min_clip_duration, max_clip_duration, min_view_count,
		       only_verified_clips, target_platforms, auto_post_enabled,
		       posting_frequency, preferred_posting_times, blocked_words,
		       allowed_games, notify_on_post, notify_on_error,
		       created_at, updated_at
		FROM automation_settings
		WHERE is_enabled = true AND auto_post_enabled = true
		ORDER BY created_at DESC
	`

	return r.queryMultiple(ctx, query)
}

func (r *automationSettingsRepository) queryMultiple(ctx context.Context, query string, args ...interface{}) ([]*domain.AutomationSettings, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Query: %w", err)
	}
	defer rows.Close()

	var settingsList []*domain.AutomationSettings
	for rows.Next() {
		var settings domain.AutomationSettings
		err := rows.Scan(
			&settings.UserTwitchID,
			&settings.IsEnabled,
			&settings.AITone,
			&settings.AIStyle,
			&settings.UseHashtags,
			&settings.MaxHashtags,
			&settings.MinClipDuration,
			&settings.MaxClipDuration,
			&settings.MinViewCount,
			&settings.OnlyVerifiedClips,
			&settings.TargetPlatforms,
			&settings.AutoPostEnabled,
			&settings.PostingFrequency,
			&settings.PreferredPostingTimes,
			&settings.BlockedWords,
			&settings.AllowedGames,
			&settings.NotifyOnPost,
			&settings.NotifyOnError,
			&settings.CreatedAt,
			&settings.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("fehler beim Scan: %w", err)
		}
		settingsList = append(settingsList, &settings)
	}

	return settingsList, rows.Err()
}

func (r *automationSettingsRepository) Enable(ctx context.Context, userTwitchID string) error {
	return r.setEnabled(ctx, userTwitchID, true)
}

func (r *automationSettingsRepository) Disable(ctx context.Context, userTwitchID string) error {
	return r.setEnabled(ctx, userTwitchID, false)
}

func (r *automationSettingsRepository) setEnabled(ctx context.Context, userTwitchID string, enabled bool) error {
	query := `
		UPDATE automation_settings
		SET is_enabled = $2, updated_at = $3
		WHERE user_twitch_id = $1
	`

	result, err := r.db.ExecContext(ctx, query, userTwitchID, enabled, time.Now())
	if err != nil {
		return fmt.Errorf("fehler beim Update: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrAutomationSettingsNotFound
	}

	return nil
}

func (r *automationSettingsRepository) UpdateTargetPlatforms(ctx context.Context, userTwitchID string, platforms []domain.Platform) error {
	query := `
		UPDATE automation_settings
		SET target_platforms = $2, updated_at = $3
		WHERE user_twitch_id = $1
	`

	platformArray := domain.PlatformArray(platforms)

	result, err := r.db.ExecContext(ctx, query, userTwitchID, platformArray, time.Now())
	if err != nil {
		return fmt.Errorf("fehler beim Update der Platforms: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrAutomationSettingsNotFound
	}

	return nil
}
