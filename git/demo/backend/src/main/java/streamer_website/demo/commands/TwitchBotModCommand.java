package streamer_website.demo.commands;

import com.github.twitch4j.TwitchClient;
import com.github.twitch4j.chat.events.channel.ChannelMessageEvent;
import com.github.twitch4j.helix.domain.CategorySearchList;
import com.github.twitch4j.helix.domain.ChannelInformation;
import com.github.twitch4j.helix.domain.ChannelInformationList;
import com.netflix.hystrix.HystrixCommand;
import lombok.Getter;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import streamer_website.demo.service.TwitchCommandService;

import java.util.List;

public enum TwitchBotModCommand {

    ADD("add", (parts, event, client, service) -> {
        if (parts.length == 3) {
            String newTrigger = parts[1].toLowerCase();
            String response = parts[2];
            service.addCommand(newTrigger, response, false);
            client.getChat().sendMessage(event.getChannel().getName(),
                    "Command !" + newTrigger + " hinzugefügt!");
        }
    }),

    EDIT("edit", (parts, event, client, service) -> {
        if (parts.length == 3) {
            String editTrigger = parts[1].toLowerCase();
            String newResponse = parts[2];
            service.editCommand(editTrigger, newResponse);
            client.getChat().sendMessage(event.getChannel().getName(), "Command !" + editTrigger + " aktualisiert!");
        }
    }),

    DELETE("delete", (parts, event, client, service) -> {
        if (parts.length >= 2) {
            String deleteTrigger = parts[1].toLowerCase();
            service.deleteCommand(deleteTrigger);
            client.getChat().sendMessage(event.getChannel().getName(),
                    "Command !" + deleteTrigger + " gelöscht!");
        }
    }),

    SET_CATEGORY("category" , (parts, event, client, service) -> {
        if (parts.length >= 2) {
            String newCategory = parts[1];
            String botToken = service.getUserToken();

            try {
                CategorySearchList result = client.getHelix().searchCategories(
                        botToken,
                        newCategory,
                        1,
                        null
                ).execute();

                if (result.getResults().isEmpty()) {
                    client.getChat().sendMessage(event.getChannel().getName(),
                            "Kategorie '" + newCategory + "' wurde nicht gefunden.");
                    return;
                }

                String gameId = result.getResults().getFirst().getId();
                String foundName = result.getResults().getFirst().getName();

                HystrixCommand<ChannelInformationList> infoList = client.getHelix().getChannelInformation(botToken, List.of(event.getChannel().getId()));

                ChannelInformation oldInfo = infoList.execute().getChannels().getFirst();
                ChannelInformation information = oldInfo.toBuilder()
                        .gameId(gameId)
                        .gameName(foundName)
                        .delay(null)
                        .build();
                client.getHelix().updateChannelInformation(botToken, event.getChannel().getId(), information).execute();
                client.getChat().sendMessage(event.getChannel().getName(),
                        "Kategorie auf \"" + foundName + "\" gesetzt!");
            } catch (Exception e){
                client.getChat().sendMessage(event.getChannel().getName(),
                        "Fehler beim Setzen der Kategorie.");
                throw new RuntimeException("Es ist ein Fehler aufgetreten {}", e);
            }
        }
    }),

    SET_TITLE("title", (parts, event, client, service) -> {
        if (parts.length >= 2) {
            String dynamicPart = parts[1];
            String fixedSuffix = " | kw-com.de | !emp !dornfinger";
            String newTitle = dynamicPart + fixedSuffix;
            String botToken = service.getUserToken();

            try {
                HystrixCommand<ChannelInformationList> infoList = client.getHelix().getChannelInformation(botToken, List.of(event.getChannel().getId()));

                ChannelInformation oldInfo = infoList.execute().getChannels().getFirst();

                ChannelInformation information = oldInfo.toBuilder()
                        .title(newTitle)
                        .delay(null)
                        .build();

                client.getHelix().updateChannelInformation(botToken, event.getChannel().getId(), information).execute();
                client.getChat().sendMessage(event.getChannel().getName(),
                        "Titel auf \"" + newTitle + "\" gesetzt!");
            } catch (Exception e){
                client.getChat().sendMessage(event.getChannel().getName(),
                        "Fehler beim Setzen des Titels.");
                throw new RuntimeException("Es ist ein Fehler aufgetreten {}", e);
            }
        }
    });

    @Getter
    private final String trigger;
    private final CommandExecutor executor;

    TwitchBotModCommand(String trigger, CommandExecutor executor) {
        this.trigger = trigger;
        this.executor = executor;
    }

    public void execute(String[] parts, ChannelMessageEvent event, TwitchClient client, TwitchCommandService service) {
        executor.execute(parts, event, client, service);
    }

    public static TwitchBotModCommand fromTrigger(String trigger) {
        for (TwitchBotModCommand cmd : values()) {
            if (cmd.trigger.equalsIgnoreCase(trigger)) {
                return cmd;
            }
        }
        return null;
    }

    @FunctionalInterface
    public interface CommandExecutor {
        void execute(String[] parts, ChannelMessageEvent event, TwitchClient client, TwitchCommandService service);
    }
}
