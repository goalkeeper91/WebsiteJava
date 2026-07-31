import { useState } from "react";
import type { ComponentType } from "react";
import { Home, Terminal, Wand2, Vote, Clock, ShieldAlert, Star, Gift, Gamepad2, Users, Menu, X } from "lucide-react";
import TwitchBotStatus from "./TwitchBotStatus";
import LiveDashboardHome from "./LiveDashboardHome";
import CommandsDashboard from "./CommandsDashboard";
import BuiltinCommandsInfo from "./BuiltinCommandsInfo";
import VoteSessionManager from "../VoteSessionManager";
import ScheduledMessagesDashboard from "./ScheduledMessagesDashboard";
import AutomodDashboard from "./AutomodDashboard";
import LoyaltyDashboard from "./LoyaltyDashboard";
import GiveawaysDashboard from "./GiveawaysDashboard";
import CS2CasterDashboard from "./CS2CasterDashboard";
import TeamDashboard from "./TeamDashboard";
import { useTeam } from "../../context/TeamContext";

type View =
  | "overview"
  | "commands"
  | "builtin"
  | "votes"
  | "scheduled"
  | "automod"
  | "loyalty"
  | "giveaways"
  | "cs2"
  | "team";

interface SectionDef {
  id: View;
  label: string;
  icon: ComponentType<{ className?: string }>;
  category: "Chat" | "Engagement" | "Erweitert" | null;
}

// "overview" (category null) is always pinned above the grouped categories -
// it renders full-width without the sidebar taking up space (see below), so
// its widget canvas keeps the full width it was designed for.
const SECTIONS: SectionDef[] = [
  { id: "overview", label: "Übersicht", icon: Home, category: null },
  { id: "commands", label: "Chat-Commands", icon: Terminal, category: "Chat" },
  { id: "builtin", label: "Eingebaute Commands", icon: Wand2, category: "Chat" },
  { id: "votes", label: "Umfragen", icon: Vote, category: "Chat" },
  { id: "scheduled", label: "Automatisierte Nachrichten", icon: Clock, category: "Chat" },
  { id: "automod", label: "Automod", icon: ShieldAlert, category: "Engagement" },
  { id: "loyalty", label: "Loyalty", icon: Star, category: "Engagement" },
  { id: "giveaways", label: "Giveaways", icon: Gift, category: "Engagement" },
  { id: "cs2", label: "CS2 Caster Tools", icon: Gamepad2, category: "Erweitert" },
  { id: "team", label: "Team", icon: Users, category: "Erweitert" },
];

const CATEGORY_ORDER: Array<"Chat" | "Engagement" | "Erweitert"> = ["Chat", "Engagement", "Erweitert"];

