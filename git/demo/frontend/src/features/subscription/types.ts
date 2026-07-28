// features/subscription/types.ts
//
// Mirrors the real backend shape (internal/domain/subscription_toer.go,
// internal/domain/user_subscription.go) - the previous SubscriptionDashboard.tsx
// used a shape that never matched the actual API response (tierID vs tierId,
// currentPeriodEnd vs expiresAt, a made-up features object) and silently
// rendered undefined values. This file is the corrected contract.

export type TierID = "free" | "pro" | "premium";

export type SubscriptionStatus = "active" | "expired" | "canceled" | "trialing" | "past_due";

export type BillingCycle = "monthly" | "yearly";

export interface TierFeatures {
  simple_commands: boolean;
  advanced_commands: boolean;
  discord_integration: boolean;
  analytics: boolean;
  workflow_templates: boolean;
  custom_workflows: boolean;
  api_access: boolean;
  priority_support: boolean;
  clip_automation: boolean;
}

export interface SubscriptionTier {
  id: TierID;
  name: string;
  priceMonthly: number;
  priceYearly: number;
  maxCommands: number | null;
  maxWorkflows: number | null;
  maxVotesPerMonth: number | null;
  features: TierFeatures;
  isActive: boolean;
  createdAt: string;
}

export interface UserSubscription {
  id: number;
  userId: string;
  tierId: TierID;
  status: SubscriptionStatus;
  billingCycle: BillingCycle | null;
  startedAt: string;
  expiresAt: string | null;
  canceledAt: string | null;
  createdAt: string;
  updatedAt: string;
  tier?: SubscriptionTier;
}
