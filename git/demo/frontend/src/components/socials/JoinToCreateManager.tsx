import React, { useEffect, useState } from "react";
import { Plus, Trash2, Save, AlertTriangle } from "lucide-react";

interface JoinToCreateConfig {
  id: number;
  guildId: string;
  joinChannelId: string;
  categoryId: string;
  channelNamePrefix: string;
  userLimit: number;
  privateChannel: boolean;
  enabled: boolean;
}

interface JoinToCreateManagerProps {
  guildId: string;
}

const JoinToCreateManager: React.FC<JoinToCreateManagerProps> = ({ guildId }) => {
  const [configs, setConfigs] = useState<JoinToCreateConfig[]>([]);
  const [loading, setLoading] = useState(true);
  const [showNewForm, setShowNewForm] = useState(false);
  const [newConfig, setNewConfig] = useState<Partial<JoinToCreateConfig>>({
    guildId,
    channelNamePrefix: "Room",
    userLimit: 0,
    privateChannel: false,
    enabled: true,
  });

  const fetchConfigs = async () => {
    try {
      setLoading(true);
      const res = await fetch("/api/discord/join-to-create", {
        credentials: "include",
      });
      if (res.ok) {
        const data = await res.json();
        // Filter configs for this guild
        setConfigs(data.filter((c: JoinToCreateConfig) => c.guildId === guildId));
      }
    } catch (err) {
      console.error("Failed to fetch configs:", err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchConfigs();
  }, [guildId]);

  const handleCreate = async () => {
    if (
      !newConfig.joinChannelId ||
      !newConfig.categoryId ||
      !newConfig.channelNamePrefix
    ) {
      alert("Please fill in all required fields");
      return;
    }

    try {
      const res = await fetch("/api/discord/join-to-create", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify(newConfig),
      });

      if (res.ok) {
        setNewConfig({
          guildId,
          channelNamePrefix: "Room",
          userLimit: 0,
          privateChannel: false,
          enabled: true,
        });
        setShowNewForm(false);
        fetchConfigs();
      } else {
        const error = await res.text();
        alert(`Failed to create config: ${error}`);
      }
    } catch (err) {
      console.error("Failed to create config:", err);
      alert("Failed to create config");
    }
  };

  const handleUpdate = async (config: JoinToCreateConfig) => {
    try {
      const res = await fetch(`/api/discord/join-to-create/${config.id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({
          joinChannelId: config.joinChannelId,
          categoryId: config.categoryId,
          channelNamePrefix: config.channelNamePrefix,
          userLimit: config.userLimit,
          privateChannel: config.privateChannel,
          enabled: config.enabled,
        }),
      });

      if (res.ok) {
        alert("Config updated successfully!");
        fetchConfigs();
      }
    } catch (err) {
      console.error("Failed to update config:", err);
      alert("Failed to update config");
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm("Delete this Join-to-Create config?")) return;

    try {
      const res = await fetch(`/api/discord/join-to-create/${id}`, {
        method: "DELETE",
        credentials: "include",
      });

      if (res.ok) {
        fetchConfigs();
      }
    } catch (err) {
      console.error("Failed to delete config:", err);
    }
  };

  const toggleEnabled = (config: JoinToCreateConfig) => {
    const updated = { ...config, enabled: !config.enabled };
    handleUpdate(updated);
  };

  if (loading) {
    return <p className="text-gray-500">Loading Join-to-Create configs...</p>;
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="text-lg font-semibold text-gray-900">
          Join-to-Create Channels
        </h3>
        <button
          onClick={() => setShowNewForm(!showNewForm)}
          className="flex items-center gap-2 px-3 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg transition-colors text-sm"
        >
          <Plus className="w-4 h-4" />
          Add Config
        </button>
      </div>

      {/* New Config Form */}
      {showNewForm && (
        <div className="p-4 bg-gray-50 rounded-lg border border-gray-200 space-y-3">
          <h4 className="font-medium text-gray-900">New Join-to-Create Config</h4>

          <div className="grid gap-3 sm:grid-cols-2">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Join Channel ID *
              </label>
              <input
                type="text"
                placeholder="123456789"
                value={newConfig.joinChannelId ?? ""}
                onChange={(e) =>
                  setNewConfig({
                    ...newConfig,
                    joinChannelId: e.target.value || undefined,
                  })
                }
                className="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Category ID *
              </label>
              <input
                type="text"
                placeholder="987654321"
                value={newConfig.categoryId ?? ""}
                onChange={(e) =>
                  setNewConfig({
                    ...newConfig,
                    categoryId: e.target.value || undefined,
                  })
                }
                className="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Channel Name Prefix *
              </label>
              <input
                type="text"
                placeholder="Room"
                value={newConfig.channelNamePrefix ?? ""}
                onChange={(e) =>
                  setNewConfig({ ...newConfig, channelNamePrefix: e.target.value })
                }
                className="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                User Limit (0 = unlimited)
              </label>
              <input
                type="number"
                min="0"
                max="99"
                value={newConfig.userLimit ?? 0}
                onChange={(e) =>
                  setNewConfig({ ...newConfig, userLimit: Number(e.target.value) })
                }
                className="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>
          </div>

          <label className="flex items-center gap-2">
            <input
              type="checkbox"
              checked={newConfig.privateChannel ?? false}
              onChange={(e) =>
                setNewConfig({ ...newConfig, privateChannel: e.target.checked })
              }
              className="w-4 h-4 text-indigo-600 rounded focus:ring-2 focus:ring-indigo-500"
            />
            <span className="text-sm text-gray-700">Private channel (only creator can see)</span>
          </label>

          <div className="flex gap-2">
            <button
              onClick={handleCreate}
              className="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg transition-colors text-sm"
            >
              Create Config
            </button>
            <button
              onClick={() => setShowNewForm(false)}
              className="px-4 py-2 bg-gray-200 hover:bg-gray-300 text-gray-700 rounded-lg transition-colors text-sm"
            >
              Cancel
            </button>
          </div>
        </div>
      )}

      {/* Config List */}
      {configs.length > 0 ? (
        <div className="space-y-3">
          {configs.map((config) => (
            <div
              key={config.id}
              className="p-4 bg-white rounded-lg border border-gray-200 space-y-3"
            >
              <div className="flex items-start justify-between">
                <div className="flex-1 grid gap-3 sm:grid-cols-2">
                  <div>
                    <label className="block text-xs font-medium text-gray-500 mb-1">
                      Join Channel ID
                    </label>
                    <input
                      type="text"
                      value={config.joinChannelId}
                      onChange={(e) => {
                        const updated = { ...config, joinChannelId: e.target.value };
                        setConfigs(configs.map((c) => (c.id === config.id ? updated : c)));
                      }}
                      className="w-full p-2 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500"
                    />
                  </div>

                  <div>
                    <label className="block text-xs font-medium text-gray-500 mb-1">
                      Category ID
                    </label>
                    <input
                      type="text"
                      value={config.categoryId}
                      onChange={(e) => {
                        const updated = { ...config, categoryId: e.target.value };
                        setConfigs(configs.map((c) => (c.id === config.id ? updated : c)));
                      }}
                      className="w-full p-2 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500"
                    />
                  </div>

                  <div>
                    <label className="block text-xs font-medium text-gray-500 mb-1">
                      Channel Prefix
                    </label>
                    <input
                      type="text"
                      value={config.channelNamePrefix}
                      onChange={(e) => {
                        const updated = { ...config, channelNamePrefix: e.target.value };
                        setConfigs(configs.map((c) => (c.id === config.id ? updated : c)));
                      }}
                      className="w-full p-2 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500"
                    />
                  </div>

                  <div>
                    <label className="block text-xs font-medium text-gray-500 mb-1">
                      User Limit
                    </label>
                    <input
                      type="number"
                      value={config.userLimit}
                      onChange={(e) => {
                        const updated = { ...config, userLimit: Number(e.target.value) };
                        setConfigs(configs.map((c) => (c.id === config.id ? updated : c)));
                      }}
                      className="w-full p-2 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500"
                    />
                  </div>
                </div>
              </div>

              <div className="flex items-center justify-between">
                <div className="flex gap-4">
                  <label className="flex items-center gap-2">
                    <input
                      type="checkbox"
                      checked={config.privateChannel}
                      onChange={(e) => {
                        const updated = { ...config, privateChannel: e.target.checked };
                        setConfigs(configs.map((c) => (c.id === config.id ? updated : c)));
                      }}
                      className="w-4 h-4 text-indigo-600 rounded"
                    />
                    <span className="text-sm text-gray-700">Private</span>
                  </label>

                  <label className="flex items-center gap-2">
                    <input
                      type="checkbox"
                      checked={config.enabled}
                      onChange={() => toggleEnabled(config)}
                      className="w-4 h-4 text-green-600 rounded"
                    />
                    <span className="text-sm text-gray-700">Enabled</span>
                  </label>
                </div>

                <div className="flex gap-2">
                  <button
                    onClick={() => handleUpdate(config)}
                    className="flex items-center gap-1 px-3 py-1 bg-blue-100 hover:bg-blue-200 text-blue-700 rounded-lg transition-colors text-sm"
                  >
                    <Save className="w-4 h-4" />
                    Save
                  </button>
                  <button
                    onClick={() => handleDelete(config.id)}
                    className="flex items-center gap-1 px-3 py-1 bg-red-100 hover:bg-red-200 text-red-700 rounded-lg transition-colors text-sm"
                  >
                    <Trash2 className="w-4 h-4" />
                    Delete
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <div className="flex items-center justify-center p-8 bg-yellow-50 border border-yellow-200 rounded-lg">
          <AlertTriangle className="w-5 h-5 mr-2 text-yellow-600" />
          <span className="text-yellow-700">No Join-to-Create configs yet</span>
        </div>
      )}
    </div>
  );
};

export default JoinToCreateManager;