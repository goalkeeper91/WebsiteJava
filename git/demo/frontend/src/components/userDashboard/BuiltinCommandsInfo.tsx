const BUILTIN_COMMANDS = [
  {
    trigger: "uptime",
    usage: "!uptime",
    description: "Zeigt, wie lange der Stream schon live ist (oder dass er offline ist).",
  },
  {
    trigger: "title",
    usage: "!title [neuer Titel]",
    description: "Ohne Argument: zeigt den aktuellen Stream-Titel. Mit Argument (nur Mods/Broadcaster): ändert ihn.",
  },
  {
    trigger: "game",
    usage: "!game [neues Spiel]",
    description: "Ohne Argument: zeigt das aktuelle Spiel. Mit Argument (nur Mods/Broadcaster): ändert es.",
  },
  {
    trigger: "followage",
    usage: "!followage [username]",
    description: "Zeigt, seit wann jemand dem Kanal folgt (ohne Argument: der aufrufende Chatter).",
  },
  {
    trigger: "shoutout",
    usage: "!shoutout <username>",
    description: "Postet eine Plug-Nachricht für einen anderen Streamer (nur Mods/Broadcaster).",
  },
];

// Purely informational reference - these commands aren't configurable, so
// unlike CommandsDashboard there's no CRUD here, just a static list.
export default function BuiltinCommandsInfo() {
  return (
    <div className="max-w-5xl mx-auto space-y-4">
      <div>
        <h2 className="text-2xl font-bold text-white">Eingebaute Commands</h2>
        <p className="text-sm text-gray-400 mt-1">
          Diese Commands sind für jeden verbundenen Kanal automatisch aktiv - kein Einrichten nötig.
          Ein eigener Chat-Command mit demselben Namen hat immer Vorrang.
        </p>
      </div>

      <div className="space-y-2">
        {BUILTIN_COMMANDS.map((cmd) => (
          <div key={cmd.trigger} className="bg-gray-800 rounded p-4">
            <div className="flex items-center gap-2 mb-1">
              <p className="font-mono text-green-400 font-semibold">{cmd.usage}</p>
            </div>
            <p className="text-sm text-gray-300">{cmd.description}</p>
          </div>
        ))}
      </div>
    </div>
  );
}
