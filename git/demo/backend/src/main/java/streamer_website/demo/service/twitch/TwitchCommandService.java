package streamer_website.demo.service.twitch;

import org.springframework.stereotype.Service;
import streamer_website.demo.entity.twitch.TwitchAuthToken;
import streamer_website.demo.entity.twitch.TwitchCommand;
import streamer_website.demo.repository.TwitchAuthTokenRepository;
import streamer_website.demo.repository.TwitchCommandRepository;

import java.time.Instant;
import java.util.Map;
import java.util.Optional;
import java.util.concurrent.ConcurrentHashMap;

@Service
public class TwitchCommandService {

    private final TwitchCommandRepository repository;
    private final TwitchAuthTokenRepository tokenRepository;
    private final Map<String, TwitchCommand> commandCache = new ConcurrentHashMap<>();

    public TwitchCommandService(TwitchCommandRepository repository, TwitchAuthTokenRepository tokenRepository) {
        this.repository = repository;
        this.tokenRepository = tokenRepository;
        loadCommandsFromDB();
    }

    private void loadCommandsFromDB() {
        repository.findAll().forEach(cmd -> commandCache.put(cmd.getTrigger().toLowerCase(), cmd));
    }

    public Optional<TwitchCommand> getCommand(String trigger) {
        return Optional.ofNullable(commandCache.get(trigger.toLowerCase()));
    }

    public TwitchCommand addCommand(String trigger, String response, boolean modOnly) {
        TwitchCommand command = new TwitchCommand();
        command.setTrigger(trigger);
        command.setResponse(response);
        command.setModOnly(modOnly);
        command.setCreatedAt(Instant.now());
        command.setUpdatedAt(Instant.now());

        TwitchCommand saved = repository.save(command);
        commandCache.put(trigger.toLowerCase(), saved);
        return saved;
    }

    public Optional<TwitchCommand> editCommand(String trigger, String newResponse) {
        return repository.findByTriggerIgnoreCase(trigger).map(cmd -> {
            cmd.setResponse(newResponse);
            cmd.setUpdatedAt(Instant.now());
            TwitchCommand saved = repository.save(cmd);
            commandCache.put(trigger.toLowerCase(), saved);
            return saved;
        });
    }

    public boolean deleteCommand(String trigger) {
        boolean deleted = repository.findByTriggerIgnoreCase(trigger).map(cmd -> {
            repository.delete(cmd);
            return true;
        }).orElse(false);

        commandCache.remove(trigger.toLowerCase());
        return deleted;
    }

    public String getBotToken() {
        return tokenRepository.findTopByUserNameOrderByCreatedAtDesc("goalkeeper91_bot")
                .map(TwitchAuthToken::getAccessToken)
                .orElseThrow(() -> new IllegalStateException("Kein aktiver Bot-Token gefunden"));
    }

    public String getUserToken() {
        return tokenRepository.findTopByUserNameOrderByCreatedAtDesc("goalkeeper91")
                .map(TwitchAuthToken::getAccessToken)
                .orElseThrow(() -> new IllegalStateException("Kein aktiver User-Token gefunden"));
    }
}
