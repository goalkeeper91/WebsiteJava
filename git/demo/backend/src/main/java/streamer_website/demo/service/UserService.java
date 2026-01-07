package streamer_website.demo.service;

import org.springframework.stereotype.Service;
import streamer_website.demo.dto.TwitchUser;
import streamer_website.demo.entity.User;
import streamer_website.demo.repository.UserRepository;

@Service
public class UserService {

    private final UserRepository userRepository;

    private static final String ADMIN_TWITCH_ID = "727153297";

    public UserService(UserRepository userRepository) {
        this.userRepository = userRepository;
    }

    public void createOrUpdate(TwitchUser twitchUser) {
        User user = userRepository.findById(twitchUser.id())
                .orElseGet(() -> {
                    User newUser = new User();
                    newUser.setTwitchId(twitchUser.id());
                    newUser.setAdmin(false);
                    return newUser;
                });

        user.setUsername(twitchUser.username());
        user.setEmail(twitchUser.email());

        if (twitchUser.id().equals(ADMIN_TWITCH_ID)) {
            user.setAdmin(true);
        }

        userRepository.save(user);
    }
}
