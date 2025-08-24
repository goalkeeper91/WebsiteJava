package streamer_website.demo.controller.twitch;

import jakarta.servlet.http.HttpServletResponse;
import jakarta.servlet.http.HttpSession;
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
import streamer_website.demo.service.twitch.TwitchService;
import streamer_website.demo.service.UserService;
import streamer_website.demo.service.twitch.TwitchTokenService;

import java.io.IOException;
import java.net.URLEncoder;
import java.nio.charset.StandardCharsets;

@Controller
@RequestMapping("/api/auth")
public class TwitchAuthController {

    @Value("${twitch.clientId}")
    private String clientId;

    @Value("${twitch.redirectUri}")
    private String redirectUri;

    @Value("${app.authorized-username}")
    private String authorizedUsername;

    @Value("${app.frontend-url}")
    private String frontendUrl;

    private final TwitchService twitchAuthService;

    private final TwitchTokenService tokenService;

    private final UserService userService;

    private static final Logger logger = LoggerFactory.getLogger(TwitchAuthController.class);

    public TwitchAuthController(TwitchService twitchAuthService, UserService userService, TwitchTokenService tokenService) {
        this.twitchAuthService = twitchAuthService;
        this.userService = userService;
        this.tokenService = tokenService;
    }

    @GetMapping("/twitch")
    public void redirectToTwitch(HttpServletResponse response) throws IOException {
        String scopes = "user:read:email bits:read chat:read chat:edit channel:manage:ads channel:read:ads " +
                "channel:manage:broadcast channel:read:editors channel:manage:extensions channel:read:goals " +
                "channel:read:charity whispers:read moderator:read:followers";

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

        TwitchAuthToken token = tokenService.exchangeCodeForToken(code, false);

        TwitchUser twitchUser = twitchAuthService.getUserInfo(token.getUserName());

        if (!twitchUser.username().equalsIgnoreCase(authorizedUsername)) {
            logger.warn("Unauthorized Twitch login attempt: {}", twitchUser.username());
            response.sendRedirect(frontendUrl + "/unauthorized");
            return;
        }

        try {
            userService.createOrUpdate(twitchUser);

            session.setAttribute("user", twitchUser);

            response.sendRedirect(frontendUrl + "/");
        } catch (Exception e) {
            logger.error("Error while updating or creating user", e);
            response.sendRedirect(frontendUrl + "/error");
        }
    }
}
