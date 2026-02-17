import { useEffect, useState } from "react";
import { Plus, BarChart3, Clock, CheckCircle, XCircle, Trophy } from "lucide-react";

interface VoteOption {
  key: string;
  label: string;
}

interface VoteResult {
  option: string;
  count: number;
  percentage: number;
}

interface VoteSession {
  id: number;
  category: string;
  title: string;
  description?: string;
  options: VoteOption[];
  status: string;
  startedAt: string;
  closedAt?: string;
  winner?: string;
  totalVotes?: number;
  results?: VoteResult[];
}

export default function VoteSessionManager() {
  const [sessions, setSessions] = useState<VoteSession[]>([]);
  const [activeSession, setActiveSession] = useState<VoteSession | null>(null);
  const [loading, setLoading] = useState(true);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showResultsModal, setShowResultsModal] = useState<VoteSession | null>(null);
  const [page, setPage] = useState(0);
  const [totalPages, setTotalPages] = useState(0);

  useEffect(() => {
    loadSessions();
    loadActiveSession();
  }, [page]);

  async function loadSessions() {
    setLoading(true);
    try {
      const res = await fetch(`/api/votes/sessions?page=${page + 1}&page_size=10`, {
        credentials: "include",
      });
      if (res.ok) {
        const data = await res.json();
        setSessions(data.data || []);
        setTotalPages(data.total_pages || 0);
      }
    } catch (err) {
      console.error("Failed to load sessions:", err);
    } finally {
      setLoading(false);
    }
  }

  async function loadActiveSession() {
    try {
      const res = await fetch("/api/votes/sessions/active", {
        credentials: "include",
      });
      if (res.ok) {
        const data = await res.json();
        setActiveSession(data);
      }
    } catch (err) {
      console.error("Failed to load active session:", err);
    }
  }

  async function handleCloseSession(sessionId: number) {
    if (!confirm("Vote Session wirklich schließen?")) {
      return;
    }

    try {
      const res = await fetch(`/api/votes/sessions/${sessionId}/close`, {
        method: "POST",
        credentials: "include",
      });

      if (res.ok) {
        const closedSession = await res.json();
        alert(`Vote geschlossen! Gewinner: ${closedSession.winner}`);
        loadSessions();
        loadActiveSession();
      } else {
        alert("Fehler beim Schließen der Session");
      }
    } catch (err) {
      console.error("Close error:", err);
      alert("Fehler beim Schließen");
    }
  }

  async function viewResults(session: VoteSession) {
    try {
      const res = await fetch(`/api/votes/sessions/${session.id}/results`, {
        credentials: "include",
      });
      if (res.ok) {
        const data = await res.json();
        setShowResultsModal({
          ...session,
          results: data.results,
          totalVotes: data.totalVotes,
        });
      }
    } catch (err) {
      console.error("Failed to load results:", err);
    }
  }

  function formatDate(dateString: string) {
    return new Date(dateString).toLocaleString("de-DE", {
      day: "2-digit",
      month: "2-digit",
      year: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  }

  function getStatusBadge(status: string) {
    switch (status) {
      case "active":
        return (
          <span className="flex items-center gap-1 text-xs bg-green-600 px-2 py-1 rounded">
            <CheckCircle className="w-3 h-3" />
            Aktiv
          </span>
        );
      case "closed":
        return (
          <span className="flex items-center gap-1 text-xs bg-gray-600 px-2 py-1 rounded">
            <XCircle className="w-3 h-3" />
            Geschlossen
          </span>
        );
      default:
        return null;
    }
  }

  if (loading && sessions.length === 0) {
    return (
      <div className="p-6 max-w-6xl mx-auto">
        <div className="animate-pulse space-y-4">
          <div className="h-8 bg-gray-800 rounded w-1/3"></div>
          <div className="h-64 bg-gray-800 rounded"></div>
        </div>
      </div>
    );
  }

  return (
    <div className="p-6 max-w-6xl mx-auto">
      {/* Header */}
      <div className="flex justify-between items-start mb-6">
        <div>
          <h1 className="text-3xl font-bold mb-2">Vote Sessions</h1>
          <p className="text-gray-400">
            Erstelle Abstimmungen für deinen Chat (Hero Voting, Giveaways, etc.)
          </p>
        </div>
        <button
          onClick={() => setShowCreateModal(true)}
          className="flex items-center gap-2 px-4 py-2 bg-green-600 hover:bg-green-500 rounded-lg transition-colors"
        >
          <Plus className="w-5 h-5" />
          Neue Session
        </button>
      </div>

      {/* Active Session Card */}
      {activeSession && (
        <div className="bg-gradient-to-r from-green-900/20 to-blue-900/20 border-2 border-green-500 rounded-xl p-6 mb-6">
          <div className="flex items-start justify-between mb-4">
            <div>
              <h2 className="text-xl font-bold mb-1">{activeSession.title}</h2>
              {activeSession.description && (
                <p className="text-gray-400 text-sm">{activeSession.description}</p>
              )}
            </div>
            <span className="flex items-center gap-1 text-xs bg-green-600 px-3 py-1 rounded font-semibold">
              <Clock className="w-3 h-3 animate-pulse" />
              LIVE
            </span>
          </div>

          <div className="grid md:grid-cols-2 gap-4 mb-4">
            <div>
              <p className="text-sm text-gray-400 mb-2">Optionen:</p>
              <div className="flex flex-wrap gap-2">
                {activeSession.options.map((opt) => (
                  <span
                    key={opt.key}
                    className="px-3 py-1 bg-gray-800 rounded text-sm"
                  >
                    {opt.label}
                  </span>
                ))}
              </div>
            </div>
            <div className="flex items-center justify-end gap-2">
              <button
                onClick={() => viewResults(activeSession)}
                className="flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-500 rounded-lg transition-colors"
              >
                <BarChart3 className="w-4 h-4" />
                Live Results
              </button>
              <button
                onClick={() => handleCloseSession(activeSession.id)}
                className="px-4 py-2 bg-red-600 hover:bg-red-500 rounded-lg transition-colors"
              >
                Schließen
              </button>
            </div>
          </div>

          <p className="text-xs text-gray-500">
            Gestartet: {formatDate(activeSession.startedAt)}
          </p>
        </div>
      )}

      {/* Sessions History */}
      <div className="bg-gray-800 rounded-xl p-6">
        <h2 className="text-xl font-bold mb-4">Alle Sessions</h2>

        {sessions.length === 0 ? (
          <div className="text-center py-12 text-gray-500">
            <BarChart3 className="w-12 h-12 mx-auto mb-3 opacity-50" />
            <p>Noch keine Vote Sessions erstellt</p>
          </div>
        ) : (
          <div className="space-y-3">
            {sessions.map((session) => (
              <div
                key={session.id}
                className="bg-gray-900 rounded-lg p-4 hover:bg-gray-850 transition-colors"
              >
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <div className="flex items-center gap-2 mb-1">
                      <h3 className="font-semibold">{session.title}</h3>
                      {getStatusBadge(session.status)}
                      {session.winner && (
                        <span className="flex items-center gap-1 text-xs bg-yellow-600 px-2 py-1 rounded">
                          <Trophy className="w-3 h-3" />
                          {session.winner}
                        </span>
                      )}
                    </div>
                    {session.description && (
                      <p className="text-sm text-gray-400 mb-2">
                        {session.description}
                      </p>
                    )}
                    <div className="flex items-center gap-4 text-xs text-gray-500">
                      <span>
                        Kategorie: <span className="text-gray-400">{session.category}</span>
                      </span>
                      <span>Gestartet: {formatDate(session.startedAt)}</span>
                      {session.closedAt && (
                        <span>Geschlossen: {formatDate(session.closedAt)}</span>
                      )}
                      {session.totalVotes !== undefined && (
                        <span>
                          Votes: <span className="text-gray-400">{session.totalVotes}</span>
                        </span>
                      )}
                    </div>
                  </div>

                  <button
                    onClick={() => viewResults(session)}
                    className="flex items-center gap-1 px-3 py-1 text-sm bg-gray-800 hover:bg-gray-700 rounded transition-colors"
                  >
                    <BarChart3 className="w-4 h-4" />
                    Ergebnisse
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="flex justify-between items-center mt-6 pt-6 border-t border-gray-700">
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
      </div>

      {/* Results Modal */}
      {showResultsModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="bg-gray-800 rounded-xl p-6 max-w-2xl w-full max-h-[80vh] overflow-y-auto">
            <div className="flex items-start justify-between mb-6">
              <div>
                <h2 className="text-2xl font-bold mb-1">
                  {showResultsModal.title}
                </h2>
                <p className="text-gray-400 text-sm">
                  Gesamt: {showResultsModal.totalVotes} Votes
                </p>
              </div>
              <button
                onClick={() => setShowResultsModal(null)}
                className="text-gray-400 hover:text-white text-2xl"
              >
                ×
              </button>
            </div>

            <div className="space-y-3">
              {showResultsModal.results?.map((result, idx) => (
                <div key={result.option} className="bg-gray-900 rounded-lg p-4">
                  <div className="flex items-center justify-between mb-2">
                    <div className="flex items-center gap-2">
                      {idx === 0 && showResultsModal.status === "closed" && (
                        <Trophy className="w-5 h-5 text-yellow-500" />
                      )}
                      <span className="font-semibold">{result.option}</span>
                    </div>
                    <div className="text-right">
                      <div className="font-bold text-lg">{result.count}</div>
                      <div className="text-sm text-gray-400">
                        {result.percentage.toFixed(1)}%
                      </div>
                    </div>
                  </div>
                  <div className="w-full bg-gray-700 rounded-full h-2">
                    <div
                      className={`h-2 rounded-full transition-all ${
                        idx === 0 ? "bg-green-500" : "bg-blue-500"
                      }`}
                      style={{ width: `${result.percentage}%` }}
                    />
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}

      {/* Create Modal Placeholder */}
      {showCreateModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="bg-gray-800 rounded-xl p-6 max-w-md w-full">
            <h2 className="text-2xl font-bold mb-4">Neue Vote Session</h2>
            <p className="text-gray-400 mb-4">
              Modal wird im nächsten Schritt implementiert (benötigt Form mit Optionen)
            </p>
            <button
              onClick={() => setShowCreateModal(false)}
              className="w-full py-2 bg-gray-700 hover:bg-gray-600 rounded-lg transition-colors"
            >
              Schließen
            </button>
          </div>
        </div>
      )}
    </div>
  );
}