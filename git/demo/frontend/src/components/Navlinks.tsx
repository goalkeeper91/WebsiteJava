import React from "react";
import { Link } from "react-router-dom";
import { isBotStorefront } from "../lib/botDomain";

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

  // bot.goalkeeper91.de carries only the Twitch Bot SaaS product for Paddle
  // review purposes - see lib/botDomain.ts. The freelance-services nav
  // links are dropped here too, not just the routes themselves, so there's
  // no dead link pointing at a page that immediately redirects away.
  const isBot = isBotStorefront();

  return (
    <div className={isMobile ? "space-y-2" : "flex space-x-6 items-center"}>

      {/* Hauptlinks */}
      <Link to={isBot ? "/pricing" : "/"} className={linkClass}>Home</Link>
      {!isBot && <Link to="/services" className={linkClass}>Lösungen</Link>}
      <Link to="/pricing" className={linkClass}>Preise</Link>
      {!isBot && <Link to="/about" className={linkClass}>About</Link>}
      {!isBot && <Link to="/contact" className={linkClass}>Kontakt</Link>}

      {/* Admin Link */}
      {isAdmin && (
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
        <a href="/auth/login" className={`${linkClass} text-green-300`}>
          Login
        </a>
      )}
    </div>
  );
};
