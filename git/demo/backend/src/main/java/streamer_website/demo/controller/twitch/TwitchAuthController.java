package streamer_website.demo.controller.twitch;

import jakarta.servlet.http.HttpServletResponse;
import jakarta.servlet.http.HttpSession;
import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import streamer_website.demo.dto.TwitchTokenResponse;
import streamer_website.demo.dto.TwitchUser;
import streamer_website.demo.entity.twitch.TwitchAuthToken;
import streamer_website.demo.service.BotSignalService;
import streamer_website.demo.service.twitch.TwitchService;
import streamer_website.demo.service.UserService;
import streamer_website.demo.service.twitch.TwitchTokenService;

import java.io.IOException;
import java.net.URLEncoder;
import java.nio.charset.StandardCharsets;

@Controller
@RequestMapping("/auth")
@RequiredArgsConstructor
public class TwitchAuthController {

    @Value("${twitch.clientId}")
    private String clientId;

    @Value("${twitch.redirectUri}")
    private String redirectUri;

    @Value("${app.frontend-url}")
    private String frontendUrl;

    private final TwitchService twitchService;
    private final TwitchTokenService tokenService;
    private final UserService userService;
    private final BotSignalService botSignalService;

    private static final Logger logger = LoggerFactory.getLogger(TwitchAuthController.class);

    @GetMapping("/twitch")
    public void redirectToTwitch(HttpServletResponse response) throws IOException {
        // Scopes für Streamer: E-Mail, Chat lesen/schreiben und Stream-Infos verwalten
        String scopes = "user:read:email chat:read chat:edit channel:manage:broadcast";

        String url = "https://id.twitch.tv/oauth2/authorize" +
                "?response_type=code" +
                "&client_id=" + clientId +
                "&redirect_uri=" + URLEncoder.encode(redirectUri, StandardCharsets.UTF_8) +
                "&scope=" + URLEncoder.encode(scopes, StandardCharsets.UTF_8);

        response.sendRedirect(url);
    }

    @GetMapping("/twitch/callback")
    public void handleTwitchCallback(@RequestParam("code") String code,
                                     HttpServletResponse response,
                                     HttpSession session) throws IOException {
        try {
            TwitchAuthToken tokenEntity = tokenService.exchangeCodeForToken(code, false);

            TwitchUser twitchUser = twitchService.getUserInfo(tokenEntity.getUserName());

            userService.createOrUpdate(twitchUser);

            session.setAttribute("user", twitchUser);

            botSignalService.sendBotJoinSignal(twitchUser.id());

            response.sendRedirect(frontendUrl + "/dashboard");

        } catch (Exception e) {
            logger.error("Fehler beim Twitch-Login Workflow", e);
            response.sendRedirect(frontendUrl + "/error?msg=auth_failed");
        }
    }
}
