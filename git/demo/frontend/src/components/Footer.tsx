import React from 'react';
import { FaTwitch, FaYoutube, FaTiktok, FaInstagram, FaDiscord } from 'react-icons/fa';
import { Link } from 'react-router-dom';
import { isBotStorefront } from '../lib/botDomain';
import { DEV_STOREFRONT_HOSTNAME, isDevStorefront } from '../lib/devDomain';

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
              &copy; {new Date().getFullYear()} Goalkeeper91. Alle Rechte vorbehalten.
            </div>

            <div className="text-sm text-gray-400 text-center sm:text-right max-w-full">
              <div className="flex flex-col sm:flex-row gap-2 sm:gap-4 justify-center sm:justify-end">
                {isDev ? (
                  <>
                    <Link to="/services" className="hover:text-white transition">Leistungen</Link>
                    <Link to="/portfolio" className="hover:text-white transition">Portfolio</Link>
                    <Link to="/about" className="hover:text-white transition">Über mich</Link>
                    <Link to="/contact" className="hover:text-white transition">Kontakt</Link>
                  </>
                ) : (
                  <>
                    <Link to="/pricing" className="hover:text-white transition">Preise</Link>
                    {!isBot && <Link to="/about" className="hover:text-white transition">About</Link>}
                    {!isBot && (
                      <a href={`https://${DEV_STOREFRONT_HOSTNAME}`} className="hover:text-white transition">
                        Entwickler gesucht?
                      </a>
                    )}
                  </>
                )}
              </div>
              <div className="mt-2 flex flex-col sm:flex-row gap-2 sm:gap-4 justify-center sm:justify-end">
                <Link to="/legal/impressum" className="hover:text-white transition">Impressum</Link>
                <Link to="/legal/datenschutz" className="hover:text-white transition">Datenschutz</Link>
                <Link to="/legal/agb" className="hover:text-white transition">AGB</Link>
                <Link to="/legal/widerruf" className="hover:text-white transition">Widerruf</Link>
                <Link to="/legal/cookies" className="hover:text-white transition">Cookies</Link>
                <Link to="/vertrag-kuendigen" className="hover:text-white transition">Kündigen</Link>
              </div>
            </div>
          </div>
        </footer>
  );
};

export default Footer;
