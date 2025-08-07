import React from 'react';

interface Props {
  activeTab: 'pages' | 'bots' | 'stats';
  onTabChange: (tab: 'pages' | 'bots' | 'stats') => void;
}

const AdminTabs: React.FC<Props> = ({ activeTab, onTabChange }) => {
  return (
    <div className="flex space-x-4 mb-6">
      <button
        onClick={() => onTabChange('pages')}
        className={`px-4 py-2 rounded ${
          activeTab === 'pages' ? 'bg-goalyBlue text-white' : 'bg-gray-600 text-gray-300'
        }`}
      >
        Seiten verwalten
      </button>
      <button
        onClick={() => onTabChange('bots')}
        className={`px-4 py-2 rounded ${
          activeTab === 'bots' ? 'bg-goalyBlue text-white' : 'bg-gray-600 text-gray-300'
        }`}
      >
        Bots
      </button>
      <button
        onClick={() => onTabChange('stats')}
        className={`px-4 py-2 rounded ${
        activeTab === 'stats' ? 'bg-goalyBlue text-white' : 'bg-gray-600 text-gray-300'
        }`}
      >
        Stats
      </button>
    </div>
  );
};

export default AdminTabs;
