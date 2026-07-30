import { useEffect, useMemo, useState } from "react";
import { getCS2LiveStatus, getCS2Notes } from "../features/cs2/api";
import type { CS2LiveStatus, CS2Note } from "../features/cs2/types";

const LIVE_STATUS_POLL_MS = 7000;
const NOTES_POLL_MS = 15000;

function findNote(notes: CS2Note[], subjectType: "team" | "player", name: string): CS2Note | undefined {
  const needle = name.trim().toLowerCase();
  return notes.find((n) => n.subject_type === subjectType && n.subject_name.trim().toLowerCase() === needle);
}

interface CardProps {
  title: string;
  content: string | undefined;
  highlighted?: boolean;
  emptyHint?: string;
}

function ModerationCard({ title, content, highlighted, emptyHint }: CardProps) {
  return (
    <div
      className={`rounded-lg p-3 border ${
        highlighted ? "border-indigo-400 bg-indigo-950/40 shadow-lg shadow-indigo-900/30" : "border-gray-700 bg-gray-800"
      }`}
    >
      <div className="flex items-center gap-1.5 mb-1">
        {highlighted && <span title="Aktuell beobachtet">👁</span>}
        <span className="font-semibold text-white text-sm">{title}</span>
      </div>
      {content ? (
        <p className="text-gray-300 text-sm whitespace-pre-wrap">{content}</p>
      ) : (
        <p className="text-gray-500 text-xs italic">{emptyHint || "Keine Notiz vorhanden."}</p>
      )}
    </div>
  );
}

function TeamColumn({
  sideLabel,
  teamName,
  players,
  notes,
  observedPlayerName,
}: {
  sideLabel: string;
  teamName: string;
  players: string[];
  notes: CS2Note[];
  observedPlayerName?: string;
}) {
  const teamNote = findNote(notes, "team", teamName);
  return (
    <div className="space-y-2 min-w-0">
      <div className="text-xs font-semibold text-gray-400 uppercase tracking-wide">{sideLabel}</div>
      <ModerationCard title={teamName || "Unbekanntes Team"} content={teamNote?.content} emptyHint="Keine Team-Notiz." />
      {players.map((name) => {
        const note = findNote(notes, "player", name);
        const isObserved = !!observedPlayerName && name.trim().toLowerCase() === observedPlayerName.trim().toLowerCase();
        return (
          <ModerationCard
            key={name}
            title={name}
            content={note?.content}
            highlighted={isObserved}
            emptyHint="Keine Notiz für diesen Spieler."
          />
        );
      })}
    </div>
  );
}

export default function CS2MatchNotesPopup() {
  const [liveStatus, setLiveStatus] = useState<CS2LiveStatus | null>(null);
  const [notes, setNotes] = useState<CS2Note[]>([]);
  const [error, setError] = useState("");
  const [manualTeam, setManualTeam] = useState("");

  useEffect(() => {
    document.title = "CS2 Match-Notizen";
  }, []);

  useEffect(() => {
    const poll = () => {
      getCS2LiveStatus()
        .then(setLiveStatus)
        .catch((err) => setError(err.message || "Fehler beim Laden des Live-Status"));
    };
    poll();
    const interval = setInterval(poll, LIVE_STATUS_POLL_MS);
    return () => clearInterval(interval);
  }, []);

  useEffect(() => {
    const poll = () => {
      getCS2Notes()
        .then(setNotes)
        .catch((err) => setError(err.message || "Fehler beim Laden der Notizen"));
    };
    poll();
    const interval = setInterval(poll, NOTES_POLL_MS);
    return () => clearInterval(interval);
  }, []);

  const knownTeamNames = useMemo(
    () => Array.from(new Set(notes.filter((n) => n.subject_type === "team").map((n) => n.subject_name))).sort(),
    [notes]
  );

  const hasLiveMatch = !!(liveStatus && (liveStatus.team_ct_name || liveStatus.team_t_name));

  return (
    <div className="min-h-screen bg-gray-900 text-white p-4">
      <div className="flex items-center justify-between gap-2 mb-4">
        <h1 className="text-lg font-bold">🗒 Match-Notizen</h1>
        <button
          onClick={() => window.close()}
          className="px-3 py-1.5 bg-gray-700 hover:bg-gray-600 rounded text-sm transition-colors"
        >
          ✕ Schließen
        </button>
      </div>

      {error && (
        <div className="bg-red-900/30 border border-red-500/50 text-red-300 rounded-lg p-3 text-sm mb-4">{error}</div>
      )}

      {hasLiveMatch && liveStatus ? (
        <>
          <div className="text-center text-sm text-gray-400 mb-3">
            {liveStatus.map_name && <span className="mr-2">{liveStatus.map_name}</span>}
            <span className="font-semibold text-white">
              {liveStatus.score_ct} : {liveStatus.score_t}
            </span>
          </div>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <TeamColumn
              sideLabel="CT"
              teamName={liveStatus.team_ct_name || "CT-Team"}
              players={liveStatus.team_ct_players || []}
              notes={notes}
              observedPlayerName={liveStatus.observed_player_name}
            />
            <TeamColumn
              sideLabel="T"
              teamName={liveStatus.team_t_name || "T-Team"}
              players={liveStatus.team_t_players || []}
              notes={notes}
              observedPlayerName={liveStatus.observed_player_name}
            />
          </div>
        </>
      ) : (
        <div className="space-y-3">
          <p className="text-sm text-gray-400">
            Kein aktives Match erkannt. Wähle ein Team aus deinen gespeicherten Notizen, um dich schon vorab
            vorzubereiten - sobald GSI verbunden ist, schaltet dieses Fenster automatisch auf den Live-Kader um.
          </p>
          <select
            value={manualTeam}
            onChange={(e) => setManualTeam(e.target.value)}
            className="w-full p-2 rounded bg-gray-800 border border-gray-700 text-white text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
          >
            <option value="">Team wählen…</option>
            {knownTeamNames.map((name) => (
              <option key={name} value={name}>
                {name}
              </option>
            ))}
          </select>
          {manualTeam && (
            <ModerationCard title={manualTeam} content={findNote(notes, "team", manualTeam)?.content} />
          )}
        </div>
      )}
    </div>
  );
}
