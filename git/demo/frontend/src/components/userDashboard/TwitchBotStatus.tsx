import { useEffect, useState } from "react";
import { CheckCircle, AlertTriangle, XCircle } from "lucide-react";
import { useAuth } from "../../context/AuthContext";
import { useTeam } from "../../context/TeamContext";

interface BotStatus {
  running: boolean;
  activeChannels?: string[];
  uptimeSeconds?: number;
}

function formatUptime(seconds: number): string {
  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  if (hours) return `${hours}h ${minutes}m`;
  return `${minutes}m`;
}

// Purely informational - unlike Discord there's no separate "connect" action
// here, Twitch login already is the platform's primary sign-in, so the bot
// automatically joins a channel once its owner is logged in.
export default function TwitchBotStatus() {
  const { username } = useAuth();
  const { actingAsChannel, managedChannels, setActingAsChannel } = useTeam();
  const [status, setStatus] = useState<BotStatus | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch("/api/bot/status", { credentials: "include" })
      .then((res) => (res.ok ? res.json() : null))
      .then(setStatus)
      .catch((err) => console.error("Failed to fetch bot status:", err))
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return (
      <div className="p-6 bg-gray-800 rounded-xl border border-gray-700">
        <p className="text-gray-400">Lade Bot-Status...</p>
      </div>
    );
  }

  const isMyChannelConnected = !!(
    username &&
    status?.activeChannels?.some((c) => c.toLowerCase() === username.toLowerCase())
  );

  return (
    <div className="space-y-3">
      {actingAsChannel && (
        <div className="flex flex-wrap items-center justify-between gap-2 p-3 bg-amber-900/30 border border-amber-500/50 rounded-lg text-sm text-amber-300">
          <span>
            Du verwaltest gerade den Kanal von <span className="font-semibold">{actingAsChannel.owner_login}</span>.
          </span>
          <button
            onClick={() => setActingAsChannel(null)}
            className="px-3 py-1 bg-amber-700/50 rounded hover:bg-amber-700 transition-colors text-xs font-semibold"
          >
            Zurück zu meinem Kanal
          </button>
        </div>
      )}

      <div className="p-6 bg-gray-800 rounded-xl border border-gray-700 space-y-1">
        <div className="flex flex-wrap items-center justify-between gap-2 mb-2">
          <h3 className="text-lg font-bold text-white">Twitch-Bot-Status</h3>
          {managedChannels.length > 0 && (
            <select
              value={actingAsChannel?.owner_twitch_id ?? ""}
              onChange={(e) => {
                const selected = managedChannels.find((c) => c.owner_twitch_id === e.target.value);
                setActingAsChannel(selected ?? null);
              }}
              className="bg-gray-900 border border-gray-700 rounded px-2 py-1 text-sm text-white"
            >
              <option value="">Mein eigener Kanal</option>
              {managedChannels.map((c) => (
                <option key={c.owner_twitch_id} value={c.owner_twitch_id}>
                  Kanal von {c.owner_login}
                </option>
              ))}
            </select>
          )}
        </div>

        {!status?.running && (
          <div className="flex items-center gap-2">
            <XCircle className="w-5 h-5 text-red-500" />
            <span className="text-sm text-gray-300">Bot ist aktuell offline.</span>
          </div>
        )}

        {status?.running && isMyChannelConnected && (
          <div className="flex items-center gap-2">
            <CheckCircle className="w-5 h-5 text-green-500" />
            <span className="text-sm text-gray-300">
              Dein Chat ist verbunden
              {status.uptimeSeconds !== undefined && (
                <span className="text-gray-500"> · seit {formatUptime(status.uptimeSeconds)}</span>
              )}
            </span>
          </div>
        )}

        {status?.running && !isMyChannelConnected && (
          <div className="flex items-center gap-2">
            <AlertTriangle className="w-5 h-5 text-yellow-500" />
            <span className="text-sm text-gray-300">
              Bot läuft, dein Kanal ist aber noch nicht verbunden. Das kann kurz nach dem Login dauern.
            </span>
          </div>
        )}
      </div>
    </div>
  );
}
