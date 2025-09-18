import { useState } from 'react';
import BotSidebar from './BotSidebar';
import BotEditor from './BotEditor';

const BotManager = () => {
  const [selectedBot, setSelectedBot] = useState<{ botId: string; subId?: string } | null>(null);

  return (
    <div className="flex h-full">
      <div className="w-1/4 border-r bg-gray-800">
        <BotSidebar
          selectedBotId={selectedBot?.subId || selectedBot?.botId || null}
          onSelectBot={(bot) => setSelectedBot(bot)}
        />
      </div>
      <div className="flex-1 p-4">
        {selectedBot ? (
          <BotEditor botId={selectedBot.botId} subId={selectedBot.subId} />
        ) : (
          <div className="text-gray-500">Bitte einen Bot auswählen.</div>
        )}
      </div>
    </div>
  );
};

export default BotManager;