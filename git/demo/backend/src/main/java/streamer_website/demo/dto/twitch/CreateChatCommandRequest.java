package streamer_website.demo.dto.twitch;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Size;
import lombok.Data;

@Data
public class CreateChatCommandRequest {

    @NotBlank
    @Size(max = 100)
    private String trigger;

    @NotBlank
    private String response;

    private Integer cooldown;
}

