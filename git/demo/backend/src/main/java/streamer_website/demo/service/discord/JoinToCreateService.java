package streamer_website.demo.service.discord;

import discord4j.common.util.Snowflake;
import discord4j.core.GatewayDiscordClient;
import discord4j.core.event.domain.VoiceStateUpdateEvent;
import discord4j.core.object.PermissionOverwrite;
import discord4j.core.object.VoiceState;
import discord4j.core.object.entity.Guild;
import discord4j.core.object.entity.Member;
import discord4j.core.object.entity.channel.VoiceChannel;
import discord4j.core.spec.GuildMemberEditSpec;
import discord4j.core.spec.VoiceChannelCreateSpec;
import discord4j.rest.util.Permission;
import discord4j.rest.util.PermissionSet;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import reactor.core.publisher.Mono;
import streamer_website.demo.entity.discord.JoinToCreateConfig;
import streamer_website.demo.repository.JoinToCreateRepository;

import java.util.*;
import java.util.concurrent.ConcurrentHashMap;
import java.util.stream.Collectors;

@Service
@RequiredArgsConstructor
public class JoinToCreateService {

    private final JoinToCreateRepository repository;
    private final Set<Snowflake> createdChannelIds = ConcurrentHashMap.newKeySet();
    private Map<Snowflake, JoinToCreateConfig> configs = new ConcurrentHashMap<>();

    /**
     * Initialisiert die Configs aus der DB.
     * Null-safe: Leere DB ist erlaubt.
     */
    public void initConfigs() {
        try {
            List<JoinToCreateConfig> allConfigs = repository.findAll();
                configs = allConfigs.stream()
                    .filter(c -> c.getJoinChannelId() != null && c.getCategoryId() != null)
                    .collect(Collectors.toConcurrentMap(
                                c -> Snowflake.of(c.getJoinChannelId()),
                                c -> c
                ));
        } catch (Exception e) {
            System.err.println("JoinToCreateService initConfigs Fehler (DB evtl. nicht bereit): " + e.getMessage());
            configs = new ConcurrentHashMap<>();
        }

    }

    /**
     * Registriert die VoiceState Listener.
     * Listener reagieren nur auf gültige Channels.
     * Fehler werden geloggt, Backend crasht nicht.
     */
    public void register(GatewayDiscordClient client) {

        // Listener: User joined / moved
        client.on(VoiceStateUpdateEvent.class)
                .flatMap(event -> {
                    if (event.getCurrent() == null) return Mono.empty();

                    Optional<Snowflake> currentChannelIdOpt = event.getCurrent().getChannelId();
                    if (currentChannelIdOpt.isEmpty()) return Mono.empty();

                    Snowflake currentChannelId = currentChannelIdOpt.get();
                    JoinToCreateConfig cfg = configs.get(currentChannelId);
                    if (cfg == null) return Mono.empty(); // Kein Join-to-Create

                    if (event.getOld()
                            .flatMap(VoiceState::getChannelId)
                            .map(id -> id.equals(currentChannelId))
                            .orElse(false)) return Mono.empty();

                    Snowflake guildId = event.getCurrent().getGuildId();
                    if (guildId == null) return Mono.empty();

                    Mono<Guild> guildMono = client.getGuildById(guildId);
                    Mono<Member> memberMono = event.getCurrent().getMember()
                            .switchIfEmpty(guildMono.flatMap(g -> g.getMemberById(event.getCurrent().getUserId())));

                    return Mono.zip(guildMono, memberMono)
                            .flatMap(tuple -> {
                                Guild guild = tuple.getT1();
                                Member member = tuple.getT2();
                                if (member == null) return Mono.empty();

                                return createTempChannel(guild, member, cfg)
                                        .flatMap(newChannel -> {
                                            createdChannelIds.add(newChannel.getId());
                                            return member.edit(GuildMemberEditSpec.builder()
                                                            .newVoiceChannel(newChannel.getId())
                                                            .build())
                                                    .thenReturn(newChannel);
                                        })
                                        .onErrorResume(e -> {
                                            System.err.println("JoinToCreate Fehler: " + e.getMessage());
                                            return Mono.empty();
                                        });
                            });
                })
                .subscribe();

        // Listener: User verlässt temporären Channel
        client.on(VoiceStateUpdateEvent.class)
                .flatMap(event -> {
                    Optional<Snowflake> oldChannelIdOpt = event.getOld().flatMap(VoiceState::getChannelId);
                    if (oldChannelIdOpt.isEmpty()) return Mono.empty();

                    Snowflake oldChannelId = oldChannelIdOpt.get();
                    if (!createdChannelIds.contains(oldChannelId)) return Mono.empty();

                    Snowflake guildId = event.getCurrent() != null ? event.getCurrent().getGuildId() : null;
                    if (guildId == null) return Mono.empty();

                    return client.getGuildById(guildId)
                            .flatMap(guild -> guild.getChannelById(oldChannelId))
                            .ofType(VoiceChannel.class)
                            .flatMap(vc -> vc.getVoiceStates().collectList()
                                    .flatMap(list -> {
                                        if (list.isEmpty()) {
                                            return vc.delete()
                                                    .then(Mono.fromRunnable(() -> createdChannelIds.remove(oldChannelId)));
                                        }
                                        return Mono.empty();
                                    }))
                            .onErrorResume(e -> {
                                System.err.println("Fehler beim Löschen temporären Channels: " + e.getMessage());
                                return Mono.empty();
                            });
                })
                .subscribe();
    }

    /**
     * Erstellt einen temporären VoiceChannel.
     * @param guild Guild
     * @param member Member
     * @param cfg JoinToCreateConfig
     */
    private Mono<VoiceChannel> createTempChannel(Guild guild, Member member, JoinToCreateConfig cfg) {
        VoiceChannelCreateSpec.Builder specBuilder = VoiceChannelCreateSpec.builder()
                .name(cfg.getChannelNamePrefix() + " - " + member.getDisplayName())
                .parentId(Snowflake.of(cfg.getCategoryId()))
                .userLimit(cfg.getUserLimit() != null ? cfg.getUserLimit() : 0);

        if (Boolean.TRUE.equals(cfg.getPrivateChannel())) {
            return guild.getEveryoneRole()
                    .flatMap(everyoneRole -> {
                        List<PermissionOverwrite> overwrites = List.of(
                                PermissionOverwrite.forRole(
                                        everyoneRole.getId(),
                                        PermissionSet.none(),
                                        PermissionSet.of(Permission.CONNECT)
                                ),
                                PermissionOverwrite.forMember(
                                        member.getId(),
                                        PermissionSet.of(Permission.CONNECT, Permission.SPEAK),
                                        PermissionSet.none()
                                )
                        );
                        specBuilder.permissionOverwrites(overwrites);
                        return guild.createVoiceChannel(specBuilder.build());
                    });
        } else {
            return guild.createVoiceChannel(specBuilder.build());
        }
    }
}