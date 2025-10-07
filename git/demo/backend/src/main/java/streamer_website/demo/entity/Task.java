package streamer_website.demo.entity;

import jakarta.persistence.*;
import lombok.Data;

import java.time.LocalDateTime;

@Entity
@Table(name="tasks")
@Data
public class Task {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    // Discord Message ID, um das Embed später zu aktualisieren oder zu löschen
    @Column(name = "message_id")
    private Long messageId;

    // Discord Channel ID, in dem das Embed gepostet wurde
    @Column(name = "channel_id")
    private Long channelId;

    private String title;
    private String description;

    // Ersteller-Informationen
    @Column(name = "creator_discord_id")
    private Long creatorDiscordId;
    @Column(name = "creator_name")
    private String creatorName;

    // Zeitstempel für Fälligkeit (Tag mit Uhrzeit)
    @Column(name = "due_date")
    private LocalDateTime dueDate;

    @Enumerated(EnumType.STRING)
    private TaskStatus status = TaskStatus.OPEN; // Standardstatus

    @Enumerated(EnumType.STRING)
    private TaskPriority priority = TaskPriority.MEDIUM; // Standardpriorität
}
