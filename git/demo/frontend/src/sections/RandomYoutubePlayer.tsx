import { useEffect, useState, useRef } from 'react';

const fetchVideoIds = async (): Promise<string[]> => {
    const res = await fetch('/api/videos');
    const data: { videoId: string }[] = await res.json();
    return data.map(v => v.videoId);
};


const getRandomVideoId = (videoIds: string[], excludeId?: string): string => {
    const filtered = excludeId ? videoIds.filter(id => id !== excludeId) : videoIds;
    return filtered[Math.floor(Math.random() * filtered.length)];
};

const RandomYoutubePlayer = () => {
    const [videoIds, setVideoIds] = useState<string[]>([]);
    const [currentVideoId, setCurrentVideoId] = useState<string | null>(null);
    const playerRef = useRef<any>(null);

    useEffect(() => {
        fetchVideoIds().then(ids => {
            setVideoIds(ids);
            if (ids.length > 0) {
                setCurrentVideoId(getRandomVideoId(ids));
            }
        });
    }, []);

    const handleVideoEnd = () => {
        if (!videoIds.length) return;
        const next = getRandomVideoId(videoIds, currentVideoId || undefined);
        setCurrentVideoId(next);
    };

    useEffect(() => {
        if (!videoIds.length || !currentVideoId) return;

        const tag = document.createElement('script');
        tag.src = 'https://www.youtube.com/iframe_api';
        document.body.appendChild(tag);

        (window as any).onYouTubeIframeAPIReady = () => {
            playerRef.current = new (window as any).YT.Player('youtube-player', {
                height: '390',
                width: '640',
                videoId: currentVideoId,
                events: {
                    onStateChange: (event: any) => {
                        if (event.data === 0) handleVideoEnd();
                    },
                },
            });
        };

        return () => {
            if (playerRef.current?.destroy) playerRef.current.destroy();
        };
    }, [videoIds, currentVideoId]);

    useEffect(() => {
        if (playerRef.current?.loadVideoById && currentVideoId) {
            playerRef.current.loadVideoById(currentVideoId);
        }
    }, [currentVideoId]);

    return (
        <section className="py-12 px-6 bg-black text-white text-center">
            <h2 className="text-3xl font-bold mb-6">📺 Zufälliges Video</h2>
            <div className="flex justify-center">
                <div id="youtube-player" className="rounded-lg overflow-hidden shadow-lg" />
            </div>
        </section>
    );
};

export default RandomYoutubePlayer;
