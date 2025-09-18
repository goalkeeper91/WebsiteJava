package streamer_website.demo.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import streamer_website.demo.entity.discord.DiscordTask;
import streamer_website.demo.entity.discord.TaskStatusEnum;

import java.util.List;

public interface DiscordTaskRepository extends JpaRepository<DiscordTask, Long> {
    List<DiscordTask> findByStatus(TaskStatusEnum status);

    List<DiscordTask> findByTitleContainingIgnoreCase(String keyword);
}
