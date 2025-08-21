import React, { useState, useEffect } from "react";
import { AlertTriangle, CheckCircle, RefreshCcw, LogOut } from "lucide-react";
import UptimeDisplay from "../../stats/UptimeDisplay";

interface DiscordStatus {
  running?: boolean;
  since: string | null;
  uptimeSeconds?: number;
  guilds?: Guild[];
}

interface Guild {
  id: string;
  name: string;
}

const DiscordBotStatus: React.FC = () => {
  const [status, setStatus] = useState<DiscordStatus | null>(null);
  const [loading, setLoading] = useState(false);
  const [inviteLink, setInviteLink] = useState<string | null>(null);

  const fetchInviteLink = async () => {
    try {
      const res = await fetch("/api/discord/guild/invite-link");
      if (!res.ok) throw new Error("Fehler beim Abrufen des Invite-Links");
      const data = await res.json();
      setInviteLink(data.inviteLink);
    } catch (err) {
      console.error(err);
    }
  };

  const fetchStatus = async () => {
    try {
      setLoading(true);
      const res = await fetch("/api/discord/status");
      if (!res.ok) throw new Error("Fehler beim Abrufen des Discord-Status");
      const data: DiscordStatus = await res.json();
      setStatus(data);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const syncGuilds = async () => {
    try {
      setLoading(true);
      const res = await fetch("/api/discord/guild/sync", {
        method: "POST",
      });
      if (!res.ok) throw new Error("Fehler beim Synchronisieren der Guilds");

      const data: Guild[] = await res.json(); // aktualisierte Guilds aus Backend
      setStatus((prev) =>
        prev
          ? { ...prev, guilds: data }
          : { running: false, since: null, guilds: data } // Fallback falls prev null
      );
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const disconnectGuild = async (guildId: string) => {
    try {
      const res = await fetch(`/api/discord/guilds/${guildId}/disconnect`, {
        method: "POST",
      });
      if (!res.ok) throw new Error("Fehler beim Trennen des Servers");
      fetchStatus(); // Status neu laden
    } catch (err) {
      console.error(err);
    }
  };

  useEffect(() => {
    fetchStatus();
    syncGuilds();
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
        onClick={fetchInviteLink}
        className="bg-blue-500 hover:bg-blue-600 text-white px-3 py-2 rounded"
      >
        Bot zu Server hinzufügen
      </button>

      {inviteLink && (
        <a
          href={inviteLink}
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
            {status.guilds.map((guild, index) => (
              <li key={`${guild.id}-${index}`} className="flex justify-between items-center p-2 rounded">
                <span>{guild.name}</span>
                <button
                  onClick={() => disconnectGuild(guild.id)}
                  className="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded flex items-center"
                >
                  <LogOut className="w-4 h-4 mr-1" /> Trennen
                </button>
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
