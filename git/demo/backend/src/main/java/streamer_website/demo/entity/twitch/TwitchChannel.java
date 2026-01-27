package streamer_website.demo.entity.twitch;

import jakarta.persistence.*;
import lombok.*;
import java.time.Instant;

@Entity
@Table(name = "twitch_channels")
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class TwitchChannel {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(name = "twitch_user_id", nullable = false, unique = true, length = 128)
    private String twitchUserId;

    @Column(name = "user_name", nullable = false, length = 128)
    private String userName;

    @Column(name = "is_active", nullable = false)
    private boolean isActive = true;

    @Column(name = "created_at", nullable = false, updatable = false)
    private Instant createdAt = Instant.now();

    @Column(name = "updated_at", nullable = false)
    private Instant updatedAt = Instant.now();

    @PreUpdate
    public void preUpdate() {
        this.updatedAt = Instant.now();
    }
}
