import { useState, useEffect } from "react";
import {
  Play,
  Square,
  AlertTriangle,
  CheckCircle,
  RefreshCcw,
  Activity,
  Users,
  MessageSquare,
  Zap,
  TrendingUp,
  Settings,
  BarChart3
} from "lucide-react";
import BotAuthCard from "./BotAuthButton";

interface BotStatus {
  running: boolean;
  tokenPresent: boolean;
  expiresAt: string | null;
  channelName: string | null;
  activeChannels?: string[];
  uptimeSeconds?: number;
}

interface BotStats {
  commands: {
    total: number;
    today: number;
    success: number;
    failed: number;
  };
  messages: {
    total: number;
    today: number;
  };
  activeChatters: number;
  topCommands?: Array<{
    name: string;
    count: number;
  }>;
}

const TwitchBotDashboard = () => {
  const [status, setStatus] = useState<BotStatus | null>(null);
  const [stats, setStats] = useState<BotStats | null>(null);
  const [loading, setLoading] = useState(false);
  const [autoRefresh, setAutoRefresh] = useState(true);

  useEffect(() => {
    fetchStatus();
    fetchStats();

    if (autoRefresh) {
      const interval = setInterval(() => {
        fetchStatus();
        fetchStats();
      }, 5000);
      return () => clearInterval(interval);
    }
  }, [autoRefresh]);

  const fetchStatus = async () => {
    try {
      const res = await fetch("/api/bot/status", {
        credentials: "include",
      });
      if (res.ok) {
        const data: BotStatus = await res.json();
        setStatus(data);
      }
    } catch (err) {
      console.error("Fehler beim Laden des Status:", err);
    }
  };

  const fetchStats = async () => {
    try {
      const res = await fetch("/api/bot/stats", {
        credentials: "include",
      });
      if (res.ok) {
        const data: BotStats = await res.json();
        setStats(data);
      }
    } catch (err) {
      console.error("Fehler beim Laden der Stats:", err);
    }
  };

  const startBot = async () => {
    setLoading(true);
    try {
      const res = await fetch("/api/bot/start", {
        method: "POST",
        credentials: "include",
      });
      if (res.ok) {
        await fetchStatus();
      } else {
        alert("Fehler beim Starten des Bots");
      }
    } catch (err) {
      console.error("Fehler beim Starten:", err);
      alert("Fehler beim Starten des Bots");
    } finally {
      setLoading(false);
    }
  };

  const stopBot = async () => {
    if (!confirm("Bot wirklich stoppen?")) return;

    setLoading(true);
    try {
      const res = await fetch("/api/bot/stop", {
        method: "POST",
        credentials: "include",
      });
      if (res.ok) {
        await fetchStatus();
      } else {
        alert("Fehler beim Stoppen des Bots");
      }
    } catch (err) {
      console.error("Fehler beim Stoppen:", err);
      alert("Fehler beim Stoppen des Bots");
    } finally {
      setLoading(false);
    }
  };

  const restartBot = async () => {
    if (!confirm("Bot neu starten?")) return;

    setLoading(true);
    try {
      const res = await fetch("/api/bot/restart", {
        method: "POST",
        credentials: "include",
      });
      if (res.ok) {
        await fetchStatus();
      } else {
        alert("Fehler beim Neustarten des Bots");
      }
    } catch (err) {
      console.error("Fehler beim Neustarten:", err);
      alert("Fehler beim Neustarten des Bots");
    } finally {
      setLoading(false);
    }
  };

  const formatUptime = (seconds?: number) => {
    if (!seconds) return "0s";

    const days = Math.floor(seconds / 86400);
    const hours = Math.floor((seconds % 86400) / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    const secs = seconds % 60;

    if (days > 0) return `${days}d ${hours}h ${minutes}m`;
    if (hours > 0) return `${hours}h ${minutes}m`;
    if (minutes > 0) return `${minutes}m ${secs}s`;
    return `${secs}s`;
  };

  const getStatusColor = () => {
    if (!status) return "bg-gray-500";
    return status.running ? "bg-green-500" : "bg-red-500";
  };

  const getStatusText = () => {
    if (!status) return "Unbekannt";
    return status.running ? "Online" : "Offline";
  };

  return (
    <div className="space-y-6 p-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold flex items-center gap-3">
            <Activity className="w-8 h-8 text-purple-500" />
            Twitch Bot Dashboard
          </h1>
          <p className="text-gray-400 mt-1">
            Verwalte und überwache deinen Twitch Bot in Echtzeit
          </p>
        </div>

        {/* Auto-Refresh Toggle */}
        <label className="flex items-center gap-2 cursor-pointer bg-gray-800 px-4 py-2 rounded-lg hover:bg-gray-700 transition-colors">
          <input
            type="checkbox"
            checked={autoRefresh}
            onChange={(e) => setAutoRefresh(e.target.checked)}
            className="w-4 h-4 accent-purple-500"
          />
          <RefreshCcw className={`w-4 h-4 ${autoRefresh ? "text-green-400" : "text-gray-400"}`} />
          <span className="text-sm">Auto-Refresh (5s)</span>
        </label>
      </div>

      {/* Bot Authentication Card */}
      <BotAuthCard />

      {/* Status Overview Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {/* Bot Status */}
        <div className="p-5 rounded-xl shadow-lg bg-gradient-to-br from-gray-800 to-gray-900 border border-gray-700">
          <div className="flex items-center justify-between mb-3">
            <span className="text-sm text-gray-400 font-medium">Bot Status</span>
            <div className={`w-3 h-3 rounded-full ${getStatusColor()} ${status?.running ? 'animate-pulse' : ''}`}></div>
          </div>
          <p className="text-3xl font-bold mb-1">{getStatusText()}</p>
          {status?.running && status.uptimeSeconds !== undefined && (
            <p className="text-xs text-gray-400">
              🕒 Uptime: {formatUptime(status.uptimeSeconds)}
            </p>
          )}
        </div>

        {/* Active Channels */}
        <div className="p-5 rounded-xl shadow-lg bg-gradient-to-br from-purple-900/30 to-gray-900 border border-purple-700/30">
          <div className="flex items-center gap-2 mb-3">
            <Users className="w-5 h-5 text-purple-400" />
            <span className="text-sm text-gray-400 font-medium">Channels</span>
          </div>
          <p className="text-3xl font-bold mb-1 text-purple-300">
            {status?.activeChannels?.length || 0}
          </p>
          <p className="text-xs text-gray-400">Verbunden</p>
        </div>

        {/* Commands Today */}
        <div className="p-5 rounded-xl shadow-lg bg-gradient-to-br from-blue-900/30 to-gray-900 border border-blue-700/30">
          <div className="flex items-center gap-2 mb-3">
            <Zap className="w-5 h-5 text-blue-400" />
            <span className="text-sm text-gray-400 font-medium">Commands</span>
          </div>
          <p className="text-3xl font-bold mb-1 text-blue-300">
            {stats?.commands.today || 0}
          </p>
          <p className="text-xs text-gray-400">Heute ausgeführt</p>
        </div>

        {/* Messages Total */}
        <div className="p-5 rounded-xl shadow-lg bg-gradient-to-br from-green-900/30 to-gray-900 border border-green-700/30">
          <div className="flex items-center gap-2 mb-3">
            <MessageSquare className="w-5 h-5 text-green-400" />
            <span className="text-sm text-gray-400 font-medium">Nachrichten</span>
          </div>
          <p className="text-3xl font-bold mb-1 text-green-300">
            {stats?.messages.total || 0}
          </p>
          <p className="text-xs text-gray-400">Verarbeitet</p>
        </div>
      </div>

      {/* Main Control Panel */}
      <div className="p-6 rounded-xl shadow-lg bg-gray-800 border border-gray-700">
        <div className="flex items-center justify-between mb-6">
          <div>
            <h2 className="text-2xl font-bold flex items-center gap-2">
              <Zap className="w-6 h-6 text-yellow-400" />
              Bot Steuerung
            </h2>
            <p className="text-sm text-gray-400 mt-1">
              Starte, stoppe oder starte den Bot neu
            </p>
          </div>

          <button
            onClick={fetchStatus}
            disabled={loading}
            className="p-2 hover:bg-gray-700 rounded-lg transition-colors disabled:opacity-50"
            title="Status aktualisieren"
          >
            <RefreshCcw className={`w-6 h-6 ${loading ? "animate-spin" : ""}`} />
          </button>
        </div>

        {/* Status Banner */}
        <div className={`mb-6 p-4 rounded-lg ${
          status?.running
            ? 'bg-green-500/10 border border-green-500/30'
            : 'bg-red-500/10 border border-red-500/30'
        }`}>
          <div className="flex items-center gap-3">
            {status?.running ? (
              <>
                <CheckCircle className="w-7 h-7 text-green-500" />
                <div className="flex-1">
                  <p className="font-semibold text-green-400 text-lg">Bot läuft</p>
                  <p className="text-sm text-gray-400">
                    {status.channelName
                      ? `Verbunden mit: ${status.channelName}`
                      : "Bereit für Channels"
                    }
                  </p>
                </div>
              </>
            ) : (
              <>
                <AlertTriangle className="w-7 h-7 text-red-500" />
                <div className="flex-1">
                  <p className="font-semibold text-red-400 text-lg">Bot ist offline</p>
                  <p className="text-sm text-gray-400">
                    {!status?.tokenPresent
                      ? "⚠️ Kein Token vorhanden - Bot authentifizieren erforderlich"
                      : "Bereit zum Starten"
                    }
                  </p>
                </div>
              </>
            )}
          </div>
        </div>

        {/* Control Buttons */}
        <div className="flex flex-wrap gap-3">
          {!status?.running ? (
            <button
              onClick={startBot}
              disabled={loading || !status?.tokenPresent}
              className="flex-1 min-w-[200px] bg-green-600 hover:bg-green-500 disabled:bg-gray-600 disabled:cursor-not-allowed text-white px-6 py-4 rounded-lg transition-all transform hover:scale-105 disabled:hover:scale-100 flex items-center justify-center gap-2 font-semibold text-lg shadow-lg"
            >
              <Play className="w-6 h-6" />
              {loading ? "Starte..." : "Bot starten"}
            </button>
          ) : (
            <>
              <button
                onClick={stopBot}
                disabled={loading}
                className="flex-1 min-w-[150px] bg-red-600 hover:bg-red-500 disabled:bg-gray-600 text-white px-6 py-4 rounded-lg transition-all transform hover:scale-105 disabled:hover:scale-100 flex items-center justify-center gap-2 font-semibold shadow-lg"
              >
                <Square className="w-5 h-5" />
                {loading ? "Stoppe..." : "Stoppen"}
              </button>

              <button
                onClick={restartBot}
                disabled={loading}
                className="flex-1 min-w-[150px] bg-yellow-600 hover:bg-yellow-500 disabled:bg-gray-600 text-white px-6 py-4 rounded-lg transition-all transform hover:scale-105 disabled:hover:scale-100 flex items-center justify-center gap-2 font-semibold shadow-lg"
              >
                <RefreshCcw className="w-5 h-5" />
                {loading ? "Starte neu..." : "Neustart"}
              </button>
            </>
          )}
        </div>

        {/* Token Warning */}
        {!status?.tokenPresent && (
          <div className="mt-4 p-4 bg-yellow-500/10 border border-yellow-500/30 rounded-lg">
            <p className="text-sm text-yellow-400 flex items-center gap-2">
              <AlertTriangle className="w-4 h-4" />
              Bot kann nicht gestartet werden: Kein Token vorhanden.
              Bitte authentifiziere den Bot oben in der "Bot Account" Card.
            </p>
          </div>
        )}
      </div>

      {/* Active Channels & Stats Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Active Channels */}
        <div className="p-6 rounded-xl shadow-lg bg-gray-800 border border-gray-700">
          <h3 className="text-xl font-bold mb-4 flex items-center gap-2">
            <Users className="w-6 h-6 text-purple-400" />
            Aktive Channels
          </h3>

          {status?.activeChannels && status.activeChannels.length > 0 ? (
            <div className="space-y-2">
              {status.activeChannels.map((channel) => (
                <div
                  key={channel}
                  className="p-4 bg-gray-700/50 rounded-lg flex items-center gap-3 hover:bg-gray-700 transition-colors"
                >
                  <div className="w-3 h-3 bg-green-500 rounded-full animate-pulse"></div>
                  <div className="flex-1">
                    <p className="font-semibold text-lg">{channel}</p>
                    <p className="text-xs text-gray-400">Verbunden</p>
                  </div>
                  <CheckCircle className="w-5 h-5 text-green-400" />
                </div>
              ))}
            </div>
          ) : (
            <div className="text-center py-12">
              <Users className="w-16 h-16 text-gray-600 mx-auto mb-4" />
              <p className="text-gray-400 font-medium mb-2">Keine aktiven Channels</p>
              <p className="text-sm text-gray-500">
                {status?.running
                  ? "Bot läuft, aber keine Channels verbunden"
                  : "Starte den Bot um Channels zu verbinden"
                }
              </p>
            </div>
          )}
        </div>

        {/* Statistics */}
        <div className="p-6 rounded-xl shadow-lg bg-gray-800 border border-gray-700">
          <h3 className="text-xl font-bold mb-4 flex items-center gap-2">
            <BarChart3 className="w-6 h-6 text-blue-400" />
            Statistiken
          </h3>

          <div className="space-y-4">
            <div className="p-4 bg-gray-700/50 rounded-lg">
              <div className="flex items-center justify-between mb-2">
                <span className="text-sm text-gray-400 flex items-center gap-2">
                  <MessageSquare className="w-4 h-4" />
                  Gesamt Commands
                </span>
                <span className="text-2xl font-bold">{stats?.commands.total || 0}</span>
              </div>
              <div className="w-full bg-gray-600 rounded-full h-2">
                <div
                  className="bg-blue-500 h-2 rounded-full transition-all duration-500"
                  style={{
                    width: `${Math.min(100, ((stats?.commands.total || 0) / 100) * 100)}%`
                  }}
                ></div>
              </div>
            </div>

            <div className="p-4 bg-gray-700/50 rounded-lg">
              <div className="flex items-center justify-between mb-2">
                <span className="text-sm text-gray-400 flex items-center gap-2">
                  <Zap className="w-4 h-4" />
                  Heute ausgeführt
                </span>
                <span className="text-2xl font-bold text-green-400">
                  {stats?.commands.today || 0}
                </span>
              </div>
              <div className="w-full bg-gray-600 rounded-full h-2">
                <div
                  className="bg-green-500 h-2 rounded-full transition-all duration-500"
                  style={{
                    width: `${Math.min(100, ((stats?.commands.today || 0) / 50) * 100)}%`
                  }}
                ></div>
              </div>
            </div>

            <div className="p-4 bg-gray-700/50 rounded-lg">
              <div className="flex items-center justify-between mb-2">
                <span className="text-sm text-gray-400 flex items-center gap-2">
                  <TrendingUp className="w-4 h-4" />
                  Nachrichten verarbeitet
                </span>
                <span className="text-2xl font-bold text-purple-400">
                  {stats?.messages.total || 0}
                </span>
              </div>
              <div className="w-full bg-gray-600 rounded-full h-2">
                <div
                  className="bg-purple-500 h-2 rounded-full transition-all duration-500"
                  style={{
                    width: `${Math.min(100, ((stats?.messages.total || 0) / 500) * 100)}%`
                  }}
                ></div>
              </div>
            </div>

            {/* Active Chatters */}
            {status?.running && (
              <div className="p-4 bg-gray-700/50 rounded-lg">
                <div className="flex items-center justify-between">
                  <span className="text-sm text-gray-400 flex items-center gap-2">
                    <Users className="w-4 h-4" />
                    Active Chatters
                  </span>
                  <span className="text-2xl font-bold text-yellow-400">
                    {stats?.activeChatters || 0}
                  </span>
                </div>
                <p className="text-xs text-gray-500 mt-1">Letzte 10 Minuten</p>
              </div>
            )}

            {/* Success Rate */}
            {stats && stats.commands.total > 0 && (
              <div className="p-4 bg-gray-700/50 rounded-lg">
                <div className="flex items-center justify-between mb-2">
                  <span className="text-sm text-gray-400 flex items-center gap-2">
                    <CheckCircle className="w-4 h-4" />
                    Erfolgsrate
                  </span>
                  <span className="text-2xl font-bold text-green-400">
                    {Math.round((stats.commands.success / stats.commands.total) * 100)}%
                  </span>
                </div>
                <div className="w-full bg-gray-600 rounded-full h-2">
                  <div
                    className="bg-green-500 h-2 rounded-full transition-all duration-500"
                    style={{
                      width: `${(stats.commands.success / stats.commands.total) * 100}%`
                    }}
                  ></div>
                </div>
                <p className="text-xs text-gray-500 mt-1">
                  {stats.commands.failed} fehlgeschlagen
                </p>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Top Commands Section */}
      {stats?.topCommands && stats.topCommands.length > 0 && (
        <div className="p-6 rounded-xl shadow-lg bg-gray-800 border border-gray-700">
          <h3 className="text-xl font-bold mb-4 flex items-center gap-2">
            <TrendingUp className="w-6 h-6 text-yellow-400" />
            Top Commands
          </h3>
          <div className="space-y-2">
            {stats.topCommands.slice(0, 5).map((cmd, idx) => (
              <div key={cmd.name} className="flex items-center gap-3 p-3 bg-gray-700/50 rounded-lg hover:bg-gray-700 transition-colors">
                <div className="flex items-center justify-center w-8 h-8 bg-yellow-500/20 rounded-full">
                  <span className="text-yellow-400 font-bold">#{idx + 1}</span>
                </div>
                <div className="flex-1">
                  <p className="font-semibold">{cmd.name}</p>
                </div>
                <span className="text-xl font-bold text-yellow-400">{cmd.count}</span>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default TwitchBotDashboard;