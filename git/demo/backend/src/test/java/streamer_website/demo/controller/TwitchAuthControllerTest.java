package streamer_website.demo.controller;

import org.junit.jupiter.api.Test;
import org.mockito.Mockito;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest;
import org.springframework.test.context.bean.override.mockito.MockitoBean;
import org.springframework.test.web.servlet.MockMvc;
import streamer_website.demo.dto.TwitchUser;
import streamer_website.demo.service.TwitchService;
import streamer_website.demo.service.UserService;

import java.io.IOException;

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
        String dummyAccessToken = "dummy_token";
        TwitchUser dummyUser = new TwitchUser("1", "goalkeeper91", "goalkeeper91", "marcel-tu@hotmail.de");

        Mockito.when(twitchAuthService.exchangeCodeForAccessToken("dummy_code")).thenReturn(dummyAccessToken);
        Mockito.when(twitchAuthService.getUserInfo(dummyAccessToken)).thenReturn(dummyUser);

        mockMvc.perform(get("/auth/twitch/callback")
                        .param("code", "dummy_code"))
                .andExpect(status().is3xxRedirection())
                .andExpect(redirectedUrl("/"));
    }
}
