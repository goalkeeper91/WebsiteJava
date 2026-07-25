import { useEffect, useState } from "react";
import { getCommands, deleteCommand, updateCommand } from "../../features/commands/api";
import type { ChatCommand } from "../../features/commands/types";
import CreateCommandModal from "../modals/dashboard/CreateCommandModal";
import EditCommandModal from "../modals/dashboard/EditCommandModal";

export default function CommandsDashboard() {
  const [commands, setCommands] = useState<ChatCommand[]>([]);
  const [page, setPage] = useState(0);
  const [totalPages, setTotalPages] = useState(0);
  const [totalElements, setTotalElements] = useState(0);
  const [loading, setLoading] = useState(true);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [editingCommand, setEditingCommand] = useState<ChatCommand | null>(null);
  const [searchTerm, setSearchTerm] = useState("");

  useEffect(() => {
    loadCommands();
  }, [page]);

  async function loadCommands() {
    setLoading(true);
    try {
      const data = await getCommands(page, 10);
      console.log("Loaded commands:", data);

      setCommands(data.data || []);
      setTotalPages(data.totalPages || 0);
      setTotalElements(data.total || 0);
    } catch (err) {
      console.error("Fehler beim Laden der Commands:", err);
      alert("Fehler beim Laden der Commands");
    } finally {
      setLoading(false);
    }
  }

  const filteredCommands = commands.filter((cmd) =>
    cmd.trigger.toLowerCase().includes(searchTerm.toLowerCase()) ||
    cmd.response.toLowerCase().includes(searchTerm.toLowerCase())
  );

  async function handleToggleEnabled(cmd: ChatCommand) {
    setCommands(prev =>
      prev.map(c => c.id === cmd.id ? { ...c, enabled: !c.enabled } : c)
    );

    try {
      await updateCommand(cmd.id, {
        enabled: !cmd.enabled
      });

      console.log(`Command ${cmd.trigger} status updated`);
    } catch (err) {
      console.error("Fehler beim Update:", err);

      setCommands(prev =>
        prev.map(c => c.id === cmd.id ? { ...c, enabled: cmd.enabled } : c)
      );

      alert("Fehler beim Aktualisieren des Status");
    }
  }

  async function handleDelete(cmd: ChatCommand) {
    if (!confirm(`Command "!${cmd.trigger}" wirklich löschen?`)) {
      return;
    }

    setCommands(prev => prev.filter(c => c.id !== cmd.id));
    setTotalElements(prev => prev - 1);

    try {
      await deleteCommand(cmd.id);
      console.log(`Command ${cmd.trigger} deleted`);
    } catch (err) {
      console.error("Fehler beim Löschen:", err);

      loadCommands();

      alert("Fehler beim Löschen des Commands");
    }
  }

  function handleCommandCreated(newCommand: ChatCommand) {
    setCommands(prev => [newCommand, ...prev]);
    setTotalElements(prev => prev + 1);
    setShowCreateModal(false);
  }

  function handleCommandUpdated(updatedCommand: ChatCommand) {
    setCommands(prev =>
      prev.map(c => c.id === updatedCommand.id ? updatedCommand : c)
    );
    setEditingCommand(null);
  }

  return (
    <div className="max-w-5xl mx-auto">
      {/* Header */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
        <div>
          <h1 className="text-2xl font-bold">Chat Commands</h1>
          <p className="text-sm text-gray-400 mt-1">
            {loading ? "Lade..." : `${totalElements} Commands insgesamt`}
          </p>
        </div>
        <button
          onClick={() => setShowCreateModal(true)}
          className="px-4 py-2 bg-green-600 rounded hover:bg-green-500 transition-colors flex items-center gap-2"
        >
          <span className="text-xl">+</span>
          Neues Command
        </button>
      </div>

      {/* Search Bar */}
      <div className="mb-4">
        <input
          type="text"
          placeholder="Commands durchsuchen..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="w-full bg-gray-800 rounded px-4 py-2 focus:outline-none focus:ring-2 focus:ring-green-500"
        />
      </div>

      {/* Loading State */}
      {loading && (
        <div className="text-center py-12">
          <div className="inline-block animate-spin rounded-full h-8 w-8 border-4 border-gray-600 border-t-green-500"></div>
          <p className="text-gray-500 mt-2">Lade Commands…</p>
        </div>
      )}

      {/* Empty State */}
      {!loading && commands.length === 0 && (
        <div className="text-center py-12 bg-gray-800 rounded">
          <p className="text-gray-500 mb-4">
            Noch keine Commands vorhanden
          </p>
          <button
            onClick={() => setShowCreateModal(true)}
            className="px-4 py-2 bg-green-600 rounded hover:bg-green-500 transition-colors"
          >
            Erstes Command erstellen
          </button>
        </div>
      )}

      {/* No Search Results */}
      {!loading && commands.length > 0 && filteredCommands.length === 0 && (
        <div className="text-center py-12 bg-gray-800 rounded">
          <p className="text-gray-500">Keine Commands gefunden für "{searchTerm}"</p>
        </div>
      )}

      {/* Commands List */}
      <div className="space-y-2">
        {filteredCommands.map((cmd) => (
          <div
            key={cmd.id}
            className="flex flex-col sm:flex-row justify-between items-start sm:items-center bg-gray-800 rounded p-4 gap-3 hover:bg-gray-750 transition-colors"
          >
            <div className="flex-1 min-w-0 w-full">
              <div className="flex items-center gap-2 mb-1">
                <p className="font-mono text-green-400 font-semibold">
                  !{cmd.trigger}
                </p>
                {cmd.cooldown > 0 && (
                  <span className="text-xs px-2 py-0.5 bg-gray-700 rounded text-gray-300">
                    ⏱️ {cmd.cooldown}s
                  </span>
                )}
              </div>
              <p className="text-sm text-gray-300 break-words">
                {cmd.response}
              </p>
            </div>

            {/* Action Buttons */}
            <div className="flex items-center gap-2 flex-shrink-0">
              {/* Toggle Enabled/Disabled */}
              <button
                onClick={() => handleToggleEnabled(cmd)}
                className={`text-xs px-3 py-1.5 rounded transition-colors ${
                  cmd.enabled
                    ? "bg-green-600 hover:bg-green-500"
                    : "bg-red-600 hover:bg-red-500"
                }`}
                title={cmd.enabled ? "Deaktivieren" : "Aktivieren"}
              >
                {cmd.enabled ? "✓ Aktiv" : "✗ Inaktiv"}
              </button>

              {/* Edit Button */}
              <button
                onClick={() => setEditingCommand(cmd)}
                className="text-xs px-3 py-1.5 rounded bg-blue-600 hover:bg-blue-500 transition-colors"
                title="Bearbeiten"
              >
                ✏️ Edit
              </button>

              {/* Delete Button */}
              <button
                onClick={() => handleDelete(cmd)}
                className="text-xs px-3 py-1.5 rounded bg-gray-700 hover:bg-red-600 transition-colors"
                title="Löschen"
              >
                🗑️
              </button>
            </div>
          </div>
        ))}
      </div>

      {/* Pagination */}
      {!loading && totalPages > 1 && (
        <div className="flex flex-wrap justify-between items-center gap-2 mt-6 pt-6 border-t border-gray-700">
          <button
            disabled={page === 0}
            onClick={() => setPage((p) => p - 1)}
            className="px-4 py-2 bg-gray-700 rounded disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-600 transition-colors"
          >
            ← Zurück
          </button>
          <span className="text-gray-400">
            Seite {page + 1} von {totalPages}
          </span>
          <button
            disabled={page + 1 >= totalPages}
            onClick={() => setPage((p) => p + 1)}
            className="px-4 py-2 bg-gray-700 rounded disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-600 transition-colors"
          >
            Weiter →
          </button>
        </div>
      )}

      {/* Modals */}
      {showCreateModal && (
        <CreateCommandModal
          onClose={() => setShowCreateModal(false)}
          onCreated={handleCommandCreated}
        />
      )}

      {editingCommand && (
        <EditCommandModal
          command={editingCommand}
          onClose={() => setEditingCommand(null)}
          onUpdated={handleCommandUpdated}
        />
      )}
    </div>
  );
}