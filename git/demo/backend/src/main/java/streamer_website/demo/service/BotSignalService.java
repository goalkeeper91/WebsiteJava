package streamer_website.demo.service;


import lombok.RequiredArgsConstructor;
import org.springframework.data.redis.core.StringRedisTemplate;
import org.springframework.stereotype.Service;

@Service
@RequiredArgsConstructor
public class BotSignalService {
    private final StringRedisTemplate redisTemplate;
    private static final String CHANNEL = "twitch_bot_signal";

    public void sendBotJoinSignal(String twitchUserId) {
        redisTemplate.convertAndSend(CHANNEL, "JOIN:" + twitchUserId);
    }

    public void sendCommandUpdateSignal(String twitchUserId) {
        redisTemplate.convertAndSend(CHANNEL, "REFRESH_COMMANDS:" + twitchUserId);
    }
}
