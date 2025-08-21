package streamer_website.demo.handler.twitch;

import org.springframework.stereotype.Component;
import streamer_website.demo.service.twitch.TwitchService;
import streamer_website.demo.entity.twitch.TwitchChannelStats;

import java.util.Optional;

@Component
public class TwitchHandler {

    private final TwitchService twitchService;

    public TwitchHandler(TwitchService twitchService) {
        this.twitchService = twitchService;
    }

    public boolean checkLiveStatus(String username) {
        if (username == null || username.isBlank()) {
            throw new IllegalArgumentException("Username must be provided");
        }
        return twitchService.isLive(username);
    }

    public Optional<TwitchChannelStats> getStats(String username) {
        return twitchService.getLatestChannelStats(username);
    }

    public Optional<TwitchChannelStats> refreshStats(String username) {
        return twitchService.refreshAndSaveChannelStats(username);
    }
}
