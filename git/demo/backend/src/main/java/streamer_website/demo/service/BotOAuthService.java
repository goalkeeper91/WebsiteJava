package streamer_website.demo.service;

import com.fasterxml.jackson.databind.JsonNode;
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

@Service
public class BotOAuthService {

    private static final Logger logger = LoggerFactory.getLogger(BotOAuthService.class);

    private final TwitchAuthTokenRepository tokenRepository;
    private final WebClient webClient = WebClient.create();

    private final String clientId;
    private final String clientSecret;
    private final String redirectUri;

    public BotOAuthService(TwitchAuthTokenRepository tokenRepository,
                           @Value("${twitch.client-id}") String clientId,
                           @Value("${twitch.client-secret}") String clientSecret,
                           @Value("${twitch.redirect-uri}") String redirectUri)
    {
        this.tokenRepository = tokenRepository;
        this.clientId = clientId;
        this.clientSecret = clientSecret;
        this.redirectUri = redirectUri;
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

    public void saveBotTokenFromCode(String code) {
        var response = webClient.post()
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

        TwitchAuthToken token = new TwitchAuthToken();
        token.setAccessToken(response.getAccessToken());
        token.setRefreshToken(response.getRefreshToken());
        token.setExpiresIn(Instant.now().getEpochSecond() + response.getExpiresIn());
        token.setTwitchUserId(fetchBotUserId(response.getAccessToken()));

        tokenRepository.save(token);
        logger.info("Bot-Token erfolgreich gespeichert für User {}", token.getTwitchUserId());
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
        String scope = "chat:read chat:edit"; // ggf. erweitern
        return "https://id.twitch.tv/oauth2/authorize" +
                "?client_id=" + clientId +
                "&redirect_uri=" + redirectUri +
                "&response_type=code" +
                "&scope=" + scope;
    }

    public TwitchAuthToken findBotToken() {
        return tokenRepository.findByTwitchUserId(getBotUserId())
                .orElse(null);
    }

    private String getBotUserId() {
        return "dein-bot-user-id";
    }

    @Getter
    private static class RefreshResponse {
        private String access_token;
        private String refresh_token;
        private long expires_in;
    }
}
