import React, { useState, useEffect } from "react";
import { AlertTriangle, CheckCircle, RefreshCcw } from "lucide-react";
import UptimeDisplay from "../../stats/UptimeDisplay";

interface DiscordStatus {
  running: boolean;
  since: string | null;
  uptimeSeconds?: number;
}

const DiscordBotStatus: React.FC = () => {
  const [status, setStatus] = useState<DiscordStatus | null>(null);
  const [loading, setLoading] = useState(false);

  const fetchStatus = async () => {
    try {
      setLoading(true);
      const res = await fetch("http://localhost:8080/api/discord/status");
      if (!res.ok) throw new Error("Fehler beim Abrufen des Discord-Status");
      const data: DiscordStatus = await res.json();
      setStatus(data);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchStatus();
    const interval = setInterval(fetchStatus, 5000);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="p-4 rounded-xl shadow flex items-center justify-between">
      <div>
        <h2 className="text-lg font-bold">Discord Bot Status</h2>
        {status ? (
          <p className="text-sm text-gray-500">
            {status.running ? (
              <span className="flex items-center text-green-600">
                <CheckCircle className="w-5 h-5 mr-1" /> Läuft
              </span>
            ) : (
              <span className="flex items-center text-red-600">
                <AlertTriangle className="w-5 h-5 mr-1" /> Gestoppt
              </span>
            )}
            {status.running && status.uptimeSeconds !== undefined && (
              <UptimeDisplay uptimeSeconds={status.uptimeSeconds} />
            )}
          </p>
        ) : (
          <p className="text-gray-400">Noch keine Daten geladen</p>
        )}
      </div>

      <div className="flex p-1">
        <button
          onClick={fetchStatus}
          disabled={loading}
          className="bg-gray-200 hover:bg-gray-300 text-gray-800 px-3 py-2 rounded"
        >
          <RefreshCcw className={`w-4 h-4 ${loading ? "animate-spin" : ""}`} />
        </button>
      </div>
    </div>
  );
};

export default DiscordBotStatus;
