import { Outlet } from "react-router-dom";
import TwitchBotStatus from "./TwitchBotStatus";

// Thin layout for every /dashboard/twitch/* route - the section switching
// itself now lives in the main DashboardSidebar (real routes/links), so this
// only needs to keep the always-visible bot status banner above whichever
// Twitch tool page is active.
export default function TwitchLayout() {
  return (
    <div className="min-h-screen bg-gray-900 text-white p-4 sm:p-6 space-y-4">
      <TwitchBotStatus />
      <Outlet />
    </div>
  );
}
