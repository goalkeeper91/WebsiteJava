package streamer_website.demo.service.discord;

import discord4j.core.GatewayDiscordClient;
import discord4j.core.object.entity.Guild;
import discord4j.rest.util.Image;
import org.springframework.stereotype.Service;
import reactor.core.publisher.Flux;
import streamer_website.demo.entity.discord.DiscordGuild;
import streamer_website.demo.repository.DiscordGuildRepository;

import java.time.Instant;
import java.util.List;
import java.util.Objects;

@Service
public class GuildService {

    private final GatewayDiscordClient gateway;
    private final DiscordGuildRepository discordGuildRepository;

    public GuildService(GatewayDiscordClient gateway,
                        DiscordGuildRepository discordGuildRepository) {
        this.gateway = gateway;
        this.discordGuildRepository = discordGuildRepository;
    }

    public void syncGuilds() {
        Objects.requireNonNull(gateway.getGuilds().collectList().block()).forEach(guildData -> {
            DiscordGuild guild = discordGuildRepository.findById(guildData.getId().asLong())
                    .orElse(new DiscordGuild());

            guild.setId(guildData.getId().asLong());
            guild.setName(guildData.getName());
            guild.setIconUrl(guildData.getIconUrl(Image.Format.PNG).map(Object::toString).orElse(null));
            guild.setUpdatedAt(Instant.now());

            discordGuildRepository.save(guild);
        });
    }

    public List<DiscordGuild> getGuilds() {
        return discordGuildRepository.findAll();
    }
}
