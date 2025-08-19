package streamer_website.demo.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import streamer_website.demo.entity.discord.DiscordGuildRole;

import java.util.List;

public interface DiscordGuildRoleRepository extends JpaRepository<DiscordGuildRole, Long> {
    List<DiscordGuildRole> findByGuildId(Long guildId);
}
