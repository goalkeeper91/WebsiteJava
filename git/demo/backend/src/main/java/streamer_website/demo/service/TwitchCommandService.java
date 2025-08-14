package streamer_website.demo.service;

import org.springframework.stereotype.Service;
import streamer_website.demo.entity.TwitchCommand;
import streamer_website.demo.repository.TwitchCommandRepository;

import java.time.Instant;
import java.util.Optional;

@Service
public class TwitchCommandService {

    private final TwitchCommandRepository repository;

    public TwitchCommandService(TwitchCommandRepository repository) {
        this.repository = repository;
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
}
