package streamer_website.demo.handler.discord;

import discord4j.common.util.Snowflake;
import discord4j.core.GatewayDiscordClient;
import discord4j.core.event.domain.message.ReactionAddEvent;
import discord4j.core.object.entity.Message;
import discord4j.core.object.entity.User;
import discord4j.core.object.reaction.ReactionEmoji;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Component;
import reactor.core.publisher.Mono;
import streamer_website.demo.service.TaskService;

@Component
@RequiredArgsConstructor
public class TaskReactionHandler {
    private final TaskService taskService;

    public void register(GatewayDiscordClient client) {
        client.on(ReactionAddEvent.class, this::handle).subscribe();
    }

    private Mono<Void> handle(ReactionAddEvent event) {
        if (event.getMember().map(User::isBot).orElse(false)) {
            return Mono.empty();
        }

        ReactionEmoji emoji = event.getEmoji();
        String emote;

        if (emoji.asUnicodeEmoji().isPresent()) {
            emote = emoji.asUnicodeEmoji().get().getRaw();
        } else if (emoji.asCustomEmoji().isPresent()) {
            ReactionEmoji.Custom custom = emoji.asCustomEmoji().get();
            emote = String.format("<:%s:%s>", custom.getName(), custom.getId().asString());
        } else {
            return Mono.empty();
        }

        final String finalEmote = emote;

        return event.getMessage()
            .flatMap(message -> {

                Long messageId = message.getId().asLong();
                Snowflake userId = event.getUserId();

                Mono<Void> updateMono = taskService.updateStatusByEmote(event.getClient(), messageId, finalEmote);

                Mono<Void> removeReactionMono = message.removeReaction(emoji, userId);

                return updateMono.then(removeReactionMono);
            });
    }
}
