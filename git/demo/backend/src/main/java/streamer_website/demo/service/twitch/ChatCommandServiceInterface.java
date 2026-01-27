package streamer_website.demo.service.twitch;

import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import streamer_website.demo.entity.twitch.ChatCommand;

import java.util.Optional;

public interface ChatCommandServiceInterface {

    /* =========================
     * Lesen (Dashboard)
     * ========================= */

    Page<ChatCommand> getCommands(
            String twitchUserId,
            Pageable pageable
    );

    Page<ChatCommand> searchCommands(
            String twitchUserId,
            String search,
            Pageable pageable
    );

    Page<ChatCommand> getCommandsByStatus(
            String twitchUserId,
            boolean enabled,
            Pageable pageable
    );

    Optional<ChatCommand> getCommand(
            String twitchUserId,
            String trigger
    );

    /* =========================
     * Schreiben (Dashboard)
     * ========================= */

    ChatCommand createCommand(
            String twitchUserId,
            String trigger,
            String response,
            Integer cooldown
    );

    ChatCommand updateCommand(
            String twitchUserId,
            String trigger,
            String newResponse,
            Integer newCooldown,
            Boolean enabled
    );

    void deleteCommand(
            String twitchUserId,
            String trigger
    );

    ChatCommand toggleCommand(
            String twitchUserId,
            String trigger,
            boolean enabled
    );
}

