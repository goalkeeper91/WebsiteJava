package streamer_website.demo.dto;

import lombok.AllArgsConstructor;
import lombok.Data;

import java.time.Instant;

@Data
@AllArgsConstructor
public class BotStatusDto {
    private boolean running;
    private boolean tokenPresent;
    private Long expiresAt;
    private String channelName;
}
