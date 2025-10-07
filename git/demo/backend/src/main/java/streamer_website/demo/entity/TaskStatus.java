package streamer_website.demo.entity;

// Enum für den Bearbeitungsstatus
public enum TaskStatus {
    OPEN("Offen", "🔴"),
    IN_PROGRESS("In Bearbeitung", "🟡"),
    ON_HOLD("Auf Eis", "🧊"),
    ACCEPTANCE("Abnahme", "🟣"),
    FINISHED("Fertiggestellt", "✅");

    public final String label;
    public final String emote;

    TaskStatus(String label, String emote) {
        this.label = label;
        this.emote = emote;
    }
}
