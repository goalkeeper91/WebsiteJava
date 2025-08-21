import { motion } from 'motion/react';
import { useEffect, useState } from 'react';
import TwitchEmbed from '../components/live/TwitchEmbed';
import HighlightVideo from '../components/live/RandomYoutubePlayer';

const LiveShowcase = () => {
    const [isLive, setIsLive] = useState(false);

    useEffect(() => {
      fetch(`/api/twitch/status`)
        .then(res => {
          if (!res.ok) {
            setIsLive(false);
            return null;
          }
          return res.json();
        })
        .then(data => {
          if (data) setIsLive(data.isLive);
        })
        .catch(() => {
          setIsLive(false);
        });
    }, []);


    return (
        <section className='relative w-full h-full py-16 px-6'>
            <div className="absolute inset-0 bg-black/50 z-0" />
            <div className='absolute z-0 inset-0 bg-gradient-to-br from-purple-500 via-pink-500 to-yellow-500 opacity-30 blur-3xl'></div>
            <div className='relative z-10 max-w-5xl mx-auto text-center'>
                <motion.h2
                    className='text-4xl font-bold mb-6'
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ duration: 0.5 }}
                >
                    {isLive ? "Schau dir meinen Stream an" : "Leider bin ich gerade offline, aber hier ein Video"}
                </motion.h2>

                <motion.div
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    transition={{ delay: 0.3 }}
                >
                    {isLive ? (
                        <TwitchEmbed channel='goalkeeper91' />
                    ) : (
                        <HighlightVideo />
                    )}
                </motion.div>
            </div>
        </section>
    );
};

export default LiveShowcase;