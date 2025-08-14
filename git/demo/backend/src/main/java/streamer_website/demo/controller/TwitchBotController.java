package streamer_website.demo.controller;

import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import streamer_website.demo.dto.BotStatusDto;
import streamer_website.demo.dto.TokenInfoDto;
import streamer_website.demo.service.BotOAuthService;
import streamer_website.demo.service.TwitchBotManagerService;

@RestController
@RequestMapping("/api/bot")
@RequiredArgsConstructor
public class TwitchBotController {

    private final TwitchBotManagerService botManagerService;
    private final BotOAuthService botOAuthService;

    @PostMapping("/start")
    public ResponseEntity<String> startBot() {
        botManagerService.startBot();
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
    public ResponseEntity<String> oauthCallback(@RequestParam String code) {
        botOAuthService.saveBotTokenFromCode(code);
        return ResponseEntity.ok("Bot-Token gespeichert");
    }
}
