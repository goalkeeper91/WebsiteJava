package streamer_website.demo.websocket;

import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Component;
import org.springframework.web.socket.CloseStatus;
import org.springframework.web.socket.TextMessage;
import org.springframework.web.socket.WebSocketSession;
import org.springframework.web.socket.handler.TextWebSocketHandler;
import streamer_website.demo.dto.twitch.ActivityDto;

import java.io.IOException;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

@Slf4j
@Component
@RequiredArgsConstructor
public class ActivityFeedWebSocketHandler extends TextWebSocketHandler {

    private final ObjectMapper objectMapper;

    private final Map<String, WebSocketSession> sessions = new ConcurrentHashMap<>();

    private final Map<String, String> sessionUserMap = new ConcurrentHashMap<>();

    @Override
    public void afterConnectionEstablished(WebSocketSession session) throws Exception {
        String sessionId = session.getId();
        sessions.put(sessionId, session);

        log.info("WebSocket verbunden: sessionId={}", sessionId);

        // TODO: User-ID aus Session extrahieren
        // String userId = extractUserId(session);
        // sessionUserMap.put(sessionId, userId);
    }

    @Override
    public void afterConnectionClosed(WebSocketSession session, CloseStatus status) throws Exception {
        String sessionId = session.getId();
        sessions.remove(sessionId);
        sessionUserMap.remove(sessionId);

        log.info("WebSocket geschlossen: sessionId={}, status={}", sessionId, status);
    }

    @Override
    protected void handleTextMessage(WebSocketSession session, TextMessage message) throws Exception {
        // Optional: Client kann Ping senden
        log.debug("WebSocket Nachricht empfangen: {}", message.getPayload());

        if ("ping".equals(message.getPayload())) {
            session.sendMessage(new TextMessage("pong"));
        }
    }

    /**
     * Sendet eine Activity an alle verbundenen Clients
     */
    public void broadcastActivity(ActivityDto activity) {
        String json;
        try {
            json = objectMapper.writeValueAsString(activity);
        } catch (Exception e) {
            log.error("Fehler beim Serialisieren der Activity", e);
            return;
        }

        TextMessage message = new TextMessage(json);

        // An alle Sessions senden
        sessions.values().forEach(session -> {
            if (session.isOpen()) {
                try {
                    session.sendMessage(message);
                    log.debug("Activity gesendet an session={}", session.getId());
                } catch (IOException e) {
                    log.error("Fehler beim Senden an session={}", session.getId(), e);
                }
            }
        });
    }

    /**
     * Sendet eine Activity nur an einen spezifischen User
     */
    public void sendActivityToUser(String twitchUserId, ActivityDto activity) {
        String json;
        try {
            json = objectMapper.writeValueAsString(activity);
        } catch (Exception e) {
            log.error("Fehler beim Serialisieren der Activity", e);
            return;
        }

        TextMessage message = new TextMessage(json);

        // Finde alle Sessions des Users
        sessionUserMap.entrySet().stream()
                .filter(entry -> twitchUserId.equals(entry.getValue()))
                .map(Map.Entry::getKey)
                .map(sessions::get)
                .filter(session -> session != null && session.isOpen())
                .forEach(session -> {
                    try {
                        session.sendMessage(message);
                        log.debug("Activity gesendet an user={}, session={}",
                                twitchUserId, session.getId());
                    } catch (IOException e) {
                        log.error("Fehler beim Senden an user={}", twitchUserId, e);
                    }
                });
    }

    /**
     * Gibt die Anzahl aktiver Verbindungen zurück
     */
    public int getActiveConnectionCount() {
        return sessions.size();
    }
}