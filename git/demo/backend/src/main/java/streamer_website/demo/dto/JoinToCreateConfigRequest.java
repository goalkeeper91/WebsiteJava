package streamer_website.demo.dto;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.PositiveOrZero;
import lombok.Data;

@Data
public class JoinToCreateConfigRequest {

    @NotNull(message = "joinChannelId darf nicht leer sein")
    private Long joinChannelId;

    @NotNull(message = "categoryId darf nicht leer sein")
    private Long categoryId;

    @NotBlank(message = "channelNamePrefix darf nicht leer sein")
    private String channelNamePrefix;

    @PositiveOrZero(message = "userLimit muss >= 0 sein")
    private Integer userLimit;

    @NotNull(message = "privateChannel muss gesetzt sein")
    private Boolean privateChannel;
}
