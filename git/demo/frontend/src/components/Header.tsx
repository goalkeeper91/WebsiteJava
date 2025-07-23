import React from 'react';
import { Link } from 'react-router-dom';

const Header: React.FC = () => {
    return (
        <header className="bg-slate-700/40 text-white w-full">
          <div className="flex items-center justify-between py-2 px-6">
            <div className="flex items-center space-x-4">
              <img src="../../../public/images/goalkeeper_logo.png" alt="Profile" className="w-15 h-15 rounded-full" />
              <h1 className="text-xl font-semibold text-goalyBlue">Goalkeeper91</h1>
            </div>
            <nav className="space-x-6">
              <Link to="/" className="text-white hover:text-blue-300">Home</Link>
              <Link to="/about" className="text-white hover:text-blue-300">About</Link>
            </nav>
          </div>
        </header>
    );
};

export default Header;
