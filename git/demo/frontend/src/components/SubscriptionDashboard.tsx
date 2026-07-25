import { useEffect, useState } from "react";
import { Crown, Zap, Rocket, Check, X } from "lucide-react";

interface SubscriptionTier {
  id: string;
  name: string;
  priceMonthly: number;
  priceYearly: number;
  features: {
    commands: number;
    advancedCommands: boolean;
    workflows: number;
    votesPerMonth: number;
    analytics: boolean;
    support: string;
  };
}

interface UserSubscription {
  id: number;
  tierID: string;
  status: string;
  billingCycle: string;
  currentPeriodEnd?: string;
  canceledAt?: string;
  tier?: SubscriptionTier;
}

const TIER_ICONS = {
  free: Crown,
  pro: Zap,
  premium: Rocket,
};

const TIER_COLORS = {
  free: "text-gray-400",
  pro: "text-blue-400",
  premium: "text-purple-400",
};

export default function SubscriptionDashboard() {
  const [subscription, setSubscription] = useState<UserSubscription | null>(null);
  const [tiers, setTiers] = useState<SubscriptionTier[]>([]);
  const [loading, setLoading] = useState(true);
  const [billingCycle, setBillingCycle] = useState<"monthly" | "yearly">("monthly");
  const [showCancelConfirm, setShowCancelConfirm] = useState(false);

  useEffect(() => {
    loadSubscriptionData();
  }, []);

  async function loadSubscriptionData() {
    setLoading(true);
    try {
      // Load user's current subscription
      const subRes = await fetch("/api/subscription", {
        credentials: "include",
      });
      if (subRes.ok) {
        const subData = await subRes.json();
        setSubscription(subData);
      }

      // Load available tiers
      const tiersRes = await fetch("/api/subscription/tiers");
      if (tiersRes.ok) {
        const tiersData = await tiersRes.json();
        setTiers(tiersData);
      }
    } catch (err) {
      console.error("Failed to load subscription data:", err);
    } finally {
      setLoading(false);
    }
  }

  async function handleUpgrade(tierID: string) {
    if (!confirm(`Möchtest du wirklich zum ${tierID.toUpperCase()} Plan upgraden?`)) {
      return;
    }

    try {
      const res = await fetch("/api/subscription/upgrade", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ tier_id: tierID }),
      });

      if (res.ok) {
        alert("Upgrade erfolgreich! 🎉");
        loadSubscriptionData();
      } else {
        const error = await res.text();
        alert(`Upgrade fehlgeschlagen: ${error}`);
      }
    } catch (err) {
      console.error("Upgrade error:", err);
      alert("Upgrade fehlgeschlagen");
    }
  }

  async function handleCancel() {
    try {
      const res = await fetch("/api/subscription/cancel", {
        method: "POST",
        credentials: "include",
      });

      if (res.ok) {
        alert("Abo wurde gekündigt. Es läuft bis zum Ende der Abrechnungsperiode.");
        loadSubscriptionData();
        setShowCancelConfirm(false);
      } else {
        alert("Kündigung fehlgeschlagen");
      }
    } catch (err) {
      console.error("Cancel error:", err);
      alert("Kündigung fehlgeschlagen");
    }
  }

  function renderFeatureList(features: SubscriptionTier["features"]) {
    return (
      <ul className="space-y-2 text-sm">
        <li className="flex items-center gap-2">
          <Check className="w-4 h-4 text-green-500" />
          <span>
            {features.commands === -1 ? "Unlimited" : features.commands} Simple Commands
          </span>
        </li>
        <li className="flex items-center gap-2">
          {features.advancedCommands ? (
            <Check className="w-4 h-4 text-green-500" />
          ) : (
            <X className="w-4 h-4 text-gray-600" />
          )}
          <span className={!features.advancedCommands ? "text-gray-500" : ""}>
            Advanced Commands (n8n)
          </span>
        </li>
        <li className="flex items-center gap-2">
          {features.workflows > 0 ? (
            <Check className="w-4 h-4 text-green-500" />
          ) : (
            <X className="w-4 h-4 text-gray-600" />
          )}
          <span className={features.workflows === 0 ? "text-gray-500" : ""}>
            {features.workflows === -1 ? "Unlimited" : features.workflows} Workflows
          </span>
        </li>
        <li className="flex items-center gap-2">
          {features.votesPerMonth > 0 ? (
            <Check className="w-4 h-4 text-green-500" />
          ) : (
            <X className="w-4 h-4 text-gray-600" />
          )}
          <span className={features.votesPerMonth === 0 ? "text-gray-500" : ""}>
            {features.votesPerMonth === -1
              ? "Unlimited"
              : features.votesPerMonth} Votes/Monat
          </span>
        </li>
        <li className="flex items-center gap-2">
          {features.analytics ? (
            <Check className="w-4 h-4 text-green-500" />
          ) : (
            <X className="w-4 h-4 text-gray-600" />
          )}
          <span className={!features.analytics ? "text-gray-500" : ""}>
            Analytics Dashboard
          </span>
        </li>
        <li className="flex items-center gap-2">
          <Check className="w-4 h-4 text-green-500" />
          <span>{features.support}</span>
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

  const currentTierID = subscription?.tierID || "free";

  return (
    <div className="min-h-screen bg-gray-900 text-white">
    <div className="p-6 max-w-7xl mx-auto">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold mb-2">Subscription & Billing</h1>
        <p className="text-gray-400">
          Aktueller Plan:{" "}
          <span className="font-semibold text-white">
            {subscription?.tier?.name || "Free"}
          </span>
        </p>
        {subscription?.currentPeriodEnd && (
          <p className="text-sm text-gray-500">
            Verlängert sich am{" "}
            {new Date(subscription.currentPeriodEnd).toLocaleDateString()}
          </p>
        )}
        {subscription?.canceledAt && (
          <p className="text-sm text-red-400">
            ⚠️ Gekündigt - läuft bis {new Date(subscription.currentPeriodEnd!).toLocaleDateString()}
          </p>
        )}
      </div>

      {/* Billing Cycle Toggle */}
      <div className="flex justify-center mb-8">
        <div className="inline-flex bg-gray-800 rounded-lg p-1">
          <button
            onClick={() => setBillingCycle("monthly")}
            className={`px-4 py-2 rounded-md transition-colors ${
              billingCycle === "monthly"
                ? "bg-gray-700 text-white"
                : "text-gray-400 hover:text-white"
            }`}
          >
            Monatlich
          </button>
          <button
            onClick={() => setBillingCycle("yearly")}
            className={`px-4 py-2 rounded-md transition-colors ${
              billingCycle === "yearly"
                ? "bg-gray-700 text-white"
                : "text-gray-400 hover:text-white"
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
          const Icon = TIER_ICONS[tier.id as keyof typeof TIER_ICONS] || Crown;
          const colorClass = TIER_COLORS[tier.id as keyof typeof TIER_COLORS];
          const isCurrent = tier.id === currentTierID;
          const price =
            billingCycle === "monthly" ? tier.priceMonthly : tier.priceYearly;

          return (
            <div
              key={tier.id}
              className={`bg-gray-800 rounded-xl p-6 border-2 transition-all ${
                isCurrent
                  ? "border-green-500 shadow-lg shadow-green-500/20"
                  : "border-gray-700 hover:border-gray-600"
              }`}
            >
              {/* Tier Header */}
              <div className="flex items-center justify-between mb-4">
                <div className="flex items-center gap-2">
                  <Icon className={`w-6 h-6 ${colorClass}`} />
                  <h3 className="text-xl font-bold">{tier.name}</h3>
                </div>
                {isCurrent && (
                  <span className="text-xs bg-green-600 px-2 py-1 rounded">
                    Aktuell
                  </span>
                )}
              </div>

              {/* Price */}
              <div className="mb-6">
                <div className="flex items-baseline gap-1">
                  <span className="text-4xl font-bold">{price}€</span>
                  <span className="text-gray-400">
                    /{billingCycle === "monthly" ? "Monat" : "Jahr"}
                  </span>
                </div>
                {billingCycle === "yearly" && price > 0 && (
                  <p className="text-sm text-gray-500 mt-1">
                    {(price / 12).toFixed(2)}€ pro Monat
                  </p>
                )}
              </div>

              {/* Features */}
              <div className="mb-6">{renderFeatureList(tier.features)}</div>

              {/* CTA Button */}
              {isCurrent ? (
                tier.id !== "free" && !subscription?.canceledAt ? (
                  <button
                    onClick={() => setShowCancelConfirm(true)}
                    className="w-full py-3 bg-gray-700 hover:bg-red-600 rounded-lg transition-colors"
                  >
                    Abo kündigen
                  </button>
                ) : (
                  <button disabled className="w-full py-3 bg-gray-700 rounded-lg cursor-not-allowed opacity-50">
                    Aktueller Plan
                  </button>
                )
              ) : tier.id !== "free" ? (
                <button
                  onClick={() => handleUpgrade(tier.id)}
                  className="w-full py-3 bg-green-600 hover:bg-green-500 rounded-lg transition-colors font-semibold"
                >
                  Upgrade auf {tier.name}
                </button>
              ) : (
                <button
                  onClick={() => handleUpgrade("free")}
                  className="w-full py-3 bg-gray-700 hover:bg-gray-600 rounded-lg transition-colors"
                >
                  Zu Free wechseln
                </button>
              )}
            </div>
          );
        })}
      </div>

      {/* Cancel Confirmation Modal */}
      {showCancelConfirm && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="bg-gray-800 rounded-xl p-6 max-w-md w-full max-h-[90vh] overflow-y-auto">
            <h3 className="text-xl font-bold mb-4">Abo wirklich kündigen?</h3>
            <p className="text-gray-400 mb-6">
              Dein Abo läuft bis zum Ende der aktuellen Abrechnungsperiode weiter.
              Danach verlierst du den Zugriff auf Premium-Features.
            </p>
            <div className="flex gap-3">
              <button
                onClick={() => setShowCancelConfirm(false)}
                className="flex-1 py-2 bg-gray-700 hover:bg-gray-600 rounded-lg transition-colors"
              >
                Abbrechen
              </button>
              <button
                onClick={handleCancel}
                className="flex-1 py-2 bg-red-600 hover:bg-red-500 rounded-lg transition-colors"
              >
                Kündigen
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
    </div>
  );
}