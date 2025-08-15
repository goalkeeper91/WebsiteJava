package streamer_website.demo.dto;

import lombok.AllArgsConstructor;
import lombok.Data;

import java.time.Instant;
import java.util.List;

@Data
@AllArgsConstructor
public class BotStatusDto {
    private boolean running;
    private boolean tokenPresent;
    private Instant expiresAt;
    private String primaryChannel;
    private List<String> activeChannels;
    private Long uptimeSeconds;
    private String version;
}
