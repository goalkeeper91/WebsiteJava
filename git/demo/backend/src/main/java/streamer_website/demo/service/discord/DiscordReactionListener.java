package streamer_website.demo.service.discord;

import discord4j.core.event.domain.message.ReactionAddEvent;
import discord4j.core.object.entity.Message;
import discord4j.core.object.entity.User;
import discord4j.core.object.reaction.ReactionEmoji;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Component;
import reactor.core.publisher.Mono;
import streamer_website.demo.entity.discord.DiscordTask;
import streamer_website.demo.entity.discord.TaskStatusEnum;
import streamer_website.demo.repository.DiscordTaskRepository;

@Component
@RequiredArgsConstructor
public class DiscordReactionListener {

    private final DiscordTaskRepository taskRepository;

    public Mono<Void> handleReaction(ReactionAddEvent event) {
        Message message = event.getMessage().block();
        User user = event.getUser().block();

        // Nur echte Benutzer reagieren lassen
        if (message == null || user == null || user.isBot()) return Mono.empty();

        ReactionEmoji emojiObj = event.getEmoji();
        if (!(emojiObj instanceof ReactionEmoji.Unicode unicodeEmoji)) return Mono.empty();

        String emoji = unicodeEmoji.getRaw();
        Long taskId = parseTaskIdFromMessage(message);
        if (taskId == null) return Mono.empty();

        DiscordTask task = taskRepository.findById(taskId).orElse(null);
        if (task == null) return Mono.empty();

        switch (emoji) {
            case "⏳" -> task.setStatus(TaskStatusEnum.NOT_STARTED);
            case "⚙️" -> task.setStatus(TaskStatusEnum.WORK_IN_PROGRESS);
            case "✅" -> task.setStatus(TaskStatusEnum.COMPLETED);
            default -> {}
        }

        taskRepository.save(task);
        return Mono.empty();
    }

    private Long parseTaskIdFromMessage(Message message) {
        String content = message.getEmbeds().getFirst().getTitle().orElse("");
        try {
            if (content.contains("ID:")) {
                String idStr = content.substring(content.indexOf("ID:") + 3).trim();
                return Long.parseLong(idStr);
            }
        } catch (NumberFormatException e) {
            return null;
        }
        return null;
    }
}
