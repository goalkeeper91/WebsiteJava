package streamer_website.demo.controller;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest;
import org.springframework.test.web.servlet.MockMvc;

import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.*;


@WebMvcTest(controllers = TwitchAuthController.class)
public class TwitchAuthControllerTest {

    @Autowired
    private MockMvc mockMvc;

    @Test
    void shouldRedirectToTwitchAuth() throws Exception {
        mockMvc.perform(get("/auth/twitch"))
                .andExpect(status().is3xxRedirection())
                .andExpect(redirectedUrlPattern("https://id.twitch.tv/oauth2/authorize*"));
    }

    @Test
    void shouldHandleTwitchCallbackAndAuthorizeUser() throws Exception {
        mockMvc.perform(get("/auth/twitch/callback")
                        .param("code", "dummy_code"))
                .andExpect(status().is3xxRedirection())
                .andExpect(redirectedUrl("/")); // später z. B. /dashboard o. ä.
    }
}
