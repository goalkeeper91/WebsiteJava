package streamer_website.demo.handler.discord;

import discord4j.core.GatewayDiscordClient;
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

        if (content.startsWith("!tasks")) {
            String argsStr = content.substring("!tasks".length()).trim();
            String[] args = new String[]{argsStr}; // Übergebe alles als ein String
            return reactor.core.publisher.Mono.fromRunnable(() -> taskCommands.execute(event, args));
        }

        return reactor.core.publisher.Mono.empty();
    }

    public void register(GatewayDiscordClient gateway) {
        // Hier startet man den Stream korrekt
        gateway.on(MessageCreateEvent.class)
                .flatMap(this::handle)  // flatMap auf die Methode, die Mono<Void> zurückgibt
                .subscribe();           // unbedingt subscribe, sonst passiert nichts
    }
}

