package streamer_website.demo.service.twitch;

import com.fasterxml.jackson.databind.JsonNode;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import org.springframework.util.LinkedMultiValueMap;
import org.springframework.util.MultiValueMap;
import org.springframework.web.reactive.function.client.WebClient;
import streamer_website.demo.controller.twitch.TwitchAuthController;
import streamer_website.demo.entity.twitch.TwitchAuthToken;
import streamer_website.demo.repository.TwitchAuthTokenRepository;

import java.time.Instant;

@Service
@RequiredArgsConstructor
@Slf4j
public class TwitchTokenService {

    private final TwitchAuthTokenRepository tokenRepository;
    private final WebClient twitchApiClient = WebClient.create();
    private static final Logger logger = LoggerFactory.getLogger(TwitchTokenService.class);

    @Value("${twitch.clientId}")
    private String clientId;

    @Value("${twitch.client-secret}")
    private String clientSecret;

    @Value("${twitch.redirect-uri}")
    private String redirectUri;

    @Value("${twitch.bot.username}")
    private String botUsername;

    public String getUserAccessToken(String username) {
        return tokenRepository.findTopByUserNameOrderByCreatedAtDesc(username)
                .map(token -> {
                    Instant expiresAt = token.getCreatedAt().plusSeconds(token.getExpiresIn());
                    if (Instant.now().isAfter(expiresAt)) {
                        token = refreshToken(token);
                    }
                    return token.getAccessToken();
                })
                .orElse(null);
    }

    public String getBotAccessToken() {
        return getValidToken(botUsername).getAccessToken();
    }

    public TwitchAuthToken findBotToken() {
        return tokenRepository.findTopByUserNameOrderByCreatedAtDesc(botUsername)
                .orElse(null);
    }

    public TwitchAuthToken saveTokenFromCode(String code, String username) {
        try {
            MultiValueMap<String, String> formData = new LinkedMultiValueMap<>();
            formData.add("client_id", clientId);
            formData.add("client_secret", clientSecret);
            formData.add("code", code);
            formData.add("grant_type", "authorization_code");
            formData.add("redirect_uri", redirectUri);

            JsonNode response = twitchApiClient.post()
                    .uri("https://id.twitch.tv/oauth2/token")
                    .header("Content-Type", "application/x-www-form-urlencoded")
                    .bodyValue(formData)
                    .retrieve()
                    .bodyToMono(JsonNode.class)
                    .block();

            if (response == null || !response.has("access_token")) {
                throw new IllegalStateException("Twitch response ohne Access-Token");
            }

            TwitchAuthToken token = TwitchAuthToken.builder()
                    .userName(username)
                    .accessToken(response.get("access_token").asText())
                    .refreshToken(response.get("refresh_token").asText())
                    .expiresIn(response.get("expires_in").asLong())
                    .createdAt(Instant.now())
                    .tokenType(response.get("token_type").asText())
                    .scope(response.get("scope").toString())
                    .build();

            return tokenRepository.save(token);

        } catch (Exception e) {
            log.error("Fehler beim Speichern des Tokens für {}: {}", username, e.getMessage(), e);
            throw new RuntimeException("Konnte Token nicht speichern", e);
        }
    }

    public String buildOAuthUrl(boolean forBot) {
        String scope = forBot
                ? "chat:read chat:edit channel:manage:broadcast"
                : "user:read:email"; // hier kannst du für User eigene Scopes definieren

        return "https://id.twitch.tv/oauth2/authorize" +
                "?client_id=" + clientId +
                "&redirect_uri=" + redirectUri +
                "&response_type=code" +
                "&scope=" + scope +
                "&force_verify=true";
    }

    public TwitchAuthToken exchangeCodeForToken(String code, boolean forBot) {
        if (forBot) {
            return saveTokenFromCode(code, botUsername);
        } else {
            JsonNode tokenResponse = fetchTokenFromCode(code);

            String accessToken = tokenResponse.get("access_token").asText();
            String refreshToken = tokenResponse.get("refresh_token").asText();
            long expiresIn = tokenResponse.get("expires_in").asLong();
            String tokenType = tokenResponse.get("token_type").asText();
            String scope = tokenResponse.get("scope").toString();

            String username = fetchUsernameFromAccessToken(accessToken);
            String userId = fetchUserIdFromAccessToken(accessToken);

            TwitchAuthToken token = TwitchAuthToken.builder()
                    .userName(username)
                    .twitchUserId(userId)
                    .accessToken(accessToken)
                    .refreshToken(refreshToken)
                    .expiresIn(expiresIn)
                    .createdAt(Instant.now())
                    .tokenType(tokenType)
                    .scope(scope)
                    .build();

            return tokenRepository.save(token);
        }
    }