export default function TwitchDashboard() {
  const [currentView, setCurrentView] = useState<View>("overview");
  const [subNavOpen, setSubNavOpen] = useState(false);
  const { actingAsChannel } = useTeam();

  const visibleSections = SECTIONS.filter((s) => s.id !== "team" || !actingAsChannel);
  const activeSection = visibleSections.find((s) => s.id === currentView) ?? visibleSections[0];

  function selectSection(id: View) {
    setCurrentView(id);
    setSubNavOpen(false);
  }

  return (
    <div className="p-4 sm:p-6 space-y-4">
      <TwitchBotStatus />

      {/* Bereich-Umschalter: auf Mobile immer sichtbar (öffnet die Sidebar als
          Drawer); ab md: nur bei den Werkzeug-Ansichten nötig, da die
          Übersicht dort schon ihre eigene, feste Sidebar-Spalte hat. */}
      <button
        onClick={() => setSubNavOpen(true)}
        className={`flex items-center gap-2 px-4 py-2.5 bg-gray-800 border border-gray-700 rounded-lg hover:bg-gray-700 transition-colors text-sm font-medium w-full sm:w-auto ${
          currentView === "overview" ? "" : "md:hidden"
        }`}
      >
        <Menu className="w-4 h-4 text-gray-400 flex-shrink-0" />
        <span>{activeSection.label}</span>
      </button>

      {currentView === "overview" ? (
        <LiveDashboardHome />
      ) : (
        <div className="flex items-start gap-6">
          <TwitchSubSidebar
            visibleSections={visibleSections}
            currentView={currentView}
            onSelect={selectSection}
            open={subNavOpen}
            onClose={() => setSubNavOpen(false)}
            pinned
          />

          <div className="flex-1 min-w-0">
            {currentView === "commands" && <CommandsDashboard />}
            {currentView === "builtin" && <BuiltinCommandsInfo />}
            {currentView === "votes" && <VoteSessionManager />}
            {currentView === "scheduled" && <ScheduledMessagesDashboard />}
            {currentView === "automod" && <AutomodDashboard />}
            {currentView === "loyalty" && <LoyaltyDashboard />}
            {currentView === "giveaways" && <GiveawaysDashboard />}
            {currentView === "cs2" && <CS2CasterDashboard />}
            {currentView === "team" && !actingAsChannel && <TeamDashboard />}
          </div>
        </div>
      )}

      {/* Auf der Übersicht gibt es keine feste Sidebar-Spalte (siehe oben) -
          hier öffnet der Umschalter-Button dieselbe Sidebar rein als Overlay. */}
      {currentView === "overview" && (
        <TwitchSubSidebar
          visibleSections={visibleSections}
          currentView={currentView}
          onSelect={selectSection}
          open={subNavOpen}
          onClose={() => setSubNavOpen(false)}
        />
      )}
    </div>
  );
}

function TwitchSubSidebar({
  visibleSections,
  currentView,
  onSelect,
  open,
  onClose,
  pinned,
}: {
  visibleSections: SectionDef[];
  currentView: View;
  onSelect: (id: View) => void;
  open: boolean;
  onClose: () => void;
  pinned?: boolean;
}) {
  return (
    <>
      {open && <div className="fixed inset-0 bg-black/50 z-40 md:hidden" onClick={onClose} />}

      <aside
        className={`w-64 flex-shrink-0 bg-gray-800 border border-gray-700 rounded-xl p-3 space-y-4
          fixed inset-y-4 left-4 z-50 overflow-y-auto transition-transform duration-200
          ${pinned ? "md:static md:z-auto md:translate-x-0" : ""}
          ${open ? "translate-x-0" : "-translate-x-[120%]"}`}
      >
        <div className={`flex items-center justify-between px-1 ${pinned ? "md:hidden" : ""}`}>
          <span className="text-sm font-semibold text-gray-400">Bereiche</span>
          <button onClick={onClose}>
            <X className="w-5 h-5 text-gray-400" />
          </button>
        </div>

        <button
          onClick={() => onSelect("overview")}
          className={`w-full flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
            currentView === "overview" ? "bg-indigo-600 text-white" : "text-gray-300 hover:bg-gray-700"
          }`}
        >
          <Home className="w-4 h-4 flex-shrink-0" />
          Übersicht
        </button>

        {CATEGORY_ORDER.map((category) => {
          const items = visibleSections.filter((s) => s.category === category);
          if (items.length === 0) return null;
          return (
            <div key={category} className="pt-3 border-t border-gray-700">
              <p className="px-3 text-xs font-semibold text-gray-500 uppercase tracking-wide mb-1">
                {category}
              </p>
              <div className="space-y-1">
                {items.map((s) => {
                  const Icon = s.icon;
                  return (
                    <button
                      key={s.id}
                      onClick={() => onSelect(s.id)}
                      className={`w-full flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                        currentView === s.id ? "bg-indigo-600 text-white" : "text-gray-300 hover:bg-gray-700"
                      }`}
                    >
                      <Icon className="w-4 h-4 flex-shrink-0" />
                      <span className="truncate">{s.label}</span>
                    </button>
                  );
                })}
              </div>
            </div>
          );
        })}
      </aside>
    </>
  );
}
