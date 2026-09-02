import React from 'react';
import { FaTwitch, FaYoutube, FaTiktok, FaInstagram, FaDiscord } from 'react-icons/fa';
import { Link } from 'react-router-dom';
import { isBotStorefront } from '../lib/botDomain';
import { isDevStorefront } from '../lib/devDomain';

const Footer: React.FC = () => {
  const isBot = isBotStorefront();
  const isDev = isDevStorefront();

  return (
    <footer className="bg-slate-900 text-white py-6 w-full box-border">
          <div className="mx-auto flex flex-col sm:flex-row items-center justify-between gap-4 px-6 max-w-full overflow-hidden">

            <div className="flex flex-wrap space-x-4">
              <a href="https://www.twitch.tv/goalkeeper91" target="_blank" rel="noopener noreferrer"
                className="bg-purple-600 hover:bg-purple-700 text-white p-2 rounded-full transition">
                <FaTwitch size={20} />
              </a>
              <a href="https://www.youtube.com/@goalkeeper91UNCUT" target="_blank" rel="noopener noreferrer"
                className="bg-red-600 hover:bg-red-700 text-white p-2 rounded-full transition">
                <FaYoutube size={20} />
              </a>
              <a href="https://www.tiktok.com/@g04lkeeper91" target="_blank" rel="noopener noreferrer"
                className="bg-black hover:bg-neutral-800 text-white p-2 rounded-full transition">
                <FaTiktok size={20} />
              </a>
              <a href="https://www.instagram.com/goalkeeper995" target="_blank" rel="noopener noreferrer"
                className="bg-gradient-to-tr from-yellow-400 via-pink-500 to-purple-600 hover:opacity-80 text-white p-2 rounded-full transition">
                <FaInstagram size={20} />
              </a>
              <a href="https://discord.gg/XE8sW56" target="_blank" rel="noopener noreferrer"
                className="bg-indigo-600 hover:bg-indigo-700 hover:opacity-80 text-white p-2 rounded-full transition">
                <FaDiscord size={20} />
              </a>
            </div>

            <div className="text-sm text-gray-400 text-center sm:text-right max-w-full">
              &copy; {new Date().getFullYear()} {isDev ? "Marcel Turlach" : "Goalkeeper91"}. Alle Rechte vorbehalten.
            </div>

            <div className="text-sm text-gray-400 text-center sm:text-right max-w-full">
              {/* Every one of these already sits in the header nav (see
                  Navlinks.tsx) - repeating them here was pure redundancy,
                  not extra discoverability. dev. no longer cross-links to/
                  from the main domain at all (this storefront is meant to
                  become fully independent later, see lib/devDomain.ts). */}
              {!isDev && (
                <div className="flex flex-col sm:flex-row gap-2 sm:gap-4 justify-center sm:justify-end">
                  <Link to="/pricing" className="hover:text-white transition">Preise</Link>
                  {!isBot && <Link to="/about" className="hover:text-white transition">About</Link>}
                </div>
              )}
              <div className="mt-2 flex flex-col sm:flex-row gap-2 sm:gap-4 justify-center sm:justify-end">
                <Link to="/legal/impressum" className="hover:text-white transition">Impressum</Link>
                <Link to="/legal/datenschutz" className="hover:text-white transition">Datenschutz</Link>
                {/* AGB und Widerrufsbelehrung handeln beide ausschließlich vom
                    Twitch-Bot-Abo (Paddle-Checkout, Twitch-Login, 14-Tage-
                    Widerrufsrecht für "kostenpflichtige Tarife") - auf dev.
                    gibt es weder einen Online-Vertragsschluss noch
                    Verbraucher-Kunden (B2B-Projektanfragen laufen über ein
                    individuelles Angebot, siehe DevContact.tsx), das
                    Verbraucher-Widerrufsrecht nach §312g BGB greift dort gar
                    nicht erst. Beide Seiten blieben sonst live erreichbar und
                    zeigen komplett falsche Inhalte (Bot-Kündigungsseite etc.). */}
                {!isDev && <Link to="/legal/agb" className="hover:text-white transition">AGB</Link>}
                {!isDev && <Link to="/legal/widerruf" className="hover:text-white transition">Widerruf</Link>}
                <Link to="/legal/cookies" className="hover:text-white transition">Cookies</Link>
                {/* §312k-BGB-Kündigungsbutton ist Twitch-Bot-Abo-spezifisch -
                    auf dev. gibt es kein Abo, das gekündigt werden könnte. */}
                {!isDev && <Link to="/vertrag-kuendigen" className="hover:text-white transition">Kündigen</Link>}
              </div>
            </div>
          </div>
        </footer>
  );
};

export default Footer;
