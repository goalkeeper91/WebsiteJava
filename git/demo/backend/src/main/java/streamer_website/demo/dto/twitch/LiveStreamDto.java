package streamer_website.demo.dto.twitch;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class LiveStreamDto {
    private boolean isLive;
    private int viewerCount;
    private String startedAt;
    private String title;
    private String gameName;
    private String thumbnailUrl;
}
