package streamer_website.demo.service.twitch;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import streamer_website.demo.dto.twitch.*;
import streamer_website.demo.entity.twitch.TwitchAuthToken;
import streamer_website.demo.repository.TwitchAuthTokenRepository;
import streamer_website.demo.security.EncryptionUtils;

import java.time.Duration;
import java.time.Instant;
import java.util.List;

@Slf4j
@Service
@RequiredArgsConstructor
public class StreamService {

    private final TwitchApiService twitchApiService;
    private final TwitchAuthTokenRepository tokenRepository;

    public StreamInfoDto getStreamInfo(String twitchUserId) {
        String accessToken = getAccessToken(twitchUserId);
        return twitchApiService.getStreamInfo(twitchUserId, accessToken);
    }

    public StreamInfoDto updateStreamInfo(
            String twitchUserId,
            String title,
            String gameId
    ) {
        String accessToken = getAccessToken(twitchUserId);

        twitchApiService.updateChannelInfo(twitchUserId, accessToken, title, gameId);

        return twitchApiService.getStreamInfo(twitchUserId, accessToken);
    }

    public LiveStreamDto getLiveStatus(String twitchUserId) {
        String accessToken = getAccessToken(twitchUserId);
        return twitchApiService.getLiveStream(twitchUserId, accessToken);
    }

    public DashboardStatsDto getDashboardStats(String twitchUserId) {
        String accessToken = getAccessToken(twitchUserId);

        LiveStreamDto liveStream = twitchApiService.getLiveStream(twitchUserId, accessToken);

        int followerCount = twitchApiService.getFollowerCount(twitchUserId, accessToken);

        String uptime = null;
        if (liveStream.isLive() && liveStream.getStartedAt() != null) {
            uptime = calculateUptime(liveStream.getStartedAt());
        }

        // TODO: Subscriber Count über EventSub oder separate API
        int subscriberCount = 0;

        // TODO: Stats aus DB holen (follows today, subs this week)
        int followsToday = 0;
        int subsThisWeek = 0;
        double avgViewers = 0.0;

        return DashboardStatsDto.builder()
                .isLive(liveStream.isLive())
                .currentViewers(liveStream.getViewerCount())
                .followerCount(followerCount)
                .subscriberCount(subscriberCount)
                .uptime(uptime)
                .followsToday(followsToday)
                .subsThisWeek(subsThisWeek)
                .avgViewers(avgViewers)
                .build();
    }

    public List<CategoryDto> searchCategories(String twitchUserId, String query) {
        String accessToken = getAccessToken(twitchUserId);
        return twitchApiService.searchCategories(query, accessToken);
    }

    private String getAccessToken(String twitchUserId) {
        TwitchAuthToken token = tokenRepository
                .findByTwitchUserId(twitchUserId)
                .orElseThrow(() -> new IllegalStateException(
                        "Kein Access Token gefunden. Bitte neu einloggen."
                ));

        return EncryptionUtils.decrypt(token.getAccessToken());
    }

    private String calculateUptime(String startedAtIso) {
        try {
            Instant startTime = Instant.parse(startedAtIso);
            Instant now = Instant.now();
            Duration duration = Duration.between(startTime, now);

            long hours = duration.toHours();
            long minutes = duration.toMinutes() % 60;

            return String.format("%d:%02d", hours, minutes);
        } catch (Exception e) {
            log.error("Fehler beim Berechnen der Uptime", e);
            return "0:00";
        }
    }
}
