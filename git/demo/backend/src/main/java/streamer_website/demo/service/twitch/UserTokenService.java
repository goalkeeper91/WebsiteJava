package streamer_website.demo.service.twitch;

import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.web.reactive.function.BodyInserters;
import org.springframework.web.reactive.function.client.WebClient;
import streamer_website.demo.dto.TwitchTokenResponse;
import streamer_website.demo.entity.twitch.TwitchAuthToken;
import streamer_website.demo.repository.TwitchAuthTokenRepository;

import java.time.Instant;
import java.util.List;

@Service
@RequiredArgsConstructor
public class UserTokenService {

    private final WebClient twitchAuthClient;
    private final TwitchAuthTokenRepository twitchAuthTokenRepository;
    private final ObjectMapper objectMapper;

    public String getAccessToken() {
        TwitchAuthToken token = twitchAuthTokenRepository.findTopByOrderByCreatedAtDesc()
                .orElseThrow(() -> new IllegalStateException("No user token found"));

        Instant expiry = token.getCreatedAt().plusSeconds(token.getExpiresIn());
        if (Instant.now().isAfter(expiry.minusSeconds(60))) {
            return refreshToken(token.getRefreshToken()).getAccessToken();
        }

        return token.getAccessToken();
    }

    private TwitchTokenResponse refreshToken(String refreshToken) {
        try {
            JsonNode response = twitchAuthClient.post()
                    .uri("/oauth2/token")
                    .body(BodyInserters.fromFormData("client_id", System.getenv("TWITCH_CLIENT_ID"))
                            .with("client_secret", System.getenv("TWITCH_CLIENT_SECRET"))
                            .with("grant_type", "refresh_token")
                            .with("refresh_token", refreshToken))
                    .retrieve()
                    .bodyToMono(JsonNode.class)
                    .block();

            if (response == null || response.get("access_token") == null) {
                throw new RuntimeException("Failed to refresh Twitch token");
            }

            List<String> scopes = objectMapper.convertValue(
                    response.get("scope"), new TypeReference<List<String>>() {}
            );

            TwitchTokenResponse tokenResponse = new TwitchTokenResponse(
                    response.get("access_token").asText(),
                    response.get("refresh_token").asText(),
                    response.get("expires_in").asLong(),
                    scopes,
                    response.get("token_type").asText()
            );

            TwitchAuthToken tokenEntity = twitchAuthTokenRepository.findTopByOrderByCreatedAtDesc()
                    .orElseThrow();

            tokenEntity.setAccessToken(tokenResponse.getAccessToken());
            tokenEntity.setRefreshToken(tokenResponse.getRefreshToken());
            tokenEntity.setExpiresIn(tokenResponse.getExpiresIn());
            tokenEntity.setScope(objectMapper.writeValueAsString(scopes));
            tokenEntity.setTokenType(tokenResponse.getTokenType());
            tokenEntity.setCreatedAt(Instant.now());

            twitchAuthTokenRepository.save(tokenEntity);

            return tokenResponse;

        } catch (Exception e) {
            throw new RuntimeException("Error refreshing user token", e);
        }
    }
}
