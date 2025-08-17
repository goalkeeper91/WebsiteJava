package streamer_website.demo.service.discord;

import discord4j.core.GatewayDiscordClient;
import discord4j.core.object.entity.Guild;
import org.springframework.stereotype.Service;
import reactor.core.publisher.Flux;

import java.util.List;

@Service
public class GuildService {

    private final GatewayDiscordClient gateway;

    public GuildService(GatewayDiscordClient gateway) {
        this.gateway = gateway;
    }

    public List<String> getGuilds() {
        Flux<Guild> guilds = gateway.getGuilds();
        return guilds.map(Guild::getName).collectList().block();
    }
}
