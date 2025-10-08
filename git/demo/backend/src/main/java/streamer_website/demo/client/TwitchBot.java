package streamer_website.demo.client;

import com.github.philippheuer.credentialmanager.domain.OAuth2Credential;
import com.github.twitch4j.TwitchClient;
import com.github.twitch4j.TwitchClientBuilder;
import com.github.twitch4j.chat.events.channel.ChannelMessageEvent;
import com.github.twitch4j.common.enums.CommandPermission;
import com.github.twitch4j.helix.domain.StreamList;
import lombok.Getter;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import streamer_website.demo.commands.twitch.TwitchBotModCommand;
import streamer_website.demo.service.twitch.TwitchCommandService;
import streamer_website.demo.service.twitch.TwitchTokenService;

import java.time.Instant;
import java.util.Collections;
import java.util.concurrent.*;

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
    private String botUserIdFromDB;

    private final ScheduledExecutorService scheduler = Executors.newScheduledThreadPool(1);
    private final ScheduledExecutorService timerScheduler = Executors.newScheduledThreadPool(4);
    private final ConcurrentHashMap<Long, ScheduledFuture<?>> activeTimers = new ConcurrentHashMap<>();
    private boolean runningTimersStarted = false;

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

        this.botUserIdFromDB = botUserIdFromDB;

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

    public void restart() {
        logger.info("Starte den Bot neu...");

        stop();

        try {
            start(this.botUserIdFromDB);
            logger.info("TwitchBot erfolgreich neu gestartet");
        } catch (Exception e) {
            logger.error("Fehler beim Neustarten des Bots", e);
        }
    }

    public boolean isStreamLive() {
        StreamList streams = client.getHelix()
                .getStreams(
                        null,
                        null,
                        null,
                        1,
                        null,
                        null,
                        null,
                        Collections.singletonList(channelName.toLowerCase())
                )
                .execute();

        return !streams.getStreams().isEmpty();
    }

    private void onMessage(ChannelMessageEvent event) {
        boolean isMod = event.getPermissions().contains(CommandPermission.MODERATOR)
                || event.getPermissions().contains(CommandPermission.BROADCASTER);
        String message = event.getMessage();

        if (!message.startsWith("!")) return;

        String[] parts = message.substring(1).split(";");
        String trigger = parts[0].toLowerCase();

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

    public void startLiveStatusPolling() {
        scheduler.scheduleAtFixedRate(() -> {
            try {
                boolean live = isStreamLive();
                if (live && !runningTimersStarted) {
                    logger.info("Stream ist live — Timer-Commands starten");
                    startAllTimerCommands();
                    runningTimersStarted = true;
                } else if (!live && runningTimersStarted) {
                    logger.info("Stream ist offline — Timer-Commands stoppen");
                    stopAllTimerCommands();
                    runningTimersStarted = false;
                }
            } catch (Exception e) {
                logger.error("Fehler beim Live-Check", e);
            }
        }, 0, 10, TimeUnit.SECONDS);
    }

    public void startAllTimerCommands() {
        commandService.getTimerCommands().forEach(cmd -> {
            if (activeTimers.containsKey(cmd.getId())) return;

            ScheduledFuture<?> future = timerScheduler.scheduleAtFixedRate(() -> {
                try {
                    if (client != null) {
                        client.getChat().sendMessage(channelName, cmd.getResponse());
                    }
                } catch (Exception e) {
                    logger.error("Fehler beim Timer-Command '{}'", cmd.getTrigger(), e);
                }
            }, cmd.getDuration(), cmd.getDuration(), TimeUnit.SECONDS);

            activeTimers.put(cmd.getId(), future);
            logger.info("Timer für Command '{}' gestartet (alle {} Sekunden)", cmd.getTrigger(), cmd.getDuration());
        });
    }

    public void stopAllTimerCommands() {
        activeTimers.forEach((id, future) -> {
            future.cancel(true);
            logger.info("Timer für Command-ID {} gestoppt", id);
        });
        activeTimers.clear();
    }
}
