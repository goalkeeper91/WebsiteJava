package streamer_website.demo.entity;

public enum TaskPriority {
    LOW("Niedrig"),
    MEDIUM("Mittel"),
    HIGH("Hoch");

    public final String label;

    TaskPriority(String label) {
        this.label = label;
    }
}
