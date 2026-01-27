package streamer_website.demo.dto.twitch;

import lombok.Data;

@Data
public class UpdateChatCommandRequest {

    private String trigger;
    private String response;
    private Integer cooldown;
    private Boolean enabled;
}

