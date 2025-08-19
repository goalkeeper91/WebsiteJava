package streamer_website.demo.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import streamer_website.demo.entity.discord.DiscordGuild;

public interface DiscordGuildRepository extends JpaRepository<DiscordGuild, Long> {
}
