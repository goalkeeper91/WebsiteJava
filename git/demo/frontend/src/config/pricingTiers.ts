// config/pricingTiers.ts
//
// Hand-edited tier list for the public pricing page - deliberately decoupled
// from the backend's subscription_tiers.name column, which stays "Free"/
// "Pro"/"Premium" internally (TierID free/pro/premium is used pervasively in
// backend feature-gating, e.g. IsFree()/IsPro()/IsPremium() - not renamed).
// This file's display names/descriptions/features are just marketing copy.
//
// Starter has no priceId - it maps to the existing free tier and never goes
// through Paddle checkout (see Pricing.tsx's "Jetzt starten" CTA instead).

export interface Tier {
  name: "Starter" | "Pro" | "Advanced";
  description: string;
  features: string[];
  priceId?: { month: string; year: string };
}

export const pricingTiers: Tier[] = [
  {
    name: "Starter",
    description: "Zum Reinschnuppern - die Grundlagen für jeden Twitch-Kanal.",
    features: ["Bis zu 10 Chat-Commands", "Automod", "Loyalty-Punkte", "Giveaways"],
  },
  {
    name: "Pro",
    description: "Für Streamer, die ihren Kanal ernsthaft ausbauen.",
    features: [
      "Alles aus Starter",
      "Erweiterte Commands (n8n)",
      "Discord-Integration",
      "Analytics-Dashboard",
      "Workflow-Templates",
    ],
    priceId: {
      month: import.meta.env.VITE_PADDLE_PRICE_PRO_MONTHLY,
      year: import.meta.env.VITE_PADDLE_PRICE_PRO_YEARLY,
    },
  },
  {
    name: "Advanced",
    description: "Das volle Paket, inklusive automatischer Clip-Erstellung.",
    features: [
      "Alles aus Pro",
      "Clip-Automatisierung",
      "Eigene Workflows",
      "API-Zugriff",
      "Priority Support",
    ],
    priceId: {
      month: import.meta.env.VITE_PADDLE_PRICE_PREMIUM_MONTHLY,
      year: import.meta.env.VITE_PADDLE_PRICE_PREMIUM_YEARLY,
    },
  },
];
