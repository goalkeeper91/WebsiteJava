package streamer_website.demo.dto.twitch;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class ActivityDto {
    private String type;        // FOLLOW, SUB, RAID, CHEER, etc.
    private String username;
    private String displayName;
    private Long timestamp;

    // Optional fields je nach Type
    private Integer viewers;    // for RAID
    private Integer bits;       // for CHEER
    private String tier;        // for SUB
    private String message;     // for CHEER
}
