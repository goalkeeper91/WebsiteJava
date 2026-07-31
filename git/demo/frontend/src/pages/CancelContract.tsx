import { useEffect, useState } from "react";
import { useAuth } from "../context/AuthContext";
import { getSubscription, getPortalLink } from "../features/subscription/api";
import type { UserSubscription } from "../features/subscription/types";
import Seo from "../components/Seo";

// Public "Kündigungsbutton" page required by § 312k BGB for online contracts
// with a recurring payment obligation - must be permanently, easily
// findable (linked from the footer on every page) and end in a clearly
// labeled "Jetzt kündigen" button. The actual cancellation itself still has
// to happen with Paddle (our Merchant of Record and the contracting party
// for the payment), so the button deep-links straight into Paddle's own
// cancel_subscription portal page rather than re-implementing cancellation
// logic here.
export default function CancelContract() {
  const { isAuthenticated, authChecked, username, email } = useAuth();
  const [subscription, setSubscription] = useState<UserSubscription | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [portalLoading, setPortalLoading] = useState(false);
  const [portalError, setPortalError] = useState("");

  useEffect(() => {
    if (!authChecked || !isAuthenticated) {
      setLoading(false);
      return;
    }
    getSubscription()
      .then(setSubscription)
      .catch((err) => setError(err.message || "Fehler beim Laden des Abos"))
      .finally(() => setLoading(false));
  }, [authChecked, isAuthenticated]);

  async function handleCancel() {
    setPortalError("");
    setPortalLoading(true);
    try {
      const url = await getPortalLink();
      window.open(url, "_blank", "noopener,noreferrer");
    } catch (err) {
      setPortalError(err instanceof Error ? err.message : "Kündigung konnte nicht gestartet werden");
    } finally {
      setPortalLoading(false);
    }
  }

  const isPaid = !!subscription && subscription.tierId !== "free";

  return (
    <div className="min-h-screen bg-slate-950 text-white flex items-start justify-center py-16 px-4">
      <Seo
        title="Vertrag kündigen"
        description="Kündige dein Goalkeeper91-Abonnement jederzeit zum Ende des Abrechnungszeitraums - direkt hier oder über das Paddle-Kundenportal."
        path="/vertrag-kuendigen"
      />
      <div className="max-w-2xl w-full bg-slate-900 rounded-xl shadow-lg p-8 space-y-6">
        <div>
          <h1 className="text-3xl font-bold mb-2">Vertrag kündigen</h1>
          <p className="text-gray-400 text-sm">
            Hier kannst du ein laufendes, kostenpflichtiges Abonnement bei Goalkeeper91 kündigen.
          </p>
        </div>

        {!authChecked ? (
          <p className="text-sm text-gray-500">Lade…</p>
        ) : !isAuthenticated ? (
          <div className="space-y-3">
            <p className="text-sm text-gray-300">
              Um dein Abonnement zu kündigen, melde dich bitte zuerst mit deinem Twitch-Account an - wir
              müssen wissen, welcher Vertrag betroffen ist.
            </p>
            <a
              href={`/auth/login?returnTo=${encodeURIComponent("/vertrag-kuendigen")}`}
              className="inline-block px-4 py-2 bg-indigo-600 hover:bg-indigo-500 rounded-lg transition-colors text-sm font-medium"
            >
              Jetzt einloggen
            </a>
          </div>
        ) : loading ? (
          <p className="text-sm text-gray-500">Lade…</p>
        ) : error ? (
          <p className="text-sm text-red-400">{error}</p>
        ) : (
          <div className="space-y-4">
            <div className="bg-slate-800 rounded-lg p-4 text-sm space-y-1">
              <p>
                <span className="text-gray-400">Vertrag:</span>{" "}
                <span className="font-semibold">
                  Goalkeeper91 {subscription?.tier?.name || "Free"}-Abonnement
                  {subscription?.billingCycle === "yearly" ? " (jährlich)" : subscription?.billingCycle === "monthly" ? " (monatlich)" : ""}
                </span>
              </p>
              <p>
                <span className="text-gray-400">Kunde:</span>{" "}
                <span className="font-semibold">{username}{email ? ` (${email})` : ""}</span>
              </p>
              {subscription?.startedAt && (
                <p>
                  <span className="text-gray-400">Vertragsbeginn:</span>{" "}
                  <span className="font-semibold">{new Date(subscription.startedAt).toLocaleDateString()}</span>
                </p>
              )}
            </div>

            {portalError && (
              <div className="bg-red-900/30 border border-red-500/50 text-red-300 rounded-lg p-3 text-sm">
                {portalError}
              </div>
            )}

            {isPaid ? (
              <>
                <p className="text-sm text-gray-300">
                  Die Kündigung wird über unseren Zahlungsdienstleister Paddle abgewickelt (öffnet in einem neuen
                  Tab). Dein Zugriff bleibt bis zum Ende des bereits bezahlten Abrechnungszeitraums bestehen.
                </p>
                <button
                  onClick={handleCancel}
                  disabled={portalLoading}
                  className="w-full sm:w-auto px-6 py-3 bg-red-600 hover:bg-red-500 rounded-lg transition-colors font-semibold disabled:opacity-50"
                >
                  {portalLoading ? "Öffne…" : "Jetzt kündigen"}
                </button>
              </>
            ) : (
              <p className="text-sm text-gray-300">
                Du hast aktuell kein kostenpflichtiges Abonnement, das gekündigt werden müsste.
              </p>
            )}
          </div>
        )}

        <p className="text-xs text-gray-500 pt-4 border-t border-slate-800">
          Weitere Informationen findest du in unseren{" "}
          <a href="/legal/agb" className="underline text-goalyBlue">Nutzungsbedingungen</a> und der{" "}
          <a href="/legal/widerruf" className="underline text-goalyBlue">Widerrufsbelehrung</a>.
        </p>
      </div>
    </div>
  );
}
