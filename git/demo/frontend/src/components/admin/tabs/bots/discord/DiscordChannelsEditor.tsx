import React, { useEffect, useState } from "react";
import { AlertTriangle, LogOut } from "lucide-react";

type DiscordChannel = {
  id: number;
  guildId: string;
  channelId: string;
  description?: string;
};

const DiscordChannelEditor: React.FC = () => {
  const [channels, setChannels] = useState<DiscordChannel[]>([]);
  const [loading, setLoading] = useState(true);
  const [newChannel, setNewChannel] = useState<Partial<DiscordChannel>>({});

  const fetchChannels = async () => {
    setLoading(true);
    try {
      const res = await fetch("/api/discord/channels");
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
    if (!newChannel.guildId || !newChannel.channelId) return;
    try {
      const res = await fetch("/api/discord/channels", {
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
      const res = await fetch(`/api/discord/channels/${id}`, { method: "DELETE" });
      if (res.ok) fetchChannels();
    } catch (err) {
      console.error("Delete failed:", err);
    }
  };

  type EditableChannelField = "guildId" | "channelId" | "description";

  const handleUpdate = async (id: number, field: EditableChannelField, value: string) => {
    const updated = channels.find((c) => c.id === id);
    if (!updated) return;

    updated[field] = value; // TypeScript weiß jetzt, dass field ein String-Feld ist

    try {
      const res = await fetch(`/api/discord/channels/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(updated),
      });
      if (res.ok) fetchChannels();
    } catch (err) {
      console.error("Update failed:", err);
    }
  };

  if (loading) return <p>Lade Channels...</p>;

  return (
    <div className="p-4 rounded-xl shadow space-y-4">
      <h3 className="text-lg font-bold">Discord Channels</h3>

      {/* Neue Channel hinzufügen */}
      <div className="flex flex-wrap gap-2 mb-4">
        <input
          type="text"
          placeholder="Guild ID"
          value={newChannel.guildId || ""}
          onChange={(e) => setNewChannel({ ...newChannel, guildId: e.target.value })}
          className="border p-2 rounded flex-1 min-w-[120px]"
        />
        <input
          type="text"
          placeholder="Channel ID"
          value={newChannel.channelId || ""}
          onChange={(e) => setNewChannel({ ...newChannel, channelId: e.target.value })}
          className="border p-2 rounded flex-1 min-w-[120px]"
        />
        <input
          type="text"
          placeholder="Description"
          value={newChannel.description || ""}
          onChange={(e) => setNewChannel({ ...newChannel, description: e.target.value })}
          className="border p-2 rounded flex-1 min-w-[150px]"
        />
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
                value={ch.guildId}
                onChange={(e) => handleUpdate(ch.id, "guildId", e.target.value)}
                className="border p-2 rounded flex-1 min-w-[100px]"
              />
              <input
                type="text"
                value={ch.channelId}
                onChange={(e) => handleUpdate(ch.id, "channelId", e.target.value)}
                className="border p-2 rounded flex-1 min-w-[100px]"
              />
              <input
                type="text"
                value={ch.description || ""}
                onChange={(e) => handleUpdate(ch.id, "description", e.target.value)}
                className="border p-2 rounded flex-1 min-w-[150px]"
              />
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

export default DiscordChannelEditor;
