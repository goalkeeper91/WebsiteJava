import { useEffect, useState } from "react";
import { getLoyaltySettings, updateLoyaltySettings, getLoyaltyLeaderboard } from "../../features/loyalty/api";
import type { LeaderboardEntry } from "../../features/loyalty/types";

function formatRelative(dateString: string): string {
  const seconds = Math.max(0, Math.floor((Date.now() - new Date(dateString).getTime()) / 1000));
  if (seconds < 60) return "gerade eben";
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return `vor ${minutes} Min.`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `vor ${hours} Std.`;
  return `vor ${Math.floor(hours / 24)} Tag(en)`;
}

export default function LoyaltyDashboard() {
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [saved, setSaved] = useState(false);

  const [enabled, setEnabled] = useState(false);
  const [pointsName, setPointsName] = useState("Punkte");
  const [pointsPerInterval, setPointsPerInterval] = useState(1);
  const [intervalMinutes, setIntervalMinutes] = useState(5);
  const [regularsThreshold, setRegularsThreshold] = useState(0);

  const [entries, setEntries] = useState<LeaderboardEntry[]>([]);
  const [entriesPage, setEntriesPage] = useState(1);
  const [entriesTotalPages, setEntriesTotalPages] = useState(0);
  const [entriesLoading, setEntriesLoading] = useState(true);

  useEffect(() => {
    getLoyaltySettings()
      .then((settings) => {
        setEnabled(settings.enabled);
        setPointsName(settings.points_name || "Punkte");
        setPointsPerInterval(settings.points_per_interval || 1);
        setIntervalMinutes(settings.interval_minutes || 5);
        setRegularsThreshold(settings.regulars_threshold || 0);
      })
      .catch((err) => setError(err.message || "Fehler beim Laden"))
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => {
    setEntriesLoading(true);
    getLoyaltyLeaderboard(entriesPage, 10)
      .then((data) => {
        setEntries(data.data || []);
        setEntriesTotalPages(data.total_pages || 0);
      })
      .catch((err) => console.error("Fehler beim Laden der Loyalty-Bestenliste:", err))
      .finally(() => setEntriesLoading(false));
  }, [entriesPage]);

  async function handleSave() {
    setSaving(true);
    setError("");
    setSaved(false);

    try {
      await updateLoyaltySettings({
        enabled,
        points_name: pointsName,
        points_per_interval: pointsPerInterval,
        interval_minutes: intervalMinutes,
        regulars_threshold: regularsThreshold,
      });
      setSaved(true);
      setTimeout(() => setSaved(false), 3000);
    } catch (err: any) {
      setError(err.message || "Fehler beim Speichern");
    } finally {
      setSaving(false);
    }
  }

  if (loading) {
    return (
      <div className="max-w-3xl mx-auto p-4 sm:p-6 animate-pulse space-y-4">
        <div className="h-8 bg-gray-800 rounded w-1/3"></div>
        <div className="h-64 bg-gray-800 rounded"></div>
      </div>
    );
  }

  return (
    <div className="max-w-3xl mx-auto space-y-6 p-4 sm:p-0">
      <div>
        <h1 className="text-2xl font-bold">Loyalty / Watchtime</h1>
        <p className="text-sm text-gray-400 mt-1">
          Schreibt Zuschauern in regelmäßigen Abständen Punkte gut, solange sie im Chat anwesend
          sind und dein Stream live ist — auch stille Zuschauer (Lurker) zählen mit, nicht nur wer
          schreibt.
        </p>
      </div>

      {error && (
        <div className="bg-red-900/30 border border-red-500/50 text-red-300 rounded-lg p-3 text-sm">
          {error}
        </div>
      )}

      <div className="bg-gray-800 border border-gray-700 rounded-xl p-4 sm:p-6 space-y-6">
        <label className="flex items-center justify-between cursor-pointer gap-3">
          <div>
            <p className="font-semibold text-white">Loyalty aktiviert</p>
            <p className="text-xs text-gray-400 mt-0.5">
              Ohne diese Option werden keine Punkte vergeben.
            </p>
          </div>
          <input
            type="checkbox"
            checked={enabled}
            onChange={(e) => setEnabled(e.target.checked)}
            className="w-5 h-5 rounded bg-gray-700 border-gray-600 text-green-600 focus:ring-2 focus:ring-green-500 flex-shrink-0"
          />
        </label>

        <div className="grid sm:grid-cols-2 gap-4 pt-4 border-t border-gray-700">
          <div>
            <label className="block text-sm font-medium mb-1">Punkte-Name</label>
            <p className="text-xs text-gray-400 mb-2">Wird im Chat angezeigt, z.B. "Goldstücke".</p>
            <input
              type="text"
              value={pointsName}
              onChange={(e) => setPointsName(e.target.value)}
              className="w-full p-2 rounded bg-gray-900 border border-gray-700 text-white text-sm focus:outline-none focus:ring-2 focus:ring-green-500"
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">Punkte pro Intervall</label>
            <p className="text-xs text-gray-400 mb-2">Pro anwesendem Zuschauer, pro Intervall.</p>
            <input
              type="number"
              min={1}
              value={pointsPerInterval}
              onChange={(e) => setPointsPerInterval(Math.max(1, Number(e.target.value)))}
              className="w-full p-2 rounded bg-gray-900 border border-gray-700 text-white text-sm focus:outline-none focus:ring-2 focus:ring-green-500"
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">Intervall (Minuten)</label>
            <p className="text-xs text-gray-400 mb-2">
              Wie oft Punkte vergeben werden, solange dein Stream live ist.
            </p>
            <input
              type="number"
              min={1}
              value={intervalMinutes}
              onChange={(e) => setIntervalMinutes(Math.max(1, Number(e.target.value)))}
              className="w-full p-2 rounded bg-gray-900 border border-gray-700 text-white text-sm focus:outline-none focus:ring-2 focus:ring-green-500"
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">Regulars-Schwelle (Punkte)</label>
            <p className="text-xs text-gray-400 mb-2">
              Ab dieser Punktzahl gelten Zuschauer als "Regular" und werden automatisch von Automod
              ausgenommen (0 = deaktiviert).
            </p>
            <input
              type="number"
              min={0}
              value={regularsThreshold}
              onChange={(e) => setRegularsThreshold(Math.max(0, Number(e.target.value)))}
              className="w-full p-2 rounded bg-gray-900 border border-gray-700 text-white text-sm focus:outline-none focus:ring-2 focus:ring-green-500"
            />
          </div>
        </div>

        <div className="flex flex-wrap items-center gap-3 pt-2">
          <button
            onClick={handleSave}
            disabled={saving}
            className="px-4 py-2 bg-green-600 rounded-lg hover:bg-green-500 disabled:opacity-50 transition-colors"
          >
            {saving ? "Speichern…" : "Speichern"}
          </button>
          {saved && <span className="text-sm text-green-400">✓ Gespeichert</span>}
        </div>
      </div>

      <div className="bg-gray-800 border border-gray-700 rounded-xl p-4 sm:p-6">
        <h3 className="font-semibold text-white mb-3">Bestenliste</h3>

        {entriesLoading && <p className="text-sm text-gray-500">Lade…</p>}

        {!entriesLoading && entries.length === 0 && (
          <p className="text-sm text-gray-500">Noch keine Punkte vergeben.</p>
        )}

        <div className="space-y-2">
          {entries.map((entry, index) => (
            <div
              key={entry.viewer_twitch_id}
              className="bg-gray-900 rounded-lg p-3 text-sm flex flex-wrap items-center justify-between gap-2"
            >
              <span className="text-white">
                <span className="text-gray-500 mr-2">
                  #{(entriesPage - 1) * 10 + index + 1}
                </span>
                {entry.viewer_login}
                {regularsThreshold > 0 && entry.points >= regularsThreshold && (
                  <span className="ml-2 text-xs text-yellow-400">⭐ Regular</span>
                )}
              </span>
              <span className="flex items-center gap-3">
                <span className="font-semibold text-green-400">{entry.points}</span>
                <span className="text-xs text-gray-500">{formatRelative(entry.updated_at)}</span>
              </span>
            </div>
          ))}
        </div>

        {!entriesLoading && entriesTotalPages > 1 && (
          <div className="flex flex-wrap justify-between items-center gap-2 mt-4 pt-4 border-t border-gray-700">
            <button
              disabled={entriesPage <= 1}
              onClick={() => setEntriesPage((p) => p - 1)}
              className="px-3 py-1.5 bg-gray-700 rounded disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-600 transition-colors text-sm"
            >
              ← Zurück
            </button>
            <span className="text-gray-400 text-sm">
              Seite {entriesPage} von {entriesTotalPages}
            </span>
            <button
              disabled={entriesPage >= entriesTotalPages}
              onClick={() => setEntriesPage((p) => p + 1)}
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
