import { useState } from "react";

const bots = [
  { id: "twitchBot", name: "Twitch Bot" },
  { id: "discordBot", name: "Discord Bot", subItems: [
      { id: "dashboard", name: "Dashboard"},
      { id: "discordChannels", name: "Channel IDs" },
    ]
  },
  { id: "widgetImplementation", name: "Widgets" },
];

const BotSidebar = ({
  selectedBotId,
  onSelectBot,
}: {
  selectedBotId: string | null;
  onSelectBot: (bot: { botId: string; subId?: string }) => void;
}) => {
  const [openSubmenu, setOpenSubmenu] = useState<string | null>(null);

  const handleClick = (botId: string, hasSubmenu?: boolean) => {
    if (hasSubmenu) {
      setOpenSubmenu(openSubmenu === botId ? null : botId);
    } else {
      onSelectBot({ botId });
    }
  };

  return (
    <ul className="p-4 space-y-2">
      {bots.map((bot) => (
        <li key={bot.id}>
          <div
            className={`cursor-pointer p-2 rounded flex justify-between ${
              selectedBotId === bot.id ? "bg-blue-500 text-white" : "hover:bg-blue-700"
            }`}
            onClick={() => handleClick(bot.id, !!bot.subItems)}
          >
            <span>{bot.name}</span>
            {bot.subItems && <span className="text-sm">{openSubmenu === bot.id ? "▲" : "▼"}</span>}
          </div>

          {bot.subItems && openSubmenu === bot.id && (
            <ul className="ml-4 mt-1 space-y-1">
              {bot.subItems.map((sub) => (
                <li
                  key={sub.id}
                  className={`cursor-pointer p-2 rounded ${
                    selectedBotId === sub.id ? "bg-blue-400 text-white" : "hover:bg-blue-600"
                  }`}
                  onClick={() => onSelectBot({ botId: bot.id, subId: sub.id })}
                >
                  {sub.name}
                </li>
              ))}
            </ul>
          )}
        </li>
      ))}
    </ul>
  );
};

export default BotSidebar;
