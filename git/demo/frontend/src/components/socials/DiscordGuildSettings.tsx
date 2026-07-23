import React, { useState, useEffect } from "react";
import { Save, Send, ArrowLeft } from "lucide-react";

interface GuildSettings {
  id: number;
  guildId: string;
  userId: string;
  notificationChannelId?: string;
  commandChannelId?: string;
  activityChannelId?: string;
  twitchNotificationsEnabled: boolean;
  joinToCreateEnabled: boolean;
  adminRoleId?: string;
}

interface Channel {
  id: string;
  name: string;
  type: string;
}

interface Role {
  id: string;
  name: string;
  color: number;
}

interface GuildDetails {
  guild: {
    id: string;
    name: string;
    iconUrl?: string;
  };
  channels: Channel[];
  roles: Role[];
  settings?: GuildSettings;
}

interface DiscordGuildSettingsProps {
  guildId: string;
  onBack: () => void;
}

const DiscordGuildSettings: React.FC<DiscordGuildSettingsProps> = ({ guildId, onBack }) => {
  const [details, setDetails] = useState<GuildDetails | null>(null);
  const [settings, setSettings] = useState<Partial<GuildSettings>>({
    twitchNotificationsEnabled: true,
    joinToCreateEnabled: true,
  });
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  const fetchGuildDetails = async () => {
    try {
      setLoading(true);
      const res = await fetch(`/api/discord/guilds/${guildId}/details`, {
        credentials: "include",
      });
      if (res.ok) {
        const data = await res.json();
        setDetails(data);
        if (data.settings) {
          setSettings(data.settings);
        }
      }
    } catch (err) {
      console.error("Failed to fetch guild details:", err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchGuildDetails();
  }, [guildId]);

  const handleSave = async () => {
    try {
      setSaving(true);
      const res = await fetch(`/api/discord/guilds/${guildId}/settings`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify(settings),
      });
      if (res.ok) {
        alert("Settings saved successfully!");
      }
    } catch (err) {
      console.error("Failed to save settings:", err);
      alert("Failed to save settings");
    } finally {
      setSaving(false);
    }
  };

  const handleTestNotification = async () => {
    if (!settings.notificationChannelId) {
      alert("Please select a notification channel first");
      return;
    }

    try {
      const res = await fetch(`/api/discord/guilds/${guildId}/test-notification`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ channelId: settings.notificationChannelId }),
      });
      if (res.ok) {
        alert("Test notification sent! Check your Discord channel.");
      }
    } catch (err) {
      console.error("Failed to send test notification:", err);
      alert("Failed to send test notification");
    }
  };

  if (loading) {
    return (
      <div className="p-6 bg-gray-800 border border-gray-700 rounded-xl">
        <p className="text-gray-400">Loading settings...</p>
      </div>
    );
  }

  if (!details) {
    return (
      <div className="p-6 bg-gray-800 border border-gray-700 rounded-xl">
        <p className="text-red-400">Failed to load guild details</p>
      </div>
    );
  }

  const textChannels = details.channels.filter((ch) => ch.type === "text" || ch.type === "GUILD_TEXT");

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-4">
        <button
          onClick={onBack}
          className="p-2 hover:bg-gray-700 rounded-lg transition-colors text-gray-300"
        >
          <ArrowLeft className="w-5 h-5" />
        </button>
        <div className="flex items-center gap-3">
          {details.guild.iconUrl && (
            <img
              src={details.guild.iconUrl}
              alt={details.guild.name}
              className="w-12 h-12 rounded-full"
            />
          )}
          <div>
            <h2 className="text-xl font-bold text-white">{details.guild.name}</h2>
            <p className="text-sm text-gray-400">Guild Settings</p>
          </div>
        </div>
      </div>

      {/* Settings Form */}
      <div className="p-6 bg-gray-800 border border-gray-700 rounded-xl space-y-6">
        {/* Feature Toggles */}
        <div className="space-y-4">
          <h3 className="font-semibold text-white">Features</h3>

          <label className="flex items-center justify-between p-4 bg-gray-900 rounded-lg cursor-pointer hover:bg-gray-900/70 transition-colors">
            <div>
              <p className="font-medium text-white">Twitch Notifications</p>
              <p className="text-sm text-gray-400">
                Receive notifications for new commands and stream activity
              </p>
            </div>
            <input
              type="checkbox"
              checked={settings.twitchNotificationsEnabled ?? true}
              onChange={(e) =>
                setSettings({ ...settings, twitchNotificationsEnabled: e.target.checked })
              }
              className="w-5 h-5 accent-indigo-600 rounded focus:ring-2 focus:ring-indigo-500"
            />
          </label>

          <label className="flex items-center justify-between p-4 bg-gray-900 rounded-lg cursor-pointer hover:bg-gray-900/70 transition-colors">
            <div>
              <p className="font-medium text-white">Join-to-Create</p>
              <p className="text-sm text-gray-400">
                Enable automatic voice channel creation
              </p>
            </div>
            <input
              type="checkbox"
              checked={settings.joinToCreateEnabled ?? true}
              onChange={(e) =>
                setSettings({ ...settings, joinToCreateEnabled: e.target.checked })
              }
              className="w-5 h-5 accent-indigo-600 rounded focus:ring-2 focus:ring-indigo-500"
            />
          </label>
        </div>

        {/* Channel Configuration */}
        <div className="space-y-4">
          <h3 className="font-semibold text-white">Notification Channels</h3>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              General Notifications
            </label>
            <select
              value={settings.notificationChannelId ?? ""}
              onChange={(e) =>
                setSettings({
                  ...settings,
                  notificationChannelId: e.target.value || undefined,
                })
              }
              className="w-full bg-gray-900 border border-gray-700 text-white p-2 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
            >
              <option value="">Select a channel...</option>
              {textChannels.map((ch) => (
                <option key={ch.id} value={ch.id}>
                  #{ch.name}
                </option>
              ))}
            </select>
            <p className="mt-1 text-xs text-gray-500">
              Where general bot notifications will be sent
            </p>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              Command Notifications
            </label>
            <select
              value={settings.commandChannelId ?? ""}
              onChange={(e) =>
                setSettings({
                  ...settings,
                  commandChannelId: e.target.value || undefined,
                })
              }
              className="w-full bg-gray-900 border border-gray-700 text-white p-2 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
            >
              <option value="">Select a channel...</option>
              {textChannels.map((ch) => (
                <option key={ch.id} value={ch.id}>
                  #{ch.name}
                </option>
              ))}
            </select>
            <p className="mt-1 text-xs text-gray-500">
              Where new Twitch command notifications will be sent
            </p>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              Activity Notifications
            </label>
            <select
              value={settings.activityChannelId ?? ""}
              onChange={(e) =>
                setSettings({
                  ...settings,
                  activityChannelId: e.target.value || undefined,
                })
              }
              className="w-full bg-gray-900 border border-gray-700 text-white p-2 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
            >
              <option value="">Select a channel...</option>
              {textChannels.map((ch) => (
                <option key={ch.id} value={ch.id}>
                  #{ch.name}
                </option>
              ))}
            </select>
            <p className="mt-1 text-xs text-gray-500">
              Where stream activity (follows, subs, etc.) will be sent
            </p>
          </div>
        </div>

        {/* Admin Role */}
        <div>
          <label className="block text-sm font-medium text-gray-300 mb-2">
            Admin Role (Optional)
          </label>
          <select
            value={settings.adminRoleId ?? ""}
            onChange={(e) =>
              setSettings({
                ...settings,
                adminRoleId: e.target.value || undefined,
              })
            }
            className="w-full bg-gray-900 border border-gray-700 text-white p-2 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
          >
            <option value="">No role required</option>
            {details.roles.map((role) => (
              <option key={role.id} value={role.id}>
                {role.name}
              </option>
            ))}
          </select>
          <p className="mt-1 text-xs text-gray-500">
            Role required to use admin bot commands
          </p>
        </div>

        {/* Action Buttons */}
        <div className="flex gap-3">
          <button
            onClick={handleSave}
            disabled={saving}
            className="flex-1 flex items-center justify-center gap-2 px-4 py-2 bg-indigo-600 hover:bg-indigo-500 disabled:bg-indigo-400/50 text-white rounded-lg transition-colors"
          >
            <Save className="w-4 h-4" />
            {saving ? "Saving..." : "Save Settings"}
          </button>
          <button
            onClick={handleTestNotification}
            disabled={!settings.notificationChannelId}
            className="flex items-center gap-2 px-4 py-2 bg-gray-700 hover:bg-gray-600 disabled:bg-gray-800 disabled:text-gray-500 text-gray-200 rounded-lg transition-colors"
          >
            <Send className="w-4 h-4" />
            Test Notification
          </button>
        </div>
      </div>
    </div>
  );
};

export default DiscordGuildSettings;