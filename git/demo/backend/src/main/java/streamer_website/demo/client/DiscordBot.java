package streamer_website.demo.client;

import discord4j.core.DiscordClient;
import discord4j.core.GatewayDiscordClient;
import discord4j.gateway.intent.IntentSet;
import lombok.Getter;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.stereotype.Component;
import streamer_website.demo.handler.discord.CommandEventHandler;
import streamer_website.demo.handler.discord.GuildEventHandler;
import streamer_website.demo.service.discord.GuildService;
import streamer_website.demo.service.discord.StatusService;

@Component
public class DiscordBot {

    private static final Logger logger = LoggerFactory.getLogger(DiscordBot.class);

    @Value("${discord.bot.token}")
    private String token;

    private final GuildEventHandler guildEventHandler;
    private final CommandEventHandler commandHandler;
    private final StatusService statusService;
    private final GuildService guildService;

    @Getter
    private GatewayDiscordClient gateway;

    public DiscordBot(GuildEventHandler guildEventHandler,
                      CommandEventHandler commandHandler,
                      StatusService statusService,
                      GuildService guildService) {
        this.guildEventHandler = guildEventHandler;
        this.commandHandler = commandHandler;
        this.statusService = statusService;
        this.guildService = guildService;
    }

    @Bean
    public GatewayDiscordClient start() {
        if (token == null || token.isEmpty()) {
            throw new IllegalStateException("Discord token is not set!");
        }

        logger.info("Starte Discord Bot …");

        try {
            DiscordClient client = DiscordClient.create(token);
            gateway = client.login().block();

            if (gateway == null) {
                throw new IllegalStateException("Failed to login to Discord. Check your token.");
            }
        } catch (Exception e) {
            logger.error("Fehler beim Discord-Login: {}", e.getMessage(), e);
            throw new IllegalStateException("Anwendung konnte keine Verbindung zu Discord herstellen.", e);
        }

        logger.info("Erfolgreich bei Discord eingeloggt. Registriere Events...");
        statusService.setRunning(true);

        // Events registrieren
        guildEventHandler.register(gateway);
        commandHandler.register(gateway);

        // Sync der Guilds beim Start
        guildService.syncGuilds(gateway).subscribe(); // Nicht-blockierender Aufruf

        gateway.onDisconnect()
                .doOnTerminate(() -> statusService.setRunning(false))
                .subscribe();

        return gateway;
    }
}