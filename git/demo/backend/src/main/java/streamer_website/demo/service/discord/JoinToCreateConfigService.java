package streamer_website.demo.service.discord;

import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;
import streamer_website.demo.entity.discord.JoinToCreateConfig;
import streamer_website.demo.repository.JoinToCreateRepository;

import java.util.Collections;
import java.util.List;

@Service
@RequiredArgsConstructor
public class JoinToCreateConfigService {

    private final JoinToCreateRepository repository;
    private static final Logger logger = LoggerFactory.getLogger(JoinToCreateConfigService.class);

    // Alle Configs abrufen, safe für leere DB
    public List<JoinToCreateConfig> getAllConfigs() {
        try {
            return repository.findAll();
        } catch (Exception e) {
            logger.error(e.getMessage());
        }
        return Collections.emptyList();
    }

    public JoinToCreateConfig getConfigById(Long id) {
        return repository.findById(id)
                .orElseThrow(() -> new RuntimeException("Config nicht gefunden: " + id));
    }

    // Create: Validierung nur bei explizitem API-Aufruf
    public JoinToCreateConfig createConfig(JoinToCreateConfig config) {
        validateConfig(config);
        return repository.save(config);
    }

    // Update: Validierung nur bei API-Aufruf
    public JoinToCreateConfig updateConfig(Long id, JoinToCreateConfig updatedConfig) {
        JoinToCreateConfig existing = repository.findById(id)
                .orElseThrow(() -> new RuntimeException("Config nicht gefunden: " + id));

        validateConfig(updatedConfig);

        existing.setJoinChannelId(updatedConfig.getJoinChannelId());
        existing.setCategoryId(updatedConfig.getCategoryId());
        existing.setChannelNamePrefix(updatedConfig.getChannelNamePrefix());
        existing.setUserLimit(updatedConfig.getUserLimit());
        existing.setPrivateChannel(updatedConfig.getPrivateChannel());

        return repository.save(existing);
    }

    public void deleteConfig(Long id) {
        repository.deleteById(id);
    }

    private void validateConfig(JoinToCreateConfig config) {
        if (config.getJoinChannelId() == null) {
            logger.error("Validation failed: Join Channel ID is empty");
            throw new RuntimeException("Join Channel ID darf nicht leer sein");
        }
        if (config.getCategoryId() == null) {
            logger.error("Validation failed: Category ID is empty");
            throw new RuntimeException("Category ID darf nicht leer sein");
        }
        if (config.getChannelNamePrefix() == null || config.getChannelNamePrefix().isBlank()) {
            logger.error("Validation failed: Channel Name Prefix is empty");
            throw new RuntimeException("Channel Name Prefix darf nicht leer sein");
        }
        if (config.getUserLimit() == null || config.getUserLimit() <= 0) {
            logger.error("Validation failed: User Limit invalid -> {}", config.getUserLimit());
            throw new RuntimeException("User Limit darf nicht leer oder <= 0 sein");
        }
        if (config.getPrivateChannel() == null) {
            logger.error("Validation failed: Private Channel is null");
            throw new RuntimeException("Private Channel darf nicht leer sein");
        }
    }
}
