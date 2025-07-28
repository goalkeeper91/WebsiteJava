package streamer_website.demo.controller;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import streamer_website.demo.service.TwitchService;

import java.util.Map;

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
}
