import React, { useEffect, useState } from 'react';
import { useVideos } from '../../hooks/useVideos';

const YoutubeGrid: React.FC = () => {
  const [videos, setVideos] = useState<{ id: string; title: string; embedUrl: string }[]>([]);
  const [featuredVideo, setFeaturedVideo] = useState<typeof videos[0] | null>(null);

  const { youtubeVideos, loading } = useVideos();

  useEffect(() => {
    // warte auf Daten und prüfe ob youtubeVideos vorhanden sind
    if (loading || youtubeVideos.length === 0) return;

    const videoConfig = youtubeVideos[0]; // wir nehmen den ersten als Konfiguration

    if (!videoConfig?.apiKey || !videoConfig?.channelId) {
      console.error("Fehlende YouTube-Konfiguration");
      return;
    }

    const API_KEY = videoConfig.apiKey;
    const CHANNEL_ID = videoConfig.channelId;
    const MAX_RESULTS = 10;

    fetch(`https://www.googleapis.com/youtube/v3/search?key=${API_KEY}&channelId=${CHANNEL_ID}&part=snippet,id&order=date&maxResults=${MAX_RESULTS}`)
      .then((res) => res.json())
      .then((data) => {
        if (!Array.isArray(data.items)) {
            return { fetchedVideos: [], loading: false };
        }
        const fetchedVideos = data.items
          .filter((item: any) => item.id.kind === "youtube#video")
          .map((item: any) => ({
            id: item.id.videoId,
            title: item.snippet.title,
            embedUrl: `https://www.youtube.com/embed/${item.id.videoId}`,
          }));

        setVideos(fetchedVideos);
      })
      .catch(console.error);
  }, [loading, youtubeVideos]);

  useEffect(() => {
    if (videos.length > 0 && !featuredVideo) {
      const random = videos[Math.floor(Math.random() * videos.length)];
      setFeaturedVideo(random);
    }
  }, [videos]);

  if (loading) return <p>Lade YouTube-Konfiguration...</p>;

  return (
    <section className="relative w-full py-20 px-4 text-white overflow-hidden">
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
                title={featuredVideo.title}
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
                title={video.title}
              ></iframe>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
};

export default YoutubeGrid;
