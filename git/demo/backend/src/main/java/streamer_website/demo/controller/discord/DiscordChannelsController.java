package streamer_website.demo.controller.discord;

import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import streamer_website.demo.entity.discord.DiscordChannels;
import streamer_website.demo.service.discord.DiscordChannelsService;

import java.util.List;

@RestController
@RequestMapping("/api/discord/channels")
@RequiredArgsConstructor
public class DiscordChannelsController {

    private final DiscordChannelsService service;

    @GetMapping
    public ResponseEntity<List<DiscordChannels>> getAllChannels() {
        return ResponseEntity.ok(service.findAll());
    }

    @PostMapping
    public ResponseEntity<DiscordChannels> createChannel(@RequestBody DiscordChannels channel) {
        DiscordChannels saved = service.save(channel);
        return ResponseEntity.ok(saved);
    }

    @PutMapping("/{id}")
    public ResponseEntity<DiscordChannels> updateChannel(@PathVariable Long id, @RequestBody DiscordChannels channel) {
        return service.findById(id).map(existing -> {
            channel.setId(id);
            DiscordChannels updated = service.save(channel);
            return ResponseEntity.ok(updated);
        }).orElse(ResponseEntity.notFound().build());
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Void> deleteChannel(@PathVariable Long id) {
        service.delete(id);
        return ResponseEntity.noContent().build();
    }
}
