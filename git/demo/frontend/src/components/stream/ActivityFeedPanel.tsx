import { useState, useEffect } from "react";
import { apiFetch } from "../../lib/apiFetch";
import type { Activity } from "../../features/stream/types";

const API_BASE = "/api/dashboard/stream";
const POLL_INTERVAL_MS = 15000;

export default function ActivityFeedPanel() {
  const [activities, setActivities] = useState<Activity[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadActivities = async () => {
      try {
        const res = await apiFetch(`${API_BASE}/activities?limit=50`);
        if (res.ok) {
          const data = await res.json();
          setActivities(data || []);
        }
      } catch (err) {
        console.error("Fehler beim Laden der Activities:", err);
      } finally {
        setLoading(false);
      }
    };

    loadActivities();
    const interval = setInterval(loadActivities, POLL_INTERVAL_MS);
    return () => clearInterval(interval);
  }, []);

  function getActivityIcon(type: Activity["type"]): string {
    switch (type) {
      case "FOLLOW":
        return "💜";
      case "SUBSCRIBE":
      case "RESUBSCRIBE":
        return "⭐";
      case "GIFT_SUB":
        return "🎁";
      case "RAID":
        return "🎉";
      case "CHEER":
        return "💎";
      case "HOSTING":
        return "📺";
      default:
        return "✨";
    }
  }

  function getActivityText(activity: Activity): string {
    switch (activity.type) {
      case "FOLLOW":
        return `${activity.display_name} folgt dir!`;
      case "SUBSCRIBE":
      case "RESUBSCRIBE":
        return `${activity.display_name} hat subscribed!${
          activity.tier === "2000" ? " (Tier 2)" :
          activity.tier === "3000" ? " (Tier 3)" : ""
        }`;
      case "GIFT_SUB":
        return `${activity.display_name} hat ein Sub verschenkt!`;
      case "RAID":
        return `${activity.display_name} raided mit ${activity.viewers} Viewern!`;
      case "CHEER":
        return `${activity.display_name} cheered ${activity.bits} Bits!`;
      case "HOSTING":
        return `${activity.display_name} hostet deinen Stream!`;
      default:
        return `${activity.display_name} - ${activity.type}`;
    }
  }

  function getTimeAgo(timestamp: string): string {
    const seconds = Math.floor((Date.now() - new Date(timestamp).getTime()) / 1000);

    if (seconds < 60) return "Gerade eben";
    if (seconds < 3600) return `vor ${Math.floor(seconds / 60)} Min`;
    if (seconds < 86400) return `vor ${Math.floor(seconds / 3600)} Std`;
    return `vor ${Math.floor(seconds / 86400)} Tagen`;
  }

  if (loading) {
    return (
      <div className="bg-gray-800 rounded-lg p-4">
        <h2 className="text-lg font-bold mb-3 flex items-center gap-2">
          <span className="text-yellow-400">⚡</span> Aktivitäten
        </h2>
        <div className="animate-pulse space-y-2">
          {[1, 2, 3].map((i) => (
            <div key={i} className="bg-gray-900 rounded p-3 h-16"></div>
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className="bg-gray-800 rounded-lg p-4">
      <div className="flex items-center justify-between mb-3">
        <h2 className="text-lg font-bold flex items-center gap-2">
          <span className="text-yellow-400">⚡</span> Aktivitäten
        </h2>
      </div>

      {activities.length === 0 ? (
        <div className="text-center py-8 text-gray-500">
          <p className="text-sm">Noch keine Aktivitäten</p>
          <p className="text-xs mt-2">Events erscheinen hier</p>
        </div>
      ) : (
        <div className="space-y-2 max-h-[400px] overflow-y-auto">
          {activities.map((activity, idx) => (
            <div
              key={`${activity.id}-${idx}`}
              className="bg-gray-900 rounded p-3 flex items-center gap-3 hover:bg-gray-850 transition-colors"
            >
              <span className="text-2xl flex-shrink-0">
                {getActivityIcon(activity.type)}
              </span>
              <div className="flex-1 min-w-0">
                <p className="text-sm font-medium break-words">
                  {getActivityText(activity)}
                </p>
                {activity.message && (
                  <p className="text-xs text-gray-400 mt-1 italic">
                    "{activity.message}"
                  </p>
                )}
                <p className="text-xs text-gray-500 mt-1">
                  {getTimeAgo(activity.timestamp)}
                </p>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
