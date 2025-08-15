package streamer_website.demo.service;

import org.springframework.stereotype.Service;
import streamer_website.demo.entity.TwitchAuthToken;
import streamer_website.demo.entity.TwitchCommand;
import streamer_website.demo.repository.TwitchAuthTokenRepository;
import streamer_website.demo.repository.TwitchCommandRepository;

import java.time.Instant;
import java.util.Optional;

@Service
public class TwitchCommandService {

    private final TwitchCommandRepository repository;
    private final TwitchAuthTokenRepository tokenRepository;

    public TwitchCommandService(TwitchCommandRepository repository, TwitchAuthTokenRepository tokenRepository) {
        this.repository = repository;
        this.tokenRepository = tokenRepository;
    }

    public Optional<TwitchCommand> getCommand(String trigger) {
        return repository.findByTrigger(trigger);
    }

    public TwitchCommand addCommand(String trigger, String response, boolean modOnly) {
        TwitchCommand command = new TwitchCommand();
        command.setTrigger(trigger);
        command.setResponse(response);
        command.setModOnly(modOnly);
        command.setCreatedAt(Instant.now());
        command.setUpdatedAt(Instant.now());

        return repository.save(command);
    }

    public Optional<TwitchCommand> editCommand(String trigger, String newResponse) {
        return repository.findByTrigger(trigger).map(cmd -> {
            cmd.setResponse(newResponse);
            cmd.setUpdatedAt(Instant.now());
            return repository.save(cmd);
        });
    }

    public boolean deleteCommand(String trigger) {
        return repository.findByTrigger(trigger).map(cmd -> {
            repository.delete(cmd);
            return true;
        }).orElse(false);
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
