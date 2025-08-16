import React from 'react';
import TwitchBotControls from "./BotActions";
import DiscordBotControls from "./discord/DiscordBotStatus"

const BotEditor = ({ botId }: { botId: string }) => {
  return (
    <div>
      <h2 className="text-xl font-bold mb-4">Bearbeite: {botId}</h2>

      {botId === "twitchBot" ? (
        <TwitchBotControls />
      ) : botId === "discordBot" ? (
        <DiscordBotControls />
      ) : (
        <textarea
            className="w-full h-64 border p-2"
            placeholder={`Bot Einstellungen "${botId}" anpassen...`}
        />
      )
      }
    </div>
  );
};

export default BotEditor;
