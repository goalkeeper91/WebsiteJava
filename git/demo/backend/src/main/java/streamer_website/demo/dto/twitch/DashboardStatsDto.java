package streamer_website.demo.dto.twitch;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class DashboardStatsDto {
    private boolean isLive;
    private int currentViewers;
    private int followerCount;
    private int subscriberCount;
    private String uptime;

    // Statistiken
    private int followsToday;
    private int subsThisWeek;
    private double avgViewers;
}
