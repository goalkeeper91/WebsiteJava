package streamer_website.demo.controller.discord;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import streamer_website.demo.entity.discord.DiscordTask;
import streamer_website.demo.entity.discord.TaskStatusEnum;
import streamer_website.demo.service.discord.DiscordTaskService;

import java.util.List;

@RestController
@RequestMapping("/api/tasks")
@RequiredArgsConstructor
@CrossOrigin(origins = "*")
public class DiscordTaskController {
    private final DiscordTaskService service;

    @GetMapping
    public ResponseEntity<List<DiscordTask>> getAllTasks() {
        return ResponseEntity.ok(service.getAllTasks());
    }

    @GetMapping("/status/{status}")
    public ResponseEntity<List<DiscordTask>> getTasksByStatus(@PathVariable TaskStatusEnum status) {
        return ResponseEntity.ok(service.getTasksByStatus(status));
    }

    @GetMapping("/search")
    public ResponseEntity<List<DiscordTask>> searchTasks(@RequestParam String q) {
        return ResponseEntity.ok(service.searchTasks(q));
    }

    @PostMapping
    public ResponseEntity<DiscordTask> createTask(@RequestBody DiscordTask task) {
        DiscordTask created = service.createTask(task);
        return ResponseEntity.ok(created);
    }

    @PutMapping("/{id}")
    public ResponseEntity<DiscordTask> updateTask(@PathVariable Long id, @RequestBody DiscordTask task) {
        return service.updateTask(id, task)
                .map(ResponseEntity::ok)
                .orElse(ResponseEntity.notFound().build());
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Void> deleteTask(@PathVariable Long id) {
        service.deleteTask(id);
        return ResponseEntity.noContent().build();
    }
}
