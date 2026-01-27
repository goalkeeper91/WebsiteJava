package streamer_website.demo.service.twitch;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import org.springframework.util.LinkedMultiValueMap;
import org.springframework.util.MultiValueMap;
import org.springframework.web.reactive.function.client.WebClient;
import streamer_website.demo.entity.twitch.TwitchAuthToken;
import streamer_website.demo.entity.twitch.TwitchTokenType;
import streamer_website.demo.repository.TwitchAuthTokenRepository;
import streamer_website.demo.service.BotSignalService;

import java.time.Instant;

@Service
@RequiredArgsConstructor
@Slf4j
public class TwitchTokenService {

    private final TwitchAuthTokenRepository tokenRepository;
    private final WebClient twitchApiClient = WebClient.create();
    private final BotSignalService botSignalService;

    @Value("${twitch.clientId}")
    private String clientId;

    @Value("${twitch.client-secret}")
    private String clientSecret;

    @Value("${twitch.redirect-uri}")
    private String redirectUri;

    @Value("${twitch.bot.username}")
    private String botUsername;

    // ---------------------- PUBLIC METHODS ----------------------

    public TwitchAuthToken exchangeCodeForToken(String code, boolean forBot) {
        JsonNode tokenResponse = fetchTokenFromCode(code);

        String accessToken = tokenResponse.get("access_token").asText();
        String[] userInfo = fetchUsernameAndId(accessToken);

        String username = forBot ? botUsername : userInfo[0];
        String userId = userInfo[1];
        TwitchTokenType owner = forBot ? TwitchTokenType.BOT : TwitchTokenType.USER;

        return upsertToken(username, userId, tokenResponse, owner);
    }

    public TwitchAuthToken getValidBotToken() {
        TwitchAuthToken token = tokenRepository
                .findByUserNameAndTokenOwner(botUsername, TwitchTokenType.BOT)
                .orElseThrow(() -> new IllegalStateException("Kein Bot-Token vorhanden"));

        if (isExpired(token)) {
            token = refreshToken(token);
        }
        return token;
    }

    public TwitchAuthToken getValidUserToken(String username) {
        TwitchAuthToken token = tokenRepository
                .findByUserNameAndTokenOwner(username, TwitchTokenType.USER)
                .orElseThrow(() -> new IllegalStateException("Kein Token für User " + username));

        if (isExpired(token)) {
            token = refreshToken(token);
        }
        return token;
    }

    // ---------------------- PRIVATE METHODS ----------------------

    private TwitchAuthToken upsertToken(String username, String userId, JsonNode response, TwitchTokenType owner) {
        TwitchAuthToken token = tokenRepository
                .findByUserNameAndTokenOwner(username, owner)
                .orElseGet(TwitchAuthToken::new);

        token.setUserName(username);
        token.setTwitchUserId(userId);
        token.setAccessToken(response.get("access_token").asText());
        token.setRefreshToken(response.get("refresh_token").asText());
        token.setExpiresIn(response.get("expires_in").asLong());
        token.setCreatedAt(Instant.now());
        token.setTokenType(response.get("token_type").asText());
        token.setScope(response.get("scope").toString());
        token.setTokenOwner(owner);

        TwitchAuthToken savedToken = tokenRepository.save(token);

        if (owner == TwitchTokenType.USER) {
            botSignalService.sendBotJoinSignal(userId);
            log.info("Bot-Signal gesendet für User: {}", userId);
        }

        return savedToken;
    }

    private boolean isExpired(TwitchAuthToken token) {
        return Instant.now().isAfter(token.getCreatedAt().plusSeconds(token.getExpiresIn()));
    }

    private TwitchAuthToken refreshToken(TwitchAuthToken oldToken) {
        try {
            MultiValueMap<String, String> formData = new LinkedMultiValueMap<>();
            formData.add("grant_type", "refresh_token");
            formData.add("refresh_token", oldToken.getRefreshToken());
            formData.add("client_id", clientId);
            formData.add("client_secret", clientSecret);

            JsonNode response = postTokenRequest(formData);

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

    private JsonNode fetchTokenFromCode(String code) {
        MultiValueMap<String, String> formData = new LinkedMultiValueMap<>();
        formData.add("client_id", clientId);
        formData.add("client_secret", clientSecret);
        formData.add("code", code);
        formData.add("grant_type", "authorization_code");
        formData.add("redirect_uri", redirectUri);

        return postTokenRequest(formData);
    }

    private JsonNode postTokenRequest(MultiValueMap<String, String> formData) {
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

    private String[] fetchUsernameAndId(String accessToken) {
        JsonNode user = fetchUserInfo(accessToken);
        return new String[]{
                user.get("login").asText(),
                user.get("id").asText()
        };
    }

    private JsonNode fetchUserInfo(String accessToken) {
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

        return userResponse.get("data").get(0);
    }
}
