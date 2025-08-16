package streamer_website.demo.client;

import discord4j.core.DiscordClient;
import discord4j.core.GatewayDiscordClient;
import discord4j.core.event.domain.message.MessageCreateEvent;
import jakarta.annotation.PostConstruct;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;
import streamer_website.demo.service.discord.CommandService;
import streamer_website.demo.service.discord.StatusService;

@Component
public class DiscordBot {

    @Value("${discord.bot.token}")
    private String token;

    private final CommandService commandService;
    private final StatusService statusService;

    private GatewayDiscordClient gateway;

    public DiscordBot(CommandService commandService, StatusService statusService) {
        this.commandService = commandService;
        this.statusService = statusService;
    }

    @PostConstruct
    public void startBot() {

        if (token == null || token.isEmpty()) {
            throw new IllegalStateException("Discord token is not set!");
        }

        DiscordClient client = DiscordClient.create(token);
        gateway = client.login().block();

        if (gateway == null) {
            throw new IllegalStateException("Failed to login to Discord. Check your token.");
        }

        statusService.setRunning(true);

        gateway.on(MessageCreateEvent.class, event -> {
            commandService.handle(event);
            return null;
        }).subscribe();

        gateway.onDisconnect()
                .doOnTerminate(() -> statusService.setRunning(false))
                .subscribe();
    }
}
