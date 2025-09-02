package streamer_website.demo.client;

import com.fasterxml.jackson.databind.JsonNode;
import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;
import org.springframework.web.reactive.function.client.WebClient;

import java.util.Optional;

@Service
@RequiredArgsConstructor
public class TwitchApiClient {

    private final WebClient twitchApiClient;
    private static final Logger logger = LoggerFactory.getLogger(TwitchApiClient.class);

    public Optional<JsonNode> getUserByLogin(String login, String accessToken, String clientId) {
        try {
            JsonNode node = twitchApiClient.get()
                    .uri("/helix/users?login=" + login)
                    .header("Authorization", "Bearer " + accessToken)
                    .header("Client-ID", clientId)
                    .retrieve()
                    .bodyToMono(JsonNode.class)
                    .block();
            return Optional.ofNullable(node);
        } catch (Exception e) {
            logger.error("Something went wrong", e);
            return Optional.empty();
        }
    }

    public Optional<JsonNode> getStreamByUserId(String userId, String accessToken, String clientId) {
        try {
            JsonNode node = twitchApiClient.get()
                    .uri("/helix/streams?user_id=" + userId)
                    .header("Authorization", "Bearer " + accessToken)
                    .header("Client-ID", clientId)
                    .retrieve()
                    .bodyToMono(JsonNode.class)
                    .block();
            return Optional.ofNullable(node);
        } catch (Exception e) {
            logger.error("Something went wrong", e);
            return Optional.empty();
        }
    }
}
