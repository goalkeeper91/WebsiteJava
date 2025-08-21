package streamer_website.demo.client;

import com.github.philippheuer.credentialmanager.domain.OAuth2Credential;
import com.github.twitch4j.TwitchClient;
import com.github.twitch4j.TwitchClientBuilder;
import com.github.twitch4j.chat.events.channel.ChannelMessageEvent;
import com.github.twitch4j.common.enums.CommandPermission;
import lombok.Getter;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import streamer_website.demo.commands.twitch.TwitchBotModCommand;
import streamer_website.demo.service.twitch.TwitchCommandService;
import streamer_website.demo.service.twitch.TwitchTokenService;

import java.time.Instant;

public class TwitchBot {

    @Getter
    private TwitchClient client;

    @Getter
    private final String channelName;

    @Getter
    private Instant startTime;

    private final TwitchCommandService commandService;
    private final TwitchTokenService twitchTokenService;

    private static final Logger logger = LoggerFactory.getLogger(TwitchBot.class);

    @Getter
    private boolean running = false;

    public TwitchBot(String channelName, TwitchCommandService commandService, TwitchTokenService twitchTokenService) {
        this.channelName = channelName;
        this.commandService = commandService;
        this.twitchTokenService = twitchTokenService;
    }

    public void start(String botUserIdFromDB) {
        if (running) {
            logger.warn("Bot ist bereits gestartet");
            return;
        }

        String accessToken = twitchTokenService.getBotAccessToken();

        OAuth2Credential credential = new OAuth2Credential("twitch", accessToken);

        client = TwitchClientBuilder.builder()
                .withEnableHelix(true)
                .withEnableChat(true)
                .withChatAccount(credential)
                .build();

        client.getChat().joinChannel(channelName);
        client.getEventManager().onEvent(ChannelMessageEvent.class, this::onMessage);

        startTime = Instant.now();
        running = true;
        logger.info("TwitchBot gestartet und Kanal {} beigetreten", channelName);
    }

    public void stop() {
        if (client != null) {
            client.close();
            client = null;
        }

        running = false;
        startTime = null;
        logger.info("TwitchBot gestoppt");
    }

    private void onMessage(ChannelMessageEvent event) {
        boolean isMod = event.getPermissions().contains(CommandPermission.MODERATOR)
                || event.getPermissions().contains(CommandPermission.BROADCASTER);
        String message = event.getMessage();

        if (!message.startsWith("!")) return;

        String[] parts = message.split(" ", 3);
        String trigger = parts[0].substring(1).toLowerCase();

        try {
            if (isMod) {
                TwitchBotModCommand cmd = TwitchBotModCommand.fromTrigger(trigger);
                if (cmd != null) {
                    try {
                        String[] args = parts;

                        if (trigger.equals("title") || trigger.equals("category")) {
                            String fullArgument = message.substring(message.indexOf(" ") + 1);
                            args = new String[] { trigger, fullArgument };
                        }

                        cmd.execute(args, event, client, commandService);
                    } catch (Exception e) {
                        logger.error("Fehler beim Ausführen des Mod-Commands {} mit Nachricht: {}",
                                trigger, message, e);
                        client.getChat().sendMessage(event.getChannel().getName(),
                                "Fehler beim Ausführen des Commands: " + trigger);
                    }
                }
            }

            commandService.getCommand(trigger).ifPresent(cmd ->
                    client.getChat().sendMessage(event.getChannel().getName(), cmd.getResponse()));

        } catch (Exception e) {
            client.getChat().sendMessage(event.getChannel().getName(), "Fehler beim Verarbeiten des Command: " + event.getMessage());
            logger.warn("Fehler beim Command", e);
        }
    }
}
