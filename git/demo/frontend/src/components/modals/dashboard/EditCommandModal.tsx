import { useState } from "react";
import { updateCommand } from "../../../features/commands/api";
import type { ChatCommand } from "../../../features/commands/types";

interface EditCommandModalProps {
  command: ChatCommand;
  onClose: () => void;
  onUpdated: (cmd: ChatCommand) => void;
}

export default function EditCommandModal({
  command,
  onClose,
  onUpdated
}: EditCommandModalProps) {
  const [trigger, setTrigger] = useState(command.trigger);
  const [response, setResponse] = useState(command.response);
  const [cooldown, setCooldown] = useState(command.cooldown);
  const [enabled, setEnabled] = useState(command.enabled);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      await updateCommand(command.id, {
        trigger,
        response,
        cooldown,
        enabled
      });

      const updatedCommand: ChatCommand = {
        ...command,
        trigger,
        response,
        cooldown,
        enabled
      };

      onUpdated(updatedCommand);
      onClose();
    } catch (err: any) {
      setError(err.message || "Fehler beim Aktualisieren");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="fixed inset-0 bg-black/50 flex justify-center items-center z-50 p-4">
      <div className="bg-gray-800 p-6 rounded w-full max-w-md">
        <h2 className="text-xl font-bold mb-4">Command bearbeiten</h2>

        {error && <p className="text-red-500 mb-2 text-sm">{error}</p>}

        <form onSubmit={handleSubmit} className="space-y-3">
          <div>
            <label className="block text-sm mb-1">Trigger (ohne !)</label>
            <input
              type="text"
              value={trigger}
              onChange={(e) => setTrigger(e.target.value)}
              required
              pattern="[a-zA-Z0-9_]+"
              title="Nur Buchstaben, Zahlen und Unterstriche"
              className="w-full p-2 rounded bg-gray-700 text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <p className="text-xs text-gray-400 mt-1">
              ✨ Trigger kann jetzt geändert werden!
            </p>
          </div>

          <div>
            <label className="block text-sm mb-1">Antwort</label>
            <textarea
              value={response}
              onChange={(e) => setResponse(e.target.value)}
              required
              rows={3}
              className="w-full p-2 rounded bg-gray-700 text-white focus:outline-none focus:ring-2 focus:ring-blue-500 resize-none"
            />
            <p className="text-xs text-gray-400 mt-1">
              {response.length}/500 Zeichen
            </p>
          </div>

          <div>
            <label className="block text-sm mb-1">Cooldown (Sekunden)</label>
            <input
              type="number"
              value={cooldown}
              onChange={(e) => setCooldown(Number(e.target.value))}
              min="0"
              max="300"
              className="w-full p-2 rounded bg-gray-700 text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div className="flex items-center gap-2">
            <input
              type="checkbox"
              id="enabled"
              checked={enabled}
              onChange={(e) => setEnabled(e.target.checked)}
              className="w-4 h-4 rounded bg-gray-700 border-gray-600 text-green-600 focus:ring-2 focus:ring-green-500"
            />
            <label htmlFor="enabled" className="text-sm cursor-pointer">
              Command aktiviert
            </label>
          </div>

          <div className="flex justify-end gap-2 mt-4 pt-4 border-t border-gray-700">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 bg-gray-600 rounded hover:bg-gray-500 transition-colors"
            >
              Abbrechen
            </button>
            <button
              type="submit"
              disabled={loading}
              className="px-4 py-2 bg-blue-600 rounded hover:bg-blue-500 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {loading ? "Speichern…" : "Speichern"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}