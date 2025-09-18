import React from "react";
import TwitchBotControls from "./BotActions";
import DiscordBotControls from "./discord/DiscordBotStatus";
import DiscordChannelEditor from "./discord/DiscordChannelsEditor";

type BotEditorProps = {
  botId: string;
  subId?: string;
};

const BotEditor: React.FC<BotEditorProps> = ({ botId, subId }) => {
  return (
    <div>
      <h2 className="text-xl font-bold mb-4">
        Bearbeite: {botId} {subId ? `> ${subId}` : ""}
      </h2>

      {botId === "twitchBot" ? (
        <TwitchBotControls />
      ) : botId === "discordBot" ? (
        subId === "dashboard" ? (
          <DiscordBotControls />
        ) : subId === "discordChannels" ? (
          <DiscordChannelEditor />
        ) : (
          <textarea
            className="w-full h-64 border p-2"
            placeholder={`Bot Einstellungen "${botId}${subId ? ` > ${subId}` : ""}" anpassen...`}
          />
        )
      ) : (
        <textarea
          className="w-full h-64 border p-2"
          placeholder={`Bot Einstellungen "${botId}${subId ? ` > ${subId}` : ""}" anpassen...`}
        />
      )}
    </div>
  );
};

export default BotEditor;
