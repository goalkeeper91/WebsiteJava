package streamer_website.demo.controller;

import org.junit.jupiter.api.Test;
import org.mockito.Mockito;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest;
import org.springframework.test.context.bean.override.mockito.MockitoBean;
import org.springframework.test.web.servlet.MockMvc;
import streamer_website.demo.controller.twitch.TwitchController;
import streamer_website.demo.service.twitch.TwitchService;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath;


@WebMvcTest(TwitchController.class)
public class TwitchControllerTest {

    @Autowired
    private MockMvc mockMvc;

    @MockitoBean
    private TwitchService twitchService;

    @Test
    void shouldReturnLiveStatusTrue() throws Exception {
        Mockito.when(twitchService.isLive()).thenReturn(true);

        mockMvc.perform(get("/twitch/status"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.live").value(true));

    }
}
