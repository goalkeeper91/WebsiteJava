package streamer_website.demo.entity.twitch;

import jakarta.persistence.*;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;
import org.hibernate.annotations.ColumnTransformer;

import java.time.Instant;

@Entity
@Data
@NoArgsConstructor
@AllArgsConstructor
@Table(
        name = "chat_commands",
        uniqueConstraints = @UniqueConstraint(columnNames = {"channel_id", "command"})
)
public class ChatCommand {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(name = "channel_id", nullable = false)
    @ColumnTransformer(
            read = "channel_id::varchar",    // Beim Lesen: BigInt zu Varchar (String)
            write = "?::bigint"               // Beim Schreiben: String zu BigInt
    )
    private String channelId;

    @Column(name = "command", nullable = false, length = 100)
    private String trigger;

    @Column(name = "response", nullable = false, columnDefinition = "TEXT")
    private String response;

    @Column(name = "cooldown", nullable = false)
    private Integer cooldown = 0;

    @Column(name = "enabled", nullable = false)
    private Boolean enabled = true;

    @Column(name = "created_at", nullable = false, updatable = false, columnDefinition = "TIMESTAMP DEFAULT CURRENT_TIMESTAMP")
    private Instant createdAt = Instant.now();

    @PreUpdate
    public void preUpdate() {
        this.updatedAt = Instant.now();
    }

    @Column(name = "updated_at", nullable = false)
    private Instant updatedAt = Instant.now();
}