    private JsonNode fetchTokenFromCode(String code) {
        MultiValueMap<String, String> formData = new LinkedMultiValueMap<>();
        formData.add("client_id", clientId);
        formData.add("client_secret", clientSecret);
        formData.add("code", code);
        formData.add("grant_type", "authorization_code");
        formData.add("redirect_uri", redirectUri);

        JsonNode response = twitchApiClient.post()
                .uri("https://id.twitch.tv/oauth2/token")
                .header("Content-Type", "application/x-www-form-urlencoded")
                .bodyValue(formData)
                .retrieve()
                .bodyToMono(JsonNode.class)
                .block();

        if (response == null || !response.has("access_token")) {
            throw new IllegalStateException("Twitch response ohne Access-Token");
        }

        return response;
    }

    private String fetchUsernameFromAccessToken(String accessToken) {
        JsonNode userResponse = twitchApiClient.get()
                .uri("https://api.twitch.tv/helix/users")
                .header("Authorization", "Bearer " + accessToken)
                .header("Client-Id", clientId)
                .retrieve()
                .bodyToMono(JsonNode.class)
                .block();

        if (userResponse == null || !userResponse.has("data") || userResponse.get("data").isEmpty()) {
            throw new IllegalStateException("Konnte Twitch-Userinfo nicht abrufen");
        }

        logger.info("Fetched Twitch username: {}", userResponse.get("data").get(0).get("login").asText());
        return userResponse.get("data").get(0).get("login").asText();
    }

    private String fetchUserIdFromAccessToken(String accessToken) {
        JsonNode userResponse = twitchApiClient.get()
                .uri("https://api.twitch.tv/helix/users")
                .header("Authorization", "Bearer " + accessToken)
                .header("Client-Id", clientId)
                .retrieve()
                .bodyToMono(JsonNode.class)
                .block();

        if (userResponse == null || !userResponse.has("data") || userResponse.get("data").isEmpty()) {
            throw new IllegalStateException("Konnte Twitch-Userinfo nicht abrufen");
        }

        return userResponse.get("data").get(0).get("id").asText();
    }

    private TwitchAuthToken getValidToken(String username) {
        TwitchAuthToken token = tokenRepository.findTopByUserNameOrderByCreatedAtDesc(username)
                .orElseThrow(() -> new IllegalStateException("Kein Token für " + username + " gefunden"));

        Instant expiresAt = token.getCreatedAt().plusSeconds(token.getExpiresIn());
        if (Instant.now().isAfter(expiresAt)) {
            log.info("Access token für {} abgelaufen, refreshe...", username);
            token = refreshToken(token);
        }
        return token;
    }

    private TwitchAuthToken refreshToken(TwitchAuthToken oldToken) {
        try {
            var response = twitchApiClient.post()
                    .uri("https://id.twitch.tv/oauth2/token" +
                            "?grant_type=refresh_token" +
                            "&refresh_token=" + oldToken.getRefreshToken() +
                            "&client_id=" + clientId +
                            "&client_secret=" + clientSecret)
                    .retrieve()
                    .bodyToMono(JsonNode.class)
                    .block();

            if (response == null || !response.has("access_token")) {
                throw new IllegalStateException("Twitch response ohne Access-Token beim Refresh");
            }

            oldToken.setAccessToken(response.get("access_token").asText());
            if (response.has("refresh_token")) {
                oldToken.setRefreshToken(response.get("refresh_token").asText());
            }
            oldToken.setExpiresIn(response.get("expires_in").asLong());
            oldToken.setCreatedAt(Instant.now());

            return tokenRepository.save(oldToken);

        } catch (Exception e) {
            log.error("Fehler beim Refresh des Tokens für {}: {}", oldToken.getUserName(), e.getMessage(), e);
            throw new RuntimeException("Token refresh fehlgeschlagen", e);
        }
    }

    private String fetchUsernameFromCode(String code) {
        JsonNode tokenResponse = twitchApiClient.post()
                .uri("https://id.twitch.tv/oauth2/token" +
                        "?client_id=" + clientId +
                        "&client_secret=" + clientSecret +
                        "&code=" + code +
                        "&grant_type=authorization_code" +
                        "&redirect_uri=" + redirectUri)
                .retrieve()
                .bodyToMono(JsonNode.class)
                .block();

        if (tokenResponse == null || !tokenResponse.has("access_token")) {
            throw new IllegalStateException("Twitch response ohne Access-Token");
        }

        String accessToken = tokenResponse.get("access_token").asText();

        JsonNode userResponse = twitchApiClient.get()
                .uri("https://api.twitch.tv/helix/users")
                .header("Authorization", "Bearer " + accessToken)
                .header("Client-Id", clientId)
                .retrieve()
                .bodyToMono(JsonNode.class)
                .block();

        if (userResponse == null || !userResponse.has("data") || userResponse.get("data").isEmpty()) {
            throw new IllegalStateException("Konnte Twitch-Userinfo nicht abrufen");
        }

        return userResponse.get("data").get(0).get("login").asText();
    }
}
