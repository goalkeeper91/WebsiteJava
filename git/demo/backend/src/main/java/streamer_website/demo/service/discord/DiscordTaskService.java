package streamer_website.demo.service.discord;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import streamer_website.demo.entity.discord.DiscordTask;
import streamer_website.demo.entity.discord.TaskStatusEnum;
import streamer_website.demo.repository.DiscordTaskRepository;

import java.util.List;
import java.util.Optional;

@Service
@RequiredArgsConstructor
public class DiscordTaskService {

    private final DiscordTaskRepository repository;

    public List<DiscordTask> getAllTasks() {
        return repository.findAll();
    }

    public List<DiscordTask> getTasksByStatus(TaskStatusEnum status) {
        return repository.findByStatus(status);
    }

    public List<DiscordTask> searchTasks(String keyword) {
        return repository.findByTitleContainingIgnoreCase(keyword);
    }

    public DiscordTask createTask(DiscordTask task) {
        return repository.save(task);
    }

    public Optional<DiscordTask> updateTask(Long id, DiscordTask updated) {
        return repository.findById(id).map(task -> {
            task.setTitle(updated.getTitle());
            task.setDescription(updated.getDescription());
            task.setStatus(updated.getStatus());
            task.setDueDate(updated.getDueDate());
            return repository.save(task);
        });
    }

    public void deleteTask(Long id) {
        repository.deleteById(id);
    }
}

