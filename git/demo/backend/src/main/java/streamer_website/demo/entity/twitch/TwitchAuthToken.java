package streamer_website.demo.entity.twitch;

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

    @Column(name = "twitch_user_id", unique = true, nullable = false)
    private String twitchUserId;

    @Column(name = "user_name")
    private String userName;

    @Convert(converter = EncryptedStringConverter.class)
    @Column(name = "access_token", length = 1024)
    private String accessToken;

    @Convert(converter = EncryptedStringConverter.class)
    @Column(name = "refresh_token", length = 1024)
    private String refreshToken;

    @Column(name = "token_type")
    private String tokenType;

    @Column(name = "expires_in")
    private Long expiresIn;

    @Column(columnDefinition = "TEXT")
    private String scope;

    @Enumerated(EnumType.STRING)
    @Column(nullable = false)
    private TwitchTokenType tokenOwner;

    @Column(name = "created_at", updatable = false)
    private Instant createdAt;

    @Column(name = "updated_at")
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
