import { useState, useEffect } from "react";
import { useWebSocket } from "../../hooks/useWebSocket";
import type { Activity } from "../../features/stream/types";

const API_BASE = "/api/dashboard/stream";

// WebSocket URL - anpassen an deine Environment
const WS_URL = window.location.protocol === "https:"
  ? `wss://${window.location.host}/ws/activities`
  : `ws://${window.location.host}/ws/activities`;

export default function ActivityFeedPanel() {
  const [activities, setActivities] = useState<Activity[]>([]);
  const [loading, setLoading] = useState(true);

  // Load initial activities
  useEffect(() => {
    loadInitialActivities();
  }, []);

  async function loadInitialActivities() {
    try {
      const res = await fetch(`${API_BASE}/activities?limit=50`, {
        credentials: "include",
      });

      if (res.ok) {
        const data = await res.json();
        setActivities(data);
      }
    } catch (err) {
      console.error("Fehler beim Laden der Activities:", err);
    } finally {
      setLoading(false);
    }
  }

  // WebSocket Connection
  const { isConnected } = useWebSocket({
    url: WS_URL,
    onMessage: (activity: Activity) => {
      console.log("Neue Activity:", activity);

      // Füge neue Activity am Anfang hinzu
      setActivities((prev) => [activity, ...prev].slice(0, 100)); // Max 100 behalten
    },
    onOpen: () => {
      console.log("Activity Feed WebSocket verbunden");
    },
    onClose: () => {
      console.log("Activity Feed WebSocket getrennt");
    },
  });

  function getActivityIcon(type: Activity["type"]): string {
    switch (type) {
      case "FOLLOW":
        return "💜";
      case "SUB":
        return "⭐";
      case "RAID":
        return "🎉";
      case "CHEER":
        return "💎";
      case "HOST":
        return "📺";
      default:
        return "✨";
    }
  }

  function getActivityText(activity: Activity): string {
    switch (activity.type) {
      case "FOLLOW":
        return `${activity.displayName} folgt dir!`;
      case "SUB":
        return `${activity.displayName} hat subscribed!${
          activity.tier === "2000" ? " (Tier 2)" :
          activity.tier === "3000" ? " (Tier 3)" : ""
        }`;
      case "RAID":
        return `${activity.displayName} raided mit ${activity.viewers} Viewern!`;
      case "CHEER":
        return `${activity.displayName} cheered ${activity.bits} Bits!`;
      case "HOST":
        return `${activity.displayName} hostet deinen Stream!`;
      default:
        return `${activity.displayName} - ${activity.type}`;
    }
  }

  function getTimeAgo(timestamp: number): string {
    const seconds = Math.floor((Date.now() - timestamp) / 1000);

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

        {/* Connection Status Indicator */}
        <div className="flex items-center gap-2 text-xs">
          <div className={`w-2 h-2 rounded-full ${
            isConnected ? "bg-green-500" : "bg-red-500"
          }`}></div>
          <span className="text-gray-400">
            {isConnected ? "Live" : "Offline"}
          </span>
        </div>
      </div>

      {activities.length === 0 ? (
        <div className="text-center py-8 text-gray-500">
          <p className="text-sm">Noch keine Aktivitäten</p>
          <p className="text-xs mt-2">Events erscheinen hier live</p>
        </div>
      ) : (
        <div className="space-y-2 max-h-[400px] overflow-y-auto">
          {activities.map((activity, idx) => (
            <div
              key={`${activity.username}-${activity.timestamp}-${idx}`}
              className="bg-gray-900 rounded p-3 flex items-center gap-3 hover:bg-gray-850 transition-colors animate-fade-in"
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