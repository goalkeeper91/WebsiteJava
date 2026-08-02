package service

import (
	"context"
	"log"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/paddle"
	"demo/backend-go/internal/repository"
	"demo/backend-go/pkg/config"
)

// PriceSyncService keeps subscription_tiers.price_monthly/price_yearly in
// sync with whatever is actually configured in the Paddle dashboard.
// /dashboard/subscription reads these DB columns directly (unlike the
// public /pricing page, which calls Paddle.PricePreview() live) - without
// this, a price change made only in Paddle silently goes stale here until
// someone remembers to write a matching DB migration by hand.
type PriceSyncService struct {
	tierRepo     repository.SubscriptionTierRepository
	paddleClient *paddle.Client
	cfg          config.PaddleConfig
}

func NewPriceSyncService(
	tierRepo repository.SubscriptionTierRepository,
	paddleClient *paddle.Client,
	cfg config.PaddleConfig,
) *PriceSyncService {
	return &PriceSyncService{
		tierRepo:     tierRepo,
		paddleClient: paddleClient,
		cfg:          cfg,
	}
}

// SyncPrices fetches the current monthly/yearly price for the two paid
// tiers from Paddle and writes them to subscription_tiers. Free has no
// Paddle price ID and is never touched, same as PaddleConfig.TierForPriceID
// only ever mapping pro/premium.
func (s *PriceSyncService) SyncPrices(ctx context.Context) {
	s.syncTier(ctx, domain.TierPro, s.cfg.ProPriceIDMonthly, s.cfg.ProPriceIDYearly)
	s.syncTier(ctx, domain.TierPremium, s.cfg.PremiumPriceIDMonthly, s.cfg.PremiumPriceIDYearly)
}

func (s *PriceSyncService) syncTier(ctx context.Context, tierID domain.TierID, monthlyPriceID, yearlyPriceID string) {
	if monthlyPriceID == "" || yearlyPriceID == "" {
		return // Paddle nicht konfiguriert (z.B. lokale Entwicklung) - nichts zu tun
	}

	monthly, err := s.paddleClient.GetPrice(ctx, monthlyPriceID)
	if err != nil {
		log.Printf("⚠️ Preis-Sync: Monatspreis für %s konnte nicht von Paddle geladen werden: %v", tierID, err)
		return
	}

	yearly, err := s.paddleClient.GetPrice(ctx, yearlyPriceID)
	if err != nil {
		log.Printf("⚠️ Preis-Sync: Jahrespreis für %s konnte nicht von Paddle geladen werden: %v", tierID, err)
		return
	}

	tier, err := s.tierRepo.GetByID(ctx, tierID)
	if err != nil {
		log.Printf("⚠️ Preis-Sync: Tier %s konnte nicht geladen werden: %v", tierID, err)
		return
	}
	if tier == nil {
		log.Printf("⚠️ Preis-Sync: Tier %s existiert nicht in der DB", tierID)
		return
	}

	if tier.PriceMonthly == monthly && tier.PriceYearly == yearly {
		return // bereits synchron - kein unnoetiges UPDATE
	}

	if err := s.tierRepo.UpdatePrices(ctx, tierID, monthly, yearly); err != nil {
		log.Printf("❌ Preis-Sync: Konnte Preise für %s nicht aktualisieren: %v", tierID, err)
		return
	}

	log.Printf(
		"💶 Preis-Sync: %s aktualisiert von %.2f€/%.2f€ auf %.2f€/Monat, %.2f€/Jahr",
		tierID, tier.PriceMonthly, tier.PriceYearly, monthly, yearly,
	)
}
