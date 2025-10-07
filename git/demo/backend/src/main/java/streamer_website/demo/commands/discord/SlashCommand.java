package streamer_website.demo.commands.discord;

import discord4j.core.GatewayDiscordClient;
import discord4j.core.event.domain.interaction.ChatInputInteractionEvent;
import discord4j.discordjson.json.ApplicationCommandRequest;
import reactor.core.publisher.Mono;

public interface SlashCommand {
    /**
     * Definiert die JSON-Struktur des Slash Commands (Name, Beschreibung, Optionen)
     * für die Registrierung bei Discord.
     * @return ApplicationCommandRequest
     */
    ApplicationCommandRequest getCommandRequest();

    /**
     * Behandelt das Ereignis, wenn der Command ausgeführt wird.
     * @param event Das ChatInputInteractionEvent.
     * @param client Der GatewayDiscordClient.
     * @return Mono<Void>
     */
    Mono<Void> handle(ChatInputInteractionEvent event, GatewayDiscordClient client);

    /**
     * Gibt den Namen des Commands zurück, der für die Zuordnung verwendet wird.
     * @return Der Name des Commands.
     */
    default String getName() {
        return getCommandRequest().name();
    }
}
