import React, { useEffect, useState } from 'react';

const YoutubeGrid: React.FC = () => {
    const [videos, setVideos] = useState<{ id: string; title: string; embedUrl: string } []>([]);
    const [featuredVideo, setFeaturedVideo] = useState<typeof videos[0] | null>(null);

    useEffect(() => {
        const API_KEY = import.meta.env.VITE_YOUTUBE_API_KEY;
        const CHANNEL_ID = 'UCXqL6NXrPmn_oFuwA611iBg';
        const MAX_RESULTS = 10;

        fetch(`https://www.googleapis.com/youtube/v3/search?key=${API_KEY}&channelId=${CHANNEL_ID}&part=snippet,id&order=date&maxResults=${MAX_RESULTS}`)
        .then((res) => res.json())
         .then((data) => {
            const videos = data.items
            .filter((item: any) => item.id.kind === "youtube#video")
            .map((item: any) => ({
                id: item.id.videoId,
                title: item.snippet.title,
                embedUrl: `https://www.youtube.com/embed/${item.id.videoId}`,
            }));

              setVideos(videos);
            })
            .catch(console.error);
    }, []);

    useEffect(() => {
        if (videos.length > 0) {
            const random = videos[Math.floor(Math.random() * videos.length)];
                    setFeaturedVideo(random);
        }
    }, [videos]);

    return (
        <section className="relative w-full py-20 px-4 text-white overflow-hidden">
          {/* Hintergrundebene */}
          <div className="absolute inset-0 -z-10">
            <div className="w-full h-full bg-black/50" />
            <div className="w-full h-full bg-gradient-to-br from-purple-500 via-pink-500 to-yellow-500 opacity-30 blur-3xl" />
          </div>

          {/* Inhaltsebene */}
          <div className="max-w-6xl mx-auto z-10 relative">
            <h1 className="text-4xl font-bold mb-8 text-center">Video Highlights</h1>

            {featuredVideo && (
              <div className="mb-12">
                <h2 className="text-2xl font-semibold mb-4 text-center">Empfohlenes Video</h2>
                <div className="aspect-video w-full max-w-3xl mx-auto">
                  <iframe
                    src={featuredVideo.embedUrl}
                    className="w-full h-full rounded-lg bg-none"
                    allowFullScreen
                    loading="lazy"
                  ></iframe>
                </div>
              </div>
            )}

            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-8">
              {videos.map((video) => (
                <div key={video.id} className="aspect-video">
                  <iframe
                    src={video.embedUrl}
                    className="w-full h-full rounded-lg"
                    allowFullScreen
                    loading="lazy"
                  ></iframe>
                </div>
              ))}
            </div>
          </div>
        </section>
    );
};

export default YoutubeGrid;