package streamer_website.demo.commands.discord;

import discord4j.core.event.domain.message.MessageCreateEvent;

public class PingCommand implements Command {
    @Override
    public String getName() {
        return "!ping";
    }

    @Override
    public void execute(MessageCreateEvent event, String[] args) {
        event.getMessage().getChannel().flatMap(channel -> channel.createMessage("Pong!"))
                .subscribe();
    }
}
