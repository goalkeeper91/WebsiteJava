package streamer_website.demo.service.twitch;

import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Service;

import java.time.Instant;

@Service
@RequiredArgsConstructor
public class TwitchBotSchedulerService {

    private final TwitchBotManagerService botManager;
    private final TwitchTokenService tokenService;
    private static final Logger logger = LoggerFactory.getLogger(TwitchBotSchedulerService.class);

    @Scheduled(fixedRate = 15 * 60 * 1000) // alle 15 Minuten prüfen
    public void refreshBotIfNeeded() {
        var token = tokenService.findBotToken();
        if (token == null) {
            logger.warn("Kein Bot-Token vorhanden, Bot kann nicht geprüft werden");
            return;
        }

        Instant expiresAt = token.getCreatedAt().plusSeconds(token.getExpiresIn());
        if (Instant.now().isAfter(expiresAt)) {
            logger.info("Bot-Token abgelaufen – starte Bot neu");
            botManager.restartBot();
        } else {
            logger.debug("Bot-Token noch gültig, kein Restart nötig");
        }
    }

}

