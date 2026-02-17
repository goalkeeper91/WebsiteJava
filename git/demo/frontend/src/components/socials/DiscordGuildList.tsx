import React, { useState, useEffect } from "react";
import { Server, Settings, Trash2, RefreshCw, ExternalLink } from "lucide-react";

interface Guild {
  id: number;
  name: string;
  iconUrl?: string;
  memberCount?: number;
  isActive: boolean;
  botAddedAt?: string;
}

interface GuildWithSettings {
  guild: Guild;
  settings?: {
    id: number;
    notificationChannelId?: number;
    commandChannelId?: number;
    activityChannelId?: number;
    twitchNotificationsEnabled: boolean;
    joinToCreateEnabled: boolean;
  };
}

interface DiscordGuildListProps {
  onSelectGuild: (guildId: number) => void;
}

const DiscordGuildList: React.FC<DiscordGuildListProps> = ({ onSelectGuild }) => {
  const [guilds, setGuilds] = useState<GuildWithSettings[]>([]);
  const [loading, setLoading] = useState(true);
  const [inviteUrl, setInviteUrl] = useState<string>("");

  const fetchGuilds = async () => {
    try {
      setLoading(true);
      const res = await fetch("/api/discord/guilds", {
        credentials: "include",
      });
      if (res.ok) {
        const data = await res.json();
        setGuilds(data);
      }
    } catch (err) {
      console.error("Failed to fetch guilds:", err);
    } finally {
      setLoading(false);
    }
  };

  const fetchInviteUrl = async () => {
    try {
      const res = await fetch("/api/discord/bot/invite-url");
      if (res.ok) {
        const data = await res.json();
        setInviteUrl(data.inviteUrl);
      }
    } catch (err) {
      console.error("Failed to fetch invite URL:", err);
    }
  };

  useEffect(() => {
    fetchGuilds();
    fetchInviteUrl();
  }, []);

  const handleRemoveBot = async (guildId: number, guildName: string) => {
    if (!confirm(`Remove bot from "${guildName}"?`)) return;

    try {
      const res = await fetch(`/api/discord/guilds/${guildId}`, {
        method: "DELETE",
        credentials: "include",
      });
      if (res.ok) {
        fetchGuilds();
      }
    } catch (err) {
      console.error("Failed to remove bot:", err);
    }
  };

  const handleSync = async (guildId: number) => {
    try {
      const res = await fetch(`/api/discord/guilds/${guildId}/sync`, {
        method: "POST",
        credentials: "include",
      });
      if (res.ok) {
        fetchGuilds();
      }
    } catch (err) {
      console.error("Failed to sync guild:", err);
    }
  };

  if (loading) {
    return (
      <div className="p-6 bg-gray-800 rounded-xl border border-gray-700">
        <p className="text-gray-400">Loading guilds...</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Invite Bot Button */}
      {inviteUrl && guilds.length === 0 && (
        <div className="p-6 bg-indigo-900/30 border border-indigo-700/50 rounded-xl">
          <h3 className="font-semibold text-white mb-2">
            No Servers Yet
          </h3>
          <p className="text-sm text-gray-300 mb-4">
            Invite the bot to your Discord server to get started
          </p>
          <a
            href={inviteUrl}
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-2 px-4 py-2 bg-indigo-600 hover:bg-indigo-500 text-white rounded-lg transition-colors"
          >
            <ExternalLink className="w-4 h-4" />
            Invite Bot to Server
          </a>
        </div>
      )}

      {/* Guild Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {guilds.map(({ guild, settings }) => (
          <div
            key={guild.id}
            className="p-4 bg-gray-800 rounded-xl border border-gray-700 hover:border-gray-600 transition-all"
          >
            <div className="flex items-start justify-between mb-3">
              <div className="flex items-center gap-3">
                {guild.iconUrl ? (
                  <img
                    src={guild.iconUrl}
                    alt={guild.name}
                    className="w-12 h-12 rounded-full"
                  />
                ) : (
                  <div className="w-12 h-12 rounded-full bg-gray-700 flex items-center justify-center">
                    <Server className="w-6 h-6 text-gray-400" />
                  </div>
                )}
                <div>
                  <h3 className="font-semibold text-white">{guild.name}</h3>
                  {guild.memberCount && (
                    <p className="text-xs text-gray-400">
                      {guild.memberCount} members
                    </p>
                  )}
                </div>
              </div>
            </div>

            {/* Settings Status */}
            {settings && (
              <div className="mb-3 space-y-1">
                <div className="flex items-center justify-between text-xs">
                  <span className="text-gray-400">Notifications:</span>
                  <span className={settings.twitchNotificationsEnabled ? "text-green-400" : "text-gray-500"}>
                    {settings.twitchNotificationsEnabled ? "Enabled" : "Disabled"}
                  </span>
                </div>
                <div className="flex items-center justify-between text-xs">
                  <span className="text-gray-400">Join-to-Create:</span>
                  <span className={settings.joinToCreateEnabled ? "text-green-400" : "text-gray-500"}>
                    {settings.joinToCreateEnabled ? "Enabled" : "Disabled"}
                  </span>
                </div>
              </div>
            )}

            {/* Actions */}
            <div className="flex gap-2">
              <button
                onClick={() => onSelectGuild(guild.id)}
                className="flex-1 flex items-center justify-center gap-2 px-3 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors text-sm"
              >
                <Settings className="w-4 h-4" />
                Settings
              </button>
              <button
                onClick={() => handleSync(guild.id)}
                className="px-3 py-2 bg-blue-600/20 hover:bg-blue-600/30 text-blue-400 rounded-lg transition-colors"
                title="Sync guild data"
              >
                <RefreshCw className="w-4 h-4" />
              </button>
              <button
                onClick={() => handleRemoveBot(guild.id, guild.name)}
                className="px-3 py-2 bg-red-600/20 hover:bg-red-600/30 text-red-400 rounded-lg transition-colors"
                title="Remove bot"
              >
                <Trash2 className="w-4 h-4" />
              </button>
            </div>
          </div>
        ))}
      </div>

      {/* Add More Guilds */}
      {inviteUrl && guilds.length > 0 && (
        <div className="text-center">
          <a
            href={inviteUrl}
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-2 px-4 py-2 bg-gray-800 hover:bg-gray-700 text-gray-300 border border-gray-700 rounded-lg transition-colors text-sm"
          >
            <ExternalLink className="w-4 h-4" />
            Add Bot to Another Server
          </a>
        </div>
      )}
    </div>
  );
};

export default DiscordGuildList;