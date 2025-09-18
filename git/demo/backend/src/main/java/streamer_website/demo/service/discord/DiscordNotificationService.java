package streamer_website.demo.service.discord;

import discord4j.common.util.Snowflake;
import discord4j.core.GatewayDiscordClient;
import discord4j.core.object.entity.Message;
import discord4j.core.object.entity.channel.MessageChannel;
import discord4j.core.object.entity.channel.TextChannel;
import discord4j.core.object.reaction.ReactionEmoji;
import discord4j.core.spec.EmbedCreateSpec;
import discord4j.core.spec.MessageCreateSpec;
import discord4j.rest.util.Color;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import reactor.core.publisher.Mono;
import streamer_website.demo.entity.ContactRequest;
import streamer_website.demo.entity.discord.DiscordChannels;
import streamer_website.demo.entity.discord.DiscordTask;
import streamer_website.demo.entity.discord.TaskStatusEnum;
import streamer_website.demo.repository.DiscordChannelsRepository;

import java.time.DayOfWeek;
import java.time.LocalDate;
import java.time.ZoneId;
import java.time.format.DateTimeFormatter;
import java.time.format.TextStyle;
import java.time.temporal.TemporalAdjusters;
import java.util.List;
import java.util.Locale;
import java.util.Map;
import java.util.stream.Collectors;

@Service
@RequiredArgsConstructor
public class DiscordNotificationService {

    private final GatewayDiscordClient client;
    private final DiscordChannelsRepository channelsRepository;
    private final DiscordReactionListener reactionListener;

    /* ---------------------- Kontaktanfragen ---------------------- */

    public Mono<Void> notifyNewContactRequest(ContactRequest request) {
        EmbedCreateSpec embed = buildContactEmbed(request);
        return sendToChannelByKeyword("kontakt", embed);
    }

    private EmbedCreateSpec buildContactEmbed(ContactRequest request) {
        EmbedCreateSpec.Builder builder = EmbedCreateSpec.builder()
                .title("📩 Neue Kontaktanfrage")
                .color(Color.BLUE)
                .addField("👤 Name", request.getName(), false)
                .addField("📧 Email", request.getEmail(), false);

        if (request.getPhone() != null && !request.getPhone().isBlank()) {
            builder.addField("📞 Telefon", request.getPhone(), false);
        }

        builder.addField("📝 Betreff", request.getSubject(), false)
                .addField("💬 Nachricht", request.getMessage(), false)
                .footer("Eingegangen am", null)
                .timestamp(request.getCreatedAt().atZone(ZoneId.systemDefault()).toInstant());

        return builder.build();
    }

    /* ---------------------- Tasks ---------------------- */

    public Mono<Void> notifyNewTask(DiscordTask task) {
        EmbedCreateSpec embed = buildTaskEmbed(task, false);

        List<DiscordChannels> channels = channelsRepository.findByDescriptionContainingIgnoreCase("task");
        if (channels.isEmpty()) {
            System.err.println("⚠️ Kein Discord-Channel gefunden für Tasks");
            return Mono.empty();
        }

        return Mono.when(
                channels.stream()
                        .map(ch -> client.getChannelById(Snowflake.of(ch.getChannelId()))
                                .ofType(TextChannel.class)
                                .flatMap(channel -> channel.createMessage(
                                        MessageCreateSpec.builder().addEmbed(embed).build()
                                ).flatMap(this::addReactionsToMessage))
                        ).toList()
        ).then();
    }

    public Mono<Void> notifyTaskCompleted(DiscordTask task) {
        EmbedCreateSpec embed = buildTaskEmbed(task, true);
        return sendToChannelByKeyword("task", embed);
    }

    private EmbedCreateSpec buildTaskEmbed(DiscordTask task, boolean completed) {
        EmbedCreateSpec.Builder builder = EmbedCreateSpec.builder()
                .title((completed ? "✅ Task abgeschlossen" : "🆕 Neuer Task") + " (ID: " + task.getId() + ")")
                .color(completed ? Color.GREEN : Color.ORANGE)
                .addField("Titel", task.getTitle(), false);

        if (task.getDescription() != null && !task.getDescription().isBlank()) {
            builder.addField("Beschreibung", task.getDescription(), false);
        }

        if (task.getDueDate() != null) {
            builder.addField("Fällig am", task.getDueDate().toString(), false);
        }

        builder.addField("Status", task.getStatus().name(), false);
        return builder.build();
    }

