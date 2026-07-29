package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
)

type UserSubscriptionRepository struct {
	db *sql.DB
}

func NewUserSubscriptionRepository(db *sql.DB) repository.UserSubscriptionRepository {
	return &UserSubscriptionRepository{db: db}
}

func (r *UserSubscriptionRepository) Create(ctx context.Context, input domain.UserSubscriptionCreateInput) (*domain.UserSubscription, error) {
	query := `
		INSERT INTO user_subscriptions
		(user_id, tier_id, status, billing_cycle, started_at, expires_at,
		 paddle_customer_id, paddle_subscription_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), $5, $6, $7, NOW(), NOW())
		RETURNING id, user_id, tier_id, status, billing_cycle,
		          started_at, expires_at, canceled_at,
		          paddle_customer_id, paddle_subscription_id,
		          created_at, updated_at
	`

	var sub domain.UserSubscription
	err := r.db.QueryRowContext(
		ctx, query,
		input.UserID,
		input.TierID,
		domain.SubscriptionActive,
		input.BillingCycle,
		input.ExpiresAt,
		input.PaddleCustomerID,
		input.PaddleSubscriptionID,
	).Scan(
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
	)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Erstellen der Subscription: %w", err)
	}

	return &sub, nil
}

func (r *UserSubscriptionRepository) GetByUserID(ctx context.Context, userID string) (*domain.UserSubscription, error) {
	query := `
		SELECT id, user_id, tier_id, status, billing_cycle,
		       started_at, expires_at, canceled_at,
		       paddle_customer_id, paddle_subscription_id,
		       created_at, updated_at
		FROM user_subscriptions
		WHERE user_id = $1
	`

	var sub domain.UserSubscription
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
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
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Subscription: %w", err)
	}

	return &sub, nil
}

func (r *UserSubscriptionRepository) GetByUserIDWithTier(ctx context.Context, userID string) (*domain.UserSubscription, error) {
	query := `
		SELECT s.id, s.user_id, s.tier_id, s.status, s.billing_cycle,
		       s.started_at, s.expires_at, s.canceled_at,
		       s.paddle_customer_id, s.paddle_subscription_id,
		       s.created_at, s.updated_at,
		       t.id, t.name, t.price_monthly, t.price_yearly,
		       t.max_commands, t.max_workflows, t.max_votes_per_month,
		       t.features, t.is_active, t.created_at
		FROM user_subscriptions s
		INNER JOIN subscription_tiers t ON s.tier_id = t.id
		WHERE s.user_id = $1
	`

	var sub domain.UserSubscription
	var tier domain.SubscriptionTier
	var featuresJSON []byte

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
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
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Subscription: %w", err)
	}

	if err := json.Unmarshal(featuresJSON, &tier.Features); err != nil {
		return nil, fmt.Errorf("fehler beim Parsen der Features: %w", err)
	}

	sub.Tier = &tier
	return &sub, nil
}

func (r *UserSubscriptionRepository) GetByPaddleCustomerID(ctx context.Context, paddleCustomerID string) (*domain.UserSubscription, error) {
	query := `
		SELECT id, user_id, tier_id, status, billing_cycle,
		       started_at, expires_at, canceled_at,
		       paddle_customer_id, paddle_subscription_id,
		       created_at, updated_at
		FROM user_subscriptions
		WHERE paddle_customer_id = $1
	`

	var sub domain.UserSubscription
	err := r.db.QueryRowContext(ctx, query, paddleCustomerID).Scan(
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
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Subscription per Paddle-Kunden-ID: %w", err)
	}

	return &sub, nil
}

func (r *UserSubscriptionRepository) Update(ctx context.Context, userID string, input domain.UserSubscriptionUpdateInput) error {
	query := `
		UPDATE user_subscriptions
		SET tier_id = COALESCE($2, tier_id),
		    status = COALESCE($3, status),
		    billing_cycle = COALESCE($4, billing_cycle),
		    expires_at = COALESCE($5, expires_at),
		    canceled_at = COALESCE($6, canceled_at),
		    paddle_customer_id = COALESCE($7, paddle_customer_id),
		    paddle_subscription_id = COALESCE($8, paddle_subscription_id),
		    last_paddle_event_at = COALESCE($9, last_paddle_event_at),
		    updated_at = NOW()
		WHERE user_id = $1
	`

	_, err := r.db.ExecContext(
		ctx, query,
		userID,
		input.TierID,
		input.Status,
		input.BillingCycle,
		input.ExpiresAt,
		input.CanceledAt,
		input.PaddleCustomerID,
		input.PaddleSubscriptionID,
		input.LastPaddleEventAt,
	)
	if err != nil {
		return fmt.Errorf("fehler beim Update der Subscription: %w", err)
	}

	return nil
}

