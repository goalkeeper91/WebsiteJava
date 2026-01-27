package streamer_website.demo.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import streamer_website.demo.entity.twitch.TwitchAuthToken;
import streamer_website.demo.entity.twitch.TwitchTokenType;

import java.util.Optional;

public interface TwitchAuthTokenRepository extends JpaRepository<TwitchAuthToken, Long> {
    Optional<TwitchAuthToken> findByTwitchUserId(String twitchUserId);
    Optional<TwitchAuthToken> findTopByUserNameOrderByCreatedAtDesc(String userName);
    Optional<TwitchAuthToken> findTopByOrderByCreatedAtDesc();
    Optional<TwitchAuthToken> findByUserNameAndTokenOwner(
            String userName,
            TwitchTokenType tokenOwner
    );
}
