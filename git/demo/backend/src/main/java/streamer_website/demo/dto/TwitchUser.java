package streamer_website.demo.dto;

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
) {}

