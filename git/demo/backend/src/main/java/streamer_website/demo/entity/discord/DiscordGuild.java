package streamer_website.demo.entity.discord;

import jakarta.persistence.*;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.time.Instant;
import java.util.List;

@Entity
@Getter
@Setter
@NoArgsConstructor
@Table(name = "discord_guilds")
public class DiscordGuild {

    @Id
    private Long id;

    private String name;

    @Column(name = "icon")
    private String iconUrl;

    @ManyToOne
    @JoinColumn(name = "owner_id")
    private DiscordUser owner;

    @Column(name = "bot_added")
    private Boolean botAdded;

    @Column(name = "created_at", updatable = false)
    private Instant createdAt;

    @Column(name = "updated_at")
    private Instant updatedAt;

    @OneToMany(mappedBy = "guild", cascade = CascadeType.ALL)
    private List<DiscordVoiceChannel> voiceChannels;

    @OneToMany(mappedBy = "guild", cascade = CascadeType.ALL)
    private List<DiscordGuildRole> roles;

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
