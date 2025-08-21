package streamer_website.demo.service;

import okhttp3.mockwebserver.MockResponse;
import okhttp3.mockwebserver.MockWebServer;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestInstance;
import org.junit.jupiter.api.extension.ExtendWith;
import org.springframework.test.context.TestPropertySource;
import org.springframework.test.context.junit.jupiter.SpringExtension;
import org.springframework.web.reactive.function.client.WebClient;
import streamer_website.demo.repository.TwitchChannelStatsRepository;
import streamer_website.demo.service.twitch.TwitchService;
import streamer_website.demo.service.twitch.TwitchTokenService;

import java.io.IOException;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

@ExtendWith(SpringExtension.class)
@TestPropertySource("classpath:application-test.properties")
@TestInstance(TestInstance.Lifecycle.PER_CLASS)
public class TwitchServiceTest {

    private MockWebServer mockWebServer;
    private TwitchService twitchService;
    private TwitchTokenService tokenService;
    private TwitchChannelStatsRepository twitchChannelStatsRepository;

    @BeforeEach
    void setup() throws IOException {
        mockWebServer = new MockWebServer();
        mockWebServer.start();

        tokenService = mock(TwitchTokenService.class);
        twitchChannelStatsRepository = mock(TwitchChannelStatsRepository.class);

        WebClient.Builder builder = WebClient.builder()
                .baseUrl(mockWebServer.url("/").toString());

        twitchService = new TwitchService(
                tokenService,
                twitchChannelStatsRepository, // falls nötig, auch mocken
                builder
        );
    }

    @Test
    void returnsTrueWhenStreamIsLive() throws Exception {
        when(tokenService.getUserAccessToken(anyString())).thenReturn("dummy_token");

        mockWebServer.enqueue(new MockResponse()
                .setBody("{\"data\":[{\"id\":\"123456\"}]}")
                .addHeader("Content-Type", "application/json"));

        mockWebServer.enqueue(new MockResponse()
                .setBody("{\"data\":[{\"type\":\"live\"}]}")
                .addHeader("Content-Type", "application/json"));

        boolean result = twitchService.isLive("dummy-user");
        assertTrue(result);
    }
}