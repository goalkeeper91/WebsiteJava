package streamer_website.demo.dto;

import lombok.Data;

@Data
public class JoinToCreateConfigResponse {

    private Long id;
    private Long joinChannelId;
    private Long categoryId;
    private String channelNamePrefix;
    private Integer userLimit;
    private Boolean privateChannel;
}
