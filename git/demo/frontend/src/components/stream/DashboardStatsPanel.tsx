import { useState, useEffect } from "react";
import { fetchDashboardStats } from "../../features/stream/api";
import type { DashboardStats } from "../../features/stream/types";

export default function DashboardStatsPanel() {
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadStats();

    // Aktualisiere Stats alle 30 Sekunden
    const interval = setInterval(loadStats, 30000);
    return () => clearInterval(interval);
  }, []);

  async function loadStats() {
    try {
      const data = await fetchDashboardStats();
      setStats(data);
      setLoading(false);
    } catch (err) {
      console.error("Fehler beim Laden der Stats:", err);
      setLoading(false);
    }
  }

  if (loading) {
    return (
      <div className="grid grid-cols-2 md:grid-cols-5 gap-4">
        {[1, 2, 3, 4, 5].map((i) => (
          <div key={i} className="bg-gray-800 rounded-lg p-4 animate-pulse">
            <div className="h-4 bg-gray-700 rounded w-1/2 mb-2"></div>
            <div className="h-8 bg-gray-700 rounded w-3/4"></div>
          </div>
        ))}
      </div>
    );
  }

  return (
    <div className="grid grid-cols-2 md:grid-cols-5 gap-4">
      {/* Viewer Count */}
      <div className="bg-gray-800 rounded-lg p-4">
        <div className="flex items-center gap-2 mb-2">
          <span className="text-2xl">👥</span>
          <p className="text-xs text-gray-400">
            {stats?.isLive ? "Aktuelle Viewer" : "Follower"}
          </p>
        </div>
        <p className="text-2xl font-bold">
          {stats?.isLive
            ? stats.currentViewers.toLocaleString()
            : stats?.followerCount.toLocaleString()
          }
        </p>
        {stats?.isLive && stats.uptime && (
          <p className="text-xs text-purple-400 mt-1">
            🔴 Live seit {stats.uptime}
          </p>
        )}
      </div>

      {/* Follower */}
      <div className="bg-gray-800 rounded-lg p-4">
        <div className="flex items-center gap-2 mb-2">
          <span className="text-2xl">💜</span>
          <p className="text-xs text-gray-400">Follower</p>
        </div>
        <p className="text-2xl font-bold">
          {stats?.followerCount.toLocaleString()}
        </p>
        {stats && stats.followsToday > 0 && (
          <p className="text-xs text-green-400 mt-1">
            +{stats.followsToday} heute
          </p>
        )}
      </div>

      {/* Subscribers */}
      <div className="bg-gray-800 rounded-lg p-4">
        <div className="flex items-center gap-2 mb-2">
          <span className="text-2xl">⭐</span>
          <p className="text-xs text-gray-400">Subs</p>
        </div>
        <p className="text-2xl font-bold">
          {stats?.subscriberCount.toLocaleString()}
        </p>
        {stats && stats.subsThisWeek > 0 && (
          <p className="text-xs text-green-400 mt-1">
            +{stats.subsThisWeek} diese Woche
          </p>
        )}
      </div>

      {/* Bits */}
      <div className="bg-gray-800 rounded-lg p-4">
        <div className="flex items-center gap-2 mb-2">
          <span className="text-2xl">💎</span>
          <p className="text-xs text-gray-400">Bits heute</p>
        </div>
        <p className="text-2xl font-bold">
          {stats?.bitsToday.toLocaleString()}
        </p>
      </div>

      {/* Average Viewers */}
      <div className="bg-gray-800 rounded-lg p-4">
        <div className="flex items-center gap-2 mb-2">
          <span className="text-2xl">📊</span>
          <p className="text-xs text-gray-400">Avg. Viewer</p>
        </div>
        <p className="text-2xl font-bold">
          {stats?.avgViewers.toFixed(0)}
        </p>
        <p className="text-xs text-gray-400 mt-1">
          Letzter Stream
        </p>
      </div>
    </div>
  );
}