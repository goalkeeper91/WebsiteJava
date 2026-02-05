package streamer_website.demo.service;

import jakarta.transaction.Transactional;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import streamer_website.demo.dto.twitch.TwitchUser;
import streamer_website.demo.entity.User;
import streamer_website.demo.entity.twitch.TwitchChannel;
import streamer_website.demo.repository.TwitchChannelRepository;
import streamer_website.demo.repository.UserRepository;

import java.time.Instant;

@Service
public class UserService {

    private final UserRepository userRepository;
    private final TwitchChannelRepository twitchChannelRepository;
    private static final Logger logger = LoggerFactory.getLogger(UserService.class);

    @Value("${twitch.admin_twitch_id}")
    private String admin_twitch_id;

    public UserService(UserRepository userRepository, TwitchChannelRepository twitchChannelRepository) {
        this.userRepository = userRepository;
        this.twitchChannelRepository = twitchChannelRepository;
    }

    public User createOrUpdate(TwitchUser twitchUser) {
        User user = userRepository.findById(twitchUser.id())
                .orElseGet(() -> {
                    User newUser = new User();
                    newUser.setTwitchId(twitchUser.id());
                    newUser.setAdmin(false);
                    return newUser;
                });

        user.setUsername(twitchUser.username());
        user.setEmail(twitchUser.email());

        if (twitchUser.id().equals(admin_twitch_id)) {
            user.setAdmin(true);
        }

        return userRepository.save(user);
    }

    @Transactional
    public void syncTwitchChannel(TwitchUser twitchUser) {
        try {
            TwitchChannel channel = twitchChannelRepository.findByTwitchUserId(twitchUser.id())
                    .orElseGet(() -> {
                        return TwitchChannel.builder()
                                .twitchUserId(twitchUser.id())
                                .isActive(true)
                                .createdAt(Instant.now())
                                .updatedAt(Instant.now())
                                .build();
                    });

            channel.setUserName(twitchUser.login());
            channel.setUpdatedAt(Instant.now());

            twitchChannelRepository.saveAndFlush(channel);
            logger.info("Twitch Channel synchronisiert: {}", twitchUser.login());

        } catch (Exception e) {
            logger.error("Kritischer Fehler beim Sync des Twitch Channels für ID: " + twitchUser.id(), e);
            throw e;
        }
    }
}
