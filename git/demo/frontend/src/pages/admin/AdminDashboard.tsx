import React from 'react';

const AdminDashboard: React.FC = () => {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-slate-900 text-white px-4">
      <h1 className="text-4xl font-bold mb-4 text-goalyBlue">Admin Dashboard</h1>
      <p className="text-lg text-gray-300">Willkommen im geschützten Adminbereich!</p>
      <div className="mt-6 bg-slate-800 p-6 rounded-2xl shadow-md w-full max-w-2xl">
        <p className="text-sm text-gray-400">
          Hier kannst du später Streams verwalten, Inhalte bearbeiten, den Chatbot steuern und vieles mehr.
        </p>
      </div>
    </div>
  );
};

export default AdminDashboard;
