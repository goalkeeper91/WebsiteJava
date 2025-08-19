package streamer_website.demo.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import streamer_website.demo.entity.twitch.TwitchCommand;

import java.util.Optional;

public interface TwitchCommandRepository extends JpaRepository<TwitchCommand, Long> {
    Optional<TwitchCommand> findByTriggerIgnoreCase(String trigger);
    void deleteByTrigger(String trigger);
}
