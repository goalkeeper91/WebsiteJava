import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { FaBars, FaTimes } from 'react-icons/fa';
import { useAuth } from '../context/AuthContext';
import { NavLinks } from './Navlinks';
import { isBotStorefront } from '../lib/botDomain';
import { isDevStorefront } from '../lib/devDomain';

const Header: React.FC = () => {
    const [isMenuOpen, setIsMenuOpen] = useState(false);
    const { isAdmin, isAuthenticated, username, logout } = useAuth();

    const handleLogout = async () => {
        await logout();
    };

    return (
            <header className="fixed top-0 left-0 z-50 w-full bg-slate-700/60 text-white">
                <div className="flex items-center justify-between py-2 px-6">
                    <Link to={isBotStorefront() ? "/pricing" : "/"} className="flex items-center space-x-4">
                        <img
                            src="/images/goalkeeper_logo.png"
                            alt={isDevStorefront() ? "Marcel Turlach Logo" : "Goalkeeper91 Logo"}
                            width={48}
                            height={48}
                            className="w-12 h-12 rounded-full"
                        />
                        {/* dev.goalkeeper91.de is meant to eventually move to its own
                            independent domain (see lib/devDomain.ts) - showing the
                            real name here instead of the gamer tag from day one
                            avoids re-branding this specific piece of copy later. */}
                        <span className="hidden md:block text-xl md:text-l font-semibold text-goalyBlue">
                            {isDevStorefront() ? "Marcel Turlach" : "Goalkeeper91"}
                        </span>
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
