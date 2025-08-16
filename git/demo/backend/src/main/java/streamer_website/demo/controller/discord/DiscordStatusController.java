package streamer_website.demo.controller.discord;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import streamer_website.demo.service.discord.StatusService;

import java.util.HashMap;
import java.util.Map;

@RestController
@RequestMapping("/api/discord")
public class DiscordStatusController {

    private final StatusService statusService;

    public DiscordStatusController(StatusService statusService) {
        this.statusService = statusService;
    }

    @GetMapping("/status")
    public Map<String, Object> getStatus() {
        Map<String, Object> status = new HashMap<>();
        status.put("running", statusService.isRunning());
        status.put("since", statusService.getStartTime());
        return status;
    }
}
