package streamer_website.demo.client;

import com.github.philippheuer.credentialmanager.domain.OAuth2Credential;
import com.github.twitch4j.TwitchClient;
import com.github.twitch4j.TwitchClientBuilder;
import com.github.twitch4j.chat.events.channel.ChannelMessageEvent;
import com.github.twitch4j.common.enums.CommandPermission;
import lombok.Getter;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import streamer_website.demo.service.BotOAuthService;
import streamer_website.demo.service.TwitchCommandService;

import java.time.Instant;

public class TwitchBot {

    @Getter
    private TwitchClient client;

    @Getter
    private final String channelName;

    @Getter
    private Instant startTime;

    private final TwitchCommandService commandService;
    private final BotOAuthService botOAuthService;

    private static final Logger logger = LoggerFactory.getLogger(TwitchBot.class);

    @Getter
    private boolean running = false;

    public TwitchBot(String channelName, TwitchCommandService commandService, BotOAuthService botOAuthService) {
        this.channelName = channelName;
        this.commandService = commandService;
        this.botOAuthService = botOAuthService;
    }

    public void start(String botUserIdFromDB) {
        if (running) {
            logger.warn("Bot ist bereits gestartet");
            return;
        }

        OAuth2Credential credential = botOAuthService.getBotCredential(botUserIdFromDB);

        client = TwitchClientBuilder.builder()
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
        String user = event.getUser().getName();
        boolean isMod = event.getPermissions().contains(CommandPermission.MODERATOR)
                || event.getPermissions().contains(CommandPermission.BROADCASTER);
        String message = event.getMessage();
        String channel = event.getChannel().getName();

        System.out.println("[" + channel + "] " + user + ": " + message);

        if (!message.startsWith("!")) return;

        String[] parts = message.split(" ", 3);
        String trigger = parts[0].substring(1).toLowerCase();

        try {
            switch (trigger) {
                case "add":
                    if (isMod && parts.length == 3) {
                        String newTrigger = parts[1].toLowerCase();
                        String response = parts[2];
                        commandService.addCommand(newTrigger, response, false);
                        client.getChat().sendMessage(channel, "Command !" + newTrigger + " hinzugefügt!");
                    } else if (!isMod) {
                        client.getChat().sendMessage(channel, "Lass den Mist sonst stille Treppe!!!111!11");
                    }
                    break;

                case "edit":
                    if (isMod && parts.length == 3) {
                        String editTrigger = parts[1].toLowerCase();
                        String newResponse = parts[2];
                        commandService.editCommand(editTrigger, newResponse);
                        client.getChat().sendMessage(channel, "Command !" + editTrigger + " aktualisiert!");
                    } else if (!isMod) {
                        client.getChat().sendMessage(channel, "Dein Ernst?");
                    }
                    break;

                case "delete":
                    if (isMod && parts.length >= 2) {
                        String deleteTrigger = parts[1].toLowerCase();
                        commandService.deleteCommand(deleteTrigger);
                        client.getChat().sendMessage(channel, "Command !" + deleteTrigger + " gelöscht!");
                    } else if (!isMod) {
                        client.getChat().sendMessage(channel, "Es reicht, ab auf die Stille Treppe");
                    }
                    break;

                default:
                    commandService.getCommand(trigger).ifPresent(cmd ->
                            client.getChat().sendMessage(channel, cmd.getResponse()));
            }
        } catch (Exception e) {
            client.getChat().sendMessage(channel, "Fehler beim Verarbeiten des Command: " + event.getMessage());
            logger.warn("Fehler beim Command", e);
        }
    }
}
