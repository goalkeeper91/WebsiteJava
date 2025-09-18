package streamer_website.demo.scheduler;

import lombok.RequiredArgsConstructor;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Component;
import streamer_website.demo.entity.discord.DiscordTask;
import streamer_website.demo.repository.DiscordTaskRepository;
import streamer_website.demo.service.discord.DiscordNotificationService;

import java.util.List;

@Component
@RequiredArgsConstructor
public class WeeklyTaskScheduler {

    private final DiscordTaskRepository taskRepository;
    private final DiscordNotificationService notificationService;

    /**
     * Läuft jeden Sonntag um 18:00 Uhr
     */
    @Scheduled(cron = "0 0 18 ? * SUN") // Sekunde Minute Stunde TagMonat Monat Wochentag
    public void postWeeklyTasks() {
        // Alle Tasks laden
        List<DiscordTask> tasks = taskRepository.findAll();
        // Wochenübersicht posten
        notificationService.notifyWeeklyTasks(tasks).subscribe();
    }
}
