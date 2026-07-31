import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { LogOut, Shield } from "lucide-react";
import { useAuth } from "../../context/AuthContext";
import { getSubscription, getPortalLink } from "../../features/subscription/api";
import type { UserSubscription } from "../../features/subscription/types";

export default function SettingsPage() {
  const { username, email, twitchId, isAdmin, logout } = useAuth();

  const [subscription, setSubscription] = useState<UserSubscription | null>(null);
  const [subLoading, setSubLoading] = useState(true);
  const [subError, setSubError] = useState("");
  const [portalLoading, setPortalLoading] = useState(false);
  const [portalError, setPortalError] = useState("");

  useEffect(() => {
    getSubscription()
      .then(setSubscription)
      .catch((err) => setSubError(err.message || "Fehler beim Laden des Abos"))
      .finally(() => setSubLoading(false));
  }, []);

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

  const isFree = !subscription || subscription.tierId === "free";
  const isCanceled = subscription?.status === "canceled";

  return (
    <div className="min-h-screen bg-gray-900 text-white">
    <div className="max-w-2xl mx-auto space-y-6 p-4 sm:p-6">
      <div>
        <h1 className="text-2xl font-bold">Einstellungen</h1>
        <p className="text-sm text-gray-400 mt-1">Account-Details und Abo-Verwaltung.</p>
      </div>

      {/* Account */}
      <div className="bg-gray-800 border border-gray-700 rounded-xl p-4 sm:p-6 space-y-4">
        <h3 className="font-semibold text-white">Account</h3>

        <div className="flex items-center gap-3">
          <div className="w-12 h-12 rounded-full bg-purple-600 flex items-center justify-center text-lg font-bold flex-shrink-0">
            {username?.charAt(0).toUpperCase() || "U"}
          </div>
          <div className="min-w-0">
            <p className="font-semibold text-white truncate">{username || "—"}</p>
            {isAdmin && (
              <span className="inline-flex items-center gap-1 text-xs text-purple-400 mt-0.5">
                <Shield className="w-3 h-3" /> Admin
              </span>
            )}
          </div>
        </div>

        <div className="grid sm:grid-cols-2 gap-3 text-sm">
          <div>
            <p className="text-gray-500">E-Mail</p>
            <p className="text-white truncate">{email || "Nicht verfügbar"}</p>
          </div>
          <div>
            <p className="text-gray-500">Twitch-ID</p>
            <p className="text-white truncate">{twitchId || "—"}</p>
          </div>
        </div>

        <button
          onClick={logout}
          className="flex items-center gap-2 px-4 py-2 bg-gray-700 hover:bg-red-600 rounded-lg transition-colors text-sm font-medium"
        >
          <LogOut className="w-4 h-4" />
          Abmelden
        </button>
      </div>

      {/* Abo & Rechnung */}
      <div className="bg-gray-800 border border-gray-700 rounded-xl p-4 sm:p-6 space-y-4">
        <h3 className="font-semibold text-white">Abo & Rechnung</h3>

        {subError && (
          <div className="bg-red-900/30 border border-red-500/50 text-red-300 rounded-lg p-3 text-sm">{subError}</div>
        )}
        {portalError && (
          <div className="bg-red-900/30 border border-red-500/50 text-red-300 rounded-lg p-3 text-sm">{portalError}</div>
        )}

        {subLoading ? (
          <p className="text-sm text-gray-500">Lade…</p>
        ) : (
          <>
            <div className="text-sm space-y-1">
              <p className="text-gray-400">
                Aktueller Plan: <span className="font-semibold text-white">{subscription?.tier?.name || "Free"}</span>
              </p>
              {subscription?.expiresAt && !isCanceled && (
                <p className="text-gray-500">
                  Verlängert sich am {new Date(subscription.expiresAt).toLocaleDateString()}
                </p>
              )}
              {isCanceled && subscription?.expiresAt && (
                <p className="text-red-400">
                  ⚠️ Gekündigt - läuft bis {new Date(subscription.expiresAt).toLocaleDateString()}
                </p>
              )}
            </div>

            <div className="flex flex-wrap gap-3">
              {!isFree && (
                <button
                  onClick={handleManageSubscription}
                  disabled={portalLoading}
                  className="px-4 py-2 bg-gray-700 hover:bg-gray-600 rounded-lg transition-colors text-sm font-medium disabled:opacity-50"
                >
                  {portalLoading ? "Öffne…" : "Abo verwalten, kündigen & Rechnungen einsehen"}
                </button>
              )}
              <Link
                to="/dashboard/subscription"
                className="px-4 py-2 bg-indigo-600 hover:bg-indigo-500 rounded-lg transition-colors text-sm font-medium"
              >
                {isFree ? "Tarif upgraden" : "Tarif wechseln"}
              </Link>
            </div>
          </>
        )}
      </div>
    </div>
    </div>
  );
}
