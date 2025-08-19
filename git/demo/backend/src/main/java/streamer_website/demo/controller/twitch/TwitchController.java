package streamer_website.demo.controller.twitch;

import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import streamer_website.demo.entity.twitch.TwitchChannelStats;
import streamer_website.demo.service.twitch.TwitchService;

import java.util.Map;
import java.util.Optional;

@RestController
@RequestMapping("/twitch")
public class TwitchController {

    private final TwitchService twitchService;

    public TwitchController(TwitchService twitchService){
        this.twitchService = twitchService;
    }

    @GetMapping("/status")
    public ResponseEntity<Map<String, Boolean>> getStreamStatus() {
        boolean isLive = twitchService.isLive();
        return ResponseEntity.ok(Map.of("live", isLive));
    }

    @GetMapping("/stats")
    public ResponseEntity<TwitchChannelStats> getStats() {
        Optional<TwitchChannelStats> stats = twitchService.getLatestChannelStats();

        return stats.map(ResponseEntity::ok)
                .orElse(ResponseEntity.badRequest().build());
    }

    @PostMapping("/stats/refresh")
    public ResponseEntity<TwitchChannelStats> refreshStats() {
        Optional<TwitchChannelStats> stats = twitchService.refreshAndSaveChannelStats();
        return stats.map(ResponseEntity::ok)
                .orElseGet(() -> ResponseEntity.status(HttpStatus.NOT_FOUND).build());
    }
}
