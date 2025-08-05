package streamer_website.demo.service;

import org.springframework.stereotype.Service;
import streamer_website.demo.dto.TwitchUser;
import streamer_website.demo.entity.User;
import streamer_website.demo.repository.UserRepository;

@Service
public class UserService {

    private final UserRepository userRepository;

    public UserService(UserRepository userRepository) {
        this.userRepository = userRepository;
    }

    public void createOrUpdate(TwitchUser twitchUser){
        User user = userRepository.findById(twitchUser.id())
                .orElse(new User());

        user.setTwitchId(twitchUser.id());
        user.setUsername(twitchUser.username());
        user.setEmail(twitchUser.email());
        user.setAdmin(true);

        userRepository.save(user);
    }
}
