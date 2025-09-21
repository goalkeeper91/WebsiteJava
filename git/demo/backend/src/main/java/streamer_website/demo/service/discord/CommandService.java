package streamer_website.demo.service.discord;

import discord4j.core.event.domain.message.MessageCreateEvent;
import org.springframework.stereotype.Service;
import streamer_website.demo.commands.discord.Command;
import streamer_website.demo.commands.discord.PingCommand;

import java.util.HashMap;
import java.util.Map;

@Service
public class CommandService {

    private final Map<String, Command> commands = new HashMap<>();

    public CommandService() {
        register(new PingCommand());
    }

    private void register(Command command) {
        commands.put(command.getName().toLowerCase(), command);
    }

    public void handle(MessageCreateEvent event) {
        String content = event.getMessage().getContent();

        if (!content.startsWith("!")) return;

        String[] parts = content.split(" ");
        String commandName = parts[0].toLowerCase();
        String[] args = parts.length > 1 ? content.substring(commandName.length()).trim().split(" ")
                : new String[0];

        Command command = commands.get(commandName);

        if (command == null) return;

        command.execute(event, args);
    }
}
