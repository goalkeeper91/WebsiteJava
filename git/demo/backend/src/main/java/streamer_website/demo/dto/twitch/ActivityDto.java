package streamer_website.demo.dto.twitch;

import jakarta.persistence.*;
import lombok.*;
import java.time.Instant;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class ActivityDto {
    private Long id;
    private String twitchUserId;
    private String type;
    private String username;
    private String displayName;
    private Integer viewers;
    private Integer bits;
    private String tier;
    private String message;
    private Instant timestamp;
}