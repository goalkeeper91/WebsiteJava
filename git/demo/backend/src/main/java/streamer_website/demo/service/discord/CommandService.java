package streamer_website.demo.service.discord;

import discord4j.core.GatewayDiscordClient;
import discord4j.core.event.domain.message.MessageCreateEvent;
import discord4j.core.object.entity.User;
import jakarta.annotation.PostConstruct;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import reactor.core.publisher.Mono;
import streamer_website.demo.commands.discord.Command;

import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

@Service
@RequiredArgsConstructor
public class CommandService {

    private final List<Command> commands; // Spring injiziert alle Command-Beans automatisch
    private final Map<String, Command> commandMap = new HashMap<>();

    @PostConstruct
    public void init() {
        for (Command cmd : commands) {
            commandMap.put(cmd.getName().toLowerCase(), cmd);
        }
    }

    public void register(GatewayDiscordClient client) {
        client.on(MessageCreateEvent.class)
            .flatMap(event -> {
                // Bots ignorieren
                if (event.getMessage().getAuthor().map(User::isBot).orElse(true)) {
                    return Mono.empty();
                }

                String content = event.getMessage().getContent().trim();
                if (!content.startsWith("!jtc")) {
                    return Mono.empty();
                }

                // z. B. !jtc add ...
                String[] parts = content.split("\\s+");
                if (parts.length < 2) {
                    return Mono.empty();
                }

                String commandName = parts[1].toLowerCase();
                String[] args = Arrays.copyOfRange(parts, 2, parts.length);

                Command cmd = commandMap.get(commandName);
                if (cmd == null) {
                    return event.getMessage().getChannel()
                        .flatMap(ch -> ch.createMessage("Unbekannter Befehl: " + commandName))
                        .then();
                }

                // Command ausführen
                cmd.execute(event, args);
                    return Mono.empty();
                })
                .subscribe();
        }
}
