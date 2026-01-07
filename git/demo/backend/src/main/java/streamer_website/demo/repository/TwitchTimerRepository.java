package streamer_website.demo.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import streamer_website.demo.entity.twitch.TwitchTimer;

import java.util.List;

public interface TwitchTimerRepository extends JpaRepository<TwitchTimer, Long> {

    List<TwitchTimer> findByTwitchUserId(String twitchUserId);

    List<TwitchTimer> findByIsEnabledTrue();

}
