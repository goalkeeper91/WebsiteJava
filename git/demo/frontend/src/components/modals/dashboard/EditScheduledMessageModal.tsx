import { useState } from "react";
import { updateScheduledMessage } from "../../../features/scheduledMessages/api";
import type { ScheduledMessage } from "../../../features/scheduledMessages/types";

interface EditScheduledMessageModalProps {
  scheduledMessage: ScheduledMessage;
  onClose: () => void;
  onUpdated: (message: ScheduledMessage) => void;
}

export default function EditScheduledMessageModal({
  scheduledMessage,
  onClose,
  onUpdated,
}: EditScheduledMessageModalProps) {
  const isCommandLinked = !!scheduledMessage.command_id;
  const [message, setMessage] = useState(scheduledMessage.message ?? "");
  const [intervalMinutes, setIntervalMinutes] = useState(
    Math.max(1, Math.round(scheduledMessage.interval_seconds / 60))
  );
  const [onlyWhenLive, setOnlyWhenLive] = useState(scheduledMessage.only_when_live);
  const [enabled, setEnabled] = useState(scheduledMessage.enabled);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      const updated = await updateScheduledMessage(scheduledMessage.id, {
        ...(isCommandLinked ? {} : { message }),
        interval_seconds: intervalMinutes * 60,
        enabled,
        only_when_live: onlyWhenLive,
      });
      onUpdated(updated);
      onClose();
    } catch (err: any) {
      setError(err.message || "Fehler beim Aktualisieren");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="fixed inset-0 bg-black/50 flex justify-center items-center z-50">
      <div className="bg-gray-800 p-6 rounded w-full max-w-md">
        <h2 className="text-xl font-bold mb-4">Automatisierte Nachricht bearbeiten</h2>

        {error && <p className="text-red-500 mb-2 text-sm">{error}</p>}

        <form onSubmit={handleSubmit} className="space-y-3">
          {isCommandLinked ? (
            <div className="bg-gray-900 border border-gray-700 rounded p-3 text-sm text-gray-400">
              🔗 Diese Nachricht ist mit einem Command verknüpft — der Text folgt automatisch
              dessen aktueller Antwort. Um den Text zu ändern, bearbeite den Command selbst.
            </div>
          ) : (
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
          </div>

          <div className="flex items-center gap-2">
            <input
              type="checkbox"
              id="only-when-live"
              checked={onlyWhenLive}
              onChange={(e) => setOnlyWhenLive(e.target.checked)}
              className="w-4 h-4 rounded bg-gray-700 border-gray-600 text-purple-600 focus:ring-2 focus:ring-purple-500"
            />
            <label htmlFor="only-when-live" className="text-sm cursor-pointer">
              Nur posten, wenn der Kanal live ist
            </label>
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
              Aktiviert
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
              className="px-4 py-2 bg-blue-600 rounded hover:bg-blue-500 transition-colors disabled:opacity-50"
            >
              {loading ? "Speichern…" : "Speichern"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
