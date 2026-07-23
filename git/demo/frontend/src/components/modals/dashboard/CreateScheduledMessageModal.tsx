import { useState } from "react";
import { createScheduledMessage } from "../../../features/scheduledMessages/api";
import type { ScheduledMessage } from "../../../features/scheduledMessages/types";

interface CreateScheduledMessageModalProps {
  onClose: () => void;
  onCreated: (message: ScheduledMessage) => void;
}

export default function CreateScheduledMessageModal({
  onClose,
  onCreated,
}: CreateScheduledMessageModalProps) {
  const [message, setMessage] = useState("");
  const [intervalMinutes, setIntervalMinutes] = useState(20);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      const created = await createScheduledMessage({
        message,
        interval_seconds: intervalMinutes * 60,
      });
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

        <form onSubmit={handleSubmit} className="space-y-3">
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
              disabled={loading}
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
