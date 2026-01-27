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

@RestController
@RequiredArgsConstructor
@RequestMapping("/api/dashboard/commands")
public class ChatCommandController {

    private final ChatCommandService chatCommandService;

    @GetMapping
    public Page<ChatCommandDto> getCommands(
            HttpSession session,
            @RequestParam(required = false) String search,
            @RequestParam(required = false) Boolean enabled,
            Pageable pageable
    ) {
        TwitchUser user = requireUser(session);
        Page<ChatCommand> page;

        if (search != null && !search.isBlank()) {
            page = chatCommandService.searchCommands(user.id(), search, pageable);
        } else if (enabled != null) {
            page = chatCommandService.getCommandsByStatus(user.id(), enabled, pageable);
        } else {
            page = chatCommandService.getCommands(user.id(), pageable);
        }

        return page.map(ChatCommandMapper::toDto);
    }

    @GetMapping("/{id}")
    public ChatCommandDto getCommand(
            HttpSession session,
            @PathVariable Long id
    ) {
        TwitchUser user = requireUser(session);

        ChatCommand command = chatCommandService
                .getCommandById(user.id(), id)
                .orElseThrow(() -> new ResponseStatusException(
                        HttpStatus.NOT_FOUND,
                        "Command nicht gefunden"
                ));

        return ChatCommandMapper.toDto(command);
    }

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

    @PutMapping("/{id}")
    public ChatCommandDto updateCommand(
            HttpSession session,
            @PathVariable Long id,
            @RequestBody @Valid UpdateChatCommandRequest request
    ) {
        TwitchUser user = requireUser(session);

        ChatCommand command = chatCommandService.updateCommand(
                user.id(),
                id,
                request.getTrigger(),
                request.getResponse(),
                request.getCooldown(),
                request.getEnabled()
        );

        return ChatCommandMapper.toDto(command);
    }

    @DeleteMapping("/{id}")
    public void deleteCommand(
            HttpSession session,
            @PathVariable Long id
    ) {
        TwitchUser user = requireUser(session);
        chatCommandService.deleteCommand(user.id(), id);
    }

    @PatchMapping("/{id}/toggle")
    public ChatCommandDto toggleCommand(
            HttpSession session,
            @PathVariable Long id,
            @RequestParam boolean enabled
    ) {
        TwitchUser user = requireUser(session);

        ChatCommand command = chatCommandService.toggleCommand(
                user.id(),
                id,
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