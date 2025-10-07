package streamer_website.demo.service.discord;

import discord4j.core.GatewayDiscordClient;
import discord4j.core.event.domain.message.MessageCreateEvent;
import discord4j.core.object.entity.User;
import discord4j.core.event.domain.interaction.ChatInputInteractionEvent;
import discord4j.rest.service.ApplicationService;
import jakarta.annotation.PostConstruct;
import lombok.RequiredArgsConstructor;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;
import streamer_website.demo.commands.discord.Command;
import streamer_website.demo.commands.discord.SlashCommand;

import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.function.Function;
import java.util.stream.Collectors;

@Service
@RequiredArgsConstructor
public class CommandService {

    // Alte Commands (für !jtc)
    private final List<Command> messageCommands;
    private final Map<String, Command> messageCommandMap = new HashMap<>();

    // Neue Slash Commands (für /task)
    private final List<SlashCommand> slashCommands; // Spring injiziert alle SlashCommand-Beans
    private final Map<String, SlashCommand> slashCommandMap = new HashMap<>();

    @Value("${discord.bot.client-id}")
    private long applicationId;

    @PostConstruct
    public void init() {
        for (Command cmd : messageCommands) {
            messageCommandMap.put(cmd.getName().toLowerCase(), cmd);
        }

        slashCommandMap.putAll(slashCommands.stream()
                .collect(Collectors.toMap(
                        command -> command.getCommandRequest().name(),
                        Function.identity()
                )));
    }

    public void register(GatewayDiscordClient client) {
        ApplicationService applicationService = client.getRestClient().getApplicationService();
        long currentApplicationId = applicationId;

        Mono.just(slashCommands)
                .flatMapMany(Flux::fromIterable)
                .map(SlashCommand::getCommandRequest)
                .collectList()
                .flatMapMany(requests -> applicationService.bulkOverwriteGlobalApplicationCommand(currentApplicationId, requests))
                .doOnNext(cmd -> System.out.println("Discord Slash Command registriert: /" + cmd.name()))
                .subscribe();

        client.on(ChatInputInteractionEvent.class)
                .flatMap(event -> {
                    SlashCommand command = slashCommandMap.get(event.getCommandName());

                    if (command == null) {
                        return event.reply("Unbekannter Slash Command: /" + event.getCommandName()).withEphemeral(true);
                    }

                    return command.handle(event, client)
                            .onErrorResume(e -> event.editReply("❌ Ein Fehler ist bei der Ausführung des Commands aufgetreten.").then());
                })
                .onErrorResume(e -> {
                    e.printStackTrace();
                    return Mono.empty();
                })
                .subscribe();

        client.on(MessageCreateEvent.class)
                .flatMap(event -> {
                    if (event.getMessage().getAuthor().map(User::isBot).orElse(true)) {
                        return Mono.empty();
                    }

                    String content = event.getMessage().getContent().trim();
                    if (!content.startsWith("!jtc")) {
                        return Mono.empty();
                    }

                    String[] parts = content.split("\\s+");
                    if (parts.length < 2) {
                        return Mono.empty();
                    }

                    String commandName = parts[1].toLowerCase();
                    String[] args = Arrays.copyOfRange(parts, 2, parts.length);

                    Command cmd = messageCommandMap.get(commandName);
                    if (cmd == null) {
                        return event.getMessage().getChannel()
                                .flatMap(ch -> ch.createMessage("Unbekannter Befehl: " + commandName))
                                .then();
                    }

                    return cmd.execute(event, args);
                })
                .subscribe();
    }
}
