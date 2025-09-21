package streamer_website.demo.handler.discord;

import discord4j.core.GatewayDiscordClient;
import discord4j.core.event.domain.message.MessageCreateEvent;
import org.springframework.stereotype.Component;
import streamer_website.demo.service.discord.CommandService;

@Component
public class CommandEventHandler {

    private final CommandService commandService;

    public CommandEventHandler(CommandService commandService) {
        this.commandService = commandService;
    }

    public void register(GatewayDiscordClient gateway) {
        gateway.on(MessageCreateEvent.class, event -> {
            commandService.handle(event);
            return null; // reactor braucht ein Publisher -> null passt hier
        }).subscribe();
    }
}

