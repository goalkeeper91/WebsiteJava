package streamer_website.demo.handler.discord;

import discord4j.core.GatewayDiscordClient;
import discord4j.core.event.domain.guild.GuildCreateEvent;
import lombok.Setter;
import org.springframework.stereotype.Component;
import streamer_website.demo.service.discord.GuildService;

@Setter
@Component
public class GuildEventHandler {

    private GuildService guildService;

    public GuildEventHandler(GuildService guildService) {
        this.guildService = guildService;
    }

    public void register(GatewayDiscordClient gateway) {
        gateway.on(GuildCreateEvent.class, event -> {
            guildService.handleGuild(event.getGuild());
            return null;
        }).subscribe();
    }
}

