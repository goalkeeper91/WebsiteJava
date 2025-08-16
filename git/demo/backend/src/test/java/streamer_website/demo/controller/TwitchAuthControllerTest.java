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

import java.util.List;

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
        mockMvc.perform(get("/auth/twitch"))
                .andExpect(status().is3xxRedirection())
                .andExpect(redirectedUrlPattern("https://id.twitch.tv/oauth2/authorize*"));
    }

    @Test
    void shouldHandleTwitchCallbackAndAuthorizeUser() throws Exception {
        TwitchUser dummyUser = new TwitchUser("1", "goalkeeper91", "goalkeeper91", "marcel-tu@hotmail.de");
        TwitchTokenResponse mockTokenResponse = TwitchTokenResponse.builder()
                .accessToken("some-access-token")
                .refreshToken("some-refresh-token")
                .expiresIn(3600L)
                .scope(List.of("user:read:email"))
                .tokenType("bearer")
                .build();

        when(twitchAuthService.exchangeCodeForAccessToken(anyString()))
                .thenReturn(mockTokenResponse);

        when(twitchAuthService.exchangeCodeForAccessToken("dummy_code")).thenReturn(mockTokenResponse);
        when(twitchAuthService.getUserInfo(mockTokenResponse.getAccessToken())).thenReturn(dummyUser);

        mockMvc.perform(get("/auth/twitch/callback")
                        .param("code", "dummy_code"))
                .andExpect(status().is3xxRedirection())
                .andExpect(redirectedUrl("/"));
    }
}
