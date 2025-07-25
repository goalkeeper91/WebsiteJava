import React from 'react';
import { FaTwitch, FaYoutube, FaTiktok, FaInstagram, FaDiscord } from 'react-icons/fa';

const FollowMe = () => {
    return (
        <section className="text-center pt-10">
            <p className="text-gray-400 mb-2">Folge mir für mehr Highlights und Live-Streams:</p>
            <div className="flex flex-wrap justify-center gap-4">
                <a
                    href="https://www.twitch.tv/goalkeeper91"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="flex items-center gap-2 bg-[#9146FF] hover:bg-[#772ce8] text-white font-semibold px-4 py-2 rounded-full transition-colors duration-200"
                >
                    <FaTwitch />
                    Twitch
                </a>

                <a
                    href="https://www.youtube.com/@Goalkeeper91"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="flex items-center gap-2 bg-[#FF0000] hover:bg-[#cc0000] text-white font-semibold px-4 py-2 rounded-full transition-colors duration-200"
                >
                    <FaYoutube />
                    YouTube
                </a>

                <a
                    href="https://www.tiktok.com/@g04lkeeper91"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="flex items-center gap-2 bg-[#000000] hover:bg-[#1c1c1c] text-white font-semibold px-4 py-2 rounded-full transition-colors duration-200"
                >
                    <FaTiktok />
                    TikTok
                </a>

                <a
                    href="https://www.instagram.com/goalkeeper995"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="flex items-center gap-2 bg-gradient-to-r from-pink-500 via-red-500 to-yellow-500 hover:opacity-90 text-white font-semibold px-4 py-2 rounded-full transition duration-200"
                >
                    <FaInstagram />
                    Instagram
                </a>

                <a
                    href="https://discord.gg/XE8sW56"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="inline-flex items-center gap-3 px-6 py-3 bg-indigo-600 hover:bg-indigo-700 text-white font-semibold rounded-full shadow-md transition"
                >
                    <FaDiscord size={24} />
                    Jetzt beitreten
                </a>
            </div>
        </section>
    );
};

export default FollowMe;