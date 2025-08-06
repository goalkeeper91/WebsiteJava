import React, { useState, useEffect } from 'react';

interface TiktokVideo {
  id?: number;
  title: string;
  videoId: string;
  platform: 'tiktok';
}

const TiktokVideoForm: React.FC = () => {
  const [videos, setVideos] = useState<TiktokVideo[]>([]);
  const [newVideo, setNewVideo] = useState<TiktokVideo>({ title: '', videoId: '', platform: 'tiktok' });
  const [loading, setLoading] = useState(true);
  const [message, setMessage] = useState('');
  const baseUrl = 'http://localhost:8080';

  // Bestehende TikTok-Videos laden
  useEffect(() => {
    fetch(`${baseUrl}/api/videos/platform/tiktok`)
      .then(res => res.json())
      .then(data => setVideos(data))
      .catch(() => setMessage('Fehler beim Laden der TikTok-Videos'))
      .finally(() => setLoading(false));
  }, []);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setNewVideo(prev => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setMessage('');

    try {
      const res = await fetch(`${baseUrl}/api/videos`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(newVideo),
      });

      if (res.ok) {
        const created = await res.json();
        setVideos(prev => [...prev, created]);
        setNewVideo({ title: '', videoId: '', platform: 'tiktok' });
        setMessage('Video erfolgreich hinzugefügt');
      } else {
        throw new Error();
      }
    } catch {
      setMessage('Fehler beim Speichern');
    }
  };

  const handleDelete = async (id: number) => {
    try {
      await fetch(`${baseUrl}/api/videos/${id}`, { method: 'DELETE' });
      setVideos(prev => prev.filter(video => video.id !== id));
    } catch {
      setMessage('Fehler beim Löschen');
    }
  };

  if (loading) return <p>Lade TikTok-Videos...</p>;

  return (
    <div className="space-y-6">
      <h2 className="text-xl font-bold">TikTok Video hinzufügen</h2>

      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block font-semibold">Titel</label>
          <input
            type="text"
            name="title"
            value={newVideo.title}
            onChange={handleChange}
            className="w-full p-2 border rounded"
            required
          />
        </div>

        <div>
          <label className="block font-semibold">TikTok Video ID</label>
          <input
            type="text"
            name="videoId"
            value={newVideo.videoId}
            onChange={handleChange}
            className="w-full p-2 border rounded"
            required
          />
        </div>

        <button
          type="submit"
          className="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700"
        >
          Speichern
        </button>
      </form>

      {message && <p className="text-sm mt-2 text-gray-700">{message}</p>}

      <div>
        <h3 className="text-lg font-semibold mt-8 mb-4">Bestehende TikTok Videos</h3>
        <ul className="space-y-2">
          {videos.map((video) => (
            <li key={video.id} className="flex justify-between items-center border p-2 rounded">
              <span>{video.title} ({video.videoId})</span>
              <button
                onClick={() => handleDelete(video.id!)}
                className="text-sm text-red-600 hover:underline"
              >
                Löschen
              </button>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
};

export default TiktokVideoForm;
