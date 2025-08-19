package streamer_website.demo.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import streamer_website.demo.entity.discord.DiscordVoiceChannel;

public interface DiscordVoiceChannelRepository extends JpaRepository<DiscordVoiceChannel, Long> {
}
