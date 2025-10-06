package streamer_website.demo.commands.discord;

import discord4j.core.event.domain.message.MessageCreateEvent;
import reactor.core.publisher.Mono;

public interface Command {
    String getName();
    Mono<Void> execute(MessageCreateEvent event, String[] args);
}
