package streamer_website.demo.service.twitch;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.Getter;
import lombok.Setter;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.reactive.function.BodyInserters;
import org.springframework.web.reactive.function.client.WebClient;
import com.fasterxml.jackson.databind.JsonNode;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.MediaType;
import org.springframework.stereotype.Service;
import reactor.core.publisher.Mono;
import streamer_website.demo.dto.TwitchTokenResponse;
import streamer_website.demo.dto.TwitchUser;
import streamer_website.demo.entity.TwitchAuthToken;
import streamer_website.demo.entity.TwitchChannelStats;
import streamer_website.demo.repository.TwitchAuthTokenRepository;
import streamer_website.demo.repository.TwitchChannelStatsRepository;

import java.time.Instant;
import java.util.List;
import java.util.Objects;
import java.util.Optional;

@Service
public class TwitchService {

    public enum TokenType { BOT, USER }

    private final String clientId;
    private final String clientSecret;
    private final String username;
    private final String apiBaseUrl;
    private final String authBaseUrl;
    private final String redirectUri;

    private String botToken;
    private long botTokenExpiry;

    private final WebClient twitchApiClient;
    private final WebClient twitchAuthClient;

    @Autowired
    private TwitchAuthTokenRepository twitchAuthTokenRepository;

    @Autowired
    private TwitchChannelStatsRepository twitchChannelStatsRepository;

    private final ObjectMapper objectMapper;

    private static final Logger logger = LoggerFactory.getLogger(TwitchService.class);

    public TwitchService(
            WebClient.Builder webClientBuilder,
            @Value("${twitch.clientId}") String clientId,
            @Value("${twitch.clientSecret}") String clientSecret,
            @Value("${twitch.username}") String username,
            @Value("${twitch.apiBaseUrl}") String apiBaseUrl,
            @Value("${twitch.auth-base-url}") String authBaseUrl,
            @Value("${twitch.redirectUri}") String redirectUri,
            ObjectMapper objectMapper
                         ) {
        this.clientId = clientId;
        this.clientSecret = clientSecret;
        this.username = username;
        this.apiBaseUrl = apiBaseUrl;
        this.authBaseUrl = authBaseUrl;
        this.twitchApiClient = webClientBuilder.baseUrl(apiBaseUrl).build();
        this.twitchAuthClient = webClientBuilder.baseUrl(authBaseUrl).build();
        this.redirectUri = redirectUri;
        this.objectMapper = objectMapper;
    }

    public boolean isLive() {
        try {
            String token = getAccessToken(TokenType.USER);

            String userId = Objects.requireNonNull(twitchApiClient.get()
                            .uri("/helix/users?login=" + username)
                            .header("Client-ID", clientId)
                            .header("Authorization", "Bearer " + token)
                            .retrieve()
                            .bodyToMono(JsonNode.class)
                            .block())
                    .get("data").get(0).get("id").asText();

            // Check Live-Status
            JsonNode streamData = twitchApiClient.get()
                    .uri("/helix/streams?user_id=" + userId)
                    .header("Client-ID", clientId)
                    .header("Authorization", "Bearer " + token)
                    .retrieve()
                    .bodyToMono(JsonNode.class)
                    .block();

            assert streamData != null;
            return !streamData.get("data").isEmpty();
        } catch (Exception e) {
            logger.error("Failed to check Twitch live status", e);
            return false;
        }
    }

    public String getAccessToken(TokenType type) {
        if (type == TokenType.BOT) {
            return getBotToken().getAccessToken();
        } else {
            return getUserAccessToken();
        }
    }

    private String getUserAccessToken() {
        Optional<TwitchAuthToken> tokenOpt = twitchAuthTokenRepository.findTopByOrderByCreatedAtDesc();

        if (tokenOpt.isEmpty()) {
            throw new IllegalStateException("No user token found, user must authenticate");
        }

        TwitchAuthToken token = tokenOpt.get();

        Instant expiryTime = token.getCreatedAt().plusSeconds(token.getExpiresIn());

        if (Instant.now().isAfter(expiryTime.minusSeconds(60))) { // 60 Sekunden vorher erneuern
            try {
                TwitchTokenResponse refreshed = refreshAccessToken(token.getRefreshToken());
                return refreshed.getAccessToken();
            } catch (Exception e) {
                logger.error("Failed to refresh Twitch token", e);
                throw new IllegalStateException("Twitch token expired and refresh failed");
            }
        }

        return token.getAccessToken();
    }

