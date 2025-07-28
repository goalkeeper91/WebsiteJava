package streamer_website.demo.service;

import okhttp3.mockwebserver.MockResponse;
import okhttp3.mockwebserver.MockWebServer;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestInstance;
import org.junit.jupiter.api.extension.ExtendWith;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.core.env.Environment;
import org.springframework.test.context.TestPropertySource;
import org.springframework.test.context.junit.jupiter.SpringExtension;
import org.springframework.web.reactive.function.client.WebClient;

import java.io.IOException;

import static org.junit.jupiter.api.Assertions.*;

@ExtendWith(SpringExtension.class)
@TestPropertySource("classpath:application-test.properties")
@TestInstance(TestInstance.Lifecycle.PER_CLASS)
public class TwitchServiceTest {

    private MockWebServer mockWebServer;

    private TwitchService twitchService;

    @BeforeEach
    void setup() throws IOException {
        mockWebServer = new MockWebServer();
        mockWebServer.start();

        WebClient.Builder builder = WebClient.builder().baseUrl(mockWebServer.url("/").toString());

        // TwitchService manuell mit Spring-Values und WebClient bauen
        twitchService = new TwitchService(
                builder,
                "dummy-client-id",
                "dummy-secret",
                "dummy-user",
                mockWebServer.url("/").toString()
        );
    }

    @AfterEach
    void tearDown() throws IOException {
        mockWebServer.shutdown();
    }

    @Test
    void returnsTrueWhenStreamIsLive() throws Exception {
        // Fake OAuth token response
        mockWebServer.enqueue(new MockResponse()
                .setBody("{\"access_token\":\"dummy_token\",\"expires_in\":3600}")
                .addHeader("Content-Type", "application/json"));

        // Fake user lookup response
        mockWebServer.enqueue(new MockResponse()
                .setBody("{\"data\":[{\"id\":\"123456\"}]}")
                .addHeader("Content-Type", "application/json"));

        // Fake stream status (live)
        mockWebServer.enqueue(new MockResponse()
                .setBody("{\"data\":[{\"type\":\"live\"}]}")
                .addHeader("Content-Type", "application/json"));

        boolean result = twitchService.isLive();
        assertTrue(result);
    }

    @Test
    void returnsFalseWhenStreamIsOffline() throws Exception {
        // Fake OAuth token
        mockWebServer.enqueue(new MockResponse()
                .setBody("{\"access_token\":\"dummy_token\",\"expires_in\":3600}")
                .addHeader("Content-Type", "application/json"));

        // Fake user ID
        mockWebServer.enqueue(new MockResponse()
                .setBody("{\"data\":[{\"id\":\"123456\"}]}")
                .addHeader("Content-Type", "application/json"));

        // Offline response (empty array)
        mockWebServer.enqueue(new MockResponse()
                .setBody("{\"data\":[]}")
                .addHeader("Content-Type", "application/json"));

        boolean result = twitchService.isLive();
        assertFalse(result);
    }

    @Test
    void returnsFalseOnApiError() {
        // Simulate an error in API call
        mockWebServer.enqueue(new MockResponse().setResponseCode(500));

        boolean result = twitchService.isLive();
        assertFalse(result);
    }
}
