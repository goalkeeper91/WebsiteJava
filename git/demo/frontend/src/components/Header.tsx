import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { FaBars, FaTimes } from 'react-icons/fa';
import useTwitchLiveStatus from '../hooks/useTwitchLiveStatus';

const Header: React.FC = () => {
    const [liveChecked, setLiveChecked] = useState(false);
    const isLive = useTwitchLiveStatus();
    const [isMenuOpen, setIsMenuOpen] = useState(false);

    useEffect(() => {
      if (isLive !== null) {
        setLiveChecked(true);
      }
    }, [isLive]);

    return (
        <header className="fixed z-50 w-full bg-slate-700/60 text-white">
            <div className="flex items-center justify-between py-2 px-6">
                <div className="flex items-center space-x-4">
                    <img
                        src="/images/goalkeeper_logo.png"
                        alt="Profile"
                        className="w-12 h-12 rounded-full"
                    />
                    <h1 className="text-xl font-semibold text-goalyBlue">Goalkeeper91</h1>
                </div>

                <nav className="hidden md:flex space-x-6 items-center">
                    {liveChecked && isLive ? (
                        <>
                        <a
                            href="https://www.twitch.tv/goalkeeper91"
                            target="_blank"
                            rel="noopener noreferrer"
                            className="relative text-white hover:text-red-400 font-bold"
                        >
                            <span className="absolute -left-4 top-1/2 transform -translate-y-1/2">
                            <span className="block w-2 h-2 bg-red-500 rounded-full animate-ping" />
                            <span className="block w-2 h-2 bg-red-500 rounded-full absolute top-0 left-0" />
                            </span>
                            Live
                        </a>
                        <Link to="/" className="hover:text-blue-300">Home</Link>
                        <Link to="/about" className="hover:text-blue-300">About</Link>
                        <Link to="/allVideos" className="hover:text-blue-300">Alle Videos</Link>
                        </>
                    ) : (
                        <>
                        <a
                            href="https://www.youtube.com/@goalkeeper91UNCUT"
                            target="_blank"
                            rel="noopener noreferrer"
                            className="relative text-white hover:text-red-400 font-bold"
                        >
                            <span className="absolute -left-4 top-1/2 transform -translate-y-1/2">
                                <span className="block w-2 h-2 bg-gray-500 rounded-full absolute top-0 left-0" />
                            </span>
                            Offline
                        </a>
                        <Link to="/" className="hover:text-blue-300">Home</Link>
                        <Link to="/about" className="hover:text-blue-300">About</Link>
                        <Link to="/allVideos" className="hover:text-blue-300">Alle Videos</Link>
                        </>
                    )}
                </nav>

                <div className="md:hidden">
                    <button onClick={() => setIsMenuOpen(!isMenuOpen)}>
                        {isMenuOpen ? <FaTimes size={24} /> : <FaBars size={24} />}
                    </button>
                </div>
            </div>

              {/* Mobile Menu */}
              {isMenuOpen && (
                <div className="md:hidden bg-slate-800/90 px-6 pb-4 space-y-4">
                  {liveChecked && isLive ? (
                    <>
                      <a
                        href="https://www.twitch.tv/goalkeeper91"
                        target="_blank"
                        rel="noopener noreferrer"
                        className="block text-white hover:text-red-400 font-bold"
                      >
                        🔴 Live
                      </a>
                      <Link to="/" className="block hover:text-blue-300">Home</Link>
                      <Link to="/about" className="block hover:text-blue-300">About</Link>
                      <Link to="/allVideos" className="block hover:text-blue-300">Alle Videos</Link>
                    </>
                  ) : (
                    <>
                      <a
                        href="https://www.youtube.com/@goalkeeper91UNCUT"
                        target="_blank"
                        rel="noopener noreferrer"
                        className="block text-white hover:text-red-400 font-bold"
                      >
                        ⚪ Offline
                      </a>
                      <Link to="/" className="block hover:text-blue-300">Home</Link>
                      <Link to="/about" className="block hover:text-blue-300">About</Link>
                      <Link to="/allVideos" className="block hover:text-blue-300">Alle Videos</Link>
                    </>
                  )}
              </div>
            )}
        </header>
    );
};

export default Header;
