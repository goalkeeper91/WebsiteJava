package streamer_website.demo.controller;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;
import streamer_website.demo.service.BotOAuthService;

@RestController
@RequestMapping("/twitch/bot")
public class TwitchBotAuthController {

    private final BotOAuthService botOAuthService;

    public TwitchBotAuthController(BotOAuthService botOAuthService) {
        this.botOAuthService = botOAuthService;
    }

    @GetMapping("/callback")
    public String botCallback(@RequestParam String code) {
        botOAuthService.saveBotTokenFromCode(code);
        return "Bot token gespeichert - Du kannst den Bot jetzt starten!";
    }
}
