package streamer_website.demo.scheduler;

import discord4j.common.util.Snowflake;
import discord4j.core.GatewayDiscordClient;
import discord4j.core.object.entity.channel.MessageChannel;
import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Component;
import reactor.core.publisher.Mono;
import streamer_website.demo.service.TaskService;

@Component
@RequiredArgsConstructor
public class DiscordWeeklyTaskScheduler {
    private static final Logger logger = LoggerFactory.getLogger(DiscordWeeklyTaskScheduler.class);

    private final GatewayDiscordClient client;
    private final TaskService taskService;

    private static final Long TASK_CHANNEL_ID = 1424810921963819059L;

    // Führt diesen Job jeden Montag um 09:00 Uhr aus (Cron-Format: Minute Stunde TagDesMonats Monat TagDerWoche)
    @Scheduled(cron = "0 0 18 * * SUN")
    public void postWeeklySummary() {
        if (client == null) {
            logger.warn("Discord Client ist nicht initialisiert. Scheduler übersprungen.");
            return;
        }

        logger.info("Starte wöchentlichen Task-Übersichts-Post.");

        // Generiere das Embed
        var summaryEmbed = taskService.createWeeklySummaryEmbed();

        // Sende das Embed in den konfigurierten Kanal
        client.getChannelById(Snowflake.of(TASK_CHANNEL_ID))
                .cast(MessageChannel.class)
                .flatMap(channel -> channel.createMessage(summaryEmbed))
                .onErrorResume(e -> {
                    logger.error("Fehler beim Posten der wöchentlichen Task-Übersicht", e);
                    return Mono.empty();
                })
                .subscribe();
    }
}
