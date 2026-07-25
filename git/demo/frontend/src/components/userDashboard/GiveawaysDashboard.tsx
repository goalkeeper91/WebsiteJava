import { useEffect, useState } from "react";
import { getGiveawayStatus, getGiveawayHistory } from "../../features/giveaways/api";
import type { Giveaway } from "../../features/giveaways/types";

function formatRelative(dateString: string): string {
  const seconds = Math.max(0, Math.floor((Date.now() - new Date(dateString).getTime()) / 1000));
  if (seconds < 60) return "gerade eben";
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return `vor ${minutes} Min.`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `vor ${hours} Std.`;
  return `vor ${Math.floor(hours / 24)} Tag(en)`;
}

export default function GiveawaysDashboard() {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [current, setCurrent] = useState<Giveaway | null>(null);
  const [entryCount, setEntryCount] = useState(0);

  const [history, setHistory] = useState<Giveaway[]>([]);
  const [historyPage, setHistoryPage] = useState(1);
  const [historyTotalPages, setHistoryTotalPages] = useState(0);
  const [historyLoading, setHistoryLoading] = useState(true);

  useEffect(() => {
    getGiveawayStatus()
      .then((data) => {
        setCurrent(data.giveaway);
        setEntryCount(data.entry_count);
      })
      .catch((err) => setError(err.message || "Fehler beim Laden"))
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => {
    setHistoryLoading(true);
    getGiveawayHistory(historyPage, 10)
      .then((data) => {
        setHistory(data.data || []);
        setHistoryTotalPages(data.total_pages || 0);
      })
      .catch((err) => console.error("Fehler beim Laden der Giveaway-Historie:", err))
      .finally(() => setHistoryLoading(false));
  }, [historyPage]);

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
        <h1 className="text-2xl font-bold">Giveaways</h1>
        <p className="text-sm text-gray-400 mt-1">
          <code className="text-gray-300">!giveaway start &lt;codewort&gt;</code> startet eine Verlosung,{" "}
          <code className="text-gray-300">!giveaway draw</code> zieht einen zufälligen Gewinner. Teilnehmen ist{" "}
          <span className="text-gray-300">kein</span> Chat-Command — Zuschauer tippen einfach das Codewort in den
          Chat.
        </p>
      </div>

      {error && (
        <div className="bg-red-900/30 border border-red-500/50 text-red-300 rounded-lg p-3 text-sm">
          {error}
        </div>
      )}

      <div className="bg-gray-800 border border-gray-700 rounded-xl p-4 sm:p-6">
        <h3 className="font-semibold text-white mb-3">Aktuelles Giveaway</h3>
        {current ? (
          <div className="flex flex-wrap items-center justify-between gap-2 text-sm">
            <span className="text-green-400 font-semibold">● Läuft</span>
            <span className="text-gray-300">
              Codewort: <code className="text-white">{current.keyword}</code>
            </span>
            <span className="text-gray-300">{entryCount} Teilnehmer</span>
            <span className="text-gray-400">{current.sub_bonus ? "Sub-Bonus aktiv" : "Kein Sub-Bonus"}</span>
          </div>
        ) : (
          <p className="text-sm text-gray-500">
            Kein Giveaway aktiv — starte eins im Chat mit <code className="text-gray-400">!giveaway start</code>.
          </p>
        )}
      </div>

      <div className="bg-gray-800 border border-gray-700 rounded-xl p-4 sm:p-6">
        <h3 className="font-semibold text-white mb-3">Verlauf</h3>

        {historyLoading && <p className="text-sm text-gray-500">Lade…</p>}

        {!historyLoading && history.length === 0 && (
          <p className="text-sm text-gray-500">Noch keine Giveaways durchgeführt.</p>
        )}

        <div className="space-y-2">
          {history.map((giveaway) => (
            <div
              key={giveaway.id}
              className="bg-gray-900 rounded-lg p-3 text-sm flex flex-wrap items-center justify-between gap-2"
            >
              <span className="text-white">
                {giveaway.winner_login ? (
                  <>
                    🏆 <span className="font-semibold">{giveaway.winner_login}</span>
                  </>
                ) : (
                  <span className="text-gray-500">Kein Gewinner</span>
                )}
              </span>
              <span className="flex items-center gap-3 text-xs text-gray-500">
                <span>
                  Codewort: <code className="text-gray-400">{giveaway.keyword}</code>
                </span>
                <span>{giveaway.entry_count} Teilnehmer</span>
                <span>{formatRelative(giveaway.started_at)}</span>
              </span>
            </div>
          ))}
        </div>

        {!historyLoading && historyTotalPages > 1 && (
          <div className="flex flex-wrap justify-between items-center gap-2 mt-4 pt-4 border-t border-gray-700">
            <button
              disabled={historyPage <= 1}
              onClick={() => setHistoryPage((p) => p - 1)}
              className="px-3 py-1.5 bg-gray-700 rounded disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-600 transition-colors text-sm"
            >
              ← Zurück
            </button>
            <span className="text-gray-400 text-sm">
              Seite {historyPage} von {historyTotalPages}
            </span>
            <button
              disabled={historyPage >= historyTotalPages}
              onClick={() => setHistoryPage((p) => p + 1)}
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
