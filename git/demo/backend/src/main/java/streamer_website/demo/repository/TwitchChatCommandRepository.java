package streamer_website.demo.repository;

import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import streamer_website.demo.entity.twitch.ChatCommand;

import java.util.List;
import java.util.Optional;

public interface TwitchChatCommandRepository extends JpaRepository<ChatCommand, Long> {

    List<ChatCommand> findAllByChannelId(String channelId);

    List<ChatCommand> findAllByChannelIdAndEnabledTrue(String channelId);

    List<ChatCommand> findAllByChannelIdAndEnabled(String channelId, Boolean enabled);

    Optional<ChatCommand> findByChannelIdAndTriggerIgnoreCase(
            String channelId,
            String trigger
    );

    Optional<ChatCommand> findByChannelIdAndTriggerIgnoreCaseAndEnabledTrue(
            String channelId,
            String trigger
    );

    Optional<ChatCommand> findByIdAndChannelId(Long id, String channelId);

    boolean existsByChannelIdAndTriggerIgnoreCase(
            String channelId,
            String trigger
    );

    void deleteByChannelIdAndTriggerIgnoreCase(
            String channelId,
            String trigger
    );

    List<ChatCommand> findByChannelIdAndTriggerContainingIgnoreCase(
            String channelId,
            String search
    );

    Page<ChatCommand> findByChannelIdAndTriggerContainingIgnoreCase(
            String channelId,
            String search,
            Pageable pageable
    );

    Page<ChatCommand> findByChannelId(
            String channelId,
            Pageable pageable
    );

    Page<ChatCommand> findByChannelIdAndEnabled(
            String channelId,
            Boolean enabled,
            Pageable pageable
    );
}
