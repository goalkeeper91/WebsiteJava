package streamer_website.demo.mapper.twitch;

import streamer_website.demo.dto.twitch.ChatCommandDto;
import streamer_website.demo.entity.twitch.ChatCommand;

public class ChatCommandMapper {

    public static ChatCommandDto toDto(ChatCommand entity) {
        ChatCommandDto dto = new ChatCommandDto();
        dto.setId(entity.getId());
        dto.setTrigger(entity.getTrigger());
        dto.setResponse(entity.getResponse());
        dto.setCooldown(entity.getCooldown());
        dto.setEnabled(entity.getEnabled());
        dto.setCreatedAt(entity.getCreatedAt());
        dto.setUpdatedAt(entity.getUpdatedAt());
        return dto;
    }
}

