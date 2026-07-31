import { useEffect, useRef, useState } from "react";
import { Link, useSearchParams } from "react-router-dom";
import { Check } from "lucide-react";
import { useAuth } from "../context/AuthContext";
import { pricingTiers, type Tier } from "../config/pricingTiers";
import { usePricePreview } from "../hooks/usePricePreview";
import { getPaddle } from "../lib/paddleSdk";
import Seo from "../components/Seo";

type Cycle = "month" | "year";

function tierKey(tier: Tier): string {
  return tier.name.toLowerCase();
}

export default function Pricing() {
  const { isAuthenticated, authChecked, email, twitchId } = useAuth();
  const [searchParams, setSearchParams] = useSearchParams();
  const [cycle, setCycle] = useState<Cycle>("month");
  const [checkoutError, setCheckoutError] = useState("");
  const [withdrawalConsent, setWithdrawalConsent] = useState(false);
  const autoCheckoutTried = useRef(false);

  const paidPriceIds = pricingTiers
    .filter((t) => t.priceId)
    .map((t) => t.priceId![cycle]);
  const { totals, loading: pricesLoading } = usePricePreview(paidPriceIds);

  // Paddle's own checkout does not surface the "loss of withdrawal right on
  // immediate performance" consent required by § 356 Abs. 5 BGB for digital
  // content/services (confirmed - no such option exists in the Paddle
  // dashboard), so we collect this ourselves as an explicit, required
  // checkbox before ever opening Paddle Checkout. Also protects the
  // post-login autocheckout resume below, since consent state doesn't (and
  // shouldn't) survive the redirect - the user must tick it again themselves.
  async function openCheckout(tier: Tier, forCycle: Cycle) {
    if (!tier.priceId) return;
    setCheckoutError("");
    if (!withdrawalConsent) {
      setCheckoutError("Bitte bestätige zuerst die Checkbox zum Widerrufsrecht unten, um fortzufahren.");
      return;
    }
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
      <Seo
        title="Preise & Leistungen"
        description="Twitch-Bot-Abo (Starter, Pro, Advanced) von Goalkeeper91 sowie individuelle Softwareentwicklung nach persönlichem Angebot - keine Pauschalpreise, jedes Projekt einzeln kalkuliert."
        path="/pricing"
      />

      {/* Intro + Sprungnavigation - orientiert sofort, dass hier zwei getrennte Angebote stehen */}
      <section className="text-center pt-16 pb-10 px-6">
        <h1 className="text-4xl sm:text-5xl font-extrabold mb-4">Preise & Leistungen</h1>
        <p className="text-gray-300 max-w-2xl mx-auto mb-8">
          Zwei getrennte Angebote: ein fertiges SaaS-Abo für den Twitch Bot mit festen Tarifen, und individuelle
          Softwareentwicklung als Freelancer - dafür gibt es keine Pauschalpreise, jedes Projekt bekommt ein
          eigenes Angebot.
        </p>
        <div className="flex flex-wrap justify-center gap-3">
          <a
            href="#twitch-bot"
            className="px-5 py-2 rounded-full border border-gray-700 hover:border-goalyBlue hover:text-goalyBlue transition-colors text-sm font-semibold"
          >
            Zum Twitch-Bot-Abo
          </a>
          <a
            href="#individuelle-entwicklung"
            className="px-5 py-2 rounded-full border border-gray-700 hover:border-goalyBlue hover:text-goalyBlue transition-colors text-sm font-semibold"
          >
            Zur individuellen Softwareentwicklung
          </a>
        </div>
      </section>

      <section id="twitch-bot" className="text-center pb-16 px-6 scroll-mt-20">
        <h2 className="text-2xl sm:text-3xl font-bold mb-2">Twitch Bot - Abo-Tarife</h2>
        <p className="text-gray-400 max-w-2xl mx-auto mb-8">
          Fertiges SaaS-Produkt mit festen monatlichen/jährlichen Tarifen - wähle den Plan, der zu deinem Kanal
          passt, jederzeit änderbar.
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

        <label className="flex items-start gap-2 max-w-md mx-auto mb-6 text-left text-sm text-gray-300 cursor-pointer">
          <input
            type="checkbox"
            checked={withdrawalConsent}
            onChange={(e) => setWithdrawalConsent(e.target.checked)}
            className="mt-1 w-4 h-4 flex-shrink-0"
          />
          <span>
            Ich bin damit einverstanden, dass mit der Ausführung des Vertrags bereits vor Ablauf der
            Widerrufsfrist begonnen wird. Mir ist bekannt, dass ich dadurch bei vollständiger
            Vertragserfüllung mein Widerrufsrecht verliere. Details in unserer{" "}
            <Link to="/legal/widerruf" className="underline text-goalyBlue">Widerrufsbelehrung</Link>.
          </span>
        </label>

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
                {tier.priceId && (
                  <p className="text-xs text-gray-500 -mt-5 mb-6">inkl. MwSt.</p>
                )}

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
                    Jetzt abonnieren
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

      {/* Bewusst visuell abgesetzt (eigener Hintergrund) vom SaaS-Abo oben,
          damit auf einen Blick klar ist: andere Leistung, andere Abrechnung -
          kein Tarif-Raster, keine Festpreise, sondern Angebot nach Anfrage. */}
      <section id="individuelle-entwicklung" className="bg-slate-900 py-16 px-6 scroll-mt-20">
        <div className="max-w-5xl mx-auto text-center">
          <h2 className="text-2xl sm:text-3xl font-bold mb-2">Individuelle Softwareentwicklung</h2>
          <p className="text-gray-400 max-w-2xl mx-auto mb-2">
            Neben dem Twitch-Bot-Abo entwickle ich als Freelancer maßgeschneiderte Software - Web, Mobile,
            Automatisierung und Streaming-Tools.
          </p>
          <p className="text-goalyBlue font-semibold max-w-2xl mx-auto mb-10">
            Dafür gibt es keine Pauschalpreise: jedes Projekt ist unterschiedlich und bekommt nach einem kurzen
            Erstgespräch ein individuelles Angebot.
          </p>

          <div className="grid sm:grid-cols-3 gap-6 mb-10 text-left">
            <div className="bg-gray-900 border border-gray-800 rounded-xl p-6">
              <h3 className="text-lg font-bold mb-2">Individuelle Softwarelösungen</h3>
              <p className="text-sm text-gray-400">
                Von Prototyp bis fertiges Produkt - maßgeschneiderte Anwendungen für Web, Mobile oder Desktop.
              </p>
            </div>
            <div className="bg-gray-900 border border-gray-800 rounded-xl p-6">
              <h3 className="text-lg font-bold mb-2">Automatisierung & Workflows</h3>
              <p className="text-sm text-gray-400">
                Bots, Integrationen und Tools, die Prozesse automatisieren - für Twitch, Discord oder interne
                Abläufe.
              </p>
            </div>
            <div className="bg-gray-900 border border-gray-800 rounded-xl p-6">
              <h3 className="text-lg font-bold mb-2">Streaming & Community Tech</h3>
              <p className="text-sm text-gray-400">
                Tools & Overlays, die deine Community stärker binden und neue Features ins Stream-Erlebnis bringen.
              </p>
            </div>
          </div>

          <Link
            to="/contact"
            className="inline-block px-8 py-3 bg-goalyBlue hover:bg-goalyCyan rounded-lg text-white font-semibold transition-colors"
          >
            Kostenloses Erstgespräch anfragen
          </Link>
          <p className="text-gray-500 text-sm mt-4">
            Mehr Details zu den Leistungen findest du auf der{" "}
            <Link to="/services" className="underline text-goalyBlue">Services-Seite</Link>.
          </p>
        </div>
      </section>
    </div>
  );
}
