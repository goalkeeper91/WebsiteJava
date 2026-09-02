import React from "react";
import { Link } from "react-router-dom";
import { isBotStorefront } from "../lib/botDomain";
import { DEV_STOREFRONT_HOSTNAME, isDevStorefront } from "../lib/devDomain";

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
  // dev.goalkeeper91.de carries only the software-development-business
  // content (see lib/devDomain.ts) - its own nav, no streamer/bot-SaaS
  // links at all, same "no dead link to gated-off content" reasoning.
  const isDev = isDevStorefront();

  return (
    <div className={isMobile ? "space-y-2" : "flex space-x-6 items-center"}>

      {/* Hauptlinks */}
      {isDev ? (
        <>
          <Link to="/" className={linkClass}>Home</Link>
          <Link to="/services" className={linkClass}>Leistungen</Link>
          <Link to="/portfolio" className={linkClass}>Portfolio</Link>
          <Link to="/about" className={linkClass}>Über mich</Link>
          <Link to="/contact" className={linkClass}>Kontakt</Link>
        </>
      ) : (
        <>
          <Link to={isBot ? "/pricing" : "/"} className={linkClass}>Home</Link>
          <Link to="/pricing" className={linkClass}>Preise</Link>
          {!isBot && <Link to="/about" className={linkClass}>About</Link>}
          {!isBot && (
            <a href={`https://${DEV_STOREFRONT_HOSTNAME}`} className={linkClass}>
              Entwickler gesucht?
            </a>
          )}
        </>
      )}

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
