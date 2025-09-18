package streamer_website.demo.handler.discord;

import discord4j.core.GatewayDiscordClient;
import discord4j.core.event.domain.message.MessageCreateEvent;
import discord4j.core.object.entity.User;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Component;
import reactor.core.publisher.Mono;
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
        System.out.println("👉 Eingehende Nachricht: " + content);

        if (content.startsWith("!tasks")) {
            System.out.println("✅ Command erkannt: " + content);
            String argsStr = content.substring("!tasks".length()).trim();
            String[] args = new String[]{argsStr}; // Übergebe alles als ein String
            return reactor.core.publisher.Mono.fromRunnable(() -> taskCommands.execute(event, args));
        }

        return Mono.empty();
    }

    public void register(GatewayDiscordClient gateway) {
        gateway.on(MessageCreateEvent.class, this::handle)
                .onErrorContinue((throwable, o) -> {
                    System.err.println("❌ Fehler im Command Handler: " + throwable.getMessage());
                    throwable.printStackTrace();
                })
                .subscribe();
    }
}

