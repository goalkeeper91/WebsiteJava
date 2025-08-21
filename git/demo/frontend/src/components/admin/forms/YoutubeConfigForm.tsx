import React, { useEffect, useState } from 'react';

interface YoutubeConfig {
  id?: number;
  apiKey: string;
  channelId: string;
}

const YoutubeConfigForm: React.FC = () => {
  const [config, setConfig] = useState<YoutubeConfig>({ apiKey: '', channelId: '' });
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [message, setMessage] = useState('');

  useEffect(() => {
    fetch(`/api/videos/platform/youtube`)
      .then(res => res.json())
      .then((data) => {
        const configEntry = data.find((v: any) => !v.videoId); // Kein echtes Video = Konfiguration
        if (configEntry) {
          setConfig({
            id: configEntry.id,
            apiKey: configEntry.apiKey || '',
            channelId: configEntry.channelId || '',
          });
        }
      })
      .catch(() => setMessage('Fehler beim Laden der Konfiguration'))
      .finally(() => setLoading(false));
  }, []);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setConfig(prev => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    setMessage('');

    const payload = {
      ...config,
      platform: 'youtube',
      title: 'YouTube Config',
      videoId: null, // Wichtig, damit das Backend zufrieden ist
    };

    const url = config.id ? `/api/videos/${config.id}` : `/api/videos`;
    const method = config.id ? 'PUT' : 'POST';

    try {
      const res = await fetch(url, {
        method,
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      });

      const result = await res.text();
        console.log(result);
      if (res.ok) {
        setMessage('Erfolgreich gespeichert');
      } else {
        setMessage(`Fehler: ${result}`);
      }
    } catch (err) {
      setMessage('Fehler beim Speichern');
    } finally {
      setSaving(false);
    }
  };

  if (loading) return <p>Lade Konfiguration...</p>;

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div>
        <label className="block font-semibold">YouTube API Key</label>
        <input
          type="text"
          name="apiKey"
          value={config.apiKey}
          onChange={handleChange}
          className="w-full p-2 border rounded"
          required
        />
      </div>

      <div>
        <label className="block font-semibold">Channel ID</label>
        <input
          type="text"
          name="channelId"
          value={config.channelId}
          onChange={handleChange}
          className="w-full p-2 border rounded"
          required
        />
      </div>

      <button
        type="submit"
        disabled={saving}
        className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
      >
        {saving ? 'Speichern...' : 'Speichern'}
      </button>

      {message && <p className="mt-2 text-sm text-gray-700">{message}</p>}
    </form>
  );
};

export default YoutubeConfigForm;
