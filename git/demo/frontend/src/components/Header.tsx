import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { FaBars, FaTimes } from 'react-icons/fa';
import useTwitchLiveStatus from '../hooks/useTwitchLiveStatus';
import { useAuth } from '../context/AuthContext';
import { NavLinks } from './Navlinks';

const Header: React.FC = () => {
    const isLive = useTwitchLiveStatus();
    const [isMenuOpen, setIsMenuOpen] = useState(false);
    const { isAdmin, isAuthenticated, username, logout } = useAuth();

    const handleLogout = async () => {
        await logout();
    };

    return (
            <header className="fixed z-50 w-full bg-slate-700/60 text-white">
                <div className="flex items-center justify-between py-2 px-6">
                    <Link to="/" className="flex items-center space-x-4">
                        <img
                            src="/images/goalkeeper_logo.png"
                            alt="Profile"
                            className="w-12 h-12 rounded-full"
                        />
                        <h1 className="hidden md:block text-xl md:text-l font-semibold text-goalyBlue">
                            Goalkeeper91
                        </h1>
                    </Link>

                    {/* Desktop Navigation */}
                    <nav className="hidden md:flex space-x-6 items-center">
                        <NavLinks
                            isAdmin={isAdmin}
                            isAuthenticated={isAuthenticated}
                            username={username}
                            handleLogout={handleLogout}
                        />
                    </nav>

                    {/* Mobile Menu Button */}
                    <div className="md:hidden">
                        <button onClick={() => setIsMenuOpen(!isMenuOpen)}>
                            {isMenuOpen ? <FaTimes size={24} /> : <FaBars size={24} />}
                        </button>
                    </div>
                </div>

                {/* Mobile Navigation */}
                {isMenuOpen && (
                    <div className="md:hidden bg-slate-800/90 px-6 pb-4 space-y-4">
                        <NavLinks
                            isAdmin={isAdmin}
                            isAuthenticated={isAuthenticated}
                            username={username}
                            handleLogout={handleLogout}
                            isMobile
                        />
                    </div>
                )}
            </header>
        );
};

export default Header;
