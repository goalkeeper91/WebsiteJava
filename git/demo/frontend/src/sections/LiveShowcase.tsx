import { motion } from 'motion/react';
import { useEffect, useState } from 'react';
import TwitchEmbed from '../components/live/TwitchEmbed';
import HighlightVideo from '../components/live/RandomYoutubePlayer';
import { Link } from 'react-router-dom';

const LiveShowcase = () => {
    const [isLive, setIsLive] = useState(false);
    const [marketingConsent, setMarketingConsent] = useState(false);

    useEffect(() => {
        fetch(`/api/twitch/status`)
            .then(res => res.ok ? res.json() : null)
            .then(data => {
                if (data) setIsLive(data.isLive);
            })
            .catch(() => setIsLive(false));

        const checkConsent = () => {
            if (window.Cookiebot?.consent?.marketing) {
                setMarketingConsent(true);
            }
        };

        window.addEventListener("CookieConsentDeclaration", checkConsent);
        checkConsent();

        return () => window.removeEventListener("CookieConsentDeclaration", checkConsent);
    }, []);

    return (
        <section className='relative w-full h-full py-16 px-6 bg-slate-950'>
            <div className='relative z-10 max-w-5xl mx-auto text-center'>
                <motion.h2
                    className='text-4xl font-bold mb-6'
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ duration: 0.5 }}
                >
                    {isLive
                        ? "Live dabei: Softwareentwicklung & Gaming"
                        : "Gerade offline – hier ein Highlight aus meinen Coding-Sessions"}
                </motion.h2>

                <motion.p
                  className="text-lg text-gray-300 mb-6"
                  initial={{ opacity: 0 }}
                  animate={{ opacity: 1 }}
                  transition={{ delay: 0.2 }}
                >
                  Ich streame nicht nur Games, sondern auch meine Arbeit als Entwickler.
                  Du bekommst live Einblicke in Projekte, Tools & Workflows – so wie sie in echten Softwareprojekten entstehen.
                </motion.p>

                <motion.div
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    transition={{ delay: 0.3 }}
                >
                    {isLive ? (
                        <TwitchEmbed channel='goalkeeper91' />
                    ) : marketingConsent ? (
                        <HighlightVideo />
                    ) : (
                        <p className="text-gray-400">
                            Videos werden nach Zustimmung zu Marketing-Cookies angezeigt.
                        </p>
                    )}
                </motion.div>
                <Link to="/contact" className="text-indigo-400 hover:underline">
                  Projektanfrage stellen
                </Link>
            </div>
        </section>
    );
};

export default LiveShowcase;
