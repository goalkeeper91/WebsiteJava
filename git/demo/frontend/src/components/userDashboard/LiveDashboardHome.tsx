import { useRef } from "react";
import { useAuth } from "../../context/AuthContext";
import { useTeam } from "../../context/TeamContext";
import { useDashboardLayout } from "../../hooks/useDashboardLayout";
import DashboardWidget from "./DashboardWidget";
import WidgetVisibilityBar from "./WidgetVisibilityBar";
import TwitchChatEmbed from "../live/TwitchChatEmbed";
import StreamInfoPanel from "../stream/StreamInfoPanel";
import DashboardStatsPanel from "../stream/DashboardStatsPanel";
import AdBreakControl from "../stream/AdBreakControl";
import ActivityFeedPanel from "../stream/ActivityFeedPanel";

interface WidgetDef {
  id: string;
  title: string;
  width: number;
  defaultPos: { x: number; y: number };
  render: () => React.ReactNode;
}

export default function LiveDashboardHome() {
  const { username } = useAuth();
  const { actingAsChannel } = useTeam();
  const containerRef = useRef<HTMLDivElement>(null);

  const channel = actingAsChannel?.owner_login ?? username;

  const widgets: WidgetDef[] = [
    { id: "stats", title: "Live-Stats", width: 700, defaultPos: { x: 20, y: 20 }, render: () => <DashboardStatsPanel /> },
    {
      id: "chat",
      title: "Twitch Chat",
      width: 360,
      defaultPos: { x: 740, y: 20 },
      render: () => (
        <div className="h-[400px]">
          <TwitchChatEmbed channel={channel!} />
        </div>
      ),
    },
    { id: "streamInfo", title: "Stream Info", width: 340, defaultPos: { x: 20, y: 280 }, render: () => <StreamInfoPanel /> },
    { id: "adBreak", title: "Werbepause", width: 340, defaultPos: { x: 380, y: 280 }, render: () => <AdBreakControl /> },
    { id: "activity", title: "Aktivitäten", width: 360, defaultPos: { x: 740, y: 440 }, render: () => <ActivityFeedPanel /> },
  ];

  const defaults = Object.fromEntries(
    widgets.map((w) => [w.id, { x: w.defaultPos.x, y: w.defaultPos.y, width: w.width }])
  );
  const { getEntry, updatePosition, updateSize, toggleVisible, bringToFront, resetLayout } =
    useDashboardLayout(defaults);
  const allIds = widgets.map((w) => w.id);

  if (!channel) {
    return (
      <div className="p-4 sm:p-6">
        <p className="text-gray-400 text-sm">Lade...</p>
      </div>
    );
  }

  const hiddenWidgets = widgets.filter((w) => !getEntry(w.id).visible).map((w) => ({ id: w.id, title: w.title }));

  return (
    <div className="p-4 sm:p-6 space-y-4">
      <WidgetVisibilityBar
        hiddenWidgets={hiddenWidgets}
        onShow={toggleVisible}
        onReset={resetLayout}
      />

      {/* Desktop: frei anordbare Fenster (siehe Live-Dashboard-Plan - Ziehen per Maus
          ist ein Desktop-Pattern, auf Mobile gibt es stattdessen die gestapelte
          Ansicht direkt darunter, gesteuert von derselben visible-Liste). */}
      <div ref={containerRef} className="hidden sm:block relative min-h-[900px]">
        {widgets.map((widget) => {
          const entry = getEntry(widget.id);
          if (!entry.visible) return null;
          return (
            <DashboardWidget
              key={widget.id}
              title={widget.title}
              x={entry.x}
              y={entry.y}
              zIndex={entry.zIndex}
              width={entry.width}
              height={entry.height}
              containerRef={containerRef}
              onDragEnd={(x, y) => updatePosition(widget.id, x, y)}
              onResizeEnd={(width, height) => updateSize(widget.id, width, height)}
              onFocus={() => bringToFront(widget.id, allIds)}
              onHide={() => toggleVisible(widget.id)}
            >
              {widget.render()}
            </DashboardWidget>
          );
        })}
      </div>

      {/* Mobile: gestapelte Fallback-Ansicht ohne Drag */}
      <div className="sm:hidden space-y-4">
        {widgets.map((widget) => {
          const entry = getEntry(widget.id);
          if (!entry.visible) return null;
          return (
            <div key={widget.id} className="bg-gray-800 border border-gray-700 rounded-xl overflow-hidden">
              <div className="flex items-center justify-between px-3 py-2 bg-gray-900 border-b border-gray-700">
                <span className="text-sm font-semibold text-white">{widget.title}</span>
                <button
                  onClick={() => toggleVisible(widget.id)}
                  className="text-gray-500 hover:text-white transition-colors text-xs"
                >
                  Ausblenden
                </button>
              </div>
              {widget.render()}
            </div>
          );
        })}
      </div>
    </div>
  );
}
