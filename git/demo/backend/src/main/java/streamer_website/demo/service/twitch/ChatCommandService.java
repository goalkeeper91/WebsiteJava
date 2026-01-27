package streamer_website.demo.service.twitch;

import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import streamer_website.demo.entity.twitch.ChatCommand;
import streamer_website.demo.entity.twitch.TwitchChannel;
import streamer_website.demo.repository.TwitchChannelRepository;
import streamer_website.demo.repository.TwitchChatCommandRepository;
import streamer_website.demo.service.BotSignalService;

import java.util.Optional;

@Service
@RequiredArgsConstructor
@Transactional
public class ChatCommandService implements ChatCommandServiceInterface {

    private final TwitchChatCommandRepository repository;
    private final TwitchChannelRepository channelRepository;
    private final BotSignalService botSignalService;

    @Override
    @Transactional(readOnly = true)
    public Page<ChatCommand> getCommands(String twitchUserId, Pageable pageable) {
        String internalChannelId = getInternalChannelId(twitchUserId);
        return repository.findByChannelId(internalChannelId, pageable);
    }

    @Override
    @Transactional(readOnly = true)
    public Page<ChatCommand> searchCommands(String twitchUserId, String search, Pageable pageable) {
        String internalChannelId = getInternalChannelId(twitchUserId);
        return repository.findByChannelIdAndTriggerContainingIgnoreCase(
                internalChannelId,
                search,
                pageable
        );
    }

    @Override
    @Transactional(readOnly = true)
    public Page<ChatCommand> getCommandsByStatus(String twitchUserId, boolean enabled, Pageable pageable) {
        String internalChannelId = getInternalChannelId(twitchUserId);
        return repository.findByChannelIdAndEnabled(internalChannelId, enabled, pageable);
    }

    @Override
    @Transactional(readOnly = true)
    public Optional<ChatCommand> getCommand(String twitchUserId, String trigger) {
        String internalChannelId = getInternalChannelId(twitchUserId);
        return repository.findByChannelIdAndTriggerIgnoreCase(internalChannelId, trigger);
    }

    @Transactional(readOnly = true)
    @Override
    public Optional<ChatCommand> getCommandById(String twitchUserId, Long id) {
        String internalChannelId = getInternalChannelId(twitchUserId);
        return repository.findByIdAndChannelId(id, internalChannelId);
    }

    @Override
    public ChatCommand createCommand(String twitchUserId, String trigger, String response, Integer cooldown) {
        validateTrigger(trigger);

        TwitchChannel channel = channelRepository.findByTwitchUserId(twitchUserId)
                .orElseThrow(() -> new IllegalArgumentException("Channel nicht registriert. Bitte erneut einloggen."));

        String internalChannelId = channel.getId().toString();

        if (repository.existsByChannelIdAndTriggerIgnoreCase(internalChannelId, trigger)) {
            throw new IllegalStateException("Command existiert bereits");
        }

        ChatCommand command = new ChatCommand();
        command.setChannelId(internalChannelId);
        command.setTrigger(normalizeTrigger(trigger));
        command.setResponse(response);
        command.setCooldown(cooldown != null ? cooldown : 0);
        command.setEnabled(true);

        ChatCommand saved = repository.save(command);

        botSignalService.sendCommandUpdateSignal(twitchUserId);

        return saved;
    }

    public ChatCommand updateCommand(
            String twitchUserId,
            Long id,
            String newTrigger,
            String newResponse,
            Integer newCooldown,
            Boolean enabled
    ) {
        String internalChannelId = getInternalChannelId(twitchUserId);

        ChatCommand command = repository
                .findByIdAndChannelId(id, internalChannelId)
                .orElseThrow(() -> new IllegalArgumentException("Command nicht gefunden"));

        // ✅ Trigger kann jetzt geändert werden!
        if (newTrigger != null && !newTrigger.isBlank()) {
            String normalized = normalizeTrigger(newTrigger);

            // Prüfe ob neuer Trigger bereits existiert (außer bei gleichem Command)
            repository.findByChannelIdAndTriggerIgnoreCase(internalChannelId, normalized)
                    .ifPresent(existing -> {
                        if (!existing.getId().equals(id)) {
                            throw new IllegalStateException("Trigger bereits vergeben");
                        }
                    });

            command.setTrigger(normalized);
        }

        if (newResponse != null) {
            command.setResponse(newResponse);
        }

        if (newCooldown != null) {
            command.setCooldown(newCooldown);
        }

        if (enabled != null) {
            command.setEnabled(enabled);
        }

        ChatCommand saved = repository.save(command);

        botSignalService.sendCommandUpdateSignal(twitchUserId);

        return saved;
    }

    public void deleteCommand(String twitchUserId, Long id) {
        String internalChannelId = getInternalChannelId(twitchUserId);

        ChatCommand command = repository
                .findByIdAndChannelId(id, internalChannelId)
                .orElseThrow(() -> new IllegalArgumentException("Command nicht gefunden"));

        repository.delete(command);

        botSignalService.sendCommandUpdateSignal(twitchUserId);
    }

    public ChatCommand toggleCommand(String twitchUserId, Long id, boolean enabled) {
        String internalChannelId = getInternalChannelId(twitchUserId);

        ChatCommand command = repository
                .findByIdAndChannelId(id, internalChannelId)
                .orElseThrow(() -> new IllegalArgumentException("Command nicht gefunden"));

        command.setEnabled(enabled);

        ChatCommand saved = repository.save(command);

        botSignalService.sendCommandUpdateSignal(twitchUserId);

        return saved;
    }

    private String getInternalChannelId(String twitchUserId) {
        TwitchChannel channel = channelRepository.findByTwitchUserId(twitchUserId)
                .orElseThrow(() -> new IllegalArgumentException(
                        "Channel nicht registriert. Bitte erneut einloggen."
                ));
        return channel.getId().toString();
    }

    private void validateTrigger(String trigger) {
        if (trigger == null || trigger.isBlank()) {
            throw new IllegalArgumentException("Trigger darf nicht leer sein");
        }

        if (trigger.length() > 100) {
            throw new IllegalArgumentException("Trigger zu lang");
        }
    }

    private String normalizeTrigger(String trigger) {
        return trigger.trim().toLowerCase();
    }
}