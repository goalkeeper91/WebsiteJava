package streamer_website.demo.service.twitch;

import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.scheduling.annotation.Async;
import org.springframework.stereotype.Service;
import streamer_website.demo.client.TwitchBot;
import streamer_website.demo.config.TwitchBotConfig;
import streamer_website.demo.dto.BotStatusDto;

import java.time.Duration;
import java.time.Instant;
import java.util.Collections;
import java.util.List;

@Service
@RequiredArgsConstructor
public class TwitchBotManagerService {

    private TwitchBot twitchBot;

    private final TwitchTokenService twitchTokenService;
    private final TwitchCommandService twitchCommandService;
    private final TwitchBotConfig twitchBotConfig;
    private static final Logger logger = LoggerFactory.getLogger(TwitchBotManagerService.class);

    @Async
    public void startBot() {
        if (twitchBot != null && twitchBot.isRunning()) {
            return;
        }

        String channelName = "goalkeeper91";
        String userId = twitchBotConfig.getBotUserId();

        twitchBot = new TwitchBot(channelName, twitchCommandService, twitchTokenService);
        twitchBot.start(userId);
    }

    public void stopBot() {
        if (twitchBot != null) {
            twitchBot.stop(); // stoppt die IRC-Verbindung
        }
    }

    public BotStatusDto getBotStatus() {
        var token = twitchTokenService.findBotToken();
        boolean tokenPresent = token != null;

        Instant expiresAt = null;
        if (tokenPresent) {
            expiresAt = token.getCreatedAt().plusSeconds(token.getExpiresIn());
        }

        boolean running = twitchBot != null && twitchBot.isRunning();
        String primaryChannel = twitchBot != null ? twitchBot.getChannelName() : null;

        List<String> activeChannels = running && twitchBot.getClient() != null
                ? twitchBot.getClient().getChat().getChannels()
                .stream()
                .sorted()
                .toList()
                : Collections.emptyList();

        Long uptimeSeconds = null;

        if (running && twitchBot.getStartTime() != null) {
            uptimeSeconds = Duration.between(twitchBot.getStartTime(), Instant.now()).getSeconds();
        }

        return new BotStatusDto(
                running,
                tokenPresent,
                expiresAt,
                primaryChannel,
                activeChannels,
                uptimeSeconds,
                "1.0.0" // Version
        );
    }

}
