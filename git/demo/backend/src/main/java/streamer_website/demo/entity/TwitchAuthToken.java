package streamer_website.demo.entity;

import jakarta.persistence.*;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import streamer_website.demo.security.EncryptedStringConverter;

import java.time.Instant;

@Entity
@Table(name = "twitch_auth_tokens")
@Data
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class TwitchAuthToken {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Convert(converter = EncryptedStringConverter.class)
    private String accessToken;
    @Convert(converter = EncryptedStringConverter.class)
    private String refreshToken;
    private String tokenType;
    private Long expiresIn;
    private String scope;
    private String twitchUserId;
    private Instant createdAt;
    private Instant updatedAt;

    @PrePersist
    protected void onCreate() {
        createdAt = Instant.now();
        updatedAt = Instant.now();
    }

    @PreUpdate
    protected void onUpdate() {
        updatedAt = Instant.now();
    }
}
