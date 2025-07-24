import React from 'react';
import { Link } from 'react-router-dom';
import useTwitchLiveStatus from '../hooks/useTwitchLiveStatus';

const Header: React.FC = () => {
    const isLive = useTwitchLiveStatus();

    return (
        <header className="fixed z-50 bg-slate-700/60 text-white w-full">
          <div className="flex items-center justify-between py-2 px-6">
            <div className="flex items-center space-x-4">
              <img src="../../../public/images/goalkeeper_logo.png" alt="Profile" className="w-15 h-15 rounded-full" />
              <h1 className="text-xl font-semibold text-goalyBlue">Goalkeeper91</h1>
            </div>
            <nav className="space-x-6">
              {isLive ? (
              <>
                <Link to="/live" className="relative text-white hover:text-red-400 font-bold">
                    <span className="absolute -left-4 top-1/2 transform -translate-y-1/2">
                        <span className="block w-2 h-2 bg-red-500 rounded-full animate-ping"></span>
                        <span className="block w-2 h-2 bg-red-500 rounded-full absolute top-0 left-0"></span>
                    </span>
                        Live
                </Link>
                <Link to="/" className="text-white hover:text-blue-300">Home</Link>
                <Link to="/about" className="text-white hover:text-blue-300">About</Link>
                <Link to="/allVideos" className="text-white hover:text-blue-300">Alle Videos</Link>
              </>
              ) : (
              <>
                <Link to="/live" className="relative text-white hover:text-red-400 font-bold">
                    <span className="absolute -left-4 top-1/2 transform -translate-y-1/2">
                        <span className="block w-2 h-2 bg-grey-500 rounded-full absolute top-0 left-0"></span>
                    </span>
                    Live
                </Link>
                <Link to="/" className="text-white hover:text-blue-300">Home</Link>
                <Link to="/about" className="text-white hover:text-blue-300">About</Link>
                <Link to="/allVideos" className="text-white hover:text-blue-300">Alle Videos</Link>
              </>
              )}
            </nav>
          </div>
        </header>
    );
};

export default Header;
