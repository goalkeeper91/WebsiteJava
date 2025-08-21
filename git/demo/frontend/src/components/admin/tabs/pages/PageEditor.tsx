import { useState } from 'react';
import YoutubeConfigForm from '../../forms/YoutubeConfigForm';
import TiktokVideoForm from '../../forms/TiktokConfigForm';

const PageEditor = ({ pageId }: { pageId: string }) => {
  const [activeTab, setActiveTab] = useState<'youtube' | 'tiktok'>('youtube');

  const renderForm = () => {
    switch (pageId) {
      case 'all videos':
        return (
          <div>
            <div className="mb-6 flex space-x-4">
              <button
                onClick={() => setActiveTab('youtube')}
                className={`px-4 py-2 rounded ${
                  activeTab === 'youtube'
                    ? 'bg-blue-600 text-white'
                    : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
                }`}
              >
                YouTube Konfiguration
              </button>
              <button
                onClick={() => setActiveTab('tiktok')}
                className={`px-4 py-2 rounded ${
                  activeTab === 'tiktok'
                    ? 'bg-blue-600 text-white'
                    : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
                }`}
              >
                TikTok Videos verwalten
              </button>
            </div>

            <div>
              {activeTab === 'youtube' && <YoutubeConfigForm />}
              {activeTab === 'tiktok' && <TiktokVideoForm />}
            </div>
          </div>
        );
      default:
        return <p>Kein Editor für diese Seite vorhanden.</p>;
    }
  };

  return (
    <div className="p-4">
      <h2 className="text-3xl font-bold mb-6">Bearbeite Seite: {pageId}</h2>
      {renderForm()}
    </div>
  );
};

export default PageEditor;
