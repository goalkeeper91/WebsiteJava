import { motion } from 'motion/react';
import { Link } from 'react-router-dom';

const Hero = () => {
    return (
        <section className="relative w-full h-full p-7  bg-cover bg-center">
            <div className="absolute inset-0 bg-black/50 z-10" />

            <div className="absolute inset-0 bg-gradient-to-br from-purple-500 via-pink-500 to-yellow-500 opacity-30 blur-3xl"></div>
            <div className="relative z-20 flex flex-col items-center justify-center text-center text-white h-full px-6">
                <motion.h1
                    className="text-4xl text-goalyBlue sm:text-6xl font-extrabold mb-4"
                    initial={{ opacity: 0, y: -30 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ duration: 0.8 }}
                    >
                        Willkommen im Streamer-Universum
                </motion.h1>
                <motion.p
                    className="text-lg text-goalyCyan sm:text-2xl mb-8 max-w-2xl"
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ duration: 0.8, delay: 0.3 }}
                    >
                        Ich bin Marcel – Tech-Enthusiast, Gamer und Live-Entertainer.
                        Tauche ein in meine Welt!
                </motion.p>

                <motion.div
                    className="flex flex-col sm:flex-row gap-4"
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    transition={{ delay: 0.6 }}
                    >
                    <Link
                        to="/about"
                        className="px-6 py-3 bg-blue-600 hover:bg-blue-700 rounded-lg text-white font-semibold transition visited:text-white active:text-white"
                        >
                            Mehr über mich
                    </Link>
                    <a
                        href="https://twitch.tv/goalkeeper91"
                        target="_blank"
                        rel="noopener noreferrer"
                        className="px-6 py-3 bg-blue-600 hover:bg-blue-700 rounded-lg text-white font-semibold transition visited:text-white active:text-white"
                        >
                            Jetzt live ansehen
                    </a>
                </motion.div>
            </div>
        </section>
    );
};

export default Hero;