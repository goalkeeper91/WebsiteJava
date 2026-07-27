import { useEffect, useState } from "react";
import { getTeamMembers, inviteTeamMember, removeTeamMember } from "../../features/team/api";
import type { TeamMember } from "../../features/team/types";

function formatRelative(dateString: string): string {
  const seconds = Math.max(0, Math.floor((Date.now() - new Date(dateString).getTime()) / 1000));
  if (seconds < 60) return "gerade eben";
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return `vor ${minutes} Min.`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `vor ${hours} Std.`;
  return `vor ${Math.floor(hours / 24)} Tag(en)`;
}

export default function TeamDashboard() {
  const [loading, setLoading] = useState(true);
  const [members, setMembers] = useState<TeamMember[]>([]);
  const [error, setError] = useState("");

  const [login, setLogin] = useState("");
  const [inviting, setInviting] = useState(false);
  const [inviteError, setInviteError] = useState("");

  const [removingID, setRemovingID] = useState<string | null>(null);

  const loadMembers = () => {
    getTeamMembers()
      .then(setMembers)
      .catch((err) => setError(err.message || "Fehler beim Laden"))
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    loadMembers();
  }, []);

  async function handleInvite(e: React.FormEvent) {
    e.preventDefault();
    setInviting(true);
    setInviteError("");
    try {
      await inviteTeamMember(login.trim());
      setLogin("");
      loadMembers();
    } catch (err: any) {
      setInviteError(err.message || "Fehler beim Einladen");
    } finally {
      setInviting(false);
    }
  }

  async function handleRemove(memberTwitchId: string) {
    setRemovingID(memberTwitchId);
    try {
      await removeTeamMember(memberTwitchId);
      setMembers((prev) => prev.filter((m) => m.member_twitch_id !== memberTwitchId));
    } catch (err) {
      alert("Fehler beim Entfernen des Team-Mitglieds");
    } finally {
      setRemovingID(null);
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

  return (
    <div className="max-w-3xl mx-auto space-y-6 p-4 sm:p-0">
      <div>
        <h1 className="text-2xl font-bold">Team</h1>
        <p className="text-sm text-gray-400 mt-1">
          Lade Personen (z.B. deine Twitch-Mods) per Twitch-Login ein, damit sie dein Dashboard mitverwalten können —
          Automod, Loyalty, Chat-Commands, Automatisierte Nachrichten und Giveaways, mit denselben Rechten wie du
          selbst.
        </p>
      </div>

      {error && (
        <div className="bg-red-900/30 border border-red-500/50 text-red-300 rounded-lg p-3 text-sm">{error}</div>
      )}

      <div className="bg-gray-800 border border-gray-700 rounded-xl p-4 sm:p-6">
        <h3 className="font-semibold text-white mb-3">Person einladen</h3>
        <form onSubmit={handleInvite} className="space-y-3">
          {inviteError && <p className="text-sm text-red-400">{inviteError}</p>}
          <div className="flex flex-wrap gap-3 items-center">
            <input
              type="text"
              value={login}
              onChange={(e) => setLogin(e.target.value)}
              placeholder="Twitch-Login (z.B. meinmod123)"
              className="bg-gray-900 border border-gray-700 rounded px-3 py-2 text-sm text-white flex-1 min-w-[150px]"
            />
            <button
              type="submit"
              disabled={inviting || !login.trim()}
              className="px-4 py-2 bg-indigo-600 rounded disabled:opacity-50 disabled:cursor-not-allowed hover:bg-indigo-500 transition-colors text-sm font-semibold"
            >
              {inviting ? "Lade ein…" : "Einladen"}
            </button>
          </div>
        </form>
      </div>

      <div className="bg-gray-800 border border-gray-700 rounded-xl p-4 sm:p-6">
        <h3 className="font-semibold text-white mb-3">Wer hat Zugriff</h3>

        {members.length === 0 && <p className="text-sm text-gray-500">Noch niemand eingeladen.</p>}

        <div className="space-y-2">
          {members.map((member) => (
            <div
              key={member.id}
              className="bg-gray-900 rounded-lg p-3 text-sm flex flex-wrap items-center justify-between gap-2"
            >
              <span className="text-white font-semibold">{member.member_login}</span>
              <div className="flex items-center gap-3">
                <span className="text-xs text-gray-500">Seit {formatRelative(member.created_at)}</span>
                <button
                  onClick={() => handleRemove(member.member_twitch_id)}
                  disabled={removingID === member.member_twitch_id}
                  className="px-3 py-1.5 bg-gray-700 rounded disabled:opacity-50 disabled:cursor-not-allowed hover:bg-red-700 transition-colors text-xs"
                >
                  {removingID === member.member_twitch_id ? "Entferne…" : "Entfernen"}
                </button>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
