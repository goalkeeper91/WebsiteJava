import { useEffect, useState } from "react";
import { RefreshCw } from "lucide-react";
import { getAdminCustomers, getAdminCustomerStats, setAdminCustomerTier } from "../../../../features/admin/customers/api";
import type { AdminCustomer, AdminCustomerStats, TierId } from "../../../../features/admin/customers/types";

const TIER_OPTIONS: { value: TierId; label: string }[] = [
  { value: "free", label: "Free" },
  { value: "pro", label: "Pro" },
  { value: "premium", label: "Premium" },
];

const PAGE_SIZE = 20;

export default function CustomerSubscriptionManager() {
  const [customers, setCustomers] = useState<AdminCustomer[]>([]);
  const [stats, setStats] = useState<AdminCustomerStats | null>(null);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [savingTwitchId, setSavingTwitchId] = useState<string | null>(null);

  useEffect(() => {
    load(page);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page]);

  async function load(pageToLoad: number) {
    setLoading(true);
    setError(null);
    try {
      const [customersRes, statsRes] = await Promise.all([
        getAdminCustomers(pageToLoad, PAGE_SIZE),
        getAdminCustomerStats(),
      ]);
      setCustomers(customersRes.data ?? []);
      setTotalPages(customersRes.total_pages || 1);
      setStats(statsRes);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Fehler beim Laden der Kundenliste");
    } finally {
      setLoading(false);
    }
  }

  async function handleTierChange(twitchId: string, tierId: string) {
    setSavingTwitchId(twitchId);
    setError(null);
    try {
      await setAdminCustomerTier(twitchId, tierId);
      await load(page);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Fehler beim Ändern des Tarifs");
    } finally {
      setSavingTwitchId(null);
    }
  }

  if (loading && customers.length === 0) {
    return (
      <div className="animate-pulse space-y-4">
        <div className="h-24 bg-gray-800 rounded"></div>
        <div className="h-64 bg-gray-800 rounded"></div>
      </div>
    );
  }

  const detectorLoad = stats?.clip_detector
    ? `${stats.clip_detector.active_channels.length}/${stats.clip_detector.max_concurrent}`
    : "keine Daten";

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-xl font-bold">Kunden/Abos</h2>
        <button
          onClick={() => load(page)}
          className="p-2 hover:bg-gray-800 rounded-lg transition-colors"
          title="Aktualisieren"
        >
          <RefreshCw className="w-5 h-5" />
        </button>
      </div>

      {error && (
        <div className="bg-red-900/30 border border-red-500/50 text-red-300 rounded-lg p-3 text-sm">
          {error}
        </div>
      )}

      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <StatCard label="MRR (geschätzt)" value={`${(stats?.mrr ?? 0).toFixed(2)} €`} />
        <StatCard label="Aktive Kunden" value={String(stats?.active_customers ?? 0)} />
        <StatCard label="Clip-Detector-Auslastung" value={detectorLoad} />
      </div>

      <div className="bg-gray-800 rounded-xl overflow-hidden">
        {customers.length === 0 ? (
          <div className="p-8 text-center text-gray-500">Keine Kunden gefunden.</div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead className="bg-gray-900 text-gray-400">
                <tr>
                  <th className="text-left px-4 py-3 whitespace-nowrap">Login</th>
                  <th className="text-left px-4 py-3 whitespace-nowrap">E-Mail</th>
                  <th className="text-left px-4 py-3 whitespace-nowrap">Tarif</th>
                  <th className="text-left px-4 py-3 whitespace-nowrap">Status</th>
                  <th className="text-left px-4 py-3 whitespace-nowrap">Läuft ab</th>
                  <th className="text-left px-4 py-3 whitespace-nowrap">Paddle-Kunde</th>
                  <th className="text-left px-4 py-3 whitespace-nowrap">Aktion</th>
                </tr>
              </thead>
              <tbody>
                {customers.map((c) => (
                  <tr key={c.twitch_id} className="border-t border-gray-700">
                    <td className="px-4 py-3 whitespace-nowrap">
                      {c.username}
                      {c.is_admin && <span className="ml-2 text-xs text-purple-400">Admin</span>}
                    </td>
                    <td className="px-4 py-3 text-gray-400 whitespace-nowrap">{c.email || "–"}</td>
                    <td className="px-4 py-3 whitespace-nowrap">{c.tier_name}</td>
                    <td className="px-4 py-3 whitespace-nowrap">
                      <span
                        className={`px-2 py-1 rounded text-xs font-medium ${
                          c.is_active ? "bg-green-900/50 text-green-400" : "bg-red-900/50 text-red-400"
                        }`}
                      >
                        {c.is_active ? "Aktiv" : "Inaktiv"}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-gray-400 whitespace-nowrap">
                      {c.expires_at ? new Date(c.expires_at).toLocaleDateString() : "–"}
                    </td>
                    <td className="px-4 py-3 text-gray-400 whitespace-nowrap font-mono text-xs">
                      {c.paddle_customer_id || "–"}
                    </td>
                    <td className="px-4 py-3 whitespace-nowrap">
                      <select
                        defaultValue=""
                        disabled={savingTwitchId === c.twitch_id}
                        onChange={(e) => {
                          if (e.target.value) {
                            handleTierChange(c.twitch_id, e.target.value);
                            e.target.value = "";
                          }
                        }}
                        className="bg-gray-700 rounded px-2 py-1 text-xs disabled:opacity-50"
                      >
                        <option value="" disabled>
                          Tarif ändern...
                        </option>
                        {TIER_OPTIONS.map((opt) => (
                          <option key={opt.value} value={opt.value}>
                            {opt.label}
                          </option>
                        ))}
                      </select>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}

        {totalPages > 1 && (
          <div className="flex flex-wrap justify-between items-center gap-2 p-4 border-t border-gray-700">
            <button
              disabled={page <= 1}
              onClick={() => setPage((p) => p - 1)}
              className="px-3 py-1.5 bg-gray-700 rounded disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-600 transition-colors text-sm"
            >
              ← Zurück
            </button>
            <span className="text-gray-400 text-sm">
              Seite {page} von {totalPages}
            </span>
            <button
              disabled={page >= totalPages}
              onClick={() => setPage((p) => p + 1)}
              className="px-3 py-1.5 bg-gray-700 rounded disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-600 transition-colors text-sm"
            >
              Weiter →
            </button>
          </div>
        )}
      </div>
    </div>
  );
}

function StatCard({ label, value }: { label: string; value: string }) {
  return (
    <div className="bg-gray-800 rounded-xl p-4">
      <div className="text-2xl font-bold text-white">{value}</div>
      <div className="text-xs text-gray-400 mt-1">{label}</div>
    </div>
  );
}
