import { useEffect, useState } from "react";
import { createScheduledMessage } from "../../../features/scheduledMessages/api";
import type { ScheduledMessage } from "../../../features/scheduledMessages/types";
import { getCommands } from "../../../features/commands/api";
import type { ChatCommand } from "../../../features/commands/types";

interface CreateScheduledMessageModalProps {
  onClose: () => void;
  onCreated: (message: ScheduledMessage) => void;
}

type Mode = "message" | "command";

export default function CreateScheduledMessageModal({
  onClose,
  onCreated,
}: CreateScheduledMessageModalProps) {
  const [mode, setMode] = useState<Mode>("message");
  const [message, setMessage] = useState("");
  const [commands, setCommands] = useState<ChatCommand[]>([]);
  const [commandId, setCommandId] = useState<number | null>(null);
  const [intervalMinutes, setIntervalMinutes] = useState(20);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    if (mode !== "command" || commands.length > 0) return;
    getCommands(1, 100)
      .then((data) => {
        const simpleCommands = (data.data || []).filter(
          (cmd) => cmd.command_type === "simple"
        );
        setCommands(simpleCommands);
        if (simpleCommands.length > 0) setCommandId(simpleCommands[0].id);
      })
      .catch((err) => console.error("Fehler beim Laden der Commands:", err));
  }, [mode, commands.length]);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      const created = await createScheduledMessage(
        mode === "command"
          ? { command_id: commandId ?? undefined, interval_seconds: intervalMinutes * 60 }
          : { message, interval_seconds: intervalMinutes * 60 }
      );
      onCreated(created);
      onClose();
    } catch (err: any) {
      setError(err.message || "Fehler beim Erstellen");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="fixed inset-0 bg-black/50 flex justify-center items-center z-50">
      <div className="bg-gray-800 p-6 rounded w-full max-w-md">
        <h2 className="text-xl font-bold mb-4">Neue automatisierte Nachricht</h2>

        {error && <p className="text-red-500 mb-2 text-sm">{error}</p>}

        <div className="flex gap-2 mb-4 bg-gray-900 p-1 rounded-lg">
          <button
            type="button"
            onClick={() => setMode("message")}
            className={`flex-1 px-3 py-1.5 rounded-md text-sm transition-colors ${
              mode === "message" ? "bg-purple-600 text-white" : "text-gray-400 hover:text-gray-200"
            }`}
          >
            Eigene Nachricht
          </button>
          <button
            type="button"
            onClick={() => setMode("command")}
            className={`flex-1 px-3 py-1.5 rounded-md text-sm transition-colors ${
              mode === "command" ? "bg-purple-600 text-white" : "text-gray-400 hover:text-gray-200"
            }`}
          >
            Bestehendes Command
          </button>
        </div>

        <form onSubmit={handleSubmit} className="space-y-3">
          {mode === "message" ? (
            <div>
              <label className="block text-sm mb-1">Nachricht</label>
              <textarea
                value={message}
                onChange={(e) => setMessage(e.target.value)}
                required
                rows={3}
                className="w-full p-2 rounded bg-gray-700 text-white focus:outline-none focus:ring-2 focus:ring-purple-500 resize-none"
              />
            </div>
          ) : (
            <div>
              <label className="block text-sm mb-1">Command</label>
              {commands.length === 0 ? (
                <p className="text-sm text-gray-400">
                  Keine einfachen Commands vorhanden. Lege zuerst ein Chat-Command an.
                </p>
              ) : (
                <select
                  value={commandId ?? ""}
                  onChange={(e) => setCommandId(Number(e.target.value))}
                  required
                  className="w-full p-2 rounded bg-gray-700 text-white focus:outline-none focus:ring-2 focus:ring-purple-500"
                >
                  {commands.map((cmd) => (
                    <option key={cmd.id} value={cmd.id}>
                      !{cmd.trigger} — {cmd.response}
                    </option>
                  ))}
                </select>
              )}
              <p className="text-xs text-gray-400 mt-1">
                Der Command bleibt normal per !trigger nutzbar und postet zusätzlich automatisch.
              </p>
            </div>
          )}

          <div>
            <label className="block text-sm mb-1">Intervall (Minuten)</label>
            <input
              type="number"
              value={intervalMinutes}
              onChange={(e) => setIntervalMinutes(Number(e.target.value))}
              min={1}
              required
              className="w-full p-2 rounded bg-gray-700 text-white focus:outline-none focus:ring-2 focus:ring-purple-500"
            />
            <p className="text-xs text-gray-400 mt-1">Mindestens 1 Minute.</p>
          </div>

          <div className="flex justify-end gap-2 mt-4">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 bg-gray-600 rounded hover:bg-gray-500"
            >
              Abbrechen
            </button>
            <button
              type="submit"
              disabled={loading || (mode === "command" && commands.length === 0)}
              className="px-4 py-2 bg-purple-600 rounded hover:bg-purple-500 disabled:opacity-50"
            >
              {loading ? "Speichern…" : "Erstellen"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
