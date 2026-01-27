package streamer_website.demo.dto.twitch;

import lombok.AllArgsConstructor;
import lombok.Getter;

import java.time.Instant;

@Getter
@AllArgsConstructor
public class TwitchCommandDTO {
    private Long id;
    private String trigger;
    private String response;
    private boolean modOnly;
    private Instant createdAt;
    private Instant updatedAt;
}

