package streamer_website.demo.handler.discord;

import discord4j.core.event.domain.message.MessageCreateEvent;
import discord4j.core.object.entity.User;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Component;
import streamer_website.demo.commands.discord.TaskCommands;

@Component
@RequiredArgsConstructor
public class CommandEventHandler {

    private final TaskCommands taskCommands;

    public reactor.core.publisher.Mono<Void> handle(MessageCreateEvent event) {
        // Bot-Nachrichten ignorieren
        if (event.getMessage().getAuthor().map(User::isBot).orElse(false)) {
            return reactor.core.publisher.Mono.empty();
        }

        String content = event.getMessage().getContent().trim();

        // Prüfe, ob es ein Task-Command ist (z.B. "!tasks …")
        if (content.startsWith("!tasks")) {
            String argsStr = content.substring("!tasks".length()).trim();
            String[] args = argsStr.isEmpty() ? new String[0] : argsStr.split("\\s+");
            taskCommands.execute(event, args);
        }

        return reactor.core.publisher.Mono.empty();
    }
}

