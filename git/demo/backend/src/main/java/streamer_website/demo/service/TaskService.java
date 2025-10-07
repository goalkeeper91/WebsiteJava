package streamer_website.demo.service;

import discord4j.common.util.Snowflake;
import discord4j.core.GatewayDiscordClient;
import discord4j.core.object.reaction.ReactionEmoji;
import discord4j.core.spec.EmbedCreateSpec;
import discord4j.rest.util.Color;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import reactor.core.publisher.Mono;
import streamer_website.demo.entity.Task;
import streamer_website.demo.entity.TaskStatus;
import streamer_website.demo.repository.TaskRepository;

import java.time.Instant;
import java.time.ZoneId;
import java.time.format.DateTimeFormatter;
import java.util.Arrays;
import java.util.Comparator;
import java.util.List;
import java.util.stream.Collectors;

@Service
@RequiredArgsConstructor
public class TaskService {

    private final TaskRepository taskRepository;

    // Formattiert einen Task in ein Discord Embed
    public EmbedCreateSpec createEmbedSpec(Task task) {
        Color color = switch (task.getPriority()) {
            case HIGH -> Color.RED;
            case MEDIUM -> Color.YELLOW;
            case LOW -> Color.GREEN;
        };

        String statusEmotes = Arrays.stream(TaskStatus.values())
                .map(status -> status.emote)
                .collect(Collectors.joining(" "));

        return EmbedCreateSpec.builder()
                .color(color)
                .title("TASK #" + task.getId() + ": " + task.getTitle())
                .description(task.getDescription())
                .addField("Status", "**" + task.getStatus().label + " " + task.getStatus().emote + "**", true)
                .addField("Priorität", task.getPriority().label, true)
                .addField("Fällig am", task.getDueDate().format(DateTimeFormatter.ofPattern("dd.MM.yyyy HH:mm")), true)
                .footer("Erstellt von: " + task.getCreatorName() + " | Reagiere zur Statusänderung:\n" + statusEmotes, null)
                .timestamp(task.getDueDate().atZone(ZoneId.of("Europe/Berlin")).toInstant())
                .build();
    }

    // Speichert den Task und fügt Reaktionen hinzu
    public Mono<Void> saveAndPublishTask(GatewayDiscordClient client, Task task, Long targetChannelId) {
        // 1. Task speichern, um ID zu erhalten - Blocking Call in Mono.fromCallable verpackt
        return Mono.fromCallable(() -> taskRepository.save(task))
                .flatMap(savedTask -> {
                    Snowflake channelId = Snowflake.of(targetChannelId);

                    // 2. Embed posten
                    return client.getChannelById(channelId)
                            .cast(discord4j.core.object.entity.channel.MessageChannel.class)
                            .flatMap(channel -> channel.createMessage(createEmbedSpec(savedTask)))
                            .flatMap(message -> {
                                // 3. Message ID im Task speichern
                                savedTask.setMessageId(message.getId().asLong());
                                savedTask.setChannelId(message.getChannelId().asLong());

                                // Erneutes Speichern - Blocking Call in Mono.fromCallable verpackt
                                return Mono.fromCallable(() -> taskRepository.save(savedTask))
                                        .flatMap(finalSavedTask -> {
                                            // 4. Reaktionen hinzufügen - Konvertiere String zu ReactionEmoji
                                            Mono<Void> addReactions = Mono.empty();
                                            for (TaskStatus status : TaskStatus.values()) {
                                                addReactions = addReactions.then(
                                                        message.addReaction(ReactionEmoji.unicode(status.emote))
                                                );
                                            }
                                            return addReactions;
                                        });
                            });
                })
                .then();
    }

    // Aktualisiert den Status eines Tasks basierend auf der Emote-Reaktion
    public Mono<Void> updateStatusByEmote(GatewayDiscordClient client, Long messageId, String emote) {
        // Task synchron aus der DB abrufen und in ein Mono verpacken.
        // HINWEIS: taskRepository.findAll() ist ineffizient. Besser wäre: taskRepository.findByMessageId(messageId).
        return Mono.fromCallable(() -> taskRepository.findAll().stream()
                        .filter(t -> messageId.equals(t.getMessageId()))
                        .findFirst())
                .flatMap(optionalTask -> {
                    if (optionalTask.isEmpty()) {
                        return Mono.empty();
                    }
                    Task task = optionalTask.get();

                    TaskStatus newStatus = Arrays.stream(TaskStatus.values())
                            .filter(s -> s.emote.equals(emote))
                            .findFirst()
                            .orElse(null);

                    if (newStatus != null && newStatus != task.getStatus()) {
                        task.setStatus(newStatus);

                        // Status in DB speichern (Blocking Call in Mono.fromCallable verpackt)
                        return Mono.fromCallable(() -> taskRepository.save(task))
                                .then(
                                        // Aktualisiere das Embed - KORRIGIERT: Nutzt die nicht-deprecierte Builder-Methode .edit().withEmbeds()
                                        client.getMessageById(Snowflake.of(task.getChannelId()), Snowflake.of(task.getMessageId()))
                                                .flatMap(message -> message.edit()
                                                        .withEmbeds(List.of(createEmbedSpec(task)))
                                                )
                                                .then()
                                );
                    }
                    return Mono.empty();
                });
    }

    // Erstellt die wöchentliche Übersicht
    public EmbedCreateSpec createWeeklySummaryEmbed() {
        // Alle nicht abgeschlossenen Tasks holen
        List<Task> openTasks = taskRepository.findByStatusNot(TaskStatus.FINISHED);

        // Nach Priorität sortieren: HIGH > MEDIUM > LOW
        openTasks.sort(Comparator.comparing(Task::getPriority, Comparator.reverseOrder()));

        if (openTasks.isEmpty()) {
            return EmbedCreateSpec.builder()
                    .color(Color.GREEN)
                    .title("🗓️ Wöchentliche Task-Übersicht")
                    .description("Alle Aufgaben sind erledigt! Perfekt!")
                    .build();
        }

        StringBuilder sb = new StringBuilder();
        for (Task task : openTasks) {
            sb.append("**")
                    .append(task.getStatus().emote)
                    .append(" #")
                    .append(task.getId())
                    .append(": ")
                    .append(task.getTitle())
                    .append("** (Prio: ")
                    .append(task.getPriority().label)
                    .append(")")
                    .append("\n> Fällig: ")
                    .append(task.getDueDate().format(DateTimeFormatter.ofPattern("dd.MM."))).append("\n");
        }

        return EmbedCreateSpec.builder()
                .color(Color.BLUE)
                .title("🗓️ Wöchentliche Task-Übersicht")
                .description("Offene und in Bearbeitung befindliche Aufgaben:\n\n" + sb.toString())
                .build();
    }
}
