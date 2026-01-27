package streamer_website.demo.controller.twitch;

import jakarta.servlet.http.HttpSession;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.server.ResponseStatusException;
import streamer_website.demo.dto.twitch.ChatCommandDto;
import streamer_website.demo.dto.twitch.CreateChatCommandRequest;
import streamer_website.demo.dto.twitch.TwitchUser;
import streamer_website.demo.dto.twitch.UpdateChatCommandRequest;
import streamer_website.demo.entity.twitch.ChatCommand;
import streamer_website.demo.mapper.twitch.ChatCommandMapper;
import streamer_website.demo.service.twitch.ChatCommandService;

import java.util.Map;

@RestController
@RequiredArgsConstructor
@RequestMapping("/api/dashboard/commands")
public class ChatCommandController {

    private final ChatCommandService chatCommandService;

    /* =========================
     * GET – Listen / Suche
     * ========================= */

    @GetMapping
    public Page<ChatCommandDto> getCommands(
            HttpSession session,
            @RequestParam(required = false) String search,
            @RequestParam(required = false) Boolean enabled,
            Pageable pageable
    ) {
        TwitchUser user = requireUser(session);
        String channelId = user.id();
        Page<ChatCommand> page;

        if (search != null && !search.isBlank()) {
            page = chatCommandService.searchCommands(channelId, search, pageable);
        } else if (enabled != null) {
            page = chatCommandService.getCommandsByStatus(channelId, enabled, pageable);
        } else {
            page = chatCommandService.getCommands(channelId, pageable);
        }

        return page.map(ChatCommandMapper::toDto);
    }

    /* =========================
     * GET – Einzelner Command
     * ========================= */

    @GetMapping("/{trigger}")
    public ChatCommandDto getCommand(
            HttpSession session,
            @PathVariable String trigger
    ) {
        TwitchUser user = requireUser(session);

        ChatCommand command = chatCommandService
                .getCommand(user.id(), trigger)
                .orElseThrow(() -> new IllegalArgumentException("Command nicht gefunden"));

        return ChatCommandMapper.toDto(command);
    }

    /* =========================
     * POST – Anlegen
     * ========================= */

    @PostMapping
    public ChatCommandDto createCommand(
            HttpSession session,
            @RequestBody @Valid CreateChatCommandRequest request
    ) {
        TwitchUser user = requireUser(session);

        ChatCommand command = chatCommandService.createCommand(
                user.id(),
                request.getTrigger(),
                request.getResponse(),
                request.getCooldown()
        );

        return ChatCommandMapper.toDto(command);
    }

    /* =========================
     * PUT – Aktualisieren
     * ========================= */

    @PutMapping("/{trigger}")
    public ChatCommandDto updateCommand(
            HttpSession session,
            @PathVariable String trigger,
            @RequestBody UpdateChatCommandRequest request
    ) {
        TwitchUser user = requireUser(session);

        ChatCommand command = chatCommandService.updateCommand(
                user.id(),
                trigger,
                request.getResponse(),
                request.getCooldown(),
                request.getEnabled()
        );

        return ChatCommandMapper.toDto(command);
    }

    /* =========================
     * DELETE – Löschen
     * ========================= */

    @DeleteMapping("/{trigger}")
    public void deleteCommand(
            HttpSession session,
            @PathVariable String trigger
    ) {
        TwitchUser user = requireUser(session);

        chatCommandService.deleteCommand(user.id(), trigger);
    }

    /* =========================
     * PATCH – Enable / Disable
     * ========================= */

    @PatchMapping("/{trigger}/toggle")
    public ChatCommandDto toggleCommand(
            HttpSession session,
            @PathVariable String trigger,
            @RequestParam boolean enabled
    ) {
        TwitchUser user = requireUser(session);

        ChatCommand command = chatCommandService.toggleCommand(
                user.id(),
                trigger,
                enabled
        );

        return ChatCommandMapper.toDto(command);
    }

    private TwitchUser requireUser(HttpSession session) {
        TwitchUser user = (TwitchUser) session.getAttribute("user");
        if (user == null) {
            throw new ResponseStatusException(
                    HttpStatus.UNAUTHORIZED,
                    "Not logged in"
            );
        }
        return user;
    }
}



