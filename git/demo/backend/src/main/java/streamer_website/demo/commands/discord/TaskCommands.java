package streamer_website.demo.commands.discord;

import discord4j.core.event.domain.message.MessageCreateEvent;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Component;
import streamer_website.demo.entity.discord.DiscordTask;
import streamer_website.demo.entity.discord.TaskStatusEnum;
import streamer_website.demo.repository.DiscordTaskRepository;
import streamer_website.demo.service.discord.DiscordNotificationService;

import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;
import java.util.Arrays;
import java.util.List;

@Component
@RequiredArgsConstructor
public class TaskCommands {

    private final DiscordNotificationService notificationService;
    private final DiscordTaskRepository taskRepository;
    private final DateTimeFormatter dateTimeFormatter = DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm");

    public void execute(MessageCreateEvent event, String[] args) {
        if (args.length == 0) {
            sendMessage(event, "Bitte einen Unterbefehl angeben: create, update, delete, weekly");
            return;
        }

        String subCommand = args[0].toLowerCase();
        String rawArgs = args.length > 1 ? String.join(" ", Arrays.copyOfRange(args, 1, args.length)) : "";

        switch (subCommand) {
            case "create" -> createTask(event, rawArgs);
            case "update" -> updateTask(event, rawArgs);
            case "delete" -> deleteTask(event, rawArgs.trim());
            case "weekly" -> weeklyTasks(event);
            default -> sendMessage(event, "Unbekannter Unterbefehl: " + subCommand);
        }
    }

    private void createTask(MessageCreateEvent event, String rawArgs) {
        // Erwartet: Titel | Beschreibung | YYYY-MM-DD HH:mm | [Priorität]
        String[] parts = rawArgs.split("\\|");

        if (parts.length < 3) {
            sendMessage(event, "Usage: /tasks create <Titel> | <Beschreibung> | <YYYY-MM-DD HH:mm> | [Priorität]");
            return;
        }

        try {
            String title = parts[0].trim();
            String description = parts[1].trim();
            LocalDateTime dueDate = LocalDateTime.parse(parts[2].trim(), dateTimeFormatter);
            String priority = parts.length > 3 ? parts[3].trim() : "medium";

            String userId = event.getMessage().getAuthor()
                    .map(user -> user.getId().asString())
                    .orElse("unbekannt");

            DiscordTask task = new DiscordTask();
            task.setTitle(title);
            task.setDescription(description);
            task.setUserId(userId);
            task.setDueDate(dueDate);
            task.setPriority(priority);
            task.setStatus(TaskStatusEnum.NOT_STARTED);

            taskRepository.save(task);
            sendMessage(event, "✅ Task erstellt: " + title);
        } catch (Exception e) {
            sendMessage(event, "Fehler beim Erstellen des Tasks: " + e.getMessage());
        }
    }

    private void updateTask(MessageCreateEvent event, String rawArgs) {
        // Erwartet: TaskID | Titel | Beschreibung | YYYY-MM-DD HH:mm | [Priorität]
        String[] parts = rawArgs.split("\\|");

        if (parts.length < 2) {
            sendMessage(event, "Usage: /tasks update <TaskID> | <Titel> | <Beschreibung> | <YYYY-MM-DD HH:mm> | [Priorität]");
            return;
        }

        try {
            Long taskId = Long.parseLong(parts[0].trim());
            DiscordTask task = taskRepository.findById(taskId).orElse(null);
            if (task == null) {
                sendMessage(event, "Task nicht gefunden: " + taskId);
                return;
            }

            task.setTitle(parts[1].trim());
            if (parts.length > 2) task.setDescription(parts[2].trim());
            if (parts.length > 3) task.setDueDate(LocalDateTime.parse(parts[3].trim(), dateTimeFormatter));
            if (parts.length > 4) task.setPriority(parts[4].trim());

            // UserID bleibt unverändert (aus Sicherheitsgründen)
            taskRepository.save(task);
            sendMessage(event, "✅ Task aktualisiert: " + task.getTitle());
        } catch (Exception e) {
            sendMessage(event, "Fehler beim Aktualisieren des Tasks: " + e.getMessage());
        }
    }

    private void deleteTask(MessageCreateEvent event, String arg) {
        try {
            Long taskId = Long.parseLong(arg);
            taskRepository.deleteById(taskId);
            sendMessage(event, "✅ Task gelöscht: " + taskId);
        } catch (Exception e) {
            sendMessage(event, "Fehler beim Löschen des Tasks: " + e.getMessage());
        }
    }

    private void weeklyTasks(MessageCreateEvent event) {
        List<DiscordTask> tasks = taskRepository.findAll();
        notificationService.notifyWeeklyTasks(tasks).subscribe();
        sendMessage(event, "✅ Wochenübersicht gepostet! (" + tasks.size() + " Tasks insgesamt)");
    }

    private void sendMessage(MessageCreateEvent event, String message) {
        event.getMessage().getChannel()
                .flatMap(channel -> channel.createMessage(message))
                .subscribe();
    }
}
