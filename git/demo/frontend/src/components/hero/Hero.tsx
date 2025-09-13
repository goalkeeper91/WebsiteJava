import { motion } from 'motion/react';
import { Link } from 'react-router-dom';

const Hero = () => {
    return (
        <section className="relative w-full h-screen p-7 bg-cover bg-center bg-slate-950">

            <div className="relative z-20 flex flex-col items-center justify-center text-center text-white h-full px-6">

                <motion.h1
                    className="text-4xl sm:text-6xl font-extrabold mb-6 leading-tight"
                    initial={{ opacity: 0, y: -30 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ duration: 0.8 }}
                >
                    Softwareentwicklung & Live Entertainment
                </motion.h1>

                <motion.p
                    className="text-lg sm:text-2xl mb-10 max-w-3xl text-goalyCyan"
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ duration: 0.8, delay: 0.3 }}
                >
                    Ich bin Marcel – Full-Stack Entwickler mit Leidenschaft für Gaming & Streaming.
                    Ob <span className="font-semibold text-goalyBlue">maßgeschneiderte Softwarelösungen</span> oder
                    <span className="font-semibold text-goalyBlue"> Community-Content</span> – hier findest du beides.
                </motion.p>

                <motion.div
                    className="flex flex-col sm:flex-row gap-4"
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    transition={{ delay: 0.6 }}
                >
                    <Link
                        to="/services"
                        className="px-8 py-3 bg-goalyBlue hover:bg-goalyCyan rounded-lg text-white font-semibold transition"
                    >
                        Entwickler anheuern
                    </Link>
                    <a
                        href="https://twitch.tv/goalkeeper91"
                        target="_blank"
                        rel="noopener noreferrer"
                        className="px-8 py-3 border-2 border-goalyBlue hover:bg-goalyCyan rounded-lg text-white font-semibold transition"
                    >
                        Zum Livestream
                    </a>
                </motion.div>
            </div>
        </section>
    );
};

export default Hero;