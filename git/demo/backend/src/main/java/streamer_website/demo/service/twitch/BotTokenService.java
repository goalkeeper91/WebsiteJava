package streamer_website.demo.service.twitch;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.RequiredArgsConstructor;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.MediaType;
import org.springframework.stereotype.Service;
import org.springframework.web.reactive.function.BodyInserters;
import org.springframework.web.reactive.function.client.WebClient;

import java.time.Instant;

@Service
@RequiredArgsConstructor
public class BotTokenService {

    private final WebClient twitchAuthClient;
    private final ObjectMapper objectMapper;

    @Value("${twitch.clientId}")
    private String clientId;

    @Value("${twitch.clientSecret}")
    private String clientSecret;

    private String cachedToken;
    private long expiryEpoch;

    public synchronized String getToken() {
        if (cachedToken != null && Instant.now().getEpochSecond() < expiryEpoch - 60) {
            return cachedToken;
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

        cachedToken = response.get("access_token").asText();
        expiryEpoch = Instant.now().getEpochSecond() + response.get("expires_in").asLong();

        return cachedToken;
    }
}
