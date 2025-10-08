import React, { useEffect, useState } from "react";
import { AlertTriangle, Trash2, Info } from "lucide-react";

interface TwitchCommand {
  id: number;
  trigger: string;
  response: string;
  modOnly: boolean;
  duration?: number; // <--- neu für Timer
  createdAt: string;
  updatedAt: string;
}

const TwitchCommandsEditor: React.FC = () => {
  const [commands, setCommands] = useState<TwitchCommand[]>([]);
  const [loading, setLoading] = useState(true);
  const [newCommand, setNewCommand] = useState<Partial<TwitchCommand>>({});

  const fetchCommands = async () => {
    setLoading(true);
    try {
      const res = await fetch("/api/bot/twitch/commands");
      const data = await res.json();
      setCommands(Array.isArray(data) ? data : []);
    } catch (err) {
      console.error("Failed to fetch commands:", err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchCommands();
  }, []);

  const handleAdd = async () => {
    if (!newCommand.trigger || !newCommand.response) return;
    try {
      const res = await fetch("/api/bot/twitch/commands", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(newCommand),
      });
      if (res.ok) {
        setNewCommand({});
        fetchCommands();
      }
    } catch (err) {
      console.error("Add failed:", err);
    }
  };

  const handleDelete = async (id: number) => {
    try {
      const res = await fetch(`/api/bot/twitch/commands/${id}`, { method: "DELETE" });
      if (res.ok) fetchCommands();
    } catch (err) {
      console.error("Delete failed:", err);
    }
  };

  const handleUpdate = async (
    id: number,
    field: "trigger" | "response" | "modOnly" | "duration",
    value: string | number | boolean | undefined  // <-- hier erweitern
  ) => {
    const updatedCommand = commands.find(c => c.id === id);
    if (!updatedCommand) return;

    const payload: Partial<TwitchCommand> = {
      ...updatedCommand,
      [field]: value,
    };

    try {
      const res = await fetch(`/api/bot/twitch/commands/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
      if (res.ok) fetchCommands();
    } catch (err) {
      console.error("Update failed:", err);
    }
  };

  if (loading) return <p>Lade Commands...</p>;

  return (
    <div className="p-4 rounded-xl shadow space-y-4">
      <h3 className="text-lg font-bold">Twitch Commands</h3>

      {/* Neue Command hinzufügen */}
      <div className="flex flex-wrap gap-2 mb-4">
        <input
          type="text"
          placeholder="Trigger"
          value={newCommand.trigger || ""}
          onChange={(e) => setNewCommand({ ...newCommand, trigger: e.target.value })}
          className="border p-2 rounded flex-1 min-w-[120px]"
        />
        <input
          type="text"
          placeholder="Response"
          value={newCommand.response || ""}
          onChange={(e) => setNewCommand({ ...newCommand, response: e.target.value })}
          className="border p-2 rounded flex-1 min-w-[200px]"
        />
        <input
          type="number"
          placeholder="Timer (Sekunden)"
          value={newCommand.duration || ""}
          onChange={(e) =>
            setNewCommand({ ...newCommand, duration: e.target.value ? Number(e.target.value) : undefined })
          }
          className="border p-2 rounded w-[120px]"
        />
        <label className="flex items-center gap-1">
          <input
            type="checkbox"
            checked={newCommand.modOnly || false}
            onChange={(e) => setNewCommand({ ...newCommand, modOnly: e.target.checked })}
          />{" "}
          Mod Only
        </label>
        <button
          onClick={handleAdd}
          className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded"
        >
          Add
        </button>
      </div>

      {/* Commands Liste */}
      {commands.length > 0 ? (
        <div className="space-y-2">
          {commands.map((cmd) => (
            <div
              key={cmd.id}
              className="flex flex-wrap items-center justify-between p-3 rounded shadow-sm gap-2"
            >
              <input
                type="text"
                value={cmd.trigger}
                onChange={(e) => handleUpdate(cmd.id, "trigger", e.target.value)}
                className="border p-2 rounded flex-1 min-w-[120px]"
              />
              <input
                type="text"
                value={cmd.response}
                onChange={(e) => handleUpdate(cmd.id, "response", e.target.value)}
                className="border p-2 rounded flex-1 min-w-[200px]"
              />
              <input
                type="number"
                value={cmd.duration || ""}
                onChange={(e) =>
                  handleUpdate(cmd.id, "duration", e.target.value ? Number(e.target.value) : undefined)
                }
                className="border p-2 rounded w-[120px]"
              />
              <label className="flex items-center gap-1">
                <input
                  type="checkbox"
                  checked={cmd.modOnly}
                  onChange={(e) => handleUpdate(cmd.id, "modOnly", e.target.checked)}
                />{" "}
                Mod Only
              </label>
              <div className="flex gap-2">
                <button
                  onClick={() => handleDelete(cmd.id)}
                  className="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded flex items-center"
                >
                  <Trash2 className="w-4 h-4 mr-1" /> Delete
                </button>
                <button
                  onClick={() => alert("Details kommen später als Popup")}
                  className="bg-gray-500 hover:bg-gray-600 text-white px-3 py-1 rounded flex items-center"
                >
                  <Info className="w-4 h-4 mr-1" /> Details
                </button>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <div className="flex items-center justify-center p-4 bg-yellow-50 text-yellow-700 rounded">
          <AlertTriangle className="w-5 h-5 mr-2" /> Keine Commands vorhanden
        </div>
      )}
    </div>
  );
};

export default TwitchCommandsEditor;
