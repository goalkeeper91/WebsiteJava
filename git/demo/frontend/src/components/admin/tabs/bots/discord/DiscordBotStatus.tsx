import React, { useState, useEffect } from "react";
import { AlertTriangle, CheckCircle, RefreshCcw, LogOut } from "lucide-react";
import UptimeDisplay from "../../stats/UptimeDisplay";

interface Guild {
  id: string;
  name: string;
}

interface DiscordStatus {
  running: boolean;
  guildCount: number;
  totalMembers: number;
  uptimeSeconds?: number;
  guilds: Guild[];
}

const DiscordBotStatus: React.FC = () => {
  const [status, setStatus] = useState<DiscordStatus | null>(null);
  const [loading, setLoading] = useState(false);
  const [inviteUrl, setInviteUrl] = useState<string | null>(null);

  const fetchInviteUrl = async () => {
    try {
      const res = await fetch("/api/discord/bot/invite-url");
      if (!res.ok) throw new Error("Fehler beim Abrufen des Invite-Links");
      const data = await res.json();
      setInviteUrl(data.inviteUrl);
    } catch (err) {
      console.error(err);
    }
  };

  const fetchStatus = async () => {
    try {
      setLoading(true);
      const res = await fetch("/api/discord/bot/status", { credentials: "include" });
      if (!res.ok) throw new Error("Fehler beim Abrufen des Discord-Status");
      const data: DiscordStatus = await res.json();
      setStatus(data);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const syncGuild = async (guildId: string) => {
    try {
      const res = await fetch(`/api/discord/guilds/${guildId}/sync`, {
        method: "POST",
        credentials: "include",
      });
      if (!res.ok) throw new Error("Fehler beim Synchronisieren des Servers");
      fetchStatus();
    } catch (err) {
      console.error(err);
    }
  };

  const removeGuild = async (guildId: string, guildName: string) => {
    if (!confirm(`Bot wirklich von "${guildName}" entfernen?`)) return;

    try {
      const res = await fetch(`/api/discord/guilds/${guildId}`, {
        method: "DELETE",
        credentials: "include",
      });
      if (!res.ok) throw new Error("Fehler beim Trennen des Servers");
      fetchStatus();
    } catch (err) {
      console.error(err);
    }
  };

  useEffect(() => {
    fetchStatus();
    const interval = setInterval(fetchStatus, 5000);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="p-4 rounded-xl shadow space-y-4">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-lg font-bold">Discord Bot Status</h2>
          {status ? (
            <div className="text-sm text-gray-500">
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
              <p className="mt-1">
                {status.guildCount} Server · {status.totalMembers} Mitglieder
              </p>
            </div>
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

      <button
        onClick={fetchInviteUrl}
        className="bg-blue-500 hover:bg-blue-600 text-white px-3 py-2 rounded"
      >
        Bot zu Server hinzufügen
      </button>

      {inviteUrl && (
        <a
          href={inviteUrl}
          target="_blank"
          rel="noopener noreferrer"
          className="ml-2 text-blue-700 underline"
        >
          Link öffnen
        </a>
      )}

      {status?.guilds && status.guilds.length > 0 ? (
        <div>
          <h3 className="font-semibold mb-2">Verbundene Server</h3>
          <ul className="space-y-2">
            {status.guilds.map((guild) => (
              <li key={guild.id} className="flex justify-between items-center p-2 rounded">
                <span>{guild.name}</span>
                <div className="flex gap-2">
                  <button
                    onClick={() => syncGuild(guild.id)}
                    className="bg-gray-200 hover:bg-gray-300 text-gray-800 px-3 py-1 rounded"
                    title="Server synchronisieren"
                  >
                    <RefreshCcw className="w-4 h-4" />
                  </button>
                  <button
                    onClick={() => removeGuild(guild.id, guild.name)}
                    className="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded flex items-center"
                  >
                    <LogOut className="w-4 h-4 mr-1" /> Trennen
                  </button>
                </div>
              </li>
            ))}
          </ul>
        </div>
      ) : (
        <div className="flex items-center justify-center p-4 bg-yellow-50 text-yellow-700 rounded">
          <AlertTriangle className="w-5 h-5 mr-2" /> Keine verbundenen Server
        </div>
      )}
    </div>
  );
};

export default DiscordBotStatus;
