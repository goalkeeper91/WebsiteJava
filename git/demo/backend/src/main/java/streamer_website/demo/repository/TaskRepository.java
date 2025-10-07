package streamer_website.demo.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import streamer_website.demo.entity.Task;
import streamer_website.demo.entity.TaskStatus;

import java.util.List;

public interface TaskRepository extends JpaRepository<Task, Long> {
    // Findet alle Tasks, die noch nicht abgeschlossen sind (für Scheduler)
    List<Task> findByStatusNot(TaskStatus status);
}