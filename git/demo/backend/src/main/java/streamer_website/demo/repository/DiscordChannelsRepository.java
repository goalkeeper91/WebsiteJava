package streamer_website.demo.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import streamer_website.demo.entity.discord.DiscordChannels;

import java.util.List;

public interface DiscordChannelsRepository extends JpaRepository<DiscordChannels, Long> {

    List<DiscordChannels> findByGuildId(String guildId);

    List<DiscordChannels> findByDescriptionContainingIgnoreCase(String keyword);

    boolean existsByGuildIdAndChannelId(String guildId, String channelId);
}
