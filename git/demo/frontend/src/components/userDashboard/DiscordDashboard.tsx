import React, { useState } from "react";
import DiscordConnect from "../socials/DiscordConnect";
import DiscordGuildList from "../socials/DiscordGuildList";
import DiscordGuildSettings from "../socials/DiscordGuildSettings";
import JoinToCreateManager from "../socials/JoinToCreateManager";

type View = "guilds" | "settings" | "join-to-create";

const DiscordDashboard: React.FC = () => {
  const [currentView, setCurrentView] = useState<View>("guilds");
  const [selectedGuildId, setSelectedGuildId] = useState<number | null>(null);

  const handleSelectGuild = (guildId: number) => {
    setSelectedGuildId(guildId);
    setCurrentView("settings");
  };

  const handleBackToGuilds = () => {
    setCurrentView("guilds");
    setSelectedGuildId(null);
  };

  return (
    <div className="p-6 space-y-6">
      {/* Discord Connection Status */}
      <DiscordConnect />

      {/* Navigation Tabs */}
      {selectedGuildId && (
        <div className="flex gap-2 bg-gray-800 p-2 rounded-lg">
          <button
            onClick={() => setCurrentView("settings")}
            className={`px-4 py-2 rounded-lg transition-colors font-medium ${
              currentView === "settings"
                ? "bg-indigo-600 text-white"
                : "bg-gray-700 text-gray-300 hover:bg-gray-600"
            }`}
          >
            Settings
          </button>
          <button
            onClick={() => setCurrentView("join-to-create")}
            className={`px-4 py-2 rounded-lg transition-colors font-medium ${
              currentView === "join-to-create"
                ? "bg-indigo-600 text-white"
                : "bg-gray-700 text-gray-300 hover:bg-gray-600"
            }`}
          >
            Join-to-Create
          </button>
        </div>
      )}

      {/* Main Content */}
      <div>
        {currentView === "guilds" && (
          <DiscordGuildList onSelectGuild={handleSelectGuild} />
        )}

        {currentView === "settings" && selectedGuildId && (
          <DiscordGuildSettings
            guildId={selectedGuildId}
            onBack={handleBackToGuilds}
          />
        )}

        {currentView === "join-to-create" && selectedGuildId && (
          <div className="p-6 bg-gray-800 rounded-xl border border-gray-700">
            <JoinToCreateManager guildId={selectedGuildId} />
          </div>
        )}
      </div>
    </div>
  );
};

export default DiscordDashboard;