    private Mono<Void> addReactionsToMessage(Message message) {
        return message.addReaction(ReactionEmoji.unicode("\u23F3")) // ⏳
                .then(message.addReaction(ReactionEmoji.unicode("\u2699\uFE0F"))) // ⚙️
                .then(message.addReaction(ReactionEmoji.unicode("\u2705"))); // ✅
    }

    /* ---------------------- Wochenübersicht ---------------------- */

    public Mono<Void> notifyWeeklyTasks(List<DiscordTask> tasks) {
        LocalDate today = LocalDate.now();
        LocalDate nextMonday = today.with(TemporalAdjusters.nextOrSame(DayOfWeek.MONDAY));
        LocalDate nextSunday = nextMonday.with(TemporalAdjusters.next(DayOfWeek.SUNDAY));

        DateTimeFormatter dateFmt = DateTimeFormatter.ofPattern("dd.MM.yyyy");
        DateTimeFormatter timeFmt = DateTimeFormatter.ofPattern("HH:mm");

        Map<DayOfWeek, List<DiscordTask>> tasksByDay = tasks.stream()
                .filter(t -> t.getDueDate() != null)
                .filter(t -> {
                    LocalDate due = t.getDueDate().toLocalDate();
                    return !due.isBefore(nextMonday) && !due.isAfter(nextSunday);
                })
                .collect(Collectors.groupingBy(t -> t.getDueDate().getDayOfWeek()));

        EmbedCreateSpec.Builder builder = EmbedCreateSpec.builder()
                .title("📆 Wochenübersicht: " + nextMonday.format(dateFmt) + " - " + nextSunday.format(dateFmt))
                .color(Color.of(0x9B59B6));

        for (DayOfWeek day : DayOfWeek.values()) {
            if (day.getValue() >= DayOfWeek.MONDAY.getValue() && day.getValue() <= DayOfWeek.FRIDAY.getValue()) {
                List<DiscordTask> dayTasks = tasksByDay.getOrDefault(day, List.of());
                if (dayTasks.isEmpty()) {
                    builder.addField(day.getDisplayName(TextStyle.FULL, Locale.GERMAN), "✅ Keine Tasks", false);
                } else {
                    StringBuilder sb = new StringBuilder();
                    for (DiscordTask t : dayTasks) {
                        sb.append("• ")
                                .append(t.getTitle())
                                .append(" – ")
                                .append(t.getDueDate().format(timeFmt))
                                .append(" (").append(mapStatus(t.getStatus())).append(")\n");
                    }
                    builder.addField(day.getDisplayName(TextStyle.FULL, Locale.GERMAN), sb.toString(), false);
                }
            }
        }

        EmbedCreateSpec embed = builder.build();
        return sendToChannelByKeyword("weekly", embed);
    }

    /* ---------------------- Generic ---------------------- */

    public void registerReactionListener() {
        client.on(discord4j.core.event.domain.message.ReactionAddEvent.class,
                reactionListener::handleReaction).subscribe();
    }

    private Mono<Void> sendToChannelByKeyword(String keyword, EmbedCreateSpec embed) {
        List<DiscordChannels> channels = channelsRepository.findByDescriptionContainingIgnoreCase(keyword);

        if (channels.isEmpty()) {
            System.err.println("⚠️ Kein Discord-Channel gefunden für Keyword: " + keyword);
            return Mono.empty();
        }

        return Mono.when(
                channels.stream()
                        .map(ch -> client.getChannelById(Snowflake.of(ch.getChannelId()))
                                .ofType(MessageChannel.class)
                                .flatMap(channel -> channel.createMessage(
                                        MessageCreateSpec.builder().addEmbed(embed).build()
                                ))
                        ).toList()
        ).then();
    }

    private String mapStatus(TaskStatusEnum status) {
        return switch (status) {
            case NOT_STARTED -> "⏳ Noch nicht begonnen";
            case WORK_IN_PROGRESS -> "⚙️ In Arbeit";
            case COMPLETED -> "✅ Erledigt";
        };
    }
}
