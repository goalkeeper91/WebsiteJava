package streamer_website.demo.controller;

import jakarta.servlet.http.HttpServletResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import streamer_website.demo.dto.TwitchUser;
import streamer_website.demo.service.TwitchService;
import streamer_website.demo.service.UserService;

import java.io.IOException;
import java.net.URLEncoder;
import java.nio.charset.StandardCharsets;

@Controller
@RequestMapping("/auth")
public class TwitchAuthController {

    @Value("${twitch.clientId}")
    private String clientId;

    @Value("${twitch.redirectUri}")
    private String redirectUri;

    @Value("${authorized-username}")
    private String authorizedUsername;

    private TwitchService twitchAuthService;

    private UserService userService;

    private static final Logger logger = LoggerFactory.getLogger(TwitchAuthController.class);

    @GetMapping("/twitch")
    public void redirectToTwitch(HttpServletResponse response) throws IOException {
        String url = "https://id.twitch.tv/oauth2/authorize" +
                "?response_type=code" +
                "&client_id=" + clientId +
                "&redirect_uri=" + URLEncoder.encode(redirectUri, StandardCharsets.UTF_8) +
                "&scope=" + URLEncoder.encode("user:read:email", StandardCharsets.UTF_8);

        response.sendRedirect(url);
    }

    @GetMapping("/twitch/callback")
    public String handleTwitchCallback(@RequestParam("code") String code) {
        String accessToken = twitchAuthService.exchangeCodeForAccessToken(code);

        TwitchUser twitchUser = twitchAuthService.getUserInfo(accessToken);

        if (!twitchUser.login().equalsIgnoreCase(authorizedUsername)) {
            logger.warn("Unauthorized Twitch login attempt: {}", twitchUser.login());
            return "redirect:/unauthorized";
        }

        userService.createOrUpdate(twitchUser);

        return "redirect:/";
    }
}
