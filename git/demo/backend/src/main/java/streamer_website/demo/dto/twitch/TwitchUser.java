package streamer_website.demo.dto.twitch;

import java.io.Serial;
import java.io.Serializable;
import java.time.Instant;

public record TwitchUser(
        String id,
        String login,
        String username,
        String email,
        String description,
        String profileImageUrl,
        String offlineImageUrl,
        String broadcasterType,
        int viewCount,
        Instant createdAt
) implements Serializable {
    @Serial
    private static final long serialVersionUID = 1L;
}

