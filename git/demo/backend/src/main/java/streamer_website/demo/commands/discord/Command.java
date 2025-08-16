package streamer_website.demo.commands.discord;

import discord4j.core.event.domain.message.MessageCreateEvent;

public interface Command {
    String getName();
    void execute(MessageCreateEvent event, String[] args);
}
