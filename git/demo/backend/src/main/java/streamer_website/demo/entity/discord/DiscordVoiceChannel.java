package streamer_website.demo.entity.discord;

import jakarta.persistence.*;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.time.Instant;

@Entity
@Getter
@Setter
@NoArgsConstructor
@Table(name = "discord_voice_channels")
public class DiscordVoiceChannel {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    private String name;

    private Boolean temporary;

    @Column(name = "user_limit")
    private Integer userLimit;

    private Integer bitrate;

    @Column(name = "category_id")
    private Long categoryId;

    @Column(name = "created_at", updatable = false)
    private Instant createdAt;

    @Column(name = "updated_at")
    private Instant updatedAt;

    @ManyToOne
    @JoinColumn(name = "guild_id")
    private DiscordGuild guild;

    @ManyToOne
    @JoinColumn(name = "created_by")
    private DiscordUser createdBy;

    @PrePersist
    protected void onCreate() {
        createdAt = Instant.now();
        updatedAt = createdAt;
    }

    @PreUpdate
    protected void onUpdate() {
        updatedAt = Instant.now();
    }
}
