package streamer_website.demo.entity.discord;

import com.fasterxml.jackson.annotation.JsonProperty;
import discord4j.common.util.Snowflake;
import jakarta.persistence.*;
import lombok.Data;

import java.time.LocalDateTime;

@Entity
@Table(name = "discord_join_to_create_channels")
@Data
public class JoinToCreateConfig {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    // Trigger Channel
    @Column(name = "join_channel_id", nullable = false)
    private Long joinChannelId;

    @Column(name = "category_id", nullable = false)
    private Long categoryId;

    @Column(name = "channel_name_prefix", nullable = false)
    private String channelNamePrefix;

    @Column(name = "user_limit")
    private Integer userLimit;

    @Column(name = "is_private")
    @JsonProperty("privateChannel")
    private Boolean privateChannel;

    @Column(name = "created_at", nullable = false, updatable = false)
    private LocalDateTime createdAt;

    @Column(name = "updated_at")
    private LocalDateTime updatedAt;

    @PrePersist
    protected void onCreate() {
        createdAt = LocalDateTime.now();
        updatedAt = LocalDateTime.now();
    }

    @PreUpdate
    protected void onUpdate() {
        updatedAt = LocalDateTime.now();
    }

    /** Hilfsmethode zum Konvertieren in Snowflake, falls du Discord4J-Objekte brauchst */
    public Snowflake joinChannelSnowflake() {
        return Snowflake.of(joinChannelId);
    }

    public Snowflake categorySnowflake() {
        return Snowflake.of(categoryId);
    }
}
