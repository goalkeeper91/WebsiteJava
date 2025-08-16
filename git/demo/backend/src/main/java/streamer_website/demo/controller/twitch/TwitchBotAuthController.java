package streamer_website.demo.controller.twitch;

import com.fasterxml.jackson.core.JsonProcessingException;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.servlet.view.RedirectView;
import streamer_website.demo.service.twitch.BotOAuthService;

@RestController
@RequestMapping("/twitch/bot")
public class TwitchBotAuthController {

    private final BotOAuthService botOAuthService;

    public TwitchBotAuthController(BotOAuthService botOAuthService) {
        this.botOAuthService = botOAuthService;
    }

    @GetMapping("/callback")
    public RedirectView botCallback(@RequestParam String code) throws JsonProcessingException {
        botOAuthService.saveBotTokenFromCode(code);
        return new RedirectView("http://localhost:5173/admin?tokenSaved=true");
    }
}
