import { motion } from 'motion/react';
import { Link } from 'react-router-dom';

const Hero = () => {
    return (
        <section className="relative w-full min-h-screen py-7 bg-cover bg-center bg-slate-950 flex items-center justify-center">
            <div className="relative z-20 text-center text-white w-full px-4 sm:px-6 md:px-8">

                <motion.h1
                    className="text-2xl sm:text-4xl md:text-5xl lg:text-6xl font-extrabold mb-6 leading-snug hyphens-auto"
                    initial={{ opacity: 0, y: -30 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ duration: 0.8 }}
                >
                    Twitch Bot SaaS & Live Entertainment
                </motion.h1>

                <motion.p
                    className="text-base sm:text-lg md:text-xl lg:text-2xl mb-10 text-white hyphens-auto"
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ duration: 0.8, delay: 0.3 }}
                >
                    Ich bin Marcel – Gamer, Streamer und Entwickler meines eigenen{' '}
                    <span className="text-indigo-400 font-semibold">Twitch Bots</span> als SaaS (Automod,
                    Loyalty-Punkte, Clip-Automatisierung). Schau beim Livestream vorbei oder werde Teil der{' '}
                    <span className="text-indigo-400 font-semibold">Community</span>.
                </motion.p>

                <motion.div
                    className="flex flex-wrap justify-center gap-4"
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    transition={{ delay: 0.6 }}
                >
                    <Link
                        to="/pricing"
                        className="w-full sm:w-auto px-6 py-3 bg-goalyBlue hover:bg-goalyCyan rounded-lg text-white font-semibold transition text-center"
                    >
                        Twitch Bot entdecken
                    </Link>
                    <a
                        href="https://twitch.tv/goalkeeper91"
                        target="_blank"
                        rel="noopener noreferrer"
                        className="w-full sm:w-auto px-6 py-3 border-2 border-goalyBlue hover:bg-goalyCyan rounded-lg text-white font-semibold transition text-center"
                    >
                        Zum Livestream
                    </a>
                </motion.div>
            </div>
        </section>
    );
};

export default Hero;
