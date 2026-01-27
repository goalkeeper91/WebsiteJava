import { useState } from "react";
import CommandsDashboard from "../../components/userDashboard/CommandsDashboard";

export default function Dashboard() {
  const [activeTab, setActiveTab] = useState<"commands" | "timers">("commands");

  return (
    <div className="min-h-screen bg-gray-900 text-white p-6">
      <h1 className="text-3xl font-bold mb-6">Twitch Dashboard</h1>

      {/* Tabs */}
      <div className="flex gap-2 mb-6">
        <button
          onClick={() => setActiveTab("commands")}
          className={`px-4 py-2 rounded ${
            activeTab === "commands"
              ? "bg-green-600"
              : "bg-gray-700 hover:bg-gray-600"
          }`}
        >
          Commands
        </button>
        <button
          onClick={() => setActiveTab("timers")}
          className={`px-4 py-2 rounded ${
            activeTab === "timers"
              ? "bg-green-600"
              : "bg-gray-700 hover:bg-gray-600"
          }`}
        >
          Timer (Demo)
        </button>
      </div>

      {/* Active Panel */}
      <div>
        {activeTab === "commands" && <CommandsDashboard />}
        {activeTab === "timers" && (
          <div className="p-4 bg-gray-800 rounded">
            <p className="text-gray-300">Hier könnte später das Timer-Dashboard stehen.</p>
          </div>
        )}
      </div>
    </div>
  );
}
