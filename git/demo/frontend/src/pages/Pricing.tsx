import { useEffect, useRef, useState } from "react";
import { Link, useSearchParams } from "react-router-dom";
import { Check } from "lucide-react";
import { useAuth } from "../context/AuthContext";
import { pricingTiers, type Tier } from "../config/pricingTiers";
import { usePricePreview } from "../hooks/usePricePreview";
import { getPaddle } from "../lib/paddleSdk";

type Cycle = "month" | "year";

function tierKey(tier: Tier): string {
  return tier.name.toLowerCase();
}

export default function Pricing() {
  const { isAuthenticated, authChecked, email, twitchId } = useAuth();
  const [searchParams, setSearchParams] = useSearchParams();
  const [cycle, setCycle] = useState<Cycle>("month");
  const [checkoutError, setCheckoutError] = useState("");
  const autoCheckoutTried = useRef(false);

  const paidPriceIds = pricingTiers
    .filter((t) => t.priceId)
    .map((t) => t.priceId![cycle]);
  const { totals, loading: pricesLoading } = usePricePreview(paidPriceIds);

  async function openCheckout(tier: Tier, forCycle: Cycle) {
    if (!tier.priceId) return;
    setCheckoutError("");
    try {
      const paddle = await getPaddle();
      if (!paddle) throw new Error("Paddle.js konnte nicht initialisiert werden");

      paddle.Checkout.open({
        items: [{ priceId: tier.priceId[forCycle], quantity: 1 }],
        ...(email ? { customer: { email } } : {}),
        customData: { twitch_user_id: twitchId ?? "" },
        settings: {
          displayMode: "overlay",
          variant: "one-page",
          successUrl: `${window.location.origin}/welcome`,
        },
      });
    } catch (err) {
      setCheckoutError(err instanceof Error ? err.message : "Checkout konnte nicht gestartet werden");
    }
  }

  function handleSubscribe(tier: Tier) {
    if (!isAuthenticated) {
      const returnTo = `/pricing?tier=${tierKey(tier)}&cycle=${cycle}&autocheckout=1`;
      window.location.href = `/auth/login?returnTo=${encodeURIComponent(returnTo)}`;
      return;
    }
    openCheckout(tier, cycle);
  }

  // Resumes a checkout after the user was sent to /auth/login mid-Subscribe
  // click and came back here with ?autocheckout=1 - only fires once
  // authChecked has settled and isAuthenticated is actually true, so it
  // never fires for a login that failed/was cancelled.
  useEffect(() => {
    if (!authChecked || !isAuthenticated || autoCheckoutTried.current) return;
    if (searchParams.get("autocheckout") !== "1") return;

    const wantedTier = searchParams.get("tier");
    const wantedCycle = searchParams.get("cycle") === "year" ? "year" : "month";
    const tier = pricingTiers.find((t) => tierKey(t) === wantedTier && t.priceId);

    autoCheckoutTried.current = true;
    setCycle(wantedCycle);
    if (tier) {
      openCheckout(tier, wantedCycle);
    }

    // Strip the one-time params so a page refresh doesn't re-trigger checkout.
    setSearchParams({}, { replace: true });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [authChecked, isAuthenticated, searchParams]);

  return (
    <div className="relative w-full min-h-screen bg-slate-950 text-white">
      <section className="text-center py-16 px-6">
        <h1 className="text-4xl sm:text-5xl font-extrabold mb-4">Preise</h1>
        <p className="text-gray-300 max-w-2xl mx-auto mb-8">
          Wähle den Plan, der zu deinem Kanal passt - jederzeit änderbar.
        </p>

        <div className="inline-flex bg-gray-800 rounded-lg p-1 mb-12">
          <button
            onClick={() => setCycle("month")}
            className={`px-4 py-2 rounded-md transition-colors ${
              cycle === "month" ? "bg-gray-700 text-white" : "text-gray-400 hover:text-white"
            }`}
          >
            Monatlich
          </button>
          <button
            onClick={() => setCycle("year")}
            className={`px-4 py-2 rounded-md transition-colors ${
              cycle === "year" ? "bg-gray-700 text-white" : "text-gray-400 hover:text-white"
            }`}
          >
            Jährlich
          </button>
        </div>

        {checkoutError && (
          <p className="text-red-400 text-sm mb-6 max-w-md mx-auto">{checkoutError}</p>
        )}

        <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-6 max-w-5xl mx-auto px-4">
          {pricingTiers.map((tier) => {
            const priceId = tier.priceId?.[cycle];
            const formattedTotal = priceId ? totals[priceId] : undefined;

            return (
              <div
                key={tier.name}
                className="bg-gray-900 border border-gray-800 rounded-xl p-6 text-left flex flex-col"
              >
                <h3 className="text-xl font-bold mb-1">{tier.name}</h3>
                <p className="text-sm text-gray-400 mb-4">{tier.description}</p>

                <div className="mb-6">
                  {tier.priceId ? (
                    <span className="text-3xl font-bold">
                      {pricesLoading ? "…" : formattedTotal ?? "-"}
                    </span>
                  ) : (
                    <span className="text-3xl font-bold">Kostenlos</span>
                  )}
                  {tier.priceId && (
                    <span className="text-gray-400 text-sm ml-1">
                      /{cycle === "month" ? "Monat" : "Jahr"}
                    </span>
                  )}
                </div>

                <ul className="space-y-2 text-sm mb-6 flex-1">
                  {tier.features.map((feature) => (
                    <li key={feature} className="flex items-start gap-2">
                      <Check className="w-4 h-4 text-green-500 mt-0.5 flex-shrink-0" />
                      <span>{feature}</span>
                    </li>
                  ))}
                </ul>

                {tier.priceId ? (
                  <button
                    onClick={() => handleSubscribe(tier)}
                    className="w-full py-3 bg-green-600 hover:bg-green-500 rounded-lg transition-colors font-semibold"
                  >
                    Subscribe
                  </button>
                ) : isAuthenticated ? (
                  <Link
                    to="/dashboard"
                    className="w-full text-center py-3 bg-gray-700 hover:bg-gray-600 rounded-lg transition-colors block"
                  >
                    Zum Dashboard
                  </Link>
                ) : (
                  <a
                    href="/auth/login"
                    className="w-full text-center py-3 bg-gray-700 hover:bg-gray-600 rounded-lg transition-colors block"
                  >
                    Jetzt starten
                  </a>
                )}
              </div>
            );
          })}
        </div>
      </section>
    </div>
  );
}
