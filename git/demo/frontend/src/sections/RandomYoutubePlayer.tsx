import { useEffect, useState, useRef } from 'react';

const videoIds = [
    'JDSXB-n1QLQ',
    'uPzpZbB4mr8&pp=0gcJCccJAYcqIYzv',
    'AA_m8rLcIx0',
    'MWrnLqYwfjE&pp=0gcJCccJAYcqIYzv',
    'PzggVsIcSuY&pp=0gcJCccJAYcqIYzv',
    'sokBPz93o5g',
    'hNwE7r25Amv8',
    'iN0P_5q8Iso',
    '6jvlln4p3o0',
    'gaZXybW9_Ew'
];

const getRandomVideoId = (excludeId?: string): string => {
    const filtered = excludeId ? videoIds.filter(id => id !== excludeId) : videoIds;
    return filtered[Math.floor(Math.random() * filtered.length)];
};

const RandomYoutubePlayer = () => {
    const [currentVideoId, setCurrentVideoId] = useState(getRandomVideoId());
    const playerRef = useRef<any>(null);

    const handleVideoEnd = () => {
        const next = getRandomVideoId(currentVideoId);
        setCurrentVideoId(next);
    };

    useEffect(() => {
        const tag = document.createElement('script');
        tag.src = 'https://www.youtube.com/iframe_api';
        document.body.appendChild(tag);

        (window as any).onYoutubeIframeAPIReady = () => {
            playerRef.current = new (window as any).YT.Player('youtube-player', {
                heigh: '390',
                width: '640',
                videoId: currentVideoId,
                events: {
                    onStateChange: (event.any) => {
                        if (event.data === 0) handleVideoEnd();
                    },
                },
            });
        };

        return () => {
            if (playerRef.current?.destroy) playerRef.current.destroy();
        };
    }, ());

    useEffect(() => {
        if (playerRef.current?.loadVideoById) {
            playerRef.current.loadVideoById(currentVideoId);
        }
    }, [currentVideoId]);

    return (
        <section className='py-12 px-6 bg-black text-white text-center'>
            <h2 className='text-3xl font-bold mb-6'> 📺 Zufälliges Video</h2>
            <div className="flex justify-center">
                <div id="youtube-player" className="rounded-lg overflow-hidden shadow-lg" />
            </div>
        </section>
    );
};

export default RandomYoutubePlayer;