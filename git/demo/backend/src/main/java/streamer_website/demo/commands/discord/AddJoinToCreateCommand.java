package streamer_website.demo.commands.discord;

import discord4j.core.event.domain.message.MessageCreateEvent;
import org.springframework.stereotype.Component;
import streamer_website.demo.entity.discord.JoinToCreateConfig;
import streamer_website.demo.repository.JoinToCreateRepository;
import streamer_website.demo.service.discord.JoinToCreateService;

@Component
public class AddJoinToCreateCommand implements Command{
    private final JoinToCreateRepository repository;
    private final JoinToCreateService joinToCreateService;

    public AddJoinToCreateCommand(JoinToCreateRepository repository,
                                  JoinToCreateService joinToCreateService) {
        this.repository = repository;
        this.joinToCreateService = joinToCreateService;
    }

    @Override
    public String getName() {
        return "add";
    }

    @Override
    public void execute(MessageCreateEvent event, String[] args) {
        if (args.length < 5) {
            event.getMessage().getChannel()
                    .flatMap(ch -> ch.createMessage("Syntax: !jtc add <joinChannelId> <categoryId> <prefix> <userLimit> <private>"))
                    .subscribe();
            return;
        }

        Long joinChannelId = Long.valueOf(args[0]);
        Long categoryId = Long.valueOf(args[1]);
        String prefix = args[2];
        Integer userLimit = Integer.valueOf(args[3]);
        Boolean priv = Boolean.valueOf(args[4]);

        JoinToCreateConfig cfg = new JoinToCreateConfig();
        cfg.setJoinChannelId(joinChannelId);
        cfg.setCategoryId(categoryId);
        cfg.setChannelNamePrefix(prefix);
        cfg.setUserLimit(userLimit);
        cfg.setPrivateChannel(priv);

        repository.save(cfg);
        joinToCreateService.initConfigs();

        event.getMessage().getChannel()
                .flatMap(ch -> ch.createMessage("✅ Join-to-Create-Config gespeichert"))
                .subscribe();
    }
}
