import { Link, useLocation } from "react-router-dom";
import {
  Home,
  CreditCard,
  Zap,
  BarChart3,
  Package,
  Settings,
  Shield
} from "lucide-react";
import { useAuth } from "../../context/AuthContext";
import AdminBadge from "../admin/AdminBadge";

export default function DashboardSidebar() {
  const location = useLocation();
  const { username, isAdmin } = useAuth();

  const isActive = (path: string) => {
    if (path === "/dashboard" && location.pathname === "/dashboard") {
      return true;
    }
    return location.pathname.startsWith(path) && path !== "/dashboard";
  };

  const navItems = [
    { path: "/dashboard", icon: Home, label: "Overview" },
    { path: "/dashboard/subscription", icon: CreditCard, label: "Subscription", badge: isAdmin ? <AdminBadge /> : null },
    { path: "/dashboard/n8n", icon: Zap, label: "n8n Integration", premium: true },
    { path: "/dashboard/votes", icon: BarChart3, label: "Vote Sessions", premium: true },
    { path: "/dashboard/workflows", icon: Package, label: "Workflows", premium: true },
  ];

  return (
    <aside className="w-64 bg-gray-800 border-r border-gray-700 min-h-screen p-4">
      {/* User Info */}
      <div className="mb-6 p-3 bg-gray-900 rounded-lg">
        <div className="flex items-center gap-2 mb-1">
          <div className="w-8 h-8 rounded-full bg-purple-600 flex items-center justify-center text-sm font-bold">
            {username?.charAt(0).toUpperCase() || "U"}
          </div>
          <div className="flex-1 min-w-0">
            <p className="font-semibold truncate">{username || "User"}</p>
          </div>
        </div>
        {isAdmin && (
          <div className="mt-2">
            <AdminBadge />
          </div>
        )}
      </div>

      {/* Navigation */}
      <nav className="space-y-1">
        {navItems.map((item) => {
          const Icon = item.icon;
          const active = isActive(item.path);

          return (
            <Link
              key={item.path}
              to={item.path}
              className={`flex items-center gap-3 px-3 py-2 rounded-lg transition-colors ${
                active
                  ? "bg-gray-900 text-white"
                  : "text-gray-400 hover:bg-gray-700 hover:text-white"
              }`}
            >
              <Icon className="w-5 h-5" />
              <span className="flex-1">{item.label}</span>
              {item.badge}
              {item.premium && !isAdmin && (
                <span className="text-xs px-1.5 py-0.5 bg-purple-600 rounded">
                  PRO
                </span>
              )}
            </Link>
          );
        })}
      </nav>

      {/* Admin Panel Link */}
      {isAdmin && (
        <div className="mt-6 pt-6 border-t border-gray-700">
          <Link
            to="/admin"
            className="flex items-center gap-3 px-3 py-2 rounded-lg text-purple-400 hover:bg-purple-900/20 hover:text-purple-300 transition-colors"
          >
            <Shield className="w-5 h-5" />
            <span>Admin Panel</span>
          </Link>
        </div>
      )}

      {/* Settings */}
      <div className="mt-auto pt-6">
        <Link
          to="/dashboard/settings"
          className="flex items-center gap-3 px-3 py-2 rounded-lg text-gray-400 hover:bg-gray-700 hover:text-white transition-colors"
        >
          <Settings className="w-5 h-5" />
          <span>Settings</span>
        </Link>
      </div>
    </aside>
  );
}