import { useState } from "react";
import TwitchBotStatus from "./TwitchBotStatus";
import CommandsDashboard from "./CommandsDashboard";
import BuiltinCommandsInfo from "./BuiltinCommandsInfo";
import VoteSessionManager from "../VoteSessionManager";

type View = "commands" | "builtin" | "votes";

export default function TwitchDashboard() {
  const [currentView, setCurrentView] = useState<View>("commands");

  return (
    <div className="p-6 space-y-6">
      <TwitchBotStatus />

      <div className="flex gap-2 bg-gray-800 p-2 rounded-lg">
        <button
          onClick={() => setCurrentView("commands")}
          className={`px-4 py-2 rounded-lg transition-colors font-medium ${
            currentView === "commands"
              ? "bg-indigo-600 text-white"
              : "bg-gray-700 text-gray-300 hover:bg-gray-600"
          }`}
        >
          Chat-Commands
        </button>
        <button
          onClick={() => setCurrentView("builtin")}
          className={`px-4 py-2 rounded-lg transition-colors font-medium ${
            currentView === "builtin"
              ? "bg-indigo-600 text-white"
              : "bg-gray-700 text-gray-300 hover:bg-gray-600"
          }`}
        >
          Eingebaute Commands
        </button>
        <button
          onClick={() => setCurrentView("votes")}
          className={`px-4 py-2 rounded-lg transition-colors font-medium ${
            currentView === "votes"
              ? "bg-indigo-600 text-white"
              : "bg-gray-700 text-gray-300 hover:bg-gray-600"
          }`}
        >
          Umfragen
        </button>
      </div>

      <div>
        {currentView === "commands" && <CommandsDashboard />}
        {currentView === "builtin" && <BuiltinCommandsInfo />}
        {currentView === "votes" && <VoteSessionManager />}
      </div>
    </div>
  );
}
