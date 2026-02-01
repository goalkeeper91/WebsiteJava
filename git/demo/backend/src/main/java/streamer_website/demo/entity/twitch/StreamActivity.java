package streamer_website.demo.entity.twitch;

import jakarta.persistence.*;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.Instant;

@Data
@NoArgsConstructor
@Entity
@Table(
        name = "stream_activities",
        indexes = {
                @Index(name = "idx_activities_user_time",
                        columnList = "twitch_user_id, timestamp DESC")
        }
)
public class StreamActivity {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(name = "twitch_user_id", nullable = false, length = 50)
    private String twitchUserId;

    @Column(nullable = false, length = 20)
    private String type; // FOLLOW, SUB, RAID, CHEER, HOST

    @Column(nullable = false, length = 100)
    private String username;

    @Column(name = "display_name", nullable = false, length = 100)
    private String displayName;

    @Column(nullable = false)
    private Instant timestamp;

    // Optional fields
    @Column
    private Integer viewers; // for RAID

    @Column
    private Integer bits; // for CHEER

    @Column(length = 10)
    private String tier; // for SUB (1000, 2000, 3000)

    @Column(length = 500)
    private String message; // for CHEER

    @PrePersist
    protected void onCreate() {
        if (timestamp == null) {
            timestamp = Instant.now();
        }
    }
}