    private TwitchTokenResponse getBotToken() {
        if (botToken != null && Instant.now().getEpochSecond() < botTokenExpiry - 60) {
            return new TwitchTokenResponse(botToken, null, botTokenExpiry, List.of(), "bearer");
        }

        JsonNode response = twitchAuthClient.post()
                .uri("/oauth2/token")
                .contentType(MediaType.APPLICATION_FORM_URLENCODED)
                .body(BodyInserters.fromFormData("client_id", clientId)
                        .with("client_secret", clientSecret)
                        .with("grant_type", "client_credentials"))
                .retrieve()
                .bodyToMono(JsonNode.class)
                .block();

        if (response == null || response.get("access_token") == null) {
            throw new RuntimeException("Failed to fetch bot token");
        }

        botToken = response.get("access_token").asText();
        botTokenExpiry = Instant.now().getEpochSecond() + response.get("expires_in").asLong();

        return new TwitchTokenResponse(
                botToken,
                null,
                response.get("expires_in").asLong(),
                objectMapper.convertValue(response.get("scope"), new TypeReference<List<String>>() {}),
                response.get("token_type").asText()
        );
    }

    public TwitchTokenResponse exchangeCodeForAccessToken(String code) throws JsonProcessingException {
        JsonNode response = twitchAuthClient.post()
            .uri("/oauth2/token")
            .contentType(MediaType.APPLICATION_FORM_URLENCODED)
            .body(BodyInserters.fromFormData("client_id", clientId)
                .with("client_secret", clientSecret)
                .with("code", code)
                .with("grant_type", "authorization_code")
                .with("redirect_uri", redirectUri))
            .retrieve()
            .bodyToMono(JsonNode.class)
            .block();

        if (response == null || response.get("access_token") == null) {
            throw new RuntimeException("Twitch token exchange failed");
        }

        List<String> scopes = objectMapper.convertValue(
                response.get("scope"), new TypeReference<List<String>>() {}
        );

        TwitchTokenResponse token = new TwitchTokenResponse(
                response.get("access_token").asText(),
                response.get("refresh_token").asText(),
                response.get("expires_in").asLong(),
                scopes,
                response.get("token_type").asText()
        );

        TwitchUser twitchUser = getUserInfo(token.getAccessToken());

        TwitchAuthToken twitchAuthToken = TwitchAuthToken.builder()
                .accessToken(token.getAccessToken())
                .refreshToken(token.getRefreshToken())
                .expiresIn(token.getExpiresIn())
                .tokenType(token.getTokenType())
                .scope(objectMapper.writeValueAsString(scopes))
                .twitchUserId(twitchUser.id())
                .userName(twitchUser.username())
                .build();

        twitchAuthTokenRepository.save(twitchAuthToken);

        return token;
    }

    public TwitchTokenResponse refreshAccessToken(String refreshToken) throws JsonProcessingException {
        JsonNode response = twitchAuthClient.post()
                .uri("/oauth2/token")
                .contentType(MediaType.APPLICATION_FORM_URLENCODED)
                .body(BodyInserters.fromFormData("client_id", clientId)
                        .with("client_secret", clientSecret)
                        .with("grant_type", "refresh_token")
                        .with("refresh_token", refreshToken))
                .retrieve()
                .bodyToMono(JsonNode.class)
                .block();

        if (response == null || response.get("access_token") == null) {
            throw new RuntimeException("Twitch token refresh failed");
        }

        List<String> scopes = objectMapper.convertValue(
                response.get("scope"), new TypeReference<List<String>>() {}
        );

        TwitchTokenResponse token = new TwitchTokenResponse(
                response.get("access_token").asText(),
                response.get("refresh_token").asText(),
                response.get("expires_in").asLong(),
                scopes,
                response.get("token_type").asText()
        );

        // Update DB token record with new token + refresh token + expiry etc.
        TwitchAuthToken twitchAuthToken = twitchAuthTokenRepository.findTopByOrderByCreatedAtDesc()
                .orElseThrow(() -> new IllegalStateException("No token in DB to refresh"));

        twitchAuthToken.setAccessToken(token.getAccessToken());
        twitchAuthToken.setRefreshToken(token.getRefreshToken());
        twitchAuthToken.setExpiresIn(token.getExpiresIn());
        twitchAuthToken.setScope(objectMapper.writeValueAsString(scopes));
        twitchAuthToken.setTokenType(token.getTokenType());
        twitchAuthToken.setCreatedAt(Instant.now());

        twitchAuthTokenRepository.save(twitchAuthToken);

        return token;
    }


    public TwitchUser getUserInfo(String accessToken) {
        JsonNode userData = twitchApiClient.get()
                .uri("/helix/users")
                .header("Client-ID", clientId)
                .header("Authorization", "Bearer " + accessToken)
                .retrieve()
                .bodyToMono(JsonNode.class)
                .block();

        if (userData == null || userData.get("data").isEmpty()) {
            throw new RuntimeException("Could not retrieve user data from Twitch.");
        }

        JsonNode userNode = userData.get("data").get(0);

        return new TwitchUser(
                userNode.get("id").asText(),
                userNode.get("login").asText(),
                userNode.get("display_name").asText(),
                userNode.get("email").asText(null)
        );
    }

    public TwitchUser fetchUserByCode(String code) throws JsonProcessingException {
        TwitchTokenResponse tokenResponse = exchangeCodeForAccessToken(code);
        String accessToken = tokenResponse.getAccessToken();

        return getUserInfo(accessToken);
    }

