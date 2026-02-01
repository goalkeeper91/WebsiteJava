package streamer_website.demo.service;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.data.redis.connection.Message;
import org.springframework.data.redis.connection.MessageListener;
import org.springframework.data.redis.core.StringRedisTemplate;
import org.springframework.data.redis.listener.ChannelTopic;
import org.springframework.data.redis.listener.RedisMessageListenerContainer;
import org.springframework.stereotype.Service;
import streamer_website.demo.service.twitch.ActivityService;

import jakarta.annotation.PostConstruct;
import jakarta.annotation.PreDestroy;
import java.util.HashMap;
import java.util.Map;

@Slf4j
@Service
@RequiredArgsConstructor
public class BotSignalService implements MessageListener {

    private final StringRedisTemplate redisTemplate;
    private final RedisMessageListenerContainer messageListenerContainer;
    private final ActivityService activityService;
    private final ObjectMapper objectMapper;

    private static final String OUTGOING_CHANNEL = "bot:events";
    private static final String INCOMING_CHANNEL = "backend:events";

    @PostConstruct
    public void init() {
        messageListenerContainer.addMessageListener(
                this,
                new ChannelTopic(INCOMING_CHANNEL)
        );
        log.info("✅ BotSignalService subscribed to channel: {}", INCOMING_CHANNEL);
    }

    @PreDestroy
    public void cleanup() {
        messageListenerContainer.removeMessageListener(this);
        log.info("🛑 BotSignalService unsubscribed");
    }

    /**
     * Sendet Signal an Bot um Channel zu joinen
     */
    public void sendBotJoinSignal(String twitchUserId) {
        Map<String, Object> payload = Map.of(
                "type", "JOIN_CHANNEL",
                "twitch_user_id", twitchUserId
        );

        publishToBot(payload);
        log.info("📤 JOIN_CHANNEL Signal gesendet: user={}", twitchUserId);
    }

    /**
     * Sendet Signal an Bot, um Commands neu zu laden
     */
    public void sendCommandUpdateSignal(String twitchUserId) {
        Map<String, Object> payload = Map.of(
                "type", "REFRESH_COMMANDS",
                "twitch_user_id", twitchUserId
        );

        publishToBot(payload);
        log.info("📤 REFRESH_COMMANDS Signal gesendet: user={}", twitchUserId);
    }

    /**
     * Allgemeine Publish-Methode
     */
    private void publishToBot(Map<String, Object> payload) {
        try {
            String json = objectMapper.writeValueAsString(payload);
            redisTemplate.convertAndSend(OUTGOING_CHANNEL, json);
        } catch (JsonProcessingException e) {
            log.error("❌ Fehler beim Serialisieren des Payloads", e);
        }
    }

    @Override
    public void onMessage(Message message, byte[] pattern) {
        String channel = new String(message.getChannel());
        String body = new String(message.getBody());

        log.debug("📥 Redis Message empfangen: channel={}, body={}", channel, body);

        try {
            JsonNode payload = objectMapper.readTree(body);
            handleIncomingMessage(payload);
        } catch (Exception e) {
            log.error("❌ Fehler beim Verarbeiten der Redis-Nachricht: {}", body, e);
        }
    }

    private void handleIncomingMessage(JsonNode payload) {
        String type = payload.has("type") ? payload.get("type").asText() : null;

        if (type == null) {
            log.warn("⚠️ Message ohne 'type' empfangen: {}", payload);
            return;
        }

        switch (type) {
            case "ACTIVITY" -> handleActivityEvent(payload);
            case "CHAT_MESSAGE" -> handleChatMessage(payload);
            case "BOT_STATUS" -> handleBotStatus(payload);
            default -> log.warn("⚠️ Unbekannter Message-Typ: {}", type);
        }
    }

    /**
     * Verarbeitet Activity Events vom Bot (Follow, Sub, Raid, etc.)
     */
    private void handleActivityEvent(JsonNode payload) {
        try {
            String twitchUserId = payload.get("twitch_user_id").asText();
            String activityType = payload.get("activity_type").asText();
            String username = payload.get("username").asText();
            String displayName = payload.get("display_name").asText();

            Integer viewers = payload.has("viewers") ? payload.get("viewers").asInt() : null;
            Integer bits = payload.has("bits") ? payload.get("bits").asInt() : null;
            String tier = payload.has("tier") ? payload.get("tier").asText() : null;
            String message = payload.has("message") ? payload.get("message").asText() : null;

            activityService.createActivity(
                    twitchUserId,
                    activityType,
                    username,
                    displayName,
                    viewers,
                    bits,
                    tier,
                    message
            );

            log.info("✅ Activity verarbeitet: type={}, user={}", activityType, username);

        } catch (Exception e) {
            log.error("❌ Fehler beim Verarbeiten der Activity", e);
        }
    }

    /**
     * Verarbeitet Chat-Messages vom Bot
     * TODO: Für Chat-Integration implementieren
     */
    private void handleChatMessage(JsonNode payload) {
        log.debug("💬 Chat Message empfangen: {}", payload);
        // TODO: WebSocket an Frontend weiterleiten
    }

    /**
     * Verarbeitet Bot-Status Updates
     */
    private void handleBotStatus(JsonNode payload) {
        String status = payload.has("status") ? payload.get("status").asText() : "unknown";
        log.info("🤖 Bot Status: {}", status);
        // TODO: Status in DB speichern oder WebSocket broadcast
    }
}