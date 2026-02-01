package streamer_website.demo.service.twitch;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import streamer_website.demo.dto.twitch.ActivityDto;
import streamer_website.demo.entity.twitch.StreamActivity;
import streamer_website.demo.repository.StreamActivityRepository;
import streamer_website.demo.websocket.ActivityFeedWebSocketHandler;

import java.time.Instant;
import java.util.List;

@Slf4j
@Service
@RequiredArgsConstructor
public class ActivityService {

    private final StreamActivityRepository activityRepository;
    private final ActivityFeedWebSocketHandler webSocketHandler;

    /**
     * Speichert eine neue Activity und sendet sie via WebSocket
     */
    public void createActivity(
            String twitchUserId,
            String type,
            String username,
            String displayName,
            Integer viewers,
            Integer bits,
            String tier,
            String message
    ) {
        // In DB speichern
        StreamActivity activity = new StreamActivity();
        activity.setTwitchUserId(twitchUserId);
        activity.setType(type);
        activity.setUsername(username);
        activity.setDisplayName(displayName);
        activity.setViewers(viewers);
        activity.setBits(bits);
        activity.setTier(tier);
        activity.setMessage(message);
        activity.setTimestamp(Instant.now());

        StreamActivity saved = activityRepository.save(activity);

        log.info("Activity erstellt: type={}, user={}", type, username);

        // Via WebSocket broadcasten
        ActivityDto dto = ActivityDto.builder()
                .type(type)
                .username(username)
                .displayName(displayName)
                .timestamp(saved.getTimestamp())
                .viewers(viewers)
                .bits(bits)
                .tier(tier)
                .message(message)
                .build();

        webSocketHandler.sendActivityToUser(twitchUserId, dto);
    }

    /**
     * Lädt letzte N Activities für einen User
     */
    public List<ActivityDto> getRecentActivities(String twitchUserId, int limit) {
        List<StreamActivity> activities = activityRepository
                .findTopByTwitchUserIdOrderByTimestampDesc(twitchUserId, limit);

        return activities.stream()
                .map(this::toDto)
                .toList();
    }

    private ActivityDto toDto(StreamActivity activity) {
        return ActivityDto.builder()
                .type(activity.getType())
                .username(activity.getUsername())
                .displayName(activity.getDisplayName())
                .timestamp(activity.getTimestamp())
                .viewers(activity.getViewers())
                .bits(activity.getBits())
                .tier(activity.getTier())
                .message(activity.getMessage())
                .build();
    }
}
