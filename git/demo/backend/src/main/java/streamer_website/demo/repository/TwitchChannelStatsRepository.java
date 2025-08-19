package streamer_website.demo.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import streamer_website.demo.entity.twitch.TwitchChannelStats;

import java.util.Optional;

public interface TwitchChannelStatsRepository extends JpaRepository<TwitchChannelStats, Long> {
    Optional<TwitchChannelStats> findTopByOrderByFetchedAtDesc();

    Optional<TwitchChannelStats> findByTwitchUserId(String twitchUserId);
}
