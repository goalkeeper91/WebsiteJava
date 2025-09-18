package streamer_website.demo.service.discord;

import discord4j.core.event.domain.message.MessageCreateEvent;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import streamer_website.demo.handler.discord.CommandEventHandler;

@Service
@RequiredArgsConstructor
public class CommandService {
    private final CommandEventHandler commandEventHandler;

    public void handle(MessageCreateEvent event) {
        commandEventHandler.handle(event).subscribe(); // wichtig: Reactor-Stream starten
    }
}
