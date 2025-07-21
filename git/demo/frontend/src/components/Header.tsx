import React from 'react';
import { Link } from 'react-router-dom';

const Header: React.FC = () => {
    return (
        <header className="bg-white shadow-md w-full">
            <div className="container mx-auto flex items-center justify-between py-4 px-6">
                <div className="flex items-center space-x-4">
                    <img
                        src="https://via.placeholder.com/40"
                        alt="Profile"
                        className="w-10 h-10 rounded-full"
                    />
                    <h1 className="text-xl font-semibold text-gray-800">Streamer Website</h1>
                </div>
                <nav className="space-x-6">
                    <Link to="/" className="text-gray-600 hover:text-blue-600">
                        Home
                    </Link>
                    <Link to="/about" className="text-gray-600 hover:text-blue-600">
                        About
                    </Link>
                </nav>
            </div>
        </header>
    );
};

export default Header;
