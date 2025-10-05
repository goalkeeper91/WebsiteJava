package streamer_website.demo.controller.discord;

import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import streamer_website.demo.dto.JoinToCreateConfigRequest;
import streamer_website.demo.dto.JoinToCreateConfigResponse;
import streamer_website.demo.entity.discord.JoinToCreateConfig;
import streamer_website.demo.service.discord.JoinToCreateConfigService;

import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

@RestController
@RequestMapping("/api/discord/join-to-create")
@RequiredArgsConstructor
public class JoinToCreateChannelController {

    private final JoinToCreateConfigService service;

    @GetMapping
    public ResponseEntity<List<JoinToCreateConfigResponse>> getAllConfigs() {
        List<JoinToCreateConfigResponse> response = service.getAllConfigs().stream()
                .map(this::toResponse)
                .collect(Collectors.toList());
        return ResponseEntity.ok(response);
    }

    @GetMapping("/{id}")
    public ResponseEntity<JoinToCreateConfigResponse> getConfigById(@PathVariable Long id) {
        return ResponseEntity.ok(toResponse(service.getConfigById(id)));
    }

    @PostMapping
    public ResponseEntity<?> createConfig(@Valid @RequestBody JoinToCreateConfigRequest request) {
        try {
            JoinToCreateConfig config = toEntity(request);
            JoinToCreateConfig saved = service.createConfig(config);
            return ResponseEntity.ok(toResponse(saved));
        } catch (Exception e) {
            e.printStackTrace(); // Container logs
            return ResponseEntity.status(500).body(Map.of("error", e.getMessage()));
        }
    }

    @PutMapping("/{id}")
    public ResponseEntity<JoinToCreateConfigResponse> updateConfig(
            @PathVariable Long id,
            @Valid @RequestBody JoinToCreateConfigRequest request) {
        JoinToCreateConfig config = toEntity(request);
        JoinToCreateConfig updated = service.updateConfig(id, config);
        return ResponseEntity.ok(toResponse(updated));
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Void> deleteConfig(@PathVariable Long id) {
        service.deleteConfig(id);
        return ResponseEntity.noContent().build();
    }

    // --- Hilfsmethoden ---
    private JoinToCreateConfigResponse toResponse(JoinToCreateConfig config) {
        JoinToCreateConfigResponse response = new JoinToCreateConfigResponse();
        response.setId(config.getId());
        response.setJoinChannelId(config.getJoinChannelId());
        response.setCategoryId(config.getCategoryId());
        response.setChannelNamePrefix(config.getChannelNamePrefix());
        response.setUserLimit(config.getUserLimit());
        response.setPrivateChannel(config.getPrivateChannel());
        return response;
    }

    private JoinToCreateConfig toEntity(JoinToCreateConfigRequest request) {
        JoinToCreateConfig config = new JoinToCreateConfig();
        config.setJoinChannelId(request.getJoinChannelId());
        config.setCategoryId(request.getCategoryId());
        config.setChannelNamePrefix(request.getChannelNamePrefix());
        config.setUserLimit(request.getUserLimit());
        config.setPrivateChannel(request.getPrivateChannel());
        return config;
    }
}
