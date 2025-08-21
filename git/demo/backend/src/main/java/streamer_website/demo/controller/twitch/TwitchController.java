package streamer_website.demo.controller.twitch;

import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import streamer_website.demo.entity.twitch.TwitchChannelStats;
import streamer_website.demo.handler.twitch.TwitchHandler;
import streamer_website.demo.service.twitch.TwitchService;

import java.util.Map;
import java.util.Optional;

@RestController
@RequestMapping("/api/twitch")
public class TwitchController {

    private final TwitchHandler twitchHandler;

    public TwitchController(TwitchHandler twitchHandler){
        this.twitchHandler = twitchHandler;
    }

    @GetMapping("/status")
    public ResponseEntity<Map<String, Boolean>> getStreamStatus(@RequestParam String username) {
        boolean isLive = twitchHandler.checkLiveStatus(username);
        return ResponseEntity.ok(Map.of("live", isLive));
    }

    @GetMapping("/stats")
    public ResponseEntity<TwitchChannelStats> getStats(@RequestParam String username) {
        return twitchHandler.getStats(username)
                .map(ResponseEntity::ok)
                .orElse(ResponseEntity.badRequest().build());
    }

    @PostMapping("/stats/refresh")
    public ResponseEntity<TwitchChannelStats> refreshStats(@RequestParam String username) {
        return twitchHandler.refreshStats(username)
                .map(ResponseEntity::ok)
                .orElseGet(() -> ResponseEntity.status(HttpStatus.NOT_FOUND).build());
    }
}
