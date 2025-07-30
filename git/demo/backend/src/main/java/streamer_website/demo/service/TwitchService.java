package streamer_website.demo.service;

import org.springframework.web.reactive.function.BodyInserters;
import org.springframework.web.reactive.function.client.WebClient;
import com.fasterxml.jackson.databind.JsonNode;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.MediaType;
import org.springframework.stereotype.Service;
import streamer_website.demo.dto.TwitchUser;

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

    private final String redirectUri;

    public TwitchService(
            WebClient.Builder webClientBuilder,
            @Value("${twitch.clientId}") String clientId,
            @Value("${twitch.clientSecret}") String clientSecret,
            @Value("${twitch.username}") String username,
            @Value("${twitch.apiBaseUrl}") String apiBaseUrl,
            @Value("${twitch.redirectUri}") String redirectUri
                         ) {
        this.clientId = clientId;
        this.clientSecret = clientSecret;
        this.username = username;
        this.apiBaseUrl = apiBaseUrl;
        this.webClient = webClientBuilder.baseUrl(apiBaseUrl).build();
        this.redirectUri = redirectUri;
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

        if(response == null) {
            throw new IllegalStateException("No response from Twitch");
        }

        cachedToken = response.get("access_token").asText();

        long expiresIn = response.get("expires_in").asLong();

        tokenExpiry = System.currentTimeMillis() + (expiresIn * 1000);

        return cachedToken;
    }

    public String exchangeCodeForAccessToken(String code) {
        JsonNode response = webClient.post()
                .uri(authBaseUrl + "/oauth2/token")
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
            throw new RuntimeException("Twitch token exchange failed.");
        }

        return response.get("access_token").asText();
    }

    public TwitchUser getUserInfo(String accessToken) {
        JsonNode userData = webClient.get()
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

    public TwitchUser fetchUserByCode(String code) {
        String accessToken = exchangeCodeForAccessToken(code);

        return getUserInfo(accessToken);
    }
}
