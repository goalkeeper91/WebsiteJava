package streamer_website.demo.dto;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;

import java.time.Instant;

@Data
@Builder
@AllArgsConstructor
public class TokenInfoDto {
    private boolean tokenPresent;
    private Instant expiresAt;
    private String twitchUserId;
}
