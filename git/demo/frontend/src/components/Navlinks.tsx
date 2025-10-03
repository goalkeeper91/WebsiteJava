import React from "react";
import { Link } from "react-router-dom";

interface NavLinksProps {
  isLive: boolean;
  isAuthenticated: boolean;
  username: string | null;
  handleLogout: () => void;
  isMobile?: boolean;
}

export const NavLinks: React.FC<NavLinksProps> = ({
  isLive,
  isAuthenticated,
  username,
  handleLogout,
  isMobile = false,
}) => {
  const linkClass = isMobile
    ? "block hover:text-blue-500"
    : "hover:text-blue-500";

  const liveClass = isMobile
    ? "block text-white font-bold"
    : "relative text-white font-bold hover:text-red-400";

  return (
    <div className={isMobile ? "space-y-2" : "flex space-x-6 items-center"}>
      {/* Live / Offline Badge */}
      {isLive ? (
        <a
          href="https://www.twitch.tv/goalkeeper91"
          target="_blank"
          rel="noopener noreferrer"
          className={liveClass}
        >
          {isMobile ? (
            "🔴 Live"
          ) : (
            <>
              <span className="absolute -left-4 top-1/2 transform -translate-y-1/2">
                <span className="block w-2 h-2 bg-red-500 rounded-full animate-ping" />
                <span className="block w-2 h-2 bg-red-500 rounded-full absolute top-0 left-0" />
              </span>
              Live
            </>
          )}
        </a>
      ) : (
        <a
          href="https://www.youtube.com/@goalkeeper91UNCUT"
          target="_blank"
          rel="noopener noreferrer"
          className={liveClass.replace("hover:text-red-400", "hover:text-gray-400")}
        >
          {isMobile ? "⚪ Offline" : "Offline"}
        </a>
      )}

      {/* Hauptlinks */}
      <Link to="/" className={linkClass}>Home</Link>
      <Link to="/contact" className={linkClass}>Kontakt</Link>
      <Link to="/services" className={linkClass}>Lösungen</Link>
      <Link to="/about" className={linkClass}>About</Link>
      <Link to="/allVideos" className={linkClass}>Alle Videos</Link>

      {/* Admin Link */}
      {isAuthenticated && (
        <Link to="/admin" className={`${linkClass} text-green-400 font-bold`}>
          Admin
        </Link>
      )}

      {/* Login / Logout */}
      {username ? (
        <button onClick={handleLogout} className={`${linkClass} text-red-300`}>
          Logout ({username})
        </button>
      ) : (
        <a href="/auth/twitch" className={`${linkClass} text-green-300`}>
          Login
        </a>
      )}
    </div>
  );
};
