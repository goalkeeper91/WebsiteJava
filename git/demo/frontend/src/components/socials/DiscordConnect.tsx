import React, { useState, useEffect } from "react";
import { CheckCircle, Link as LinkIcon, LogOut } from "lucide-react";

interface ConnectionStatus {
  connected: boolean;
  discordUsername?: string;
  discordDiscriminator?: string;
  discordUserId?: string;
}

const DiscordConnect: React.FC = () => {
  const [status, setStatus] = useState<ConnectionStatus | null>(null);
  const [loading, setLoading] = useState(true);

  const fetchStatus = async () => {
    try {
      const res = await fetch("/api/discord/auth/status", {
        credentials: "include",
      });
      if (res.ok) {
        const data = await res.json();
        setStatus(data);
      }
    } catch (err) {
      console.error("Failed to fetch Discord status:", err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchStatus();
  }, []);

  const handleConnect = async () => {
    try {
      const res = await fetch("/api/discord/auth/url", {
        credentials: "include",
      });
      if (res.ok) {
        const data = await res.json();
        localStorage.setItem("discord_oauth_state", data.state);
        window.location.href = data.authUrl;
      }
    } catch (err) {
      console.error("Failed to get auth URL:", err);
    }
  };

  const handleDisconnect = async () => {
    if (!confirm("Are you sure you want to disconnect Discord?")) return;

    try {
      const res = await fetch("/api/discord/auth/disconnect", {
        method: "DELETE",
        credentials: "include",
      });
      if (res.ok) {
        setStatus({ connected: false });
      }
    } catch (err) {
      console.error("Failed to disconnect:", err);
    }
  };

  if (loading) {
    return (
      <div className="p-6 bg-gray-800 rounded-xl border border-gray-700">
        <p className="text-gray-400">Loading Discord connection...</p>
      </div>
    );
  }

  return (
    <div className="p-6 bg-gray-800 rounded-xl border border-gray-700 space-y-4">
      <div className="flex items-center justify-between">
        <div>
          <h3 className="text-lg font-bold text-white">Discord Connection</h3>
          {status?.connected ? (
            <div className="flex items-center gap-2 mt-1">
              <CheckCircle className="w-5 h-5 text-green-500" />
              <span className="text-sm text-gray-300">
                Connected as{" "}
                <span className="font-semibold text-white">
                  {status.discordUsername}#{status.discordDiscriminator}
                </span>
              </span>
            </div>
          ) : (
            <p className="text-sm text-gray-400 mt-1">
              Connect your Discord account to enable bot features
            </p>
          )}
        </div>

        {status?.connected ? (
          <button
            onClick={handleDisconnect}
            className="flex items-center gap-2 px-4 py-2 bg-red-600 hover:bg-red-500 text-white rounded-lg transition-colors"
          >
            <LogOut className="w-4 h-4" />
            Disconnect
          </button>
        ) : (
          <button
            onClick={handleConnect}
            className="flex items-center gap-2 px-4 py-2 bg-indigo-600 hover:bg-indigo-500 text-white rounded-lg transition-colors"
          >
            <LinkIcon className="w-4 h-4" />
            Connect Discord
          </button>
        )}
      </div>

      {status?.connected && (
        <div className="pt-4 border-t border-gray-700">
          <p className="text-xs text-gray-500">
            Discord User ID: {status.discordUserId}
          </p>
        </div>
      )}
    </div>
  );
};

export default DiscordConnect;