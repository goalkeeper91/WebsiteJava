package streamer_website.demo.client;

import discord4j.core.DiscordClient;
import discord4j.core.GatewayDiscordClient;
import lombok.Getter;
import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.stereotype.Component;
import streamer_website.demo.handler.discord.CommandEventHandler;
import streamer_website.demo.handler.discord.GuildEventHandler;
import streamer_website.demo.service.discord.GuildService;
import streamer_website.demo.service.discord.JoinToCreateService;
import streamer_website.demo.service.discord.StatusService;

@Component
@RequiredArgsConstructor
public class DiscordBot {

    private static final Logger logger = LoggerFactory.getLogger(DiscordBot.class);

    @Value("${discord.bot.token:}")
    private String token; // Default leer, damit null-safe

    private final GuildEventHandler guildEventHandler;
    private final CommandEventHandler commandHandler;
    private final StatusService statusService;
    private final GuildService guildService;
    private final JoinToCreateService joinToCreateService;

    @Bean
    public GatewayDiscordClient gatewayDiscordClient() {
        if (token == null || token.isEmpty()) {
            logger.warn("Kein Discord-Token gesetzt – Bot bleibt inaktiv");
            return null;
        }

        try {
            logger.info("Starte Discord Bot…");
            DiscordClient client = DiscordClient.create(token);
            GatewayDiscordClient gateway = client.login().block();
            if (gateway == null) {
                logger.error("Discord Login fehlgeschlagen");
                return null;
            }

            // Null-safe: DB kann leer sein
            joinToCreateService.initConfigs();
            joinToCreateService.register(gateway);

            guildEventHandler.register(gateway);
            commandHandler.register(gateway);
            guildService.syncGuilds(gateway).subscribe();

            statusService.setRunning(true);
            logger.info("Discord Bot erfolgreich gestartet");
            return gateway;
        } catch (Exception e) {
            logger.error("Fehler beim Start des DiscordBots", e);
            return null; // Kein Throw, damit Spring Container nicht crasht
        }
    }
}