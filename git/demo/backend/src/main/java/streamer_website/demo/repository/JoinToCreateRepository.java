package streamer_website.demo.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import streamer_website.demo.entity.discord.JoinToCreateConfig;

public interface JoinToCreateRepository extends JpaRepository<JoinToCreateConfig, Long> {
}
