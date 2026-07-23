import { useEffect, useState } from "react";
import {
  getScheduledMessages,
  deleteScheduledMessage,
  updateScheduledMessage,
} from "../../features/scheduledMessages/api";
import type { ScheduledMessage } from "../../features/scheduledMessages/types";
import CreateScheduledMessageModal from "../modals/dashboard/CreateScheduledMessageModal";
import EditScheduledMessageModal from "../modals/dashboard/EditScheduledMessageModal";

function formatRelative(dateString?: string): string | null {
  if (!dateString) return null;
  const seconds = Math.max(0, Math.floor((Date.now() - new Date(dateString).getTime()) / 1000));
  if (seconds < 60) return "gerade eben";
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return `vor ${minutes} Min.`;
  const hours = Math.floor(minutes / 60);
  return `vor ${hours} Std.`;
}

export default function ScheduledMessagesDashboard() {
  const [messages, setMessages] = useState<ScheduledMessage[]>([]);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(0);
  const [totalElements, setTotalElements] = useState(0);
  const [loading, setLoading] = useState(true);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [editingMessage, setEditingMessage] = useState<ScheduledMessage | null>(null);

  useEffect(() => {
    loadMessages();
  }, [page]);

  async function loadMessages() {
    setLoading(true);
    try {
      const data = await getScheduledMessages(page, 10);
      setMessages(data.data || []);
      setTotalPages(data.total_pages || 0);
      setTotalElements(data.total || 0);
    } catch (err) {
      console.error("Fehler beim Laden der automatisierten Nachrichten:", err);
    } finally {
      setLoading(false);
    }
  }

  async function handleToggleEnabled(m: ScheduledMessage) {
    setMessages((prev) => prev.map((x) => (x.id === m.id ? { ...x, enabled: !x.enabled } : x)));

    try {
      await updateScheduledMessage(m.id, { enabled: !m.enabled });
    } catch (err) {
      console.error("Fehler beim Update:", err);
      setMessages((prev) => prev.map((x) => (x.id === m.id ? { ...x, enabled: m.enabled } : x)));
      alert("Fehler beim Aktualisieren des Status");
    }
  }

  async function handleDelete(m: ScheduledMessage) {
    if (!confirm("Automatisierte Nachricht wirklich löschen?")) return;

    setMessages((prev) => prev.filter((x) => x.id !== m.id));
    setTotalElements((prev) => prev - 1);

    try {
      await deleteScheduledMessage(m.id);
    } catch (err) {
      console.error("Fehler beim Löschen:", err);
      loadMessages();
      alert("Fehler beim Löschen");
    }
  }

  function handleCreated(newMessage: ScheduledMessage) {
    setMessages((prev) => [newMessage, ...prev]);
    setTotalElements((prev) => prev + 1);
    setShowCreateModal(false);
  }

  function handleUpdated(updated: ScheduledMessage) {
    setMessages((prev) => prev.map((x) => (x.id === updated.id ? updated : x)));
    setEditingMessage(null);
  }

  return (
    <div className="max-w-5xl mx-auto">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
        <div>
          <h1 className="text-2xl font-bold">Automatisierte Nachrichten</h1>
          <p className="text-sm text-gray-400 mt-1">
            {loading ? "Lade..." : `${totalElements} Nachrichten insgesamt`}
          </p>
        </div>
        <button
          onClick={() => setShowCreateModal(true)}
          className="px-4 py-2 bg-green-600 rounded hover:bg-green-500 transition-colors flex items-center gap-2"
        >
          <span className="text-xl">+</span>
          Neue Nachricht
        </button>
      </div>

      {loading && (
        <div className="text-center py-12">
          <div className="inline-block animate-spin rounded-full h-8 w-8 border-4 border-gray-600 border-t-green-500"></div>
          <p className="text-gray-500 mt-2">Lade…</p>
        </div>
      )}

      {!loading && messages.length === 0 && (
        <div className="text-center py-12 bg-gray-800 rounded">
          <p className="text-gray-500 mb-4">Noch keine automatisierten Nachrichten vorhanden</p>
          <button
            onClick={() => setShowCreateModal(true)}
            className="px-4 py-2 bg-green-600 rounded hover:bg-green-500 transition-colors"
          >
            Erste Nachricht erstellen
          </button>
        </div>
      )}

      <div className="space-y-2">
        {messages.map((m) => {
          const lastSent = formatRelative(m.last_sent_at);
          return (
            <div
              key={m.id}
              className="flex flex-col sm:flex-row justify-between items-start sm:items-center bg-gray-800 rounded p-4 gap-3"
            >
              <div className="flex-1 min-w-0 w-full">
                <div className="flex items-center gap-2 mb-1 flex-wrap">
                  <span className="text-xs px-2 py-0.5 bg-gray-700 rounded text-gray-300">
                    ⏱️ alle {Math.round(m.interval_seconds / 60)} Min.
                  </span>
                  {m.only_when_live && (
                    <span className="text-xs px-2 py-0.5 bg-gray-700 rounded text-gray-300">
                      🔴 nur wenn live
                    </span>
                  )}
                  {lastSent && (
                    <span className="text-xs text-gray-500">Zuletzt gepostet: {lastSent}</span>
                  )}
                </div>
                <p className="text-sm text-gray-300 break-words">{m.message}</p>
              </div>

              <div className="flex items-center gap-2 flex-shrink-0">
                <button
                  onClick={() => handleToggleEnabled(m)}
                  className={`text-xs px-3 py-1.5 rounded transition-colors ${
                    m.enabled ? "bg-green-600 hover:bg-green-500" : "bg-red-600 hover:bg-red-500"
                  }`}
                  title={m.enabled ? "Deaktivieren" : "Aktivieren"}
                >
                  {m.enabled ? "✓ Aktiv" : "✗ Inaktiv"}
                </button>
                <button
                  onClick={() => setEditingMessage(m)}
                  className="text-xs px-3 py-1.5 rounded bg-blue-600 hover:bg-blue-500 transition-colors"
                  title="Bearbeiten"
                >
                  ✏️ Edit
                </button>
                <button
                  onClick={() => handleDelete(m)}
                  className="text-xs px-3 py-1.5 rounded bg-gray-700 hover:bg-red-600 transition-colors"
                  title="Löschen"
                >
                  🗑️
                </button>
              </div>
            </div>
          );
        })}
      </div>

      {!loading && totalPages > 1 && (
        <div className="flex justify-between items-center mt-6 pt-6 border-t border-gray-700">
          <button
            disabled={page <= 1}
            onClick={() => setPage((p) => p - 1)}
            className="px-4 py-2 bg-gray-700 rounded disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-600 transition-colors"
          >
            ← Zurück
          </button>
          <span className="text-gray-400">
            Seite {page} von {totalPages}
          </span>
          <button
            disabled={page >= totalPages}
            onClick={() => setPage((p) => p + 1)}
            className="px-4 py-2 bg-gray-700 rounded disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-600 transition-colors"
          >
            Weiter →
          </button>
        </div>
      )}

      {showCreateModal && (
        <CreateScheduledMessageModal
          onClose={() => setShowCreateModal(false)}
          onCreated={handleCreated}
        />
      )}

      {editingMessage && (
        <EditScheduledMessageModal
          scheduledMessage={editingMessage}
          onClose={() => setEditingMessage(null)}
          onUpdated={handleUpdated}
        />
      )}
    </div>
  );
}
