import { useAuth } from "../../context/AuthContext";
import { useTeam } from "../../context/TeamContext";
import TwitchEmbed from "../live/TwitchEmbed";
import TwitchChatEmbed from "../live/TwitchChatEmbed";
import StreamInfoPanel from "../stream/StreamInfoPanel";
import DashboardStatsPanel from "../stream/DashboardStatsPanel";
import AdBreakControl from "../stream/AdBreakControl";
import ActivityFeedPanel from "../stream/ActivityFeedPanel";

export default function LiveDashboardHome() {
  const { username } = useAuth();
  const { actingAsChannel } = useTeam();

  const channel = actingAsChannel?.owner_login ?? username;

  if (!channel) {
    return (
      <div className="p-4 sm:p-6">
        <p className="text-gray-400 text-sm">Lade...</p>
      </div>
    );
  }

  return (
    <div className="p-4 sm:p-6 space-y-6">
      <DashboardStatsPanel />

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
        <div className="lg:col-span-2">
          <TwitchEmbed channel={channel} />
        </div>
        <div className="h-[300px] lg:h-auto">
          <TwitchChatEmbed channel={channel} />
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <StreamInfoPanel />
        <AdBreakControl />
      </div>

      <ActivityFeedPanel />
    </div>
  );
}
