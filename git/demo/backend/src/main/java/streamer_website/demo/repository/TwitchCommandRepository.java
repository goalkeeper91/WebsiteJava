package streamer_website.demo.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import streamer_website.demo.entity.TwitchCommand;

import java.util.Optional;

public interface TwitchCommandRepository extends JpaRepository<TwitchCommand, Long> {
    Optional<TwitchCommand> findByTrigger(String trigger);
    void deleteByTrigger(String trigger);
}
