package streamer_website.demo.entity;

import jakarta.persistence.*;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.Instant;

@Entity
@Table(name = "twitch_channel_stats")
@Data
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class TwitchChannelStats {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(unique = true)
    private String twitchUserId;
    private String displayName;
    private String description;
    private String profileImageUrl;
    private String offlineImageUrl;
    private String broadcasterType;
    private int viewCount;
    private int followers;

    private Instant accountCreatedAt;
    private Instant fetchedAt;

    @PrePersist
    protected void onCreate() {
        fetchedAt = Instant.now();
    }
}
