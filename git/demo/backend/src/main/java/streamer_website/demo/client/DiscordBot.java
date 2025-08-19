package streamer_website.demo.client;

import discord4j.core.DiscordClient;
import discord4j.core.GatewayDiscordClient;
import discord4j.core.event.domain.guild.GuildCreateEvent;
import discord4j.core.event.domain.message.MessageCreateEvent;
import discord4j.rest.util.Image;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.stereotype.Component;
import streamer_website.demo.controller.discord.GuildController;
import streamer_website.demo.entity.discord.DiscordGuild;
import streamer_website.demo.repository.DiscordGuildRepository;
import streamer_website.demo.service.discord.CommandService;
import streamer_website.demo.service.discord.StatusService;

import java.time.Instant;

@Component
public class DiscordBot {

    @Value("${discord.bot.token}")
    private String token;

    private final CommandService commandService;
    private final StatusService statusService;
    private final DiscordGuildRepository discordGuildRepository;

    private static final Logger logger = LoggerFactory.getLogger(DiscordBot.class);

    private GatewayDiscordClient gateway;

    public DiscordBot(CommandService commandService,
                      StatusService statusService,
                      DiscordGuildRepository discordGuildRepository) {
        this.commandService = commandService;
        this.statusService = statusService;
        this.discordGuildRepository = discordGuildRepository;
    }

    @Bean
    public GatewayDiscordClient gatewayDiscordClient() {
        if (token == null || token.isEmpty()) {
            throw new IllegalStateException("Discord token is not set!");
        }

        DiscordClient client = DiscordClient.create(token);
        gateway = client.login().block();

        if (gateway == null) {
            throw new IllegalStateException("Failed to login to Discord. Check your token.");
        }

        statusService.setRunning(true);

        setupListeners();

        return gateway;
    }

    private void setupListeners() {
        gateway.on(MessageCreateEvent.class, event -> {
            commandService.handle(event);
            return null;
        }).subscribe();

        gateway.on(GuildCreateEvent.class, event -> {
            saveGuild(event.getGuild());
            return null;
        }).subscribe();

        gateway.onDisconnect()
                .doOnTerminate(() -> statusService.setRunning(false))
                .subscribe();
    }

    private void saveGuild(discord4j.core.object.entity.Guild guildData) {
        DiscordGuild guild = new DiscordGuild();
        guild.setId(Long.parseLong(guildData.getId().asString()));
        guild.setName(guildData.getName());
        guild.setIconUrl(guildData.getIconUrl(Image.Format.PNG).map(Object::toString).orElse(null));
        guild.setCreatedAt(Instant.now());
        guild.setUpdatedAt(Instant.now());

        discordGuildRepository.save(guild);
        System.out.println("Guild gespeichert: " + guild.getName());
    }

    public void syncGuilds() {
        if (gateway == null) {
            throw new IllegalStateException("GatewayDiscordClient ist nicht initialisiert!");
        }

        gateway.getGuilds()
                .collectList()
                .subscribe(guilds -> guilds.forEach(this::saveGuild));
    }
}
