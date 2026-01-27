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

    /**
     * Prüft, ob der Twitch-Stream für den angegebenen Benutzer live ist.
     */
    public boolean checkLiveStatus(String username) {
        if (username == null || username.isBlank()) {
            throw new IllegalArgumentException("Username must be provided");
        }
        return twitchService.isLive(username);
    }

    /**
     * Liefert die aktuell gespeicherten Twitch-Statistiken für den Benutzer.
     */
    public Optional<TwitchChannelStats> getStats(String username) {
        return twitchService.twitchChannelStatsRepository
                .findByTwitchUserId(
                        twitchService.twitchChannelStatsRepository.findIdByDisplayName(username)
                                .orElse(null)
                );
    }

    /**
     * Holt die aktuellen Twitch-Statistiken von der API und speichert sie in der DB.
     */
    public Optional<TwitchChannelStats> refreshStats(String username) {
        try {
            return Optional.of(twitchService.fetchAndSaveChannelStats(username));
        } catch (Exception e) {
            return Optional.empty();
        }
    }
}
