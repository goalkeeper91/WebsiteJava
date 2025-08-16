package streamer_website.demo.controller.twitch;

import com.fasterxml.jackson.core.JsonProcessingException;
import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import streamer_website.demo.dto.BotStatusDto;
import streamer_website.demo.dto.TokenInfoDto;
import streamer_website.demo.service.twitch.BotOAuthService;
import streamer_website.demo.service.twitch.TwitchBotManagerService;

@RestController
@RequestMapping("/api/bot")
@RequiredArgsConstructor
public class TwitchBotController {

    private final TwitchBotManagerService botManagerService;
    private final BotOAuthService botOAuthService;
    private static final Logger logger = LoggerFactory.getLogger(TwitchBotController.class);

    @PostMapping("/start")
    public ResponseEntity<String> startBot() {
        System.out.println("🎯 Start-Endpoint wurde aufgerufen");
        botManagerService.startBot();
        System.out.println("✅ startBot() wurde ausgeführt");
        logger.info("[DEBUG] POST /api/bot/start called");
        return ResponseEntity.ok("Bot gestartet");
    }

    @PostMapping("/stop")
    public ResponseEntity<String> stopBot() {
        botManagerService.stopBot();
        return ResponseEntity.ok("Bot gestoppt");
    }

    @GetMapping("/status")
    public ResponseEntity<BotStatusDto> getBotStatus() {
        return ResponseEntity.ok(botManagerService.getBotStatus());
    }

    @GetMapping("/token")
    public ResponseEntity<TokenInfoDto> getTokenInfo() {
        return ResponseEntity.ok(botOAuthService.getTokenInfo());
    }

    @GetMapping("/oauth/url")
    public ResponseEntity<String> getOAuthUrl() {
        return ResponseEntity.ok(botOAuthService.buildOAuthUrl());
    }

    @GetMapping("/oauth/callback")
    public ResponseEntity<String> oauthCallback(@RequestParam String code) throws JsonProcessingException {
        botOAuthService.saveBotTokenFromCode(code);
        return ResponseEntity.ok("Bot-Token gespeichert");
    }
}
