package streamer_website.demo.commands.discord;

import discord4j.core.GatewayDiscordClient;
import discord4j.core.event.domain.interaction.ChatInputInteractionEvent;
import discord4j.core.object.command.ApplicationCommandInteractionOptionValue;
import discord4j.core.object.command.ApplicationCommandOption;
import discord4j.discordjson.json.ApplicationCommandOptionChoiceData;
import discord4j.discordjson.json.ApplicationCommandOptionData;
import discord4j.discordjson.json.ApplicationCommandRequest;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Component;
import reactor.core.publisher.Mono;
import streamer_website.demo.entity.Task;
import streamer_website.demo.entity.TaskPriority;
import streamer_website.demo.service.TaskService;

import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;

@Component
@RequiredArgsConstructor
public class CreateTaskCommand implements SlashCommand {

    private final TaskService taskService;

    private static final Long TASK_CHANNEL_ID = 1418259428288757760L;

    @Override
    public ApplicationCommandRequest getCommandRequest() {
        return ApplicationCommandRequest.builder()
                .name("task")
                .description("Erstellt eine neue Aufgabe (Task) mit Priorität und Fälligkeitsdatum.")
                .addOption(discord4j.discordjson.json.ApplicationCommandOptionData.builder()
                        .name("titel").description("Titel der Aufgabe").type(ApplicationCommandOption.Type.STRING.getValue()).required(true).build())
                .addOption(discord4j.discordjson.json.ApplicationCommandOptionData.builder()
                        .name("beschreibung").description("Detaillierte Beschreibung der Aufgabe").type(ApplicationCommandOption.Type.STRING.getValue()).required(true).build())
                .addOption(discord4j.discordjson.json.ApplicationCommandOptionData.builder()
                        .name("prioritaet").description("Priorität der Aufgabe (NIEDRIG, MITTEL, HOCH)").type(ApplicationCommandOption.Type.STRING.getValue()).required(true)
                        .addChoice(ApplicationCommandOptionChoiceData.builder().name("HOCH").value(TaskPriority.HIGH.name()).build())
                        .addChoice(ApplicationCommandOptionChoiceData.builder().name("MITTEL").value(TaskPriority.MEDIUM.name()).build())
                        .addChoice(ApplicationCommandOptionChoiceData.builder().name("NIEDRIG").value(TaskPriority.LOW.name()).build())
                        .build())
                .addOption(discord4j.discordjson.json.ApplicationCommandOptionData.builder()
                        .name("faellig").description("Fälligkeitsdatum und Uhrzeit (Format: TT.MM.JJJJ HH:MM)").type(ApplicationCommandOption.Type.STRING.getValue()).required(true).build())
                .build();
    }

    @Override
    public Mono<Void> handle(ChatInputInteractionEvent event, GatewayDiscordClient client) {
        String title = event.getOption("titel").flatMap(option -> option.getValue().map(ApplicationCommandInteractionOptionValue::asString)).orElse("Unbekannt");
        String description = event.getOption("beschreibung").flatMap(option -> option.getValue().map(ApplicationCommandInteractionOptionValue::asString)).orElse("Keine Beschreibung");
        String priorityStr = event.getOption("prioritaet").flatMap(option -> option.getValue().map(ApplicationCommandInteractionOptionValue::asString)).orElse("MEDIUM");
        String dueDateStr = event.getOption("faellig").flatMap(option -> option.getValue().map(ApplicationCommandInteractionOptionValue::asString)).orElse(null);

        TaskPriority priority = TaskPriority.valueOf(priorityStr);
        LocalDateTime dueDate;

        try {
            // Beispiel: 06.10.2025 18:00
            DateTimeFormatter formatter = DateTimeFormatter.ofPattern("dd.MM.yyyy HH:mm");
            assert dueDateStr != null;
            dueDate = LocalDateTime.parse(dueDateStr, formatter);
        } catch (Exception e) {
            return event.reply("Fehler: Ungültiges Fälligkeitsdatum-Format. Bitte verwende TT.MM.JJJJ HH:MM.").withEphemeral(true);
        }

        return event.deferReply()
                .then(Mono.justOrEmpty(event.getInteraction().getMember()))
                .map(member -> {
                    Task task = new Task();
                    task.setTitle(title);
                    task.setDescription(description);
                    task.setPriority(priority);
                    task.setDueDate(dueDate);
                    task.setCreatorDiscordId(member.getId().asLong());
                    task.setCreatorName(member.getDisplayName());
                    return task;
                })
                .flatMap(task -> taskService.saveAndPublishTask(client, task, TASK_CHANNEL_ID)
                        .then(event.editReply("✅ Task **'" + title + "'** erfolgreich erstellt und in den Sammelkanal gepostet!"))
                )
                .onErrorResume(e -> {
                    e.printStackTrace();
                    return event.editReply("❌ Fehler beim Erstellen des Tasks: " + e.getMessage());
                })
                .then();
    }
}
