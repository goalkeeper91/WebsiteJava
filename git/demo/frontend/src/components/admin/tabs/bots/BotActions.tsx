import React, { useState, useEffect } from "react";
import { Play, Square, AlertTriangle, CheckCircle, RefreshCcw } from "lucide-react";
import UptimeDisplay from "../stats/UptimeDisplay";

interface BotStatus {
  running: boolean;
  tokenPresent: boolean;
  expiresAt: string | null;
  channelName: string | null;
  activeChannels?: string[];
  uptimeSeconds?: number;
}

const TwitchBotControls = () => {
  const [status, setStatus] = useState<BotStatus | null>(null);
  const [loading, setLoading] = useState(false);

  const fetchStatus = async () => {
    try {
      const res = await fetch("http://localhost:8080/api/bot/status");
      if (!res.ok) throw new Error("Fehler beim Abrufen des Status");
      const data: BotStatus = await res.json();
      setStatus(data);
    } catch (err) {
      console.error(err);
    }
  };

  const startBot = async () => {
    setLoading(true);
    await fetch("http://localhost:8080/api/bot/start", { method: "POST" });
    await fetchStatus();
    setLoading(false);
  };

  const stopBot = async () => {
    setLoading(true);
    await fetch("http://localhost:8080/api/bot/stop", { method: "POST" });
    await fetchStatus();
    setLoading(false);
  };

  useEffect(() => {
    fetchStatus();
    const interval = setInterval(fetchStatus, 5000); // Live-Refresh alle 5s
    return () => clearInterval(interval);
  }, []);

  const renderTokenProgress = () => {
    if (!status?.expiresAt) return <p>Kein Token vorhanden</p>;
    const expiryDate = new Date(status.expiresAt);
    const now = new Date();
    const totalMs = expiryDate.getTime() - now.getTime();
    const daysLeft = totalMs / (1000 * 60 * 60 * 24);
    const percent = Math.max(0, Math.min(100, (daysLeft / 60) * 100)); // 60 Tage max

    return (
      <div className="w-full rounded-full h-3 mt-2">
        <div
          className={`h-3 rounded-full ${percent > 30 ? "bg-green-500" : "bg-red-500"}`}
          style={{ width: `${percent}%` }}
        />
        <p className="text-sm mt-1">
          Token läuft ab am {expiryDate.toLocaleString()} ({daysLeft.toFixed(1)} Tage)
        </p>
      </div>
    );
  };

  return (
    <div className="space-y-4">
      {/* Status Card */}
      <div className="p-4 rounded-xl shadow flex items-center justify-between">
        <div>
          <h2 className="text-lg font-bold">Bot Status</h2>
          <p className="text-sm text-gray-500">
            {status?.running ? (
              <span className="flex items-center text-green-600">
                <CheckCircle className="w-5 h-5 mr-1" /> Läuft
              </span>
            ) : (
              <span className="flex items-center text-red-600">
                <AlertTriangle className="w-5 h-5 mr-1" /> Gestoppt
              </span>
            )}
            {status?.running && status.uptimeSeconds !== undefined && (
              <UptimeDisplay uptimeSeconds={status.uptimeSeconds} />
            )}
          </p>
        </div>
        <div className="space-x-2">
          {!status?.running ? (
            <button
              onClick={startBot}
              disabled={loading}
              className="bg-green-500 hover:bg-green-600 text-white px-4 py-2 rounded flex items-center"
            >
              <Play className="w-4 h-4 mr-1" /> Starten
            </button>
          ) : (
            <button
              onClick={stopBot}
              disabled={loading}
              className="bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded flex items-center"
            >
              <Square className="w-4 h-4 mr-1" /> Stoppen
            </button>
          )}
          <button
            onClick={fetchStatus}
            className="bg-gray-200 hover:bg-gray-300 text-gray-800 px-3 py-2 rounded"
          >
            <RefreshCcw className="w-4 h-4" />
          </button>
        </div>
      </div>

      {/* Token Card */}
      <div className="p-4 rounded-xl shadow">
        <h3 className="font-bold mb-2">Token Info</h3>
        {status?.tokenPresent ? (
          renderTokenProgress()
        ) : (
          <p className="text-red-600">Kein Token vorhanden</p>
        )}
      </div>

      {/* Channel Card */}
      <div className="p-4 rounded-xl shadow">
        <h3 className="font-bold mb-2">Aktive Channels</h3>
        {status?.activeChannels?.length ? (
          <ul className="list-disc pl-5">
            {status.activeChannels.map((ch) => (
              <li key={ch}>{ch}</li>
            ))}
          </ul>
        ) : (
          <p className="text-gray-500">Keine Channels verbunden</p>
        )}
      </div>
    </div>
  );
};

export default TwitchBotControls;
