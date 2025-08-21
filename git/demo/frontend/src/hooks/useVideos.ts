import { useEffect, useState } from 'react';

export type Video = {
  id: number;
  platform: 'youtube' | 'tiktok';
  videoId?: string;
  title?: string;
  apiKey?: string;
  channelId?: string;
};

export const useVideos = () => {
  const [videos, setVideos] = useState<Video[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchVideos = async () => {
      try {
        const response = await fetch('/api/videos');
        const data = await response.json();
        setVideos(data);
      } catch (error) {
        console.error('Fehler beim Laden der Videos:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchVideos();
  }, []);

  if (!Array.isArray(videos)) {
    console.error("useVideos: videos ist kein Array:", videos);
    return { youtubeVideos: [], loading: false };
  }

  const youtubeVideos = videos.filter(v => v.platform === 'youtube');
  const tiktokVideos = videos.filter(v => v.platform === 'tiktok');

  return { videos, youtubeVideos, tiktokVideos, loading };
};
