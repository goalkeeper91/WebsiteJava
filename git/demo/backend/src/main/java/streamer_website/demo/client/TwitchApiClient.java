package streamer_website.demo.service.twitch;

import com.fasterxml.jackson.databind.JsonNode;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.web.reactive.function.client.WebClient;

@Service
@RequiredArgsConstructor
public class TwitchApiClient {

    private final WebClient twitchApiClient;

    public JsonNode getUserByLogin(String login, String accessToken, String clientId) {
        return twitchApiClient.get()
                .uri("/helix/users?login=" + login)
                .header("Authorization", "Bearer " + accessToken)
                .header("Client-ID", clientId)
                .retrieve()
                .bodyToMono(JsonNode.class)
                .block();
    }

    public JsonNode getStreamByUserId(String userId, String accessToken, String clientId) {
        return twitchApiClient.get()
                .uri("/helix/streams?user_id=" + userId)
                .header("Authorization", "Bearer " + accessToken)
                .header("Client-ID", clientId)
                .retrieve()
                .bodyToMono(JsonNode.class)
                .block();
    }
}
