package streamer_website.demo.commands.discord;

import discord4j.core.event.domain.message.MessageCreateEvent;
import reactor.core.publisher.Mono;

public class PingCommand implements Command {
    @Override
    public String getName() {
        return "!ping";
    }

    @Override
    public Mono<Void> execute(MessageCreateEvent event, String[] args) {
        // ⭐ KORREKTUR: subscribe() entfernen und das Mono zurückgeben.
        return event.getMessage().getChannel()
                .flatMap(channel -> channel.createMessage("Pong!"))
                .then(); // ⭐ Konvertiert das Mono<Message> zu Mono<Void> und gibt es zurück.
    }
}
