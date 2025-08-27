package streamer_website.demo.service.discord;

import discord4j.common.util.Snowflake;
import discord4j.core.GatewayDiscordClient;
import discord4j.core.object.entity.channel.MessageChannel;
import discord4j.core.spec.EmbedCreateSpec;
import discord4j.core.spec.MessageCreateSpec;
import discord4j.rest.util.Color;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import reactor.core.publisher.Mono;
import streamer_website.demo.entity.ContactRequest;

import java.time.ZoneId;

@Service
public class DiscordNotificationService {

    private final GatewayDiscordClient client;

    @Value("${discord.contact.channel-id}")
    private String channelId;

    public DiscordNotificationService(GatewayDiscordClient client) {
        this.client = client;
    }

    public Mono<Void> notifyNewContactRequest(ContactRequest request) {
        return client.getChannelById(Snowflake.of(channelId))
                .ofType(MessageChannel.class)
                .flatMap(channel -> channel.createMessage(
                        MessageCreateSpec.builder()
                                .addEmbed(buildEmbed(request))
                                .build()
                ))
                .then();
    }

    private EmbedCreateSpec buildEmbed(ContactRequest request) {
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
}
