import React from "react";
import { Link } from "react-router-dom";

interface NavLinksProps {
  isAuthenticated: boolean;
  isAdmin: boolean;
  username: string | null;
  handleLogout: () => void;
  isMobile?: boolean;
}

export const NavLinks: React.FC<NavLinksProps> = ({
  isAuthenticated,
  isAdmin,
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

      {/* Hauptlinks */}
      <Link to="/" className={linkClass}>Home</Link>
      <Link to="/contact" className={linkClass}>Kontakt</Link>
      <Link to="/services" className={linkClass}>Lösungen</Link>
      <Link to="/about" className={linkClass}>About</Link>
      <Link to="/allVideos" className={linkClass}>Alle Videos</Link>

      {/* Admin Link */}
      {isAuthenticated && isAdmin && (
        <Link to="/admin" className={`${linkClass} text-green-400 font-bold`}>
          Admin
        </Link>
      )}

      {isAuthenticated && (
        <Link to="/dashboard" className={`${linkClass} text-green-400 font-bold`}>
            Dashboard
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
