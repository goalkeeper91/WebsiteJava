package streamer_website.demo.service.discord;

import discord4j.core.GatewayDiscordClient;
import discord4j.core.object.entity.Guild;
import discord4j.rest.util.Image;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import reactor.core.publisher.Mono;
import streamer_website.demo.entity.discord.DiscordGuild;
import streamer_website.demo.entity.discord.DiscordGuildRole;
import streamer_website.demo.repository.DiscordGuildRepository;
import streamer_website.demo.repository.DiscordGuildRoleRepository;

import java.util.List;
import java.util.Optional;

@Service
public class GuildService {

    private final DiscordGuildRepository discordGuildRepository;
    private final DiscordGuildRoleRepository discordGuildRoleRepository;

    public GuildService(DiscordGuildRepository discordGuildRepository,
                        DiscordGuildRoleRepository discordGuildRoleRepository) {
        this.discordGuildRepository = discordGuildRepository;
        this.discordGuildRoleRepository = discordGuildRoleRepository;
    }

    @Transactional
    public void saveOrUpdateGuild(DiscordGuild guild) {
        discordGuildRepository.save(guild);
    }

    public List<DiscordGuild> getAllGuilds() {
        return discordGuildRepository.findAll();
    }

    public Optional<DiscordGuild> getGuild(Long id) {
        return discordGuildRepository.findById(id);
    }

    public Mono<Void> syncGuilds(GatewayDiscordClient gateway) {
        return gateway.getGuilds()
                .flatMap(guild -> Mono.fromRunnable(() -> handleGuild(guild)))
                .then();
    }

    public Mono<Void> handleGuild(Guild guildData) {
        return Mono.fromCallable(() -> {
            DiscordGuild guild = getGuild(guildData.getId().asLong())
                    .orElse(new DiscordGuild());

            guild.setId(guildData.getId().asLong());
            guild.setName(guildData.getName());
            guild.setIconUrl(guildData.getIconUrl(Image.Format.PNG).map(Object::toString).orElse(null));
            guild.setBotAdded(true);

            saveOrUpdateGuild(guild);
            return Mono.empty();
        }).then();
    }

    public void deleteGuild(Long guildId) {
        List<DiscordGuildRole> roles = discordGuildRoleRepository.findByGuildId(guildId);
        discordGuildRoleRepository.deleteAll(roles);
        discordGuildRepository.deleteById(guildId);
    }
}