    public TwitchChannelStats fetchAndSaveChannelStats(String accessToken, String twitchUserId) {
        JsonNode userResponse = twitchApiClient.get()
                .uri(uriBuilder -> uriBuilder
                        .path("/helix/users")
                        .queryParam("id", twitchUserId)
                        .build()
                )
                .header("Client-ID", clientId)
                .header("Authorization", "Bearer " + accessToken)
                .retrieve()
                .bodyToMono(JsonNode.class)
                .block();

        if (userResponse == null || !userResponse.has("data") || userResponse.get("data").isEmpty()) {
            logger.error("Twitch API returned no userResponse data or empty data node");
            return null;
        }

        JsonNode user = userResponse.get("data").get(0);

        Integer followerCount = getFollowerCount(twitchUserId, accessToken).block();

        if (followerCount == null) {
            logger.warn("Follower count konnte nicht abgerufen werden, setze auf 0.");
        }
        int followers = (followerCount != null) ? followerCount : 0;

        TwitchChannelStats stats = TwitchChannelStats.builder()
                .twitchUserId(user.get("id").asText())
                .displayName(user.get("display_name").asText())
                .description(user.get("description").asText())
                .profileImageUrl(user.get("profile_image_url").asText())
                .offlineImageUrl(user.get("offline_image_url").asText(null))
                .broadcasterType(user.get("broadcaster_type").asText())
                .viewCount(user.get("view_count").asInt())
                .followers(followers)
                .accountCreatedAt(Instant.parse(user.get("created_at").asText()))
                .build();

        return twitchChannelStatsRepository.save(stats);
    }

    public Optional<TwitchChannelStats> getLatestChannelStats() {
        return twitchChannelStatsRepository.findTopByOrderByFetchedAtDesc();
    }

    public Optional<TwitchChannelStats> refreshAndSaveChannelStats() {
        Optional<TwitchAuthToken> latestToken = twitchAuthTokenRepository.findTopByOrderByCreatedAtDesc();

        if (latestToken.isEmpty()) {
            logger.warn("Kein Twitch Auth Token gefunden, abbrechen.");
            return Optional.empty();
        }

        String accessToken = latestToken.get().getAccessToken();

        JsonNode userData;

        try {
            userData = twitchApiClient.get()
                    .uri("/helix/users")
                    .header("Client-ID", clientId)
                    .header("Authorization", "Bearer " + accessToken)
                    .retrieve()
                    .bodyToMono(JsonNode.class)
                    .block();
        } catch (Exception e) {
            logger.error("Fehler bei Twitch API Request /helix/users", e);
            return Optional.empty();
        }

        if (userData == null || !userData.has("data") || userData.get("data").isEmpty()) {
            logger.error("Twitch API returned no user data or empty data node");
            return Optional.empty();
        }

        JsonNode user = userData.get("data").get(0);

        Integer followerCount = null;
        try {
            followerCount = getFollowerCount(user.get("id").asText(), accessToken).block();
        } catch (Exception e) {
            logger.info("User Id ist" + user.get("id").asText());
            logger.warn("Fehler beim Abrufen der Follower-Anzahl", e);
        }

        if (followerCount == null) {
            logger.warn("Follower count konnte nicht abgerufen werden, setze auf 0.");
        }
        int followers = (followerCount != null) ? followerCount : 0;

        String twitchUserId = user.get("id").asText();

        Optional<TwitchChannelStats> existingOpt = twitchChannelStatsRepository.findByTwitchUserId(twitchUserId);

        TwitchChannelStats stats = existingOpt.orElseGet(() -> {
            TwitchChannelStats newStats = new TwitchChannelStats();
            newStats.setTwitchUserId(twitchUserId);
            return newStats;
        });

        stats.setDisplayName(user.get("display_name").asText());
        stats.setDescription(user.get("description").asText());
        stats.setProfileImageUrl(user.get("profile_image_url").asText());
        stats.setOfflineImageUrl(user.get("offline_image_url").asText(null));
        stats.setBroadcasterType(user.get("broadcaster_type").asText());
        stats.setViewCount(user.get("view_count").asInt());
        stats.setFollowers(followers);
        stats.setAccountCreatedAt(Instant.parse(user.get("created_at").asText()));

        TwitchChannelStats savedStats = twitchChannelStatsRepository.save(stats);

        return Optional.of(savedStats);
    }

    public Mono<Integer> getFollowerCount(String broadcasterId, String accessToken) {
        return twitchApiClient.get()
                .uri(uriBuilder -> uriBuilder
                        .path("/helix/channels/followers")
                        .queryParam("broadcaster_id", broadcasterId)
                        .build())
                .header("Authorization", "Bearer " + accessToken)  // <--- wichtig!
                .header("Client-ID", clientId)
                .retrieve()
                .bodyToMono(FollowersResponse.class)
                .map(FollowersResponse::getTotal);
    }


    @Setter
    @Getter
    public static class FollowersResponse {
        private int total;
    }
}
