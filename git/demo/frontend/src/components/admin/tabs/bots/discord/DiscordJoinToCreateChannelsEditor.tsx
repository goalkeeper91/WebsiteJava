import React, { useEffect, useState } from "react";
import { AlertTriangle, LogOut } from "lucide-react";

type JoinToCreateChannel = {
  id: number;
  joinChannelId: string;
  categoryId: string;
  channelNamePrefix: string;
  userLimit?: number;
  privateChannel: boolean;
};

const DiscordJoinToCreateChannelsEditor: React.FC = () => {
  const [channels, setChannels] = useState<JoinToCreateChannel[]>([]);
  const [loading, setLoading] = useState(true);
  const [newChannel, setNewChannel] = useState<Partial<JoinToCreateChannel>>({});

  const fetchChannels = async () => {
    setLoading(true);
    try {
      const res = await fetch("/api/discord/join-to-create");
      const data = await res.json();
      setChannels(Array.isArray(data) ? data : []);
    } catch (err) {
      console.error("Failed to fetch channels:", err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchChannels();
  }, []);

  const handleAdd = async () => {
    if (!newChannel.joinChannelId || !newChannel.categoryId || !newChannel.channelNamePrefix) return;
    try {
      const res = await fetch("/api/discord/join-to-create", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(newChannel),
      });
      if (res.ok) {
        setNewChannel({});
        fetchChannels();
      }
    } catch (err) {
      console.error("Add failed:", err);
    }
  };

  const handleDelete = async (id: number) => {
    try {
      const res = await fetch(`/api/discord/join-to-create/${id}`, { method: "DELETE" });
      if (res.ok) fetchChannels();
    } catch (err) {
      console.error("Delete failed:", err);
    }
  };

  const handleUpdate = async (id: number, field: keyof JoinToCreateChannel, value: any) => {
    const updated = channels.find((c) => c.id === id);
    if (!updated) return;

    updated[field] = value;

    try {
      const res = await fetch(`/api/discord/join-to-create/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(updated),
      });
      if (res.ok) fetchChannels();
    } catch (err) {
      console.error("Update failed:", err);
    }
  };

  if (loading) return <p>Lade JoinToCreate Channels...</p>;

  return (
    <div className="p-4 rounded-xl shadow space-y-4">
      <h3 className="text-lg font-bold">Join To Create Channels</h3>

      {/* Neue Channel hinzufügen */}
      <div className="flex flex-wrap gap-2 mb-4">
        <input
          type="text"
          placeholder="Join Channel ID"
          value={newChannel.joinChannelId || ""}
          onChange={(e) => setNewChannel({ ...newChannel, joinChannelId: e.target.value })}
          className="border p-2 rounded flex-1 min-w-[120px]"
        />
        <input
          type="text"
          placeholder="Category ID"
          value={newChannel.categoryId || ""}
          onChange={(e) => setNewChannel({ ...newChannel, categoryId: e.target.value })}
          className="border p-2 rounded flex-1 min-w-[120px]"
        />
        <input
          type="text"
          placeholder="Channel Name Prefix"
          value={newChannel.channelNamePrefix || ""}
          onChange={(e) => setNewChannel({ ...newChannel, channelNamePrefix: e.target.value })}
          className="border p-2 rounded flex-1 min-w-[150px]"
        />
        <input
          type="number"
          placeholder="User Limit"
          value={newChannel.userLimit ?? ""}
          onChange={(e) => setNewChannel({ ...newChannel, userLimit: Number(e.target.value) })}
          className="border p-2 rounded flex-1 min-w-[80px]"
        />
        <label className="flex items-center gap-1">
          <input
            type="checkbox"
            checked={newChannel.privateChannel ?? false}
            onChange={(e) => setNewChannel({ ...newChannel, privateChannel: e.target.checked })}
          />
          Private
        </label>
        <button
          onClick={handleAdd}
          className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded"
        >
          Add
        </button>
      </div>

      {/* Channels Liste */}
      {channels.length > 0 ? (
        <div className="space-y-2">
          {channels.map((ch) => (
            <div
              key={ch.id}
              className="flex flex-wrap items-center justify-between p-3 rounded shadow-sm gap-2"
            >
              <input
                type="text"
                value={ch.joinChannelId}
                onChange={(e) => handleUpdate(ch.id, "joinChannelId", e.target.value)}
                className="border p-2 rounded flex-1 min-w-[100px]"
              />
              <input
                type="text"
                value={ch.categoryId}
                onChange={(e) => handleUpdate(ch.id, "categoryId", e.target.value)}
                className="border p-2 rounded flex-1 min-w-[100px]"
              />
              <input
                type="text"
                value={ch.channelNamePrefix}
                onChange={(e) => handleUpdate(ch.id, "channelNamePrefix", e.target.value)}
                className="border p-2 rounded flex-1 min-w-[150px]"
              />
              <input
                type="number"
                value={ch.userLimit ?? ""}
                onChange={(e) => handleUpdate(ch.id, "userLimit", Number(e.target.value))}
                className="border p-2 rounded flex-1 min-w-[80px]"
              />
              <label className="flex items-center gap-1">
                <input
                  type="checkbox"
                  checked={ch.privateChannel}
                  onChange={(e) => handleUpdate(ch.id, "privateChannel", e.target.checked)}
                />
                Private
              </label>
              <button
                onClick={() => handleDelete(ch.id)}
                className="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded flex items-center"
              >
                <LogOut className="w-4 h-4 mr-1" /> Delete
              </button>
            </div>
          ))}
        </div>
      ) : (
        <div className="flex items-center justify-center p-4 bg-yellow-50 text-yellow-700 rounded">
          <AlertTriangle className="w-5 h-5 mr-2" /> Keine Channels vorhanden
        </div>
      )}
    </div>
  );
};

export default DiscordJoinToCreateChannelsEditor;
