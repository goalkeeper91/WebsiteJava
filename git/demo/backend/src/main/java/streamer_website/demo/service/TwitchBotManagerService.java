package streamer_website.demo.service;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import streamer_website.demo.client.TwitchBot;
import streamer_website.demo.config.TwitchBotConfig;
import streamer_website.demo.dto.BotStatusDto;

@Service
@RequiredArgsConstructor
public class TwitchBotManagerService {

    private TwitchBot twitchBot;

    private final BotOAuthService botOAuthService;
    private final TwitchCommandService twitchCommandService;
    private final TwitchBotConfig twitchBotConfig;

    public void startBot() {
        if (twitchBot == null) {
            String channelName = "goalkeeper91_bot";
            String userId = twitchBotConfig.getBotUserId();

            twitchBot = new TwitchBot(channelName, twitchCommandService, botOAuthService);
            twitchBot.start(userId);
        }
    }

    public void stopBot() {
        if (twitchBot != null) {
            twitchBot.stop();
            twitchBot = null;
        }
    }

    public BotStatusDto getBotStatus() {
        var token = botOAuthService.findBotToken(); // Methode muss in BotOAuthService existieren
        boolean tokenPresent = token != null;
        long expiry = tokenPresent ? token.getExpiresIn() : 0;

        return new BotStatusDto(
                twitchBot.isRunning(),
                tokenPresent,
                expiry,
                twitchBot.getChannelName()
        );
    }
}
