package streamer_website.demo.service.discord;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import streamer_website.demo.entity.discord.DiscordChannels;
import streamer_website.demo.repository.DiscordChannelsRepository;

import java.util.List;
import java.util.Optional;

@Service
@RequiredArgsConstructor
public class DiscordChannelsService {

    private final DiscordChannelsRepository repository;

    public List<DiscordChannels> findAll() {
        return repository.findAll();
    }

    public Optional<DiscordChannels> findById(Long id) {
        return repository.findById(id);
    }

    public List<DiscordChannels> findByGuildId(String guildId) {
        return repository.findByGuildId(guildId);
    }

    @Transactional
    public DiscordChannels save(DiscordChannels channel) {
        return repository.save(channel);
    }

    @Transactional
    public void delete(Long id) {
        repository.deleteById(id);
    }
}
