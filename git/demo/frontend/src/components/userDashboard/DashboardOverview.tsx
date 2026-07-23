import React, { useState } from "react";
import { useSearchParams } from "react-router-dom";
import { MessageSquare, Hash, Radio } from "lucide-react";
import CommandsDashboard from "./CommandsDashboard";
import DiscordDashboard from "./DiscordDashboard";
import SubathonPage from "./SubathonPage";

type DashboardTab = "twitch" | "discord" | "subathon";

const DashboardOverview: React.FC = () => {
  const [searchParams] = useSearchParams();
  const [activeTab, setActiveTab] = useState<DashboardTab>(
    searchParams.get("tab") === "discord"
      ? "discord"
      : searchParams.get("tab") === "subathon"
      ? "subathon"
      : "twitch"
  );

  return (
    <div className="min-h-screen bg-gray-900 text-white">
      {/* Header with Tab Navigation */}
      <div className="bg-gray-800 border-b border-gray-700">
        <div className="max-w-7xl mx-auto px-6 py-4">
          <div className="flex items-center justify-between mb-4">
            <div>
              <h1 className="text-2xl font-bold">Dashboard</h1>
              <p className="text-sm text-gray-400 mt-1">
                Manage your Twitch and Discord bot settings
              </p>
            </div>
          </div>

          {/* Tab Buttons */}
          <div className="flex gap-2">
            <button
              onClick={() => setActiveTab("twitch")}
              className={`flex items-center gap-2 px-6 py-3 rounded-t-lg font-medium transition-all ${
                activeTab === "twitch"
                  ? "bg-gray-900 text-white border-t-2 border-purple-500"
                  : "bg-gray-700 text-gray-400 hover:bg-gray-600 hover:text-gray-200"
              }`}
            >
              <MessageSquare className="w-5 h-5" />
              <span>Twitch Bot</span>
            </button>

            <button
              onClick={() => setActiveTab("discord")}
              className={`flex items-center gap-2 px-6 py-3 rounded-t-lg font-medium transition-all ${
                activeTab === "discord"
                  ? "bg-gray-900 text-white border-t-2 border-indigo-500"
                  : "bg-gray-700 text-gray-400 hover:bg-gray-600 hover:text-gray-200"
              }`}
            >
              <Hash className="w-5 h-5" />
              <span>Discord Bot</span>
            </button>

            <button
              onClick={() => setActiveTab("subathon")}
              className={`flex items-center gap-2 px-6 py-3 rounded-t-lg font-medium transition-all ${
                activeTab === "subathon"
                  ? "bg-gray-900 text-white border-t-2 border-purple-500"
                  : "bg-gray-700 text-gray-400 hover:bg-gray-600 hover:text-gray-200"
              }`}
            >
              <Radio className="w-5 h-5" />
              <span>Subathon Timer</span>
            </button>
          </div>
        </div>
      </div>

      {/* Dashboard Content */}
      <div className="max-w-7xl mx-auto">
        {activeTab === "twitch" && <CommandsDashboard />}
        {activeTab === "discord" && <DiscordDashboard />}
        {activeTab === "subathon" && <SubathonPage />}
      </div>
    </div>
  );
};

export default DashboardOverview;