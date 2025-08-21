package streamer_website.demo.service.twitch;

import com.fasterxml.jackson.databind.JsonNode;
import jakarta.annotation.PostConstruct;
import lombok.Getter;
import lombok.RequiredArgsConstructor;
import lombok.Setter;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import org.springframework.web.reactive.function.client.WebClient;
import streamer_website.demo.dto.TwitchTokenResponse;
import streamer_website.demo.dto.TwitchUser;
import streamer_website.demo.entity.twitch.TwitchAuthToken;
import streamer_website.demo.entity.twitch.TwitchChannelStats;
import streamer_website.demo.repository.TwitchChannelStatsRepository;

import java.time.Instant;
import java.util.Arrays;
import java.util.List;
import java.util.Optional;

@Service
@Slf4j
public class TwitchService {

    @Value("${twitch.clientId}")
    private String clientId;

    @Value("${twitch.apiBaseUrl}")
    private String baseUri;

    private final TwitchTokenService tokenService;
    private final TwitchChannelStatsRepository twitchChannelStatsRepository;
    private WebClient twitchApiClient;

    private final WebClient.Builder webClientBuilder;

    public TwitchService(TwitchTokenService tokenService,
                         TwitchChannelStatsRepository twitchChannelStatsRepository,
                         WebClient.Builder webClientBuilder) {
        this.tokenService = tokenService;
        this.twitchChannelStatsRepository = twitchChannelStatsRepository;
        this.webClientBuilder = webClientBuilder;
    }

    @PostConstruct
    private void init() {
        this.twitchApiClient = webClientBuilder.baseUrl(baseUri).build();
    }

    public TwitchTokenResponse exchangeCodeForAccessToken(String code, boolean forBot) {
        TwitchAuthToken token = tokenService.exchangeCodeForToken(code, forBot);

        TwitchTokenResponse response = new TwitchTokenResponse();
        response.setAccessToken(token.getAccessToken());
        response.setRefreshToken(token.getRefreshToken());
        response.setExpiresIn(token.getExpiresIn());

        List<String> scopes = Arrays.asList(token.getScope().split(" "));
        response.setScope(scopes);

        return response;
    }

    public boolean isLive(String username) {
        try {
            String token = tokenService.getUserAccessToken(username);

            String userId = getFromTwitch("/helix/users?login=" + username, token)
                    .get("data").get(0).get("id").asText();

            JsonNode streamData = getFromTwitch("/helix/streams?user_id=" + userId, token);

            return !streamData.get("data").isEmpty();
        } catch (Exception e) {
            log.error("Failed to check Twitch live status for {}: {}", username, e.getMessage());
            return false;
        }
    }

    public TwitchUser getUserInfo(String username) {
        // 1️⃣ Holen des Tokens
        String accessToken = tokenService.getUserAccessToken(username);
        if (accessToken == null) {
            throw new RuntimeException("Kein Access-Token für Benutzer " + username + " gefunden. OAuth-Code nötig.");
        }

        // 2️⃣ Abruf der Twitch-Userinfos
        JsonNode userData = twitchApiClient.get()
                .uri("https://api.twitch.tv/helix/users?login=" + username)
                .header("Authorization", "Bearer " + accessToken)
                .header("Client-Id", clientId)
                .retrieve()
                .bodyToMono(JsonNode.class)
                .block();

        if (userData == null || !userData.has("data") || userData.get("data").isEmpty()) {
            throw new RuntimeException("Konnte Twitch-Userinfo nicht abrufen für " + username);
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

    public TwitchChannelStats fetchAndSaveChannelStats(String accessToken, String twitchUserId) {
        JsonNode userResponse = getFromTwitch("/helix/users?id=" + twitchUserId, accessToken);
        if (userResponse == null || userResponse.get("data").isEmpty()) {
            throw new RuntimeException("Could not retrieve user data from Twitch.");
        }

        JsonNode userNode = userResponse.get("data").get(0);

        TwitchUser twitchUser = new TwitchUser(
                userNode.get("id").asText(),
                userNode.get("login").asText(),
                userNode.get("display_name").asText(),
                userNode.get("email").asText(null),
                userNode.get("description").asText(null),
                userNode.get("profile_image_url").asText(null),
                userNode.get("offline_image_url").asText(null),
                userNode.get("broadcaster_type").asText(null),
                userNode.get("view_count").asInt(0),
                Instant.parse(userNode.get("created_at").asText())
        );

        int followers = Optional.ofNullable(getFollowerCount(twitchUserId, accessToken))
                .orElse(0);

        return updateChannelStats(twitchUser, followers);
    }


    public Integer getFollowerCount(String broadcasterId, String accessToken) {
        try {
            FollowersResponse response = twitchApiClient.get()
                    .uri(uriBuilder -> uriBuilder.path("/helix/channels/followers")
                            .queryParam("broadcaster_id", broadcasterId)
                            .build())
                    .header("Authorization", "Bearer " + accessToken)
                    .header("Client-ID", clientId)
                    .retrieve()
                    .bodyToMono(FollowersResponse.class)
                    .block();

            return response != null ? response.getTotal() : 0;
        } catch (Exception e) {
            log.warn("Failed to fetch follower count for broadcaster {}", broadcasterId, e);
            return 0;
        }
    }

    private JsonNode getFromTwitch(String uri, String accessToken) {
        return twitchApiClient.get()
                .uri(uri)
                .header("Client-ID", clientId)
                .header("Authorization", "Bearer " + accessToken)
                .retrieve()
                .bodyToMono(JsonNode.class)
                .block();
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

    public Optional<TwitchChannelStats> getLatestChannelStats(String username) {
        return twitchChannelStatsRepository.findByTwitchUserId(
                twitchChannelStatsRepository.findIdByDisplayName(username)
                        .orElse(null) // Optional-Anpassung je nach Repo-Methode
        );
    }

    public Optional<TwitchChannelStats> refreshAndSaveChannelStats(String username) {
        try {
            String token = tokenService.getUserAccessToken(username);
            TwitchUser user = getUserInfo(username);
            TwitchChannelStats stats = fetchAndSaveChannelStats(token, user.id());
            return Optional.of(stats);
        } catch (Exception e) {
            log.error("Failed to refresh Twitch stats for {}: {}", username, e.getMessage());
            return Optional.empty();
        }
    }

    @Setter
    @Getter
    public static class FollowersResponse {
        private int total;
    }
}