func (r *UserSubscriptionRepository) Cancel(ctx context.Context, userID string) error {
	query := `
		UPDATE user_subscriptions
		SET status = $2, canceled_at = NOW(), updated_at = NOW()
		WHERE user_id = $1
	`

	_, err := r.db.ExecContext(ctx, query, userID, domain.SubscriptionCanceled)
	if err != nil {
		return fmt.Errorf("fehler beim Kündigen der Subscription: %w", err)
	}

	return nil
}

func (r *UserSubscriptionRepository) Exists(ctx context.Context, userID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM user_subscriptions WHERE user_id = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("fehler beim Prüfen der Subscription: %w", err)
	}

	return exists, nil
}

func (r *UserSubscriptionRepository) ListAllWithUserAndTier(ctx context.Context, limit, offset int) ([]*domain.AdminCustomerRow, int64, error) {
	countQuery := `SELECT COUNT(*) FROM users WHERE is_bot = false`
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("fehler beim Zählen der Kunden: %w", err)
	}

	query := `
		SELECT u.twitch_id, u.username, u.email, u.discord_id, u.is_admin, u.created_at, u.updated_at,
		       s.id, s.user_id, s.tier_id, s.status, s.billing_cycle,
		       s.started_at, s.expires_at, s.canceled_at,
		       s.paddle_customer_id, s.paddle_subscription_id,
		       s.created_at, s.updated_at,
		       t.id, t.name, t.price_monthly, t.price_yearly,
		       t.max_commands, t.max_workflows, t.max_votes_per_month,
		       t.features, t.is_active, t.created_at
		FROM users u
		LEFT JOIN user_subscriptions s ON s.user_id = u.twitch_id
		LEFT JOIN subscription_tiers t ON s.tier_id = t.id
		WHERE u.is_bot = false
		ORDER BY u.created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("fehler beim Laden der Kundenliste: %w", err)
	}
	defer rows.Close()

	var result []*domain.AdminCustomerRow
	for rows.Next() {
		var user domain.User
		var sID sql.NullInt64
		var sUserID, sTierID, sStatus, sBillingCycle, sPaddleCustomerID, sPaddleSubscriptionID sql.NullString
		var sStartedAt, sExpiresAt, sCanceledAt, sCreatedAt, sUpdatedAt sql.NullTime
		var tID, tName sql.NullString
		var tPriceMonthly, tPriceYearly sql.NullFloat64
		var tMaxCommands, tMaxWorkflows, tMaxVotesPerMonth sql.NullInt64
		var tFeaturesJSON []byte
		var tIsActive sql.NullBool
		var tCreatedAt sql.NullTime

		err := rows.Scan(
			&user.TwitchID, &user.Username, &user.Email, &user.DiscordID, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt,
			&sID, &sUserID, &sTierID, &sStatus, &sBillingCycle,
			&sStartedAt, &sExpiresAt, &sCanceledAt,
			&sPaddleCustomerID, &sPaddleSubscriptionID,
			&sCreatedAt, &sUpdatedAt,
			&tID, &tName, &tPriceMonthly, &tPriceYearly,
			&tMaxCommands, &tMaxWorkflows, &tMaxVotesPerMonth,
			&tFeaturesJSON, &tIsActive, &tCreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("fehler beim Scan der Kundenliste: %w", err)
		}

		row := &domain.AdminCustomerRow{User: &user}

		if sID.Valid {
			billingCycle := domain.BillingCycle(sBillingCycle.String)
			var expiresAt, canceledAt *time.Time
			if sExpiresAt.Valid {
				expiresAt = &sExpiresAt.Time
			}
			if sCanceledAt.Valid {
				canceledAt = &sCanceledAt.Time
			}
			var billingCyclePtr *domain.BillingCycle
			if sBillingCycle.Valid {
				billingCyclePtr = &billingCycle
			}
			var paddleCustomerID, paddleSubscriptionID *string
			if sPaddleCustomerID.Valid {
				paddleCustomerID = &sPaddleCustomerID.String
			}
			if sPaddleSubscriptionID.Valid {
				paddleSubscriptionID = &sPaddleSubscriptionID.String
			}

			sub := &domain.UserSubscription{
				ID:                   sID.Int64,
				UserID:               sUserID.String,
				TierID:               domain.TierID(sTierID.String),
				Status:               domain.SubscriptionStatus(sStatus.String),
				BillingCycle:         billingCyclePtr,
				StartedAt:            sStartedAt.Time,
				ExpiresAt:            expiresAt,
				CanceledAt:           canceledAt,
				PaddleCustomerID:     paddleCustomerID,
				PaddleSubscriptionID: paddleSubscriptionID,
				CreatedAt:            sCreatedAt.Time,
				UpdatedAt:            sUpdatedAt.Time,
			}

			if tID.Valid {
				tier := &domain.SubscriptionTier{
					ID:           domain.TierID(tID.String),
					Name:         tName.String,
					PriceMonthly: tPriceMonthly.Float64,
					PriceYearly:  tPriceYearly.Float64,
					IsActive:     tIsActive.Bool,
					CreatedAt:    tCreatedAt.Time,
				}
				if tMaxCommands.Valid {
					v := int(tMaxCommands.Int64)
					tier.MaxCommands = &v
				}
				if tMaxWorkflows.Valid {
					v := int(tMaxWorkflows.Int64)
					tier.MaxWorkflows = &v
				}
				if tMaxVotesPerMonth.Valid {
					v := int(tMaxVotesPerMonth.Int64)
					tier.MaxVotesPerMonth = &v
				}
				if tFeaturesJSON != nil {
					if err := json.Unmarshal(tFeaturesJSON, &tier.Features); err != nil {
						return nil, 0, fmt.Errorf("fehler beim Parsen der Tier-Features: %w", err)
					}
				}
				sub.Tier = tier
			}

			row.Subscription = sub
		}

		result = append(result, row)
	}

	return result, total, rows.Err()
}

func (r *UserSubscriptionRepository) GetAllActiveForMRR(ctx context.Context) ([]*domain.UserSubscription, error) {
	query := `
		SELECT s.id, s.user_id, s.tier_id, s.status, s.billing_cycle,
		       s.started_at, s.expires_at, s.canceled_at,
		       s.paddle_customer_id, s.paddle_subscription_id,
		       s.created_at, s.updated_at,
		       t.id, t.name, t.price_monthly, t.price_yearly,
		       t.max_commands, t.max_workflows, t.max_votes_per_month,
		       t.features, t.is_active, t.created_at
		FROM user_subscriptions s
		INNER JOIN subscription_tiers t ON s.tier_id = t.id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Subscriptions für MRR: %w", err)
	}
	defer rows.Close()

	var result []*domain.UserSubscription
	for rows.Next() {
		var sub domain.UserSubscription
		var tier domain.SubscriptionTier
		var featuresJSON []byte

		err := rows.Scan(
			&sub.ID, &sub.UserID, &sub.TierID, &sub.Status, &sub.BillingCycle,
			&sub.StartedAt, &sub.ExpiresAt, &sub.CanceledAt,
			&sub.PaddleCustomerID, &sub.PaddleSubscriptionID,
			&sub.CreatedAt, &sub.UpdatedAt,
			&tier.ID, &tier.Name, &tier.PriceMonthly, &tier.PriceYearly,
			&tier.MaxCommands, &tier.MaxWorkflows, &tier.MaxVotesPerMonth,
			&featuresJSON, &tier.IsActive, &tier.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("fehler beim Scan für MRR: %w", err)
		}

		if err := json.Unmarshal(featuresJSON, &tier.Features); err != nil {
			return nil, fmt.Errorf("fehler beim Parsen der Features: %w", err)
		}

		sub.Tier = &tier
		result = append(result, &sub)
	}

	return result, rows.Err()
}

func (r *UserSubscriptionRepository) AdminOverrideTier(ctx context.Context, userID string, tierID domain.TierID) error {
	query := `
		INSERT INTO user_subscriptions (user_id, tier_id, status, expires_at, started_at, created_at, updated_at)
		VALUES ($1, $2, $3, NULL, NOW(), NOW(), NOW())
		ON CONFLICT (user_id) DO UPDATE SET
			tier_id = EXCLUDED.tier_id,
			status = EXCLUDED.status,
			expires_at = NULL,
			updated_at = NOW()
	`

	_, err := r.db.ExecContext(ctx, query, userID, tierID, domain.SubscriptionActive)
	if err != nil {
		return fmt.Errorf("fehler beim manuellen Setzen des Tarifs: %w", err)
	}

	return nil
}