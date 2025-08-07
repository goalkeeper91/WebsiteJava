import React, { useEffect, useState } from 'react';

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

    const fetchStats = async () => {
        setLoading(true);
        try {
            const res = await fetch('http://localhost:8080/twitch/stats');
            if (!res.ok) throw new Error('Fehler beim Laden');
            const data = await res.json();
            setStats(data);
        } catch (err) {
            console.error('Ladefehler:', err);
            setStats(null);
        } finally {
            setLoading(false);
        }
    };

    const refreshStats = async () => {
        setLoading(true);
        try {
            const res = await fetch('http://localhost:8080/twitch/stats/refresh', {
               method: 'POST',
            });
            if (!res.ok) throw new Error('Fehler beim Aktualisieren');
            const data = await res.json();
            setStats(data);
        } catch (err) {
            console.error('Aktualisierungsfehler:', err)
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
    </div>
  );
};

export default TwitchStatsCard;