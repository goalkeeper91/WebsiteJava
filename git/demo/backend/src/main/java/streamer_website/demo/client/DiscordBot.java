package streamer_website.demo.client;

import discord4j.core.DiscordClient;
import discord4j.core.GatewayDiscordClient;
import discord4j.gateway.intent.Intent;
import discord4j.gateway.intent.IntentSet;
import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.stereotype.Component;
import streamer_website.demo.handler.discord.GuildEventHandler;
import streamer_website.demo.service.discord.CommandService;
import streamer_website.demo.service.discord.GuildService;
import streamer_website.demo.service.discord.JoinToCreateService;
import streamer_website.demo.service.discord.StatusService;

@Component
@RequiredArgsConstructor
public class DiscordBot {

    private static final Logger logger = LoggerFactory.getLogger(DiscordBot.class);

    private final GuildEventHandler guildEventHandler;
    private final CommandService commandService;
    private final StatusService statusService;
    private final GuildService guildService;
    private final JoinToCreateService joinToCreateService;

    @Bean
    public GatewayDiscordClient gatewayDiscordClient(@Value("${discord.bot.token:}") String token) {
        if (token == null || token.isEmpty()) {
            logger.warn("Kein Discord-Token gesetzt – Bot bleibt inaktiv");
            return null;
        }

        try {
            logger.info("Starte Discord Bot…");

            // ⭐ KORREKTUR: Intents beim Login setzen!
            GatewayDiscordClient gateway = DiscordClient.create(token).gateway()
                    .setEnabledIntents(IntentSet.of(
                            Intent.GUILD_MESSAGES,
                            Intent.MESSAGE_CONTENT,
                            Intent.GUILD_VOICE_STATES,
                            Intent.GUILDS,
                            Intent.GUILD_MEMBERS
                    ))
                    .login()
                    .block();

            if (gateway == null) {
                logger.error("Discord Login fehlgeschlagen");
                return null;
            }

            // ⭐ BEHALTEN: Service Registrierung erfolgt DIREKT HIER
            joinToCreateService.initConfigs();
            joinToCreateService.register(gateway);
            guildEventHandler.register(gateway);
            commandService.register(gateway);
            guildService.syncGuilds(gateway).subscribe();

            statusService.setRunning(true);
            logger.info("Discord Bot erfolgreich gestartet");
            return gateway;
        } catch (Exception e) {
            logger.error("Fehler beim Start des DiscordBots", e);
            return null;
        }
    }
}