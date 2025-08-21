package streamer_website.demo.controller.twitch;

import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import streamer_website.demo.dto.BotStatusDto;
import streamer_website.demo.handler.twitch.TwitchBotHandler;
import streamer_website.demo.service.twitch.TwitchBotManagerService;
import streamer_website.demo.service.twitch.TwitchTokenService;

@RestController
@RequestMapping("/api/bot")
@RequiredArgsConstructor
public class TwitchBotController {

    private final TwitchBotManagerService botManagerService;
    private final TwitchTokenService twitchTokenService;
    private final TwitchBotHandler twitchBotHandler;
    private static final Logger logger = LoggerFactory.getLogger(TwitchBotController.class);

    @PostMapping("/start")
    public ResponseEntity<String> startBot() {
        logger.info("[POST] /api/bot/start called");
        botManagerService.startBot();
        return ResponseEntity.ok("Bot wird gestartet");
    }

    @PostMapping("/stop")
    public ResponseEntity<String> stopBot() {
        logger.info("[POST] /api/bot/stop called");
        botManagerService.stopBot();
        return ResponseEntity.ok("Bot gestoppt");
    }

    @GetMapping("/status")
    public ResponseEntity<BotStatusDto> getBotStatus() {
        logger.info("[GET] /api/bot/status called");
        return ResponseEntity.ok(botManagerService.getBotStatus());
    }

    @GetMapping("/oauth/url")
    public ResponseEntity<String> getOAuthUrl(@RequestParam boolean forBot) {
        logger.info("[GET] /api/bot/oauth/url called with forBot={}", forBot);
        String url = twitchTokenService.buildOAuthUrl(forBot);
        return ResponseEntity.ok(url);
    }

    @GetMapping("/oauth/callback")
    public ResponseEntity<String> oauthCallback(@RequestParam String code) {
        return twitchBotHandler.handleOAuthCallback(code);
    }
}
