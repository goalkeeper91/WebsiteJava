package streamer_website.demo.config;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;
import streamer_website.demo.repository.TwitchAuthTokenRepository;

@Component
public class TwitchBotConfig {

    private final TwitchAuthTokenRepository botRepo;

    @Autowired
    public TwitchBotConfig(TwitchAuthTokenRepository botRepo) {
        this.botRepo = botRepo;
    }

    public String getBotUserId() {
        var bot = botRepo.findTopByUserNameOrderByCreatedAtDesc("goalkeeper91_bot").orElseThrow(() -> new IllegalStateException("User not found"));
        return bot.getTwitchUserId();
    }
}
