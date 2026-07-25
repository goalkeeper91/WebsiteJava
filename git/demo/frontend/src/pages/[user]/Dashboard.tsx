import { useState } from "react";
import { Outlet } from "react-router-dom";
import { Menu } from "lucide-react";
import DashboardSidebar from "../../components/dashboard/DashboardSidebar";

const UserDashboard = () => {
  const [sidebarOpen, setSidebarOpen] = useState(false);

  return (
    <div className="flex flex-col md:flex-row">
      {/* Mobiler Header mit Hamburger - ab md: ausgeblendet, Sidebar ist dann permanent sichtbar.
          Sticky statt fixed, da der globale Site-Header bereits "fixed top-0 z-50" ist -
          ein zweiter fixed Header an derselben Stelle würde dessen Klickbereich blockieren. */}
      <div className="md:hidden sticky top-17 z-30 flex items-center gap-3 bg-gray-800 border-b border-gray-700 px-4 py-3">
        <button
          onClick={() => setSidebarOpen(true)}
          aria-label="Menü öffnen"
          className="text-gray-300 hover:text-white"
        >
          <Menu className="w-6 h-6" />
        </button>
        <span className="font-semibold text-white">Dashboard</span>
      </div>

      {sidebarOpen && (
        <div
          className="fixed inset-0 bg-black/50 z-40 md:hidden"
          onClick={() => setSidebarOpen(false)}
        />
      )}

      <DashboardSidebar isOpen={sidebarOpen} onClose={() => setSidebarOpen(false)} />

      <main className="flex-grow min-w-0">
        <Outlet />
      </main>
    </div>
  );
};

export default UserDashboard;
