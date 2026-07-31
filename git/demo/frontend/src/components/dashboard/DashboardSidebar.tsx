import { Link, useLocation } from "react-router-dom";
import type { ComponentType } from "react";
import {
  CreditCard,
  Zap,
  Radio,
  Settings,
  Shield,
  Film,
  MessageSquare,
  Hash,
  Terminal,
  Wand2,
  Vote,
  Clock,
  ShieldAlert,
  Star,
  Gift,
  Gamepad2,
  Users,
} from "lucide-react";
import { useAuth } from "../../context/AuthContext";
import { useTeam } from "../../context/TeamContext";
import AdminBadge from "../admin/AdminBadge";

interface DashboardSidebarProps {
  isOpen: boolean;
  onClose: () => void;
}

interface ToolLink {
  path: string;
  icon: ComponentType<{ className?: string }>;
  label: string;
  category: "Chat" | "Engagement" | "Erweitert";
  teamOnly?: boolean;
}

// The 9 Twitch-Chatbot tool pages, grouped exactly like the old TwitchDashboard
// sub-sidebar - now rendered directly inside the one main sidebar instead of a
// second, nested sidebar, so there's only ever one navigation surface on screen.
const TWITCH_TOOLS: ToolLink[] = [
  { path: "/dashboard/twitch/commands", icon: Terminal, label: "Chat-Commands", category: "Chat" },
  { path: "/dashboard/twitch/builtin", icon: Wand2, label: "Eingebaute Commands", category: "Chat" },
  { path: "/dashboard/twitch/votes", icon: Vote, label: "Umfragen", category: "Chat" },
  { path: "/dashboard/twitch/scheduled", icon: Clock, label: "Automatisierte Nachrichten", category: "Chat" },
  { path: "/dashboard/twitch/automod", icon: ShieldAlert, label: "Automod", category: "Engagement" },
  { path: "/dashboard/twitch/loyalty", icon: Star, label: "Loyalty", category: "Engagement" },
  { path: "/dashboard/twitch/giveaways", icon: Gift, label: "Giveaways", category: "Engagement" },
  { path: "/dashboard/twitch/cs2", icon: Gamepad2, label: "CS2 Caster Tools", category: "Erweitert" },
  { path: "/dashboard/twitch/team", icon: Users, label: "Team", category: "Erweitert", teamOnly: true },
];

const TWITCH_CATEGORY_ORDER: Array<ToolLink["category"]> = ["Chat", "Engagement", "Erweitert"];

