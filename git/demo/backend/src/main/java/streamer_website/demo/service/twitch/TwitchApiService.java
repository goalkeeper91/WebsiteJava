package streamer_website.demo.service.twitch;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.*;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;
import org.springframework.web.util.UriComponentsBuilder;
import streamer_website.demo.dto.twitch.*;
import streamer_website.demo.repository.TwitchAuthTokenRepository;

import java.util.ArrayList;
import java.util.List;

@Slf4j
@Service
@RequiredArgsConstructor
public class TwitchApiService {

    @Value("${twitch.clientId}")
    private String clientId;

    private final TwitchAuthTokenRepository tokenRepository;
    private final RestTemplate restTemplate = new RestTemplate();

    private static final String TWITCH_API_BASE = "https://api.twitch.tv/helix";

    /**
     * Holt aktuelle Stream-Informationen
     */
    public StreamInfoDto getStreamInfo(String broadcasterId, String accessToken) {
        String url = UriComponentsBuilder
                .fromUriString(TWITCH_API_BASE + "/channels")
                .queryParam("broadcaster_id", broadcasterId)
                .toUriString();

        try {
            HttpHeaders headers = createHeaders(accessToken);
            HttpEntity<Void> entity = new HttpEntity<>(headers);

            ResponseEntity<JsonNode> response = restTemplate.exchange(
                    url,
                    HttpMethod.GET,
                    entity,
                    JsonNode.class
            );

            if (response.getBody() != null && response.getBody().has("data")) {
                JsonNode data = response.getBody().get("data");
                if (data.isArray() && !data.isEmpty()) {
                    JsonNode channel = data.get(0);
                    return StreamInfoDto.builder()
                            .broadcasterId(channel.get("broadcaster_id").asText())
                            .broadcasterName(channel.get("broadcaster_name").asText())
                            .title(channel.get("title").asText())
                            .gameId(channel.get("game_id").asText())
                            .gameName(channel.get("game_name").asText())
                            .build();
                }
            }

            return null;
        } catch (Exception e) {
            log.error("Fehler beim Abrufen der Stream-Info", e);
            throw new RuntimeException("Konnte Stream-Info nicht abrufen", e);
        }
    }

    /**
     * Aktualisiert Titel und/oder Kategorie des Streams
     */
    public void updateChannelInfo(
            String broadcasterId,
            String accessToken,
            String title,
            String gameId
    ) {
        String url = UriComponentsBuilder
                .fromUriString(TWITCH_API_BASE + "/channels")
                .queryParam("broadcaster_id", broadcasterId)
                .toUriString();

        try {
            HttpHeaders headers = createHeaders(accessToken);
            headers.setContentType(MediaType.APPLICATION_JSON);

            // Body zusammenstellen
            var body = new java.util.HashMap<String, String>();
            if (title != null) body.put("title", title);
            if (gameId != null) body.put("game_id", gameId);

            HttpEntity<Object> entity = new HttpEntity<>(body, headers);

            restTemplate.exchange(
                    url,
                    HttpMethod.PATCH,
                    entity,
                    Void.class
            );

            log.info("Channel-Info aktualisiert: broadcaster={}, title={}, gameId={}",
                    broadcasterId, title, gameId);

        } catch (Exception e) {
            log.error("Fehler beim Aktualisieren der Channel-Info", e);
            throw new RuntimeException("Konnte Channel-Info nicht aktualisieren", e);
        }
    }

    /**
     * Prüft, ob der Stream live ist und liefert Live-Daten
     */
    public LiveStreamDto getLiveStream(String broadcasterId, String accessToken) {
        String url = UriComponentsBuilder
                .fromUriString(TWITCH_API_BASE + "/streams")
                .queryParam("user_id", broadcasterId)
                .toUriString();

        try {
            HttpHeaders headers = createHeaders(accessToken);
            HttpEntity<Void> entity = new HttpEntity<>(headers);

            ResponseEntity<JsonNode> response = restTemplate.exchange(
                    url,
                    HttpMethod.GET,
                    entity,
                    JsonNode.class
            );

            if (response.getBody() != null && response.getBody().has("data")) {
                JsonNode data = response.getBody().get("data");
                if (data.isArray() && !data.isEmpty()) {
                    JsonNode stream = data.get(0);
                    return LiveStreamDto.builder()
                            .isLive(true)
                            .viewerCount(stream.get("viewer_count").asInt())
                            .startedAt(stream.get("started_at").asText())
                            .title(stream.get("title").asText())
                            .gameName(stream.get("game_name").asText())
                            .thumbnailUrl(stream.get("thumbnail_url").asText())
                            .build();
                }
            }

            // Stream ist offline
            return LiveStreamDto.builder()
                    .isLive(false)
                    .viewerCount(0)
                    .build();

        } catch (Exception e) {
            log.error("Fehler beim Abrufen des Stream-Status", e);
            throw new RuntimeException("Konnte Stream-Status nicht abrufen", e);
        }
    }

    /**
     * Holt Follower-Anzahl
     */
    public int getFollowerCount(String broadcasterId, String accessToken) {
        String url = UriComponentsBuilder
                .fromUriString(TWITCH_API_BASE + "/channels/followers")
                .queryParam("broadcaster_id", broadcasterId)
                .queryParam("first", "1")
                .toUriString();

        try {
            HttpHeaders headers = createHeaders(accessToken);
            HttpEntity<Void> entity = new HttpEntity<>(headers);

            ResponseEntity<JsonNode> response = restTemplate.exchange(
                    url,
                    HttpMethod.GET,
                    entity,
                    JsonNode.class
            );

            if (response.getBody() != null && response.getBody().has("total")) {
                return response.getBody().get("total").asInt();
            }

            return 0;
        } catch (Exception e) {
            log.error("Fehler beim Abrufen der Follower-Anzahl", e);
            return 0;
        }
    }

    /**
     * Sucht nach Kategorien/Spielen
     */
    public List<CategoryDto> searchCategories(String query, String accessToken) {
        String url = UriComponentsBuilder
                .fromUriString(TWITCH_API_BASE + "/search/categories")
                .queryParam("query", query)
                .queryParam("first", "10")
                .toUriString();

        try {
            HttpHeaders headers = createHeaders(accessToken);
            HttpEntity<Void> entity = new HttpEntity<>(headers);

            ResponseEntity<JsonNode> response = restTemplate.exchange(
                    url,
                    HttpMethod.GET,
                    entity,
                    JsonNode.class
            );

            List<CategoryDto> categories = new ArrayList<>();

            if (response.getBody() != null && response.getBody().has("data")) {
                JsonNode data = response.getBody().get("data");
                if (data.isArray()) {
                    for (JsonNode cat : data) {
                        categories.add(CategoryDto.builder()
                                .id(cat.get("id").asText())
                                .name(cat.get("name").asText())
                                .boxArtUrl(cat.has("box_art_url") ?
                                        cat.get("box_art_url").asText() : null)
                                .build());
                    }
                }
            }

            return categories;

        } catch (Exception e) {
            log.error("Fehler beim Suchen von Kategorien", e);
            throw new RuntimeException("Konnte Kategorien nicht suchen", e);
        }
    }

    private HttpHeaders createHeaders(String accessToken) {
        HttpHeaders headers = new HttpHeaders();
        headers.set("Authorization", "Bearer " + accessToken);
        headers.set("Client-Id", clientId);
        return headers;
    }
}
