package streamer_website.demo.dto.twitch;

import lombok.Data;

import java.time.Instant;

@Data
public class ChatCommandDto {

    private Long id;
    private String trigger;
    private String response;
    private Integer cooldown;
    private Boolean enabled;
    private Instant createdAt;
    private Instant updatedAt;
}

