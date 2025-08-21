import { useEffect, useState } from 'react';

type TwitchStats = {
  twitchUserId: string;
  displayName: string;
  description: string;
  profileImageUrl: string;
  offlineImageUrl: string;
  broadcasterType: string;
  viewCount: number;
  followers: number;
  accountCreatedAt: string;
  fetchedAt: string;
};

const TwitchStatsCard: React.FC = () => {
  const [stats, setStats] = useState<TwitchStats | null>(null);
  const [loading, setLoading] = useState(false);

  const refreshStats = async () => {
    setLoading(true);
    try {
      const username = "goalkeeper91"; // <- vorerst hardcoded
      const res = await fetch(`${import.meta.env.VITE_BACKEND_URL}/api/twitch/stats/refresh?username=${username}`, {
        method: 'POST',
      });
      if (!res.ok) throw new Error('Fehler beim Aktualisieren');
      const data = await res.json();
      setStats(data);
    } catch (err) {
      console.error('Aktualisierungsfehler:', err);
    } finally {
      setLoading(false);
    }
  };


  useEffect(() => {
    refreshStats();
  }, []);

  if (loading) return <p>⏳ Lädt…</p>;

  if (!stats) return <p>⚠️ Keine Daten vorhanden</p>;

  return (
    <div className="p-4 bg-gray-900 text-white rounded-xl shadow-md max-w-md mx-auto">
      <h2 className="text-2xl font-bold mb-2">{stats.displayName}</h2>
      <img
        src={stats.profileImageUrl}
        alt="Profile"
        className="w-24 h-24 rounded-full mb-2"
      />
      <p className="text-sm italic mb-2">{stats.description}</p>
      <ul>
        <li><strong>Follower:</strong> {stats.followers}</li>
        <li><strong>Views:</strong> {stats.viewCount}</li>
        <li><strong>Typ:</strong> {stats.broadcasterType}</li>
        <li><strong>Erstellt:</strong> {new Date(stats.accountCreatedAt).toLocaleDateString()}</li>
        <li><strong>Zuletzt aktualisiert:</strong> {new Date(stats.fetchedAt).toLocaleString()}</li>
      </ul>

      {/* Optional Reload Button */}
      <button
        onClick={refreshStats}
        className="mt-4 px-4 py-2 bg-blue-600 rounded hover:bg-blue-500"
      >
        🔄 Aktualisieren
      </button>
    </div>
  );
};

export default TwitchStatsCard;
