import { useState } from 'react';
import BotSidebar from './BotSidebar';
import BotEditor from './BotEditor';

const BotManager = () => {
  const [selectedBotId, setSelectedBotId] = useState<string | null>(null);

  return (
    <div className="flex h-full">
      <div className="w-1/4 border-r bg-gray-800">
        <BotSidebar onSelectBot={setSelectedBotId} selectedBotId={selectedBotId} />
      </div>
      <div className="flex-1 p-4">
        {selectedBotId ? (
          <BotEditor botId={selectedBotId} />
        ) : (
          <div className="text-gray-500">Bitte einen Bot auswählen.</div>
        )}
      </div>
    </div>
  );
};

export default BotManager;