package streamer_website.demo.entity;

import jakarta.persistence.*;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;
import streamer_website.demo.security.EncryptedStringConverter;

@Entity
@Table(name="videos")
@Data
@NoArgsConstructor
@AllArgsConstructor
public class Video {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    private String platform;

    @Column(nullable = true)
    private String title;
    @Column(nullable = true)
    private String videoId;

    @Convert(converter = EncryptedStringConverter.class)
    @Column(nullable = true)
    private String apiKey;
    @Convert(converter = EncryptedStringConverter.class)
    @Column(nullable = true)
    private String channelId;
}
