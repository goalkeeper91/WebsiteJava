import React from 'react';
import { FaTwitch, FaYoutube, FaTiktok, FaInstagram, FaDiscord } from 'react-icons/fa';

const Footer: React.FC = () => {
  return (
    <footer className="bg-slate-900 text-white py-6 px-4">
          <div className="max-w-6xl mx-auto flex flex-col sm:flex-row items-center justify-between gap-4">

            <div className="flex space-x-4">
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

            <div className="text-sm text-gray-400 text-center sm:text-right">
              &copy; {new Date().getFullYear()} Goalkeeper91. Alle Rechte vorbehalten.
            </div>

            <div className="text-sm text-gray-400 text-center sm:text-right">
              <div className="mt-2 flex flex-col sm:flex-row gap-2 sm:gap-4 justify-center sm:justify-end">
                <a href="/legal/impressum" className="hover:text-white transition">Impressum</a>
                <a href="/legal/datenschutz" className="hover:text-white transition">Datenschutz</a>
                <a href="/legal/agb" className="hover:text-white transition">AGB</a>
              </div>
            </div>


          </div>
        </footer>
  );
};

export default Footer;
