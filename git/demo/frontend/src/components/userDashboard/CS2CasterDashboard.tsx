import { useEffect, useState } from "react";
import {
  getCS2Settings,
  updateCS2Settings,
  getCS2LiveStatus,
  getCS2Notes,
  createCS2Note,
  updateCS2Note,
  deleteCS2Note,
} from "../../features/cs2/api";
import type { CS2CasterSettings, CS2LiveStatus, CS2Note, CS2NoteSubjectType } from "../../features/cs2/types";

function buildGSIConfig(gsiToken: string): string {
  const uri = `${window.location.origin}/api/cs2/gsi/${gsiToken}`;
  return `"Punishers CS2 Caster Tools"
{
 "uri" "${uri}"
 "timeout" "5.0"
 "buffer" "0.1"
 "throttle" "0.5"
 "heartbeat" "30.0"
 "data"
 {
  "provider" "1"
  "map" "1"
  "round" "1"
  "player_id" "1"
  "player_state" "1"
  "allplayers_id" "1"
  "allplayers_state" "1"
  "allplayers_match_stats" "1"
 }
}
`;
}

export default function CS2CasterDashboard() {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [saving, setSaving] = useState(false);
  const [saved, setSaved] = useState(false);
  const [copied, setCopied] = useState(false);

  const [settings, setSettings] = useState<CS2CasterSettings | null>(null);
  const [liveStatus, setLiveStatus] = useState<CS2LiveStatus | null>(null);

  const [notes, setNotes] = useState<CS2Note[]>([]);
  const [notesLoading, setNotesLoading] = useState(true);
  const [notesError, setNotesError] = useState("");

  const [newSubjectType, setNewSubjectType] = useState<CS2NoteSubjectType>("team");
  const [newSubjectName, setNewSubjectName] = useState("");
  const [newContent, setNewContent] = useState("");
  const [creating, setCreating] = useState(false);

  const [editingId, setEditingId] = useState<number | null>(null);
  const [editContent, setEditContent] = useState("");

  useEffect(() => {
    getCS2Settings()
      .then(setSettings)
      .catch((err) => setError(err.message || "Fehler beim Laden"))
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => {
    loadNotes();
  }, []);

  useEffect(() => {
    const poll = () => {
      getCS2LiveStatus()
        .then(setLiveStatus)
        .catch(() => {});
    };
    poll();
    const interval = setInterval(poll, 10000);
    return () => clearInterval(interval);
  }, []);

  function loadNotes() {
    setNotesLoading(true);
    getCS2Notes()
      .then(setNotes)
      .catch((err) => setNotesError(err.message || "Fehler beim Laden der Notizen"))
      .finally(() => setNotesLoading(false));
  }

  async function handleToggle(field: keyof CS2CasterSettings, value: boolean) {
    if (!settings) return;
    const previous = settings;
    setSettings({ ...settings, [field]: value });
    setSaving(true);
    setSaved(false);
    setError("");

    try {
      const updated = await updateCS2Settings({ [field]: value });
      setSettings(updated);
      setSaved(true);
      setTimeout(() => setSaved(false), 3000);
    } catch (err: any) {
      setSettings(previous);
      setError(err.message || "Fehler beim Speichern");
    } finally {
      setSaving(false);
    }
  }

  function openMatchNotesPopup() {
    window.open(
      "/cs2/match-notes-popup",
      "cs2_match_notes",
      "width=460,height=820,resizable=yes,scrollbars=yes"
    );
  }

  async function handleCopy(text: string) {
    try {
      await navigator.clipboard.writeText(text);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch {
      // Clipboard-Zugriff kann in manchen Browser-Kontexten fehlschlagen - kein harter Fehler nötig.
    }
  }

  async function handleCreateNote(e: React.FormEvent) {
    e.preventDefault();
    if (!newSubjectName.trim() || !newContent.trim()) return;

    setCreating(true);
    setNotesError("");
    try {
      await createCS2Note({
        subject_type: newSubjectType,
        subject_name: newSubjectName.trim(),
        content: newContent.trim(),
      });
      setNewSubjectName("");
      setNewContent("");
      loadNotes();
    } catch (err: any) {
      setNotesError(err.message || "Fehler beim Anlegen der Notiz");
    } finally {
      setCreating(false);
    }
  }

  function startEdit(note: CS2Note) {
    setEditingId(note.id);
    setEditContent(note.content);
  }

  async function handleSaveEdit(id: number) {
    try {
      await updateCS2Note(id, { content: editContent });
      setEditingId(null);
      loadNotes();
    } catch (err: any) {
      setNotesError(err.message || "Fehler beim Aktualisieren der Notiz");
    }
  }

  async function handleDeleteNote(id: number) {
    try {
      await deleteCS2Note(id);
      loadNotes();
    } catch (err: any) {
      setNotesError(err.message || "Fehler beim Löschen der Notiz");
    }
  }

  if (loading) {
    return (
      <div className="max-w-3xl mx-auto p-4 sm:p-6 animate-pulse space-y-4">
        <div className="h-8 bg-gray-800 rounded w-1/3"></div>
        <div className="h-64 bg-gray-800 rounded"></div>
      </div>
    );
  }

  const teamNotes = notes.filter((n) => n.subject_type === "team");
  const playerNotes = notes.filter((n) => n.subject_type === "player");

  return (
    <div className="max-w-3xl mx-auto space-y-6 p-4 sm:p-0">
      <div>
        <h1 className="text-2xl font-bold">CS2 Caster Tools</h1>
        <p className="text-sm text-gray-400 mt-1">
          Automatisierungen fürs CS2-Casting (DACH CS / ESEA) auf Basis der Game State
          Integration (GSI) sowie strukturierte Notizen pro Team/Spieler.
        </p>
      </div>

      {error && (
        <div className="bg-red-900/30 border border-red-500/50 text-red-300 rounded-lg p-3 text-sm">
          {error}
        </div>
      )}

      {/* GSI-Einrichtung */}
      {settings && (
        <div className="bg-gray-800 border border-gray-700 rounded-xl p-4 sm:p-6 space-y-3">
          <h3 className="font-semibold text-white">GSI-Einrichtung</h3>
          <p className="text-xs text-gray-400">
            Lege diese Datei als <span className="font-mono">cs2_gsi_punishers.cfg</span> in deinem
            CS2-Ordner unter <span className="font-mono">game/csgo/cfg/</span> ab (Observer-/Spectator-Slot
            reicht, keine Server-Admin-Rechte nötig).
          </p>
          <div className="relative">
            <pre className="bg-gray-900 border border-gray-700 rounded-lg p-3 text-xs text-gray-300 overflow-x-auto whitespace-pre">
              {buildGSIConfig(settings.gsi_token)}
            </pre>
            <button
              onClick={() => handleCopy(buildGSIConfig(settings.gsi_token))}
              className="absolute top-2 right-2 px-2 py-1 bg-gray-700 hover:bg-gray-600 rounded text-xs transition-colors"
            >
              {copied ? "✓ Kopiert" : "Kopieren"}
            </button>
          </div>
        </div>
      )}

      {/* Automatisierungen */}
      {settings && (
        <div className="bg-gray-800 border border-gray-700 rounded-xl p-4 sm:p-6 space-y-4">
          <div className="flex flex-wrap items-center justify-between gap-2">
            <h3 className="font-semibold text-white">Automatisierungen</h3>
            {saved && <span className="text-sm text-green-400">✓ Gespeichert</span>}
          </div>

          <label className="flex items-center justify-between cursor-pointer gap-3">
            <div>
              <p className="font-medium text-white text-sm">Kanalpunkte-Wette bei Match-Start</p>
              <p className="text-xs text-gray-400 mt-0.5">
                Erstellt automatisch eine Twitch-Prediction "Wer gewinnt die Map?" mit den beiden
                Team-Namen, sobald die Map live geht.
              </p>
            </div>
            <input
              type="checkbox"
              checked={settings.predictions_enabled}
              disabled={saving}
              onChange={(e) => handleToggle("predictions_enabled", e.target.checked)}
              className="w-5 h-5 rounded bg-gray-700 border-gray-600 text-indigo-600 focus:ring-2 focus:ring-indigo-500 flex-shrink-0"
            />
          </label>

          <label className="flex items-center justify-between cursor-pointer gap-3 pt-3 border-t border-gray-700">
            <div>
              <p className="font-medium text-white text-sm">Multikill-Ankündigung im Chat</p>
              <p className="text-xs text-gray-400 mt-0.5">
                Postet automatisch bei 3/4/5 Kills eines Spielers in einer Runde (Ace).
              </p>
            </div>
            <input
              type="checkbox"
              checked={settings.multikill_announce_enabled}
              disabled={saving}
              onChange={(e) => handleToggle("multikill_announce_enabled", e.target.checked)}
              className="w-5 h-5 rounded bg-gray-700 border-gray-600 text-indigo-600 focus:ring-2 focus:ring-indigo-500 flex-shrink-0"
            />
          </label>

          <label className="flex items-center justify-between cursor-pointer gap-3 pt-3 border-t border-gray-700">
            <div>
              <p className="font-medium text-white text-sm">"GG WP"-Ankündigung am Map-Ende</p>
              <p className="text-xs text-gray-400 mt-0.5">
                Postet automatisch eine Glückwunsch-Nachricht ans Sieger-Team, sobald die Map endet.
              </p>
            </div>
            <input
              type="checkbox"
              checked={settings.map_end_announce_enabled}
              disabled={saving}
              onChange={(e) => handleToggle("map_end_announce_enabled", e.target.checked)}
              className="w-5 h-5 rounded bg-gray-700 border-gray-600 text-indigo-600 focus:ring-2 focus:ring-indigo-500 flex-shrink-0"
            />
          </label>

          <label className="flex items-center justify-between cursor-pointer gap-3 pt-3 border-t border-gray-700">
            <div>
              <p className="font-medium text-white text-sm">Automatisches Stream-Titel-Update</p>
              <p className="text-xs text-gray-400 mt-0.5">
                Aktualisiert deinen Stream-Titel automatisch mit dem aktuellen Punktestand.
              </p>
            </div>
            <input
              type="checkbox"
              checked={settings.title_update_enabled}
              disabled={saving}
              onChange={(e) => handleToggle("title_update_enabled", e.target.checked)}
              className="w-5 h-5 rounded bg-gray-700 border-gray-600 text-indigo-600 focus:ring-2 focus:ring-indigo-500 flex-shrink-0"
            />
          </label>
        </div>
      )}

      {/* Live-Status */}
      <div className="bg-gray-800 border border-gray-700 rounded-xl p-4 sm:p-6">
        <h3 className="font-semibold text-white mb-3">Live-Status</h3>
        {liveStatus?.active ? (
          <div className="space-y-1 text-sm">
            <p className="text-white">
              <span className="text-gray-400">Map:</span> {liveStatus.map_name || "—"}
            </p>
            <p className="text-white">
              <span className="text-gray-400">Stand:</span> {liveStatus.team_ct_name || "CT"}{" "}
              {liveStatus.score_ct} : {liveStatus.score_t} {liveStatus.team_t_name || "T"}
            </p>
            {liveStatus.observed_player_name && (
              <p className="text-white">
                <span className="text-gray-400">Beobachtet:</span> {liveStatus.observed_player_name}
              </p>
            )}
          </div>
        ) : (
          <p className="text-sm text-gray-500">Keine aktive GSI-Session.</p>
        )}
      </div>

      {/* Notizen */}
      <div className="bg-gray-800 border border-gray-700 rounded-xl p-4 sm:p-6 space-y-4">
        <div className="flex flex-wrap items-center justify-between gap-2">
          <h3 className="font-semibold text-white">Notizen</h3>
          <button
            onClick={openMatchNotesPopup}
            className="px-3 py-1.5 bg-indigo-600 hover:bg-indigo-500 rounded-lg text-sm transition-colors"
          >
            🗒 Match-Notizen als Popup öffnen
          </button>
        </div>
        <p className="text-xs text-gray-400 -mt-2">
          Öffnet ein eigenständiges Fenster mit den Notizen zu den aktuell spielenden Teams/Spielern -
          ideal zum Ablegen auf einem zweiten Monitor, jederzeit minimier- oder schließbar.
        </p>

        {notesError && (
          <div className="bg-red-900/30 border border-red-500/50 text-red-300 rounded-lg p-3 text-sm">
            {notesError}
          </div>
        )}

        <form onSubmit={handleCreateNote} className="space-y-2 pb-4 border-b border-gray-700">
          <div className="grid sm:grid-cols-2 gap-2">
            <select
              value={newSubjectType}
              onChange={(e) => setNewSubjectType(e.target.value as CS2NoteSubjectType)}
              className="w-full p-2 rounded bg-gray-900 border border-gray-700 text-white text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
            >
              <option value="team">Team</option>
              <option value="player">Spieler</option>
            </select>
            <input
              type="text"
              placeholder={newSubjectType === "team" ? "Team-Name" : "Spieler-Name"}
              value={newSubjectName}
              onChange={(e) => setNewSubjectName(e.target.value)}
              className="w-full p-2 rounded bg-gray-900 border border-gray-700 text-white text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
            />
          </div>
          <textarea
            placeholder="Notiz..."
            value={newContent}
            onChange={(e) => setNewContent(e.target.value)}
            rows={2}
            className="w-full p-2 rounded bg-gray-900 border border-gray-700 text-white text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
          />
          <button
            type="submit"
            disabled={creating || !newSubjectName.trim() || !newContent.trim()}
            className="px-4 py-2 bg-indigo-600 rounded-lg hover:bg-indigo-500 disabled:opacity-50 transition-colors text-sm"
          >
            {creating ? "Anlegen…" : "Notiz anlegen"}
          </button>
        </form>

        {notesLoading ? (
          <p className="text-sm text-gray-500">Lade…</p>
        ) : (
          <>
            <NoteGroup
              title="Teams"
              items={teamNotes}
              editingId={editingId}
              editContent={editContent}
              setEditContent={setEditContent}
              onStartEdit={startEdit}
              onSaveEdit={handleSaveEdit}
              onCancelEdit={() => setEditingId(null)}
              onDelete={handleDeleteNote}
            />
            <NoteGroup
              title="Spieler"
              items={playerNotes}
              editingId={editingId}
              editContent={editContent}
              setEditContent={setEditContent}
              onStartEdit={startEdit}
              onSaveEdit={handleSaveEdit}
              onCancelEdit={() => setEditingId(null)}
              onDelete={handleDeleteNote}
            />
          </>
        )}
      </div>
    </div>
  );
}

function NoteGroup({
  title,
  items,
  editingId,
  editContent,
  setEditContent,
  onStartEdit,
  onSaveEdit,
  onCancelEdit,
  onDelete,
}: {
  title: string;
  items: CS2Note[];
  editingId: number | null;
  editContent: string;
  setEditContent: (v: string) => void;
  onStartEdit: (note: CS2Note) => void;
  onSaveEdit: (id: number) => void;
  onCancelEdit: () => void;
  onDelete: (id: number) => void;
}) {
  if (items.length === 0) return null;

  return (
    <div className="pt-3 first:pt-0">
      <h4 className="text-xs font-semibold text-gray-400 uppercase mb-2">{title}</h4>
      <div className="space-y-2">
        {items.map((note) => (
          <div key={note.id} className="bg-gray-900 rounded-lg p-3 text-sm space-y-2">
            <div className="flex flex-wrap items-center justify-between gap-2">
              <span className="font-medium text-white">{note.subject_name}</span>
              {editingId !== note.id && (
                <div className="flex gap-2">
                  <button
                    onClick={() => onStartEdit(note)}
                    className="text-xs text-indigo-400 hover:text-indigo-300"
                  >
                    Bearbeiten
                  </button>
                  <button
                    onClick={() => onDelete(note.id)}
                    className="text-xs text-red-400 hover:text-red-300"
                  >
                    Löschen
                  </button>
                </div>
              )}
            </div>
            {editingId === note.id ? (
              <div className="space-y-2">
                <textarea
                  value={editContent}
                  onChange={(e) => setEditContent(e.target.value)}
                  rows={2}
                  className="w-full p-2 rounded bg-gray-800 border border-gray-700 text-white text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
                />
                <div className="flex gap-2">
                  <button
                    onClick={() => onSaveEdit(note.id)}
                    className="px-3 py-1 bg-indigo-600 rounded hover:bg-indigo-500 text-xs transition-colors"
                  >
                    Speichern
                  </button>
                  <button
                    onClick={onCancelEdit}
                    className="px-3 py-1 bg-gray-700 rounded hover:bg-gray-600 text-xs transition-colors"
                  >
                    Abbrechen
                  </button>
                </div>
              </div>
            ) : (
              <p className="text-gray-300 whitespace-pre-wrap">{note.content}</p>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
