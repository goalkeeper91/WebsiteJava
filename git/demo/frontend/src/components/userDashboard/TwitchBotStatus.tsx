import { useEffect, useState } from "react";
import { CheckCircle, AlertTriangle, XCircle } from "lucide-react";
import { useAuth } from "../../context/AuthContext";

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
    <div className="p-6 bg-gray-800 rounded-xl border border-gray-700 space-y-1">
      <h3 className="text-lg font-bold text-white mb-2">Twitch-Bot-Status</h3>

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
  );
}
