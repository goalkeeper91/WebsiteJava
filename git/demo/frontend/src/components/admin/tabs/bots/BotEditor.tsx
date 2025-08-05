import React from 'react';

const BotEditor = ({ botId }: { botId: string }) => {
  return (
    <div>
      <h2 className="text-xl font-bold mb-4">Bearbeite: {botId}</h2>
      <textarea
        className="w-full h-64 border p-2"
        placeholder={`Bot Einstellungen "${botId}" anpassen...`}
      />
    </div>
  );
};

export default BotEditor;
