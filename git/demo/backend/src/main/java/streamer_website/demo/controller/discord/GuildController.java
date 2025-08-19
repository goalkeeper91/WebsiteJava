package streamer_website.demo.controller.discord;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.*;
import org.springframework.web.bind.annotation.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import streamer_website.demo.client.DiscordBot;
import streamer_website.demo.entity.discord.DiscordGuild;
import streamer_website.demo.entity.discord.DiscordGuildRole;
import streamer_website.demo.repository.DiscordGuildRepository;
import streamer_website.demo.repository.DiscordGuildRoleRepository;

import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api/discord/guild")
public class GuildController {

    private final DiscordGuildRepository discordGuildRepository;
    private final DiscordGuildRoleRepository discordGuildRoleRepository;
    private final DiscordBot discordBot;
    private static final Logger logger = LoggerFactory.getLogger(GuildController.class);

    @Value("${discord.bot.client-id}")
    private String clientId;

    @Value("${discord.bot.permissions}")
    private String permissions;

    public GuildController(DiscordGuildRepository discordGuildRepository,
                           DiscordGuildRoleRepository discordGuildRoleRepository,
                           DiscordBot discordBot) {
        this.discordGuildRepository = discordGuildRepository;
        this.discordGuildRoleRepository = discordGuildRoleRepository;
        this.discordBot = discordBot;
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
        return discordGuildRepository.findAll();
    }

    @GetMapping("/{guildId}")
    public ResponseEntity<DiscordGuild> getGuild(@PathVariable Long guildId) {
        return discordGuildRepository.findById(guildId)
                .map(ResponseEntity::ok)
                .orElse(ResponseEntity.notFound().build());
    }

    @DeleteMapping("/{guildId}")
    public ResponseEntity<String> disconnectGuild(@PathVariable Long guildId) {
        var guild = discordGuildRepository.findById(guildId).orElse(null);
        if (guild == null) {
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body("Guild nicht gefunden");
        }

        try {
            List<DiscordGuildRole> roles = discordGuildRoleRepository.findByGuildId(guildId);
            discordGuildRoleRepository.deleteAll(roles);
            discordGuildRepository.delete(guild);
            logger.info("Guild {} und zugehörige Rollen gelöscht", guild.getName());

            return ResponseEntity.ok("Guild erfolgreich getrennt");
        } catch (Exception e) {
            logger.error("Fehler beim Disconnect der Guild {}", guildId, e);
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body("Fehler beim Disconnect");
        }
    }

    @PostMapping("/sync")
    public List<DiscordGuild> syncGuild() {
        discordBot.syncGuilds();
        return discordGuildRepository.findAll();
    }
}
