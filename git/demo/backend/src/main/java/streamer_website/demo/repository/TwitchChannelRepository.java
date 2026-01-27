package streamer_website.demo.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import streamer_website.demo.entity.twitch.TwitchChannel;

import java.util.Optional;

public interface TwitchChannelRepository extends JpaRepository<TwitchChannel, Long> {

    Optional<TwitchChannel> findByTwitchUserId(String twitchUserId);
}
