package repository

import (
	"context"
	"demo/backend-go/internal/domain"
)

type UserSubscriptionRepository interface {
	Create(ctx context.Context, input domain.UserSubscriptionCreateInput) (*domain.UserSubscription, error)

	GetByUserID(ctx context.Context, userID string) (*domain.UserSubscription, error)

	GetByUserIDWithTier(ctx context.Context, userID string) (*domain.UserSubscription, error)

	// GetByPaddleCustomerID resolves a webhook event back to a local user
	// when the event only carries the Paddle customer ID (not customData).
	GetByPaddleCustomerID(ctx context.Context, paddleCustomerID string) (*domain.UserSubscription, error)

	Update(ctx context.Context, userID string, input domain.UserSubscriptionUpdateInput) error

	Cancel(ctx context.Context, userID string) error

	Exists(ctx context.Context, userID string) (bool, error)

	// ListAllWithUserAndTier lädt alle Nicht-Bot-Nutzer (LEFT JOIN gegen
	// Subscription+Tier - ein Nutzer ohne bisherige Subscription-Zeile hat
	// AdminCustomerRow.Subscription == nil) für die Admin-Kundenliste,
	// paginiert. Gibt zusätzlich die Gesamtzahl für die Pagination zurück.
	ListAllWithUserAndTier(ctx context.Context, limit, offset int) ([]*domain.AdminCustomerRow, int64, error)

	// GetAllActiveForMRR lädt alle existierenden Subscription-Zeilen (inkl.
	// Tier) ungepaginiert - für die MRR-Berechnung im Service, der IsActive()
	// und den Tarifpreis pro Zeile in Go auswertet statt die Logik in SQL zu
	// duplizieren.
	GetAllActiveForMRR(ctx context.Context) ([]*domain.UserSubscription, error)

	// AdminOverrideTier setzt den Tarif eines Nutzers manuell (Support-Fall),
	// unabhängig vom Paddle-Webhook-Pfad - legt bei Bedarf eine neue Zeile an.
	// expires_at wird dabei bewusst genullt (siehe Aufrufer).
	AdminOverrideTier(ctx context.Context, userID string, tierID domain.TierID) error
}