export default function DashboardSidebar({ isOpen, onClose }: DashboardSidebarProps) {
  const location = useLocation();
  const { username, isAdmin } = useAuth();
  const { actingAsChannel } = useTeam();

  const isActive = (path: string) => location.pathname === path;
  const visibleTools = TWITCH_TOOLS.filter((t) => !t.teamOnly || !actingAsChannel);

  // Reihenfolge von oben nach unten: Twitch/Discord (Kernbereiche, taeglich
  // genutzt) vor Subathon/Clips (Feature-Bereiche), Subscription/n8n ganz
  // unten (Account-/Erweiterungs-Verwaltung, seltener aufgerufen).
  const otherItems = [
    { path: "/dashboard/subathon", icon: Radio, label: "Subathon Timer" },
    { path: "/dashboard/clips", icon: Film, label: "Clip-Automatisierung", premium: true },
    { path: "/dashboard/subscription", icon: CreditCard, label: "Subscription", badge: isAdmin ? <AdminBadge /> : null },
    { path: "/dashboard/n8n", icon: Zap, label: "n8n Integration", premium: true },
  ];

  return (
    <aside
      className={`w-64 bg-gray-800 border-r border-gray-700 min-h-screen p-4 flex flex-col overflow-y-auto
        fixed inset-y-0 left-0 z-50 transform transition-transform duration-200 ease-in-out
        md:static md:translate-x-0 md:z-auto
        ${isOpen ? "translate-x-0" : "-translate-x-full"}`}
    >
      {/* User Info + Settings */}
      <div className="mb-6 p-3 bg-gray-900 rounded-lg">
        <div className="flex items-center gap-2 mb-1">
          <div className="w-8 h-8 rounded-full bg-purple-600 flex items-center justify-center text-sm font-bold flex-shrink-0">
            {username?.charAt(0).toUpperCase() || "U"}
          </div>
          <div className="flex-1 min-w-0">
            <p className="font-semibold truncate">{username || "User"}</p>
          </div>
          <Link
            to="/dashboard/settings"
            onClick={onClose}
            title="Einstellungen"
            className={`p-1.5 rounded-lg transition-colors flex-shrink-0 ${
              isActive("/dashboard/settings") ? "bg-gray-800 text-white" : "text-gray-400 hover:bg-gray-800 hover:text-white"
            }`}
          >
            <Settings className="w-5 h-5" />
          </Link>
        </div>
        {isAdmin && (
          <div className="mt-2">
            <AdminBadge />
          </div>
        )}
      </div>

      {/* Navigation - jeder Menue-Oberpunkt bekommt seine eigene Trennlinie
          darunter, statt einer einzigen Trennlinie in der Mitte der Liste. */}
      <nav className="flex-1">
        {/* Twitch */}
        <div className="pb-3 mb-3 border-b border-gray-700 space-y-1">
          <Link
            to="/dashboard/twitch"
            onClick={onClose}
            className={`flex items-center gap-3 px-3 py-2 rounded-lg transition-colors ${
              isActive("/dashboard/twitch") ? "bg-gray-900 text-white" : "text-gray-400 hover:bg-gray-700 hover:text-white"
            }`}
          >
            <MessageSquare className="w-5 h-5" />
            <span className="flex-1">Twitch</span>
          </Link>

          <div className="pl-3 space-y-1">
            {TWITCH_CATEGORY_ORDER.map((category) => {
              const items = visibleTools.filter((t) => t.category === category);
              if (items.length === 0) return null;
              return (
                <div key={category} className="pt-1.5">
                  <p className="px-3 text-[11px] font-semibold text-gray-500 uppercase tracking-wide mb-0.5">
                    {category}
                  </p>
                  {items.map((item) => {
                    const Icon = item.icon;
                    const active = isActive(item.path);
                    return (
                      <Link
                        key={item.path}
                        to={item.path}
                        onClick={onClose}
                        className={`flex items-center gap-3 px-3 py-1.5 rounded-lg text-sm transition-colors ${
                          active ? "bg-gray-900 text-white" : "text-gray-400 hover:bg-gray-700 hover:text-white"
                        }`}
                      >
                        <Icon className="w-4 h-4 flex-shrink-0" />
                        <span className="truncate">{item.label}</span>
                      </Link>
                    );
                  })}
                </div>
              );
            })}
          </div>
        </div>

        {/* Discord */}
        <div className="pb-3 mb-3 border-b border-gray-700">
          <Link
            to="/dashboard/discord"
            onClick={onClose}
            className={`flex items-center gap-3 px-3 py-2 rounded-lg transition-colors ${
              isActive("/dashboard/discord") ? "bg-gray-900 text-white" : "text-gray-400 hover:bg-gray-700 hover:text-white"
            }`}
          >
            <Hash className="w-5 h-5" />
            <span className="flex-1">Discord</span>
          </Link>
        </div>

        {/* Subathon Timer / Clip-Automatisierung / Subscription / n8n Integration -
            jeder Punkt einzeln durch eine eigene Trennlinie abgesetzt (letzter
            Punkt ohne, da danach ohnehin ggf. das Admin-Panel folgt). */}
        {otherItems.map((item, index) => {
          const Icon = item.icon;
          const active = isActive(item.path);
          const isLast = index === otherItems.length - 1;
          return (
            <div key={item.path} className={isLast ? "" : "pb-3 mb-3 border-b border-gray-700"}>
              <Link
                to={item.path}
                onClick={onClose}
                className={`flex items-center gap-3 px-3 py-2 rounded-lg transition-colors ${
                  active ? "bg-gray-900 text-white" : "text-gray-400 hover:bg-gray-700 hover:text-white"
                }`}
              >
                <Icon className="w-5 h-5" />
                <span className="flex-1">{item.label}</span>
                {item.badge}
                {item.premium && !isAdmin && (
                  <span className="text-xs px-1.5 py-0.5 bg-purple-600 rounded">PRO</span>
                )}
              </Link>
            </div>
          );
        })}
      </nav>

      {/* Admin Panel Link */}
      {isAdmin && (
        <div className="mt-6 pt-6 border-t border-gray-700">
          <Link
            to="/admin"
            onClick={onClose}
            className="flex items-center gap-3 px-3 py-2 rounded-lg text-purple-400 hover:bg-purple-900/20 hover:text-purple-300 transition-colors"
          >
            <Shield className="w-5 h-5" />
            <span>Admin Panel</span>
          </Link>
        </div>
      )}
    </aside>
  );
}
