package streamer_website.demo.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import streamer_website.demo.entity.discord.DiscordUser;

public interface DiscordUserRepository extends JpaRepository<DiscordUser, Long> {
}
