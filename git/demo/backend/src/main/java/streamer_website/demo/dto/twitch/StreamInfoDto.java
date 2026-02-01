package streamer_website.demo.dto.twitch;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class StreamInfoDto {
    private String broadcasterId;
    private String broadcasterName;
    private String title;
    private String gameId;
    private String gameName;
}
