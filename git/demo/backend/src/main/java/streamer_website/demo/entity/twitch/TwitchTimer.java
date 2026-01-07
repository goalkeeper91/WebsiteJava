package streamer_website.demo.entity.twitch;

import jakarta.persistence.*;
import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
@Entity
@Table(name = "twitch_timers")
public class TwitchTimer {

    // 🔹 Getter & Setter
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(name = "twitch_user_id", nullable = false)
    private String twitchUserId;

    @Column(name = "name")
    private String name;

    @Column(name = "response", nullable = false, columnDefinition = "TEXT")
    private String response;

    @Column(name = "interval_minutes", nullable = false)
    private Integer intervalMinutes;

    @Column(name = "min_chat_messages", nullable = false)
    private Integer minChatMessages = 0;

    @Column(name = "is_enabled", nullable = false)
    private Boolean isEnabled = true;

    public TwitchTimer() {
    }

    public TwitchTimer(String twitchUserId, String name, String response, Integer intervalMinutes) {
        this.twitchUserId = twitchUserId;
        this.name = name;
        this.response = response;
        this.intervalMinutes = intervalMinutes;
        this.minChatMessages = 0;
        this.isEnabled = true;
    }

}

