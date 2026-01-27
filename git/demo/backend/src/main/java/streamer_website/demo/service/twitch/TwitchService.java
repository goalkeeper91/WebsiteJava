package streamer_website.demo.service.twitch;

import com.fasterxml.jackson.databind.JsonNode;
import jakarta.annotation.PostConstruct;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import org.springframework.web.reactive.function.client.WebClient;
import streamer_website.demo.dto.twitch.TwitchUser;
import streamer_website.demo.entity.twitch.TwitchChannelStats;
import streamer_website.demo.repository.TwitchChannelStatsRepository;

import java.time.Instant;
import java.util.Optional;

@Service
@RequiredArgsConstructor
@Slf4j
public class TwitchService {

    @Value("${twitch.clientId}")
    private String clientId;

    @Value("${twitch.apiBaseUrl}")
    private String baseUri;

    private final TwitchTokenService tokenService;
    public final TwitchChannelStatsRepository twitchChannelStatsRepository;
    private final WebClient.Builder webClientBuilder;

    private WebClient twitchApiClient;

    @PostConstruct
    private void init() {
        this.twitchApiClient = webClientBuilder.baseUrl(baseUri).build();
    }

    // ---------------- USER / BOT TOKEN USAGE ----------------

    private String getAccessTokenForUser(String username) {
        return tokenService.getValidUserToken(username).getAccessToken();
    }

    // ---------------- TWITCH API CALLS ----------------

    public TwitchUser getUserInfo(String username) {
        String accessToken = getAccessTokenForUser(username);
        if (accessToken == null) {
            throw new IllegalStateException("Kein Access-Token für Benutzer " + username + " vorhanden.");
        }

        JsonNode userData = twitchApiClient.get()
                .uri("/helix/users?login=" + username)
                .header("Authorization", "Bearer " + accessToken)
                .header("Client-Id", clientId)
                .retrieve()
                .bodyToMono(JsonNode.class)
                .block();

        if (userData == null || !userData.has("data") || userData.get("data").isEmpty()) {
            throw new IllegalStateException("Konnte Twitch-Userinfo nicht abrufen für " + username);
        }

        JsonNode userNode = userData.get("data").get(0);

        return new TwitchUser(
                userNode.get("id").asText(),
                userNode.get("login").asText(),
                userNode.get("display_name").asText(),
                userNode.has("email") ? userNode.get("email").asText() : null,
                userNode.has("description") ? userNode.get("description").asText() : null,
                userNode.has("profile_image_url") ? userNode.get("profile_image_url").asText() : null,
                userNode.has("offline_image_url") ? userNode.get("offline_image_url").asText() : null,
                userNode.has("broadcaster_type") ? userNode.get("broadcaster_type").asText() : null,
                userNode.has("view_count") ? userNode.get("view_count").asInt() : 0,
                Instant.parse(userNode.get("created_at").asText())
        );
    }

    public boolean isLive(String username) {
        try {
            String accessToken = getAccessTokenForUser(username);
            String userId = getUserInfo(username).id();

            JsonNode streamData = twitchApiClient.get()
                    .uri("/helix/streams?user_id=" + userId)
                    .header("Authorization", "Bearer " + accessToken)
                    .header("Client-Id", clientId)
                    .retrieve()
                    .bodyToMono(JsonNode.class)
                    .block();

            return streamData != null && !streamData.get("data").isEmpty();
        } catch (Exception e) {
            log.error("Fehler beim Prüfen des Live-Status für {}: {}", username, e.getMessage());
            return false;
        }
    }

    public TwitchChannelStats fetchAndSaveChannelStats(String username) {
        String accessToken = getAccessTokenForUser(username);
        TwitchUser user = getUserInfo(username);

        JsonNode userResponse = twitchApiClient.get()
                .uri("/helix/users?id=" + user.id())
                .header("Authorization", "Bearer " + accessToken)
                .header("Client-Id", clientId)
                .retrieve()
                .bodyToMono(JsonNode.class)
                .block();

        if (userResponse == null || userResponse.get("data").isEmpty()) {
            throw new IllegalStateException("Konnte User-Daten von Twitch nicht abrufen.");
        }

        int followers = Optional.ofNullable(getFollowerCount(user.id(), accessToken)).orElse(0);

        return updateChannelStats(user, followers);
    }

    public Integer getFollowerCount(String broadcasterId, String accessToken) {
        try {
            JsonNode response = twitchApiClient.get()
                    .uri(uriBuilder -> uriBuilder.path("/helix/channels/followers")
                            .queryParam("broadcaster_id", broadcasterId)
                            .build())
                    .header("Authorization", "Bearer " + accessToken)
                    .header("Client-Id", clientId)
                    .retrieve()
                    .bodyToMono(JsonNode.class)
                    .block();

            return response != null && response.has("total") ? response.get("total").asInt() : 0;
        } catch (Exception e) {
            log.warn("Fehler beim Abrufen der Follower für {}: {}", broadcasterId, e.getMessage());
            return 0;
        }
    }

    private TwitchChannelStats updateChannelStats(TwitchUser user, int followers) {
        TwitchChannelStats stats = twitchChannelStatsRepository.findByTwitchUserId(user.id())
                .orElseGet(() -> {
                    TwitchChannelStats newStats = new TwitchChannelStats();
                    newStats.setTwitchUserId(user.id());
                    return newStats;
                });

        stats.setDisplayName(user.username());
        stats.setDescription(user.description());
        stats.setProfileImageUrl(user.profileImageUrl());
        stats.setOfflineImageUrl(user.offlineImageUrl());
        stats.setBroadcasterType(user.broadcasterType());
        stats.setViewCount(user.viewCount());
        stats.setFollowers(followers);
        stats.setAccountCreatedAt(user.createdAt());

        return twitchChannelStatsRepository.save(stats);
    }
}
