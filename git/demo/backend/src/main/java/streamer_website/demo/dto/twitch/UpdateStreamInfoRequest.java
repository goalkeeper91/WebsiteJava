package streamer_website.demo.dto.twitch;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class UpdateStreamInfoRequest {
    private String title;
    private String gameId;
}
