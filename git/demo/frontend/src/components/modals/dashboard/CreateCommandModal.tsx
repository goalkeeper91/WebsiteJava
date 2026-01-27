import { useState } from "react";
import { createCommand } from "../../../features/commands/api";
import type { ChatCommand } from "../../../features/commands/types";

interface CreateCommandModalProps {
  onClose: () => void;
  onCreated: (cmd: ChatCommand) => void;
}

export default function CreateCommandModal({ onClose, onCreated }: CreateCommandModalProps) {
  const [trigger, setTrigger] = useState("");
  const [response, setResponse] = useState("");
  const [cooldown, setCooldown] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      const newCommand = await createCommand({ trigger, response, cooldown });
      onCreated(newCommand);
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
        <h2 className="text-xl font-bold mb-4">Neues Command erstellen</h2>

        {error && <p className="text-red-500 mb-2">{error}</p>}

        <form onSubmit={handleSubmit} className="space-y-3">
          <div>
            <label className="block text-sm mb-1">Trigger (ohne !)</label>
            <input
              type="text"
              value={trigger}
              onChange={(e) => setTrigger(e.target.value)}
              required
              className="w-full p-2 rounded bg-gray-700 text-white"
            />
          </div>

          <div>
            <label className="block text-sm mb-1">Antwort</label>
            <input
              type="text"
              value={response}
              onChange={(e) => setResponse(e.target.value)}
              required
              className="w-full p-2 rounded bg-gray-700 text-white"
            />
          </div>

          <div>
            <label className="block text-sm mb-1">Cooldown (Sekunden)</label>
            <input
              type="number"
              value={cooldown}
              onChange={(e) => setCooldown(Number(e.target.value))}
              className="w-full p-2 rounded bg-gray-700 text-white"
            />
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
              disabled={loading}
              className="px-4 py-2 bg-green-600 rounded hover:bg-green-500"
            >
              {loading ? "Speichern…" : "Erstellen"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
