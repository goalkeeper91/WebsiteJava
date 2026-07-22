import { useEffect, useState } from "react";
import { Film, ExternalLink, RefreshCw } from "lucide-react";
import { getClips, getClipStats } from "../../features/clips/api";
import type { ClipLog, ClipLogStats, ClipStatus } from "../../features/clips/types";

const STATUS_STYLES: Record<ClipStatus, string> = {
  pending: "bg-gray-700 text-gray-300",
  processing: "bg-yellow-900/50 text-yellow-400",
  completed: "bg-green-900/50 text-green-400",
  failed: "bg-red-900/50 text-red-400",
  cancelled: "bg-gray-700 text-gray-400",
};

export default function ClipHistoryPanel() {
  const [clips, setClips] = useState<ClipLog[]>([]);
  const [stats, setStats] = useState<ClipLogStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    load();
  }, []);

  async function load() {
    setLoading(true);
    setError(null);
    try {
      const [clipsRes, statsRes] = await Promise.all([getClips(1, 20), getClipStats()]);
      setClips(clipsRes.data ?? []);
      setStats(statsRes);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Fehler beim Laden der Clips");
    } finally {
      setLoading(false);
    }
  }

  if (loading) {
    return (
      <div className="animate-pulse space-y-4">
        <div className="h-24 bg-gray-800 rounded"></div>
        <div className="h-64 bg-gray-800 rounded"></div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Film className="w-5 h-5 text-purple-400" />
          <h2 className="text-xl font-bold">Clip-Historie</h2>
        </div>
        <button
          onClick={load}
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

      {stats && (
        <div className="grid grid-cols-2 md:grid-cols-5 gap-4">
          <StatCard label="Gesamt" value={stats.total_clips} />
          <StatCard label="Ausstehend" value={stats.pending_clips} />
          <StatCard label="In Bearbeitung" value={stats.processing_clips} />
          <StatCard label="Fertig" value={stats.completed_clips} accent="text-green-500" />
          <StatCard label="Fehlgeschlagen" value={stats.failed_clips} accent="text-red-500" />
        </div>
      )}

      <div className="bg-gray-800 rounded-xl overflow-hidden">
        {clips.length === 0 ? (
          <div className="p-8 text-center text-gray-500">
            Noch keine Clips erkannt. Sobald die Erkennung live ist, tauchen deine Highlights hier auf.
          </div>
        ) : (
          <table className="w-full text-sm">
            <thead className="bg-gray-900 text-gray-400">
              <tr>
                <th className="text-left px-4 py-3">Titel</th>
                <th className="text-left px-4 py-3">Spiel</th>
                <th className="text-left px-4 py-3">Dauer</th>
                <th className="text-left px-4 py-3">Status</th>
                <th className="text-left px-4 py-3">Erstellt</th>
                <th className="text-left px-4 py-3"></th>
              </tr>
            </thead>
            <tbody>
              {clips.map((clip) => (
                <tr key={clip.id} className="border-t border-gray-700">
                  <td className="px-4 py-3">{clip.clip_title || "Ohne Titel"}</td>
                  <td className="px-4 py-3 text-gray-400">{clip.game_name || "–"}</td>
                  <td className="px-4 py-3 text-gray-400">{clip.duration_seconds}s</td>
                  <td className="px-4 py-3">
                    <span className={`px-2 py-1 rounded text-xs font-medium ${STATUS_STYLES[clip.status]}`}>
                      {clip.status}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-gray-400">
                    {new Date(clip.created_at).toLocaleString()}
                  </td>
                  <td className="px-4 py-3">
                    <a
                      href={clip.clip_url}
                      target="_blank"
                      rel="noreferrer"
                      className="text-purple-400 hover:text-purple-300 inline-flex items-center gap-1"
                    >
                      Ansehen <ExternalLink className="w-3 h-3" />
                    </a>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}

function StatCard({ label, value, accent }: { label: string; value: number; accent?: string }) {
  return (
    <div className="bg-gray-800 rounded-xl p-4">
      <div className={`text-2xl font-bold ${accent || "text-white"}`}>{value}</div>
      <div className="text-xs text-gray-400 mt-1">{label}</div>
    </div>
  );
}
