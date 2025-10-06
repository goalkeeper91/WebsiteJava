package streamer_website.demo.controller.discord;

import discord4j.core.GatewayDiscordClient;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Lazy;
import org.springframework.http.*;
import org.springframework.web.bind.annotation.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import reactor.core.publisher.Mono;
import streamer_website.demo.entity.discord.DiscordGuild;
import streamer_website.demo.service.discord.GuildService;

import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api/discord/guild")
public class GuildController {

    private final GuildService guildService;
    private final GatewayDiscordClient gateway;
    private static final Logger logger = LoggerFactory.getLogger(GuildController.class);

    @Value("${discord.bot.client-id}")
    private String clientId;

    @Value("${discord.bot.permissions}")
    private String permissions;

    public GuildController(GuildService guildService, GatewayDiscordClient gateway) {
        this.guildService = guildService;
        this.gateway = gateway;
    }

    @GetMapping("/invite-link")
    public Map<String, String> getInviteLink() {
        String link = String.format(
                "https://discord.com/oauth2/authorize?client_id=%s&scope=bot%%20applications.commands&permissions=%s",
                clientId, permissions
        );
        return Map.of("inviteLink", link);
    }

    @GetMapping
    public List<DiscordGuild> getAllGuilds() {
        return guildService.getAllGuilds();
    }

    @GetMapping("/{guildId}")
    public ResponseEntity<DiscordGuild> getGuild(@PathVariable Long guildId) {
        return guildService.getGuild(guildId)
                .map(ResponseEntity::ok)
                .orElse(ResponseEntity.notFound().build());
    }

    @DeleteMapping("/{guildId}")
    public ResponseEntity<String> disconnectGuild(@PathVariable Long guildId) {
        try {
            guildService.deleteGuild(guildId);
            logger.info("Guild {} gelöscht", guildId);
            return ResponseEntity.ok("Guild erfolgreich getrennt");
        } catch (Exception e) {
            logger.error("Fehler beim Disconnect der Guild {}", guildId, e);
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body("Fehler beim Disconnect");
        }
    }

    @PostMapping("/sync")
    public Mono<ResponseEntity<String>> syncGuilds() {
        return guildService.syncGuilds(gateway) // <- Gateway direkt nutzen
                .then(Mono.just(ResponseEntity.ok("Guild-Synchronisierung gestartet")));
    }
}
