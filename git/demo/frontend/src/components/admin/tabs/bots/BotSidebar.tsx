const bots = [
  { id: 'twitchBot', name: 'Twitch Bot' },
  { id: 'discordBot', name: 'Discord Bot' },
  { id: 'widgetImplementation', name: 'Widgets' },
];

const BotSidebar = ({ selectedBotId, onSelectBot }: {
  selectedBotId: string | null;
  onSelectBot: (id: string) => void;
}) => {
  return (
    <ul className="p-4 space-y-2">
      {bots.map((bot) => (
        <li
          key={bot.id}
          className={`cursor-pointer p-2 rounded ${
            selectedBotId === bot.id ? 'bg-blue-500 text-white' : 'hover:bg-blue-700'
          }`}
          onClick={() => onSelectBot(bot.id)}
        >
          {bot.name}
        </li>
      ))}
    </ul>
  );
};

export default BotSidebar;
