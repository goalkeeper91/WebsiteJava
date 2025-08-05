import React, { useState } from 'react';
import AdminTabs from '../../components/admin/AdminTabs';
import PageManager from '../../components/admin/tabs/pages/PageManager';
import BotManager from '../../components/admin/tabs/bots/BotManager';

const AdminDashboard: React.FC = () => {
    const [activeTab, setActiveTab] = useState<'pages' | 'bots'>('pages');

  return (
    <div className="flex flex-col items-center justify-center bg-slate-900 text-white p-4">
      <AdminTabs activeTab={activeTab} onTabChange={setActiveTab} />
      <h1 className="text-4xl font-bold mb-4 text-goalyBlue">Admin Dashboard</h1>
      <div className="flex-1 bg-slate-800 p-6 rounded-2xl shadow-md mt-4 min-w-full">
        {activeTab === 'pages' && <PageManager />}
        {activeTab === 'bots' && <BotManager />}
      </div>
    </div>
  );
};

export default AdminDashboard;
