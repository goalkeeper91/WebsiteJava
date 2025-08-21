package streamer_website.demo.handler.twitch;

import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Service;
import streamer_website.demo.entity.twitch.TwitchAuthToken;
import streamer_website.demo.service.twitch.TwitchTokenService;

@Service
@RequiredArgsConstructor
public class TwitchBotHandler {

    private final TwitchTokenService twitchTokenService;

    public ResponseEntity<String> handleOAuthCallback(String code) {
        // Code gegen Token tauschen
        TwitchAuthToken token = twitchTokenService.exchangeCodeForToken(code, true);
        return ResponseEntity.ok("Bot-Token gespeichert für User: " + token.getUserName());
    }
}
