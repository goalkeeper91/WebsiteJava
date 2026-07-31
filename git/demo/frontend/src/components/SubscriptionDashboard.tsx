import { useEffect, useState } from "react";
import { Crown, Zap, Rocket, Check, X } from "lucide-react";
import { getSubscription, getTiers, downgradeToFree, getPortalLink } from "../features/subscription/api";
import type { SubscriptionTier, UserSubscription, TierID } from "../features/subscription/types";
import { openPaddleCheckout, getPaddlePriceID, setCheckoutCompletedHandler } from "../lib/paddle";

const TIER_ICONS: Record<TierID, typeof Crown> = {
  free: Crown,
  pro: Zap,
  premium: Rocket,
};

const TIER_COLORS: Record<TierID, string> = {
  free: "text-gray-400",
  pro: "text-blue-400",
  premium: "text-purple-400",
};

export default function SubscriptionDashboard() {
  const [subscription, setSubscription] = useState<UserSubscription | null>(null);
  const [tiers, setTiers] = useState<SubscriptionTier[]>([]);
  const [loading, setLoading] = useState(true);
  const [billingCycle, setBillingCycle] = useState<"monthly" | "yearly">("monthly");
  const [checkoutError, setCheckoutError] = useState("");
  const [portalLoading, setPortalLoading] = useState(false);
  const [portalError, setPortalError] = useState("");
  const [downgrading, setDowngrading] = useState(false);

  useEffect(() => {
    loadSubscriptionData();

    // Best-effort auto-refresh once the Paddle overlay reports the checkout
    // finished - the tier itself only actually changes once Paddle's
    // webhook lands, this just re-polls sooner than a manual page reload.
    setCheckoutCompletedHandler(() => {
      setTimeout(loadSubscriptionData, 2000);
    });
    return () => setCheckoutCompletedHandler(null);
  }, []);

  async function loadSubscriptionData() {
    setLoading(true);
    try {
      const [subData, tiersData] = await Promise.all([getSubscription(), getTiers()]);
      setSubscription(subData);
      setTiers(tiersData);
    } catch (err) {
      console.error("Failed to load subscription data:", err);
    } finally {
      setLoading(false);
    }
  }

  async function handleUpgrade(tier: SubscriptionTier) {
    setCheckoutError("");

    if (tier.id === "free") {
      setDowngrading(true);
      try {
        await downgradeToFree();
        await loadSubscriptionData();
      } catch (err) {
        setCheckoutError(err instanceof Error ? err.message : "Wechsel auf Free fehlgeschlagen");
      } finally {
        setDowngrading(false);
      }
      return;
    }

    if (!subscription) return;

    const priceId = getPaddlePriceID(tier.id as "pro" | "premium", billingCycle);
    if (!priceId) {
      setCheckoutError("Für diesen Tarif ist aktuell kein Checkout verfügbar.");
      return;
    }

    try {
      await openPaddleCheckout(priceId, subscription.userId);
    } catch (err) {
      setCheckoutError(err instanceof Error ? err.message : "Checkout konnte nicht gestartet werden");
    }
  }

  async function handleManageSubscription() {
    setPortalError("");
    setPortalLoading(true);
    try {
      const url = await getPortalLink();
      window.open(url, "_blank", "noopener,noreferrer");
    } catch (err) {
      setPortalError(err instanceof Error ? err.message : "Kundenportal konnte nicht geöffnet werden");
    } finally {
      setPortalLoading(false);
    }
  }

  function renderFeatureList(tier: SubscriptionTier) {
    const f = tier.features;
    return (
      <ul className="space-y-2 text-sm">
        <li className="flex items-center gap-2">
          <Check className="w-4 h-4 text-green-500" />
          <span>{tier.maxCommands === null ? "Unlimited" : tier.maxCommands} Simple Commands</span>
        </li>
        <li className="flex items-center gap-2">
          {f.advanced_commands ? (
            <Check className="w-4 h-4 text-green-500" />
          ) : (
            <X className="w-4 h-4 text-gray-600" />
          )}
          <span className={!f.advanced_commands ? "text-gray-500" : ""}>Advanced Commands (n8n)</span>
        </li>
        <li className="flex items-center gap-2">
          {tier.maxWorkflows === null || tier.maxWorkflows > 0 ? (
            <Check className="w-4 h-4 text-green-500" />
          ) : (
            <X className="w-4 h-4 text-gray-600" />
          )}
          <span className={tier.maxWorkflows === 0 ? "text-gray-500" : ""}>
            {tier.maxWorkflows === null ? "Unlimited" : tier.maxWorkflows} Workflows
          </span>
        </li>
        <li className="flex items-center gap-2">
          {tier.maxVotesPerMonth === null || tier.maxVotesPerMonth > 0 ? (
            <Check className="w-4 h-4 text-green-500" />
          ) : (
            <X className="w-4 h-4 text-gray-600" />
          )}
          <span className={tier.maxVotesPerMonth === 0 ? "text-gray-500" : ""}>
            {tier.maxVotesPerMonth === null ? "Unlimited" : tier.maxVotesPerMonth} Votes/Monat
          </span>
        </li>
        <li className="flex items-center gap-2">
          {f.discord_integration ? (
            <Check className="w-4 h-4 text-green-500" />
          ) : (
            <X className="w-4 h-4 text-gray-600" />
          )}
          <span className={!f.discord_integration ? "text-gray-500" : ""}>Discord-Integration</span>
        </li>
        <li className="flex items-center gap-2">
          {f.analytics ? (
            <Check className="w-4 h-4 text-green-500" />
          ) : (
            <X className="w-4 h-4 text-gray-600" />
          )}
          <span className={!f.analytics ? "text-gray-500" : ""}>Analytics Dashboard</span>
        </li>
        <li className="flex items-center gap-2">
          {f.clip_automation ? (
            <Check className="w-4 h-4 text-green-500" />
          ) : (
            <X className="w-4 h-4 text-gray-600" />
          )}
          <span className={!f.clip_automation ? "text-gray-500" : ""}>Clip-Automatisierung</span>
        </li>
        <li className="flex items-center gap-2">
          {f.priority_support ? (
            <Check className="w-4 h-4 text-green-500" />
          ) : (
            <X className="w-4 h-4 text-gray-600" />
          )}
          <span className={!f.priority_support ? "text-gray-500" : ""}>Priority Support</span>
        </li>
      </ul>
    );
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-900 text-white">
        <div className="p-6 max-w-7xl mx-auto">
          <div className="animate-pulse space-y-6">
            <div className="h-8 bg-gray-800 rounded w-1/3"></div>
            <div className="grid md:grid-cols-3 gap-6">
              {[1, 2, 3].map((i) => (
                <div key={i} className="h-96 bg-gray-800 rounded"></div>
              ))}
            </div>
          </div>
        </div>
      </div>
    );
  }

  const currentTierID = subscription?.tierId || "free";
  const isCanceled = subscription?.status === "canceled";

  return (
    <div className="min-h-screen bg-gray-900 text-white">
      <div className="p-6 max-w-7xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold mb-2">Subscription & Billing</h1>
          <p className="text-gray-400">
            Aktueller Plan: <span className="font-semibold text-white">{subscription?.tier?.name || "Free"}</span>
          </p>
          {subscription?.expiresAt && !isCanceled && (
            <p className="text-sm text-gray-500">
              Verlängert sich am {new Date(subscription.expiresAt).toLocaleDateString()}
            </p>
          )}
          {isCanceled && subscription?.expiresAt && (
            <p className="text-sm text-red-400">
              ⚠️ Gekündigt - läuft bis {new Date(subscription.expiresAt).toLocaleDateString()}
            </p>
          )}
          {checkoutError && <p className="text-sm text-red-400 mt-2">{checkoutError}</p>}
          {portalError && <p className="text-sm text-red-400 mt-2">{portalError}</p>}
        </div>

        {/* Billing Cycle Toggle */}
        <div className="flex justify-center mb-8">
          <div className="inline-flex bg-gray-800 rounded-lg p-1">
            <button
              onClick={() => setBillingCycle("monthly")}
              className={`px-4 py-2 rounded-md transition-colors ${
                billingCycle === "monthly" ? "bg-gray-700 text-white" : "text-gray-400 hover:text-white"
              }`}
            >
              Monatlich
            </button>
            <button
              onClick={() => setBillingCycle("yearly")}
              className={`px-4 py-2 rounded-md transition-colors ${
                billingCycle === "yearly" ? "bg-gray-700 text-white" : "text-gray-400 hover:text-white"
              }`}
            >
              Jährlich
              <span className="ml-1 text-xs text-green-400">(20% sparen)</span>
            </button>
          </div>
        </div>

        {/* Pricing Cards */}
        <div className="grid md:grid-cols-3 gap-6 mb-8">
          {tiers.map((tier) => {
            const Icon = TIER_ICONS[tier.id] || Crown;
            const colorClass = TIER_COLORS[tier.id];
            const isCurrent = tier.id === currentTierID;
            const price = billingCycle === "monthly" ? tier.priceMonthly : tier.priceYearly;

            return (
              <div
                key={tier.id}
                className={`bg-gray-800 rounded-xl p-6 border-2 transition-all ${
                  isCurrent ? "border-green-500 shadow-lg shadow-green-500/20" : "border-gray-700 hover:border-gray-600"
                }`}
              >
                <div className="flex items-center justify-between mb-4">
                  <div className="flex items-center gap-2">
                    <Icon className={`w-6 h-6 ${colorClass}`} />
                    <h3 className="text-xl font-bold">{tier.name}</h3>
                  </div>
                  {isCurrent && <span className="text-xs bg-green-600 px-2 py-1 rounded">Aktuell</span>}
                </div>

                <div className="mb-6">
                  <div className="flex items-baseline gap-1">
                    <span className="text-4xl font-bold">{price}€</span>
                    <span className="text-gray-400">/{billingCycle === "monthly" ? "Monat" : "Jahr"}</span>
                  </div>
                  {billingCycle === "yearly" && price > 0 && (
                    <p className="text-sm text-gray-500 mt-1">{(price / 12).toFixed(2)}€ pro Monat</p>
                  )}
                </div>

                <div className="mb-6">{renderFeatureList(tier)}</div>

                {isCurrent ? (
                  tier.id !== "free" ? (
                    <button
                      onClick={handleManageSubscription}
                      disabled={portalLoading}
                      className="w-full py-3 bg-gray-700 hover:bg-red-600 rounded-lg transition-colors disabled:opacity-50"
                    >
                      {portalLoading ? "Öffne..." : "Abo verwalten / kündigen"}
                    </button>
                  ) : (
                    <button disabled className="w-full py-3 bg-gray-700 rounded-lg cursor-not-allowed opacity-50">
                      Aktueller Plan
                    </button>
                  )
                ) : tier.id !== "free" ? (
                  <button
                    onClick={() => handleUpgrade(tier)}
                    className="w-full py-3 bg-green-600 hover:bg-green-500 rounded-lg transition-colors font-semibold"
                  >
                    Upgrade auf {tier.name}
                  </button>
                ) : (
                  <button
                    onClick={() => handleUpgrade(tier)}
                    disabled={downgrading}
                    className="w-full py-3 bg-gray-700 hover:bg-gray-600 rounded-lg transition-colors disabled:opacity-50"
                  >
                    {downgrading ? "Wechsle..." : "Zu Free wechseln"}
                  </button>
                )}
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}
