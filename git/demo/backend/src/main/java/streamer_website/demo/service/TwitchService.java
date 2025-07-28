package streamer_website.demo.service;

import org.springframework.web.reactive.function.BodyInserters;
import org.springframework.web.reactive.function.client.WebClient;
import com.fasterxml.jackson.databind.JsonNode;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.MediaType;
import org.springframework.stereotype.Service;

import java.util.Objects;

@Service
public class TwitchService {

    private final String clientId;

    private final String clientSecret;

    private final String username;

    private final String apiBaseUrl;

    @Value("${twitch.auth-base-url}")
    private String authBaseUrl;

    private String cachedToken;
    private long tokenExpiry;

    private final WebClient webClient;

    public TwitchService(
            WebClient.Builder webClientBuilder,
            @Value("${twitch.clientId}") String clientId,
            @Value("${twitch.clientSecret}") String clientSecret,
            @Value("${twitch.username}") String username,
            @Value("${twitch.apiBaseUrl}") String apiBaseUrl
                         ) {
        this.clientId = clientId;
        this.clientSecret = clientSecret;
        this.username = username;
        this.apiBaseUrl = apiBaseUrl;
        this.webClient = webClientBuilder.baseUrl(apiBaseUrl).build();
    }

    public boolean isLive() {
        try {
            String token = getAccessToken();

            String userId = Objects.requireNonNull(webClient.get()
                            .uri("/helix/users?login=" + username)
                            .header("Client-ID", clientId)
                            .header("Authorization", "Bearer " + token)
                            .retrieve()
                            .bodyToMono(JsonNode.class)
                            .block())
                    .get("data").get(0).get("id").asText();

            // Check Live-Status
            JsonNode streamData = webClient.get()
                    .uri("/helix/streams?user_id=" + userId)
                    .header("Client-ID", clientId)
                    .header("Authorization", "Bearer " + token)
                    .retrieve()
                    .bodyToMono(JsonNode.class)
                    .block();

            assert streamData != null;
            return !streamData.get("data").isEmpty();
        } catch (Exception e) {
            return false;
        }
    }

    private String getAccessToken() {
        if(cachedToken != null && System.currentTimeMillis() < tokenExpiry) {
            return cachedToken;
        }

        JsonNode response = webClient.post()
                .uri(authBaseUrl + "/oauth2/token")
                .contentType(MediaType.APPLICATION_FORM_URLENCODED)
                .body(BodyInserters.fromFormData("client_id", clientId)
                        .with("client_secret", clientSecret)
                        .with("grant_type", "client_credentials"))
                .retrieve()
                .bodyToMono(JsonNode.class)
                .block();

        assert response != null;
        cachedToken = response.get("access_token").asText();
        long expiresIn = response.get("expires_in").asLong();
        tokenExpiry = System.currentTimeMillis() + (expiresIn * 1000);

        return cachedToken;
    }
}
