package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
)

type SubscriptionTierRepository struct {
	db *sql.DB
}

func NewSubscriptionTierRepository(db *sql.DB) repository.SubscriptionTierRepository {
	return &SubscriptionTierRepository{db: db}
}

func (r *SubscriptionTierRepository) GetAll(ctx context.Context) ([]*domain.SubscriptionTier, error) {
	query := `
		SELECT id, name, price_monthly, price_yearly,
		       max_commands, max_workflows, max_votes_per_month,
		       features, is_active, created_at
		FROM subscription_tiers
		ORDER BY price_monthly ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Tiers: %w", err)
	}
	defer rows.Close()

	var tiers []*domain.SubscriptionTier
	for rows.Next() {
		tier, err := scanTier(rows)
		if err != nil {
			return nil, err
		}
		tiers = append(tiers, tier)
	}

	return tiers, nil
}

func (r *SubscriptionTierRepository) GetByID(ctx context.Context, tierID domain.TierID) (*domain.SubscriptionTier, error) {
	query := `
		SELECT id, name, price_monthly, price_yearly,
		       max_commands, max_workflows, max_votes_per_month,
		       features, is_active, created_at
		FROM subscription_tiers
		WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, tierID)
	tier, err := scanTierRow(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden des Tiers: %w", err)
	}

	return tier, nil
}

func (r *SubscriptionTierRepository) GetActive(ctx context.Context) ([]*domain.SubscriptionTier, error) {
	query := `
		SELECT id, name, price_monthly, price_yearly,
		       max_commands, max_workflows, max_votes_per_month,
		       features, is_active, created_at
		FROM subscription_tiers
		WHERE is_active = true
		ORDER BY price_monthly ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der aktiven Tiers: %w", err)
	}
	defer rows.Close()

	var tiers []*domain.SubscriptionTier
	for rows.Next() {
		tier, err := scanTier(rows)
		if err != nil {
			return nil, err
		}
		tiers = append(tiers, tier)
	}

	return tiers, nil
}

// Helper: scan from rows
func scanTier(rows *sql.Rows) (*domain.SubscriptionTier, error) {
	var tier domain.SubscriptionTier
	var featuresJSON []byte

	err := rows.Scan(
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
	if err != nil {
		return nil, fmt.Errorf("fehler beim Scannen des Tiers: %w", err)
	}

	if err := json.Unmarshal(featuresJSON, &tier.Features); err != nil {
		return nil, fmt.Errorf("fehler beim Parsen der Features: %w", err)
	}

	return &tier, nil
}

// Helper: scan from single row
func scanTierRow(row *sql.Row) (*domain.SubscriptionTier, error) {
	var tier domain.SubscriptionTier
	var featuresJSON []byte

	err := row.Scan(
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
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(featuresJSON, &tier.Features); err != nil {
		return nil, fmt.Errorf("fehler beim Parsen der Features: %w", err)
	}

	return &tier, nil
}