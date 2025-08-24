package streamer_website.demo.controller;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest;
import org.springframework.test.context.bean.override.mockito.MockitoBean;
import org.springframework.test.web.servlet.MockMvc;
import streamer_website.demo.controller.twitch.TwitchAuthController;
import streamer_website.demo.dto.TwitchTokenResponse;
import streamer_website.demo.dto.TwitchUser;
import streamer_website.demo.service.twitch.TwitchService;
import streamer_website.demo.service.UserService;

import java.time.Instant;
import java.util.List;

import static org.mockito.ArgumentMatchers.anyBoolean;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.when;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.*;


@WebMvcTest(controllers = TwitchAuthController.class)
public class TwitchAuthControllerTest {

    @Autowired
    private MockMvc mockMvc;

    @MockitoBean
    private TwitchService twitchAuthService;

    @MockitoBean
    private UserService userService;

    @Value("${app.authorized-username}")
    private String authorizedUsername;


    @Test
    void shouldRedirectToTwitchAuth() throws Exception {
        mockMvc.perform(get("/api/auth/twitch"))
                .andExpect(status().is3xxRedirection())
                .andExpect(redirectedUrlPattern("https://id.twitch.tv/oauth2/authorize*"));
    }

    @Test
    void shouldHandleTwitchCallbackAndAuthorizeUser() throws Exception {
        TwitchUser dummyUser = new TwitchUser(
                "1",                     // id
                "goalkeeper91",          // login
                "goalkeeper91",          // username
                "marcel-tu@hotmail.de",  // email
                "dummy description",     // description
                "http://example.com/profile.jpg", // profileImageUrl
                "http://example.com/offline.jpg", // offlineImageUrl
                "affiliate",             // broadcasterType
                1234,                    // viewCount
                Instant.now()            // createdAt
        );
        TwitchTokenResponse mockTokenResponse = TwitchTokenResponse.builder()
                .accessToken("some-access-token")
                .refreshToken("some-refresh-token")
                .expiresIn(3600L)
                .scope(List.of("user:read:email"))
                .tokenType("bearer")
                .build();

        when(twitchAuthService.exchangeCodeForAccessToken(anyString(), anyBoolean()))
                .thenReturn(mockTokenResponse);

        when(twitchAuthService.exchangeCodeForAccessToken("dummy_code", true)).thenReturn(mockTokenResponse);
        when(twitchAuthService.getUserInfo(mockTokenResponse.getAccessToken())).thenReturn(dummyUser);

        mockMvc.perform(get("/api/auth/twitch/callback")
                        .param("code", "dummy_code"))
                .andExpect(status().is3xxRedirection())
                .andExpect(redirectedUrl("/"));
    }
}
