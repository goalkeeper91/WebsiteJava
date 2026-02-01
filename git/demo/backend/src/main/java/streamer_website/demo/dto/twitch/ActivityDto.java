package streamer_website.demo.dto.twitch;

import jakarta.persistence.*;
import lombok.*;
import java.time.Instant;

@Entity
@Data
@NoArgsConstructor
@AllArgsConstructor
public class ActivityDto {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
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