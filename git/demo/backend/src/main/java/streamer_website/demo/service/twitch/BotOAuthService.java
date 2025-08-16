package streamer_website.demo.service.twitch;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.github.philippheuer.credentialmanager.domain.OAuth2Credential;
import lombok.Getter;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import org.springframework.web.reactive.function.client.WebClient;
import streamer_website.demo.dto.TokenInfoDto;
import streamer_website.demo.dto.TwitchTokenResponse;
import streamer_website.demo.entity.TwitchAuthToken;
import streamer_website.demo.repository.TwitchAuthTokenRepository;

import java.time.Instant;
import java.util.List;

@Service
public class BotOAuthService {

    private static final Logger logger = LoggerFactory.getLogger(BotOAuthService.class);

    private final TwitchAuthTokenRepository tokenRepository;
    private final WebClient webClient = WebClient.create();

    private final String clientId;
    private final String clientSecret;
    private final String redirectUri;

    private final ObjectMapper objectMapper;

    public BotOAuthService(TwitchAuthTokenRepository tokenRepository,
                           @Value("${twitch.bot.client-id}") String clientId,
                           @Value("${twitch.bot.clientSecret}") String clientSecret,
                           @Value("${twitch.bot.redirectUri}") String redirectUri,
                           ObjectMapper objectMapper
    )
    {
        this.tokenRepository = tokenRepository;
        this.clientId = clientId;
        this.clientSecret = clientSecret;
        this.redirectUri = redirectUri;
        this.objectMapper = objectMapper;
    }

    public OAuth2Credential getBotCredential(String twitchUserId) {
        TwitchAuthToken token = tokenRepository.findByTwitchUserId(twitchUserId)
                .orElseThrow(() -> new IllegalStateException("Token für TwitchUser " + twitchUserId + " nicht gefunden"));

        if (token.getExpiresIn() < Instant.now().getEpochSecond()) {
            logger.info("Access token für {} abgelaufen, refreshen...", twitchUserId);
            refreshToken(token);
        }

        return new OAuth2Credential("twitch", token.getAccessToken());
    }

    private void refreshToken(TwitchAuthToken token) {
        var response = webClient.post()
                .uri("https://id.twitch.tv/oauth2/token" +
                        "?grant_type=refresh_token" +
                        "&refresh_token=" + token.getRefreshToken() +
                        "&client_id=" + clientId +
                        "&client_secret=" + clientSecret)
                .retrieve()
                .bodyToMono(RefreshResponse.class)
                .block();

        if (response == null || response.getAccess_token() == null) {
            logger.error("Bot Token konnte nicht erneuert werden (Response ist null)");
            throw new IllegalStateException("Failed to refresh token für " + token.getTwitchUserId());
        }

        token.setAccessToken(response.getAccess_token());
        token.setRefreshToken(response.getRefresh_token() != null ? response.getRefresh_token() : token.getRefreshToken());
        token.setExpiresIn(Instant.now().getEpochSecond() + response.getExpires_in());

        tokenRepository.save(token);
        logger.info("Token für {} erfolgreich refreshed", token.getTwitchUserId());
    }

    public void saveBotTokenFromCode(String code) throws JsonProcessingException {
        // Token von Twitch holen
        TwitchTokenResponse response = webClient.post()
                .uri("https://id.twitch.tv/oauth2/token" +
                        "?client_id=" + clientId +
                        "&client_secret=" + clientSecret +
                        "&code=" + code +
                        "&grant_type=authorization_code" +
                        "&redirect_uri=" + redirectUri)
                .retrieve()
                .bodyToMono(TwitchTokenResponse.class)
                .block();

        if (response == null || response.getAccessToken() == null) {
            logger.error("Konnte Bot-Token nicht speichern: Response leer oder AccessToken fehlt");
            throw new IllegalStateException("Bot-Token konnte nicht geholt werden");
        }

        List<String> scopes = response.getScope();

        // Bot-Userdaten abfragen
        String twitchUserId = fetchBotUserId(response.getAccessToken());
        String botUsername = fetchBotUsername(response.getAccessToken());

        // Token in DB speichern oder aktualisieren
        TwitchAuthToken token = tokenRepository.findByTwitchUserId(twitchUserId)
                .orElse(TwitchAuthToken.builder().build());

        token.setTwitchUserId(twitchUserId);
        token.setUserName(botUsername);
        token.setAccessToken(response.getAccessToken());
        token.setRefreshToken(response.getRefreshToken());
        token.setExpiresIn(Instant.now().getEpochSecond() + response.getExpiresIn());
        token.setScope(objectMapper.writeValueAsString(scopes));
        token.setTokenType(response.getTokenType());
        token.setCreatedAt(Instant.now());

        tokenRepository.save(token);
        logger.info("Bot-Token erfolgreich gespeichert für User {}", twitchUserId);
    }


    private String fetchBotUsername(String accessToken) {
        var userResponse = webClient.get()
                .uri("https://api.twitch.tv/helix/users")
                .header("Authorization", "Bearer " + accessToken)
                .header("Client-Id", clientId)
                .retrieve()
                .bodyToMono(JsonNode.class)
                .block();

        if (userResponse == null || !userResponse.has("data") || userResponse.get("data").isEmpty()) {
            throw new IllegalStateException("Konnte Bot-Username nicht abrufen");
        }

        return userResponse.get("data").get(0).get("login").asText();
    }


    private String fetchBotUserId(String accessToken) {
        var resp = webClient.get()
                .uri("https://api.twitch.tv/helix/users?login=goalkeeper91_bot")
                .header("Client-ID", clientId)
                .header("Authorization", "Bearer " + accessToken)
                .retrieve()
                .bodyToMono(JsonNode.class)
                .block();

        if (resp == null || !resp.has("data") || resp.get("data").isEmpty()) {
            logger.error("Konnte Bot-UserID nicht ermitteln – leere Antwort von Twitch");
            throw new IllegalStateException("Bot-UserID konnte nicht geholt werden");
        }

        return resp.get("data").get(0).get("id").asText();
    }

    public TokenInfoDto getTokenInfo() {
        TwitchAuthToken token = findBotToken();

        if (token == null) {
            return TokenInfoDto.builder()
                    .tokenPresent(false)
                    .expiresAt(null)
                    .twitchUserId(null)
                    .build();
        }

        return TokenInfoDto.builder()
                .tokenPresent(true)
                .expiresAt(Instant.ofEpochSecond(token.getExpiresIn()))
                .twitchUserId(token.getTwitchUserId())
                .build();
    }

    public String buildOAuthUrl() {
        String scope = "chat:read chat:edit channel:manage:broadcast";
        return "https://id.twitch.tv/oauth2/authorize" +
                "?client_id=" + clientId +
                "&redirect_uri=" + redirectUri +
                "&response_type=code" +
                "&scope=" + scope +
                "&force_verify=true";
    }

    public TwitchAuthToken findBotToken() {
        return tokenRepository.findByTwitchUserId(getBotUserId())
                .orElse(null);
    }

    private String getBotUserId() {
        return tokenRepository.findTopByUserNameOrderByCreatedAtDesc("goalkeeper91_bot")
                .map(TwitchAuthToken::getTwitchUserId)
                .orElse("Not found");
    }

    @Getter
    private static class RefreshResponse {
        private String access_token;
        private String refresh_token;
        private long expires_in;
    }
}
