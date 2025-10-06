package streamer_website.demo.service.discord;

import discord4j.common.util.Snowflake;
import discord4j.core.GatewayDiscordClient;
import discord4j.core.object.entity.channel.MessageChannel;
import discord4j.core.spec.EmbedCreateSpec;
import discord4j.core.spec.MessageCreateSpec;
import discord4j.rest.util.Color;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;
import reactor.core.publisher.Mono;
import streamer_website.demo.entity.ContactRequest;
import streamer_website.demo.entity.discord.DiscordChannels;
import streamer_website.demo.entity.twitch.TwitchCommand;
import streamer_website.demo.repository.DiscordChannelsRepository;

import java.time.ZoneId;
import jakarta.inject.Provider;


@Service
public class DiscordNotificationService {

    private final Provider<GatewayDiscordClient> client;
    private final DiscordChannelsRepository channelsRepo;
    private static final Logger logger = LoggerFactory.getLogger(DiscordNotificationService.class);

    public DiscordNotificationService(Provider<GatewayDiscordClient> client, DiscordChannelsRepository channelsRepo) {
        this.client = client;
        this.channelsRepo = channelsRepo;
    }

    // Hilfsmethode für den einfachen Zugriff auf den Client
    private GatewayDiscordClient getClient() {
        return client.get();
    }

    public Mono<Void> notifyNewContactRequest(ContactRequest request) {
        String channelId = channelsRepo.findByDescriptionContainingIgnoreCase("kontakt")
                .stream().findFirst()
                .map(DiscordChannels::getChannelId)
                .orElseGet(() -> {
                    logger.warn("Kein passender Discord-Channel gefunden (beschreibung enthält 'kontakt')");
                    return null;
                });

        if (channelId == null) {
            return Mono.empty();
        }

        // 3. ÄNDERUNG: Zugriff über getClient()
        return getClient().getChannelById(Snowflake.of(channelId))
                .ofType(MessageChannel.class)
                .flatMap(channel -> channel.createMessage(
                        MessageCreateSpec.builder()
                                .addEmbed(buildContactEmbed(request))
                                .build()
                ))
                .then();
    }

    public Mono<Void> notifyNewTwitchCommandRequest(TwitchCommand command) {
        String channelId = channelsRepo.findByDescriptionContainingIgnoreCase("twitch command")
                .stream().findFirst()
                .map(DiscordChannels::getChannelId)
                .orElseGet(() -> {
                    logger.warn("Kein passender Discord-Channel gefunden (beschreibung enthält 'twitch command')");
                    return null;
                });

        if (channelId == null) {
            return Mono.empty();
        }

        // 4. ÄNDERUNG: Zugriff über getClient()
        return getClient().getChannelById(Snowflake.of(channelId))
                .ofType(MessageChannel.class)
                .flatMap(channel -> channel.createMessage(
                        MessageCreateSpec.builder()
                                .addEmbed(buildCommandEmbed(command))
                                .build()
                ))
                .then();
    }

    private EmbedCreateSpec buildContactEmbed(ContactRequest request) {
        EmbedCreateSpec.Builder builder = EmbedCreateSpec.builder()
                .title("📩 Neue Kontaktanfrage")
                .color(Color.BLUE)
                .addField("👤 Name", request.getName(), false)
                .addField("📧 Email", request.getEmail(), false);

        if (request.getPhone() != null && !request.getPhone().isBlank()) {
            builder.addField("📞 Telefon", request.getPhone(), false);
        }

        builder.addField("📝 Betreff", request.getSubject(), false)
                .addField("💬 Nachricht", request.getMessage(), false)
                .footer("Eingegangen am", null)
                .timestamp(request.getCreatedAt().atZone(ZoneId.systemDefault()).toInstant());

        return builder.build();
    }

    private EmbedCreateSpec buildCommandEmbed(TwitchCommand command) {
        EmbedCreateSpec.Builder builder = EmbedCreateSpec.builder()
                .title("📩 Neuer Command")
                .color(Color.BLUE)
                .addField("👤 Trigger", command.getTrigger(), false);

        builder.addField("💬 Nachricht", command.getResponse(), false)
                .footer("Erstellt am", null)
                .timestamp(command.getCreatedAt().atZone(ZoneId.systemDefault()).toInstant());

        return builder.build();
    }
}
