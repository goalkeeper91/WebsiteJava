package streamer_website.demo.controller.twitch;

import jakarta.servlet.http.HttpSession;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.server.ResponseStatusException;
import streamer_website.demo.dto.twitch.*;
import streamer_website.demo.service.twitch.ActivityService;
import streamer_website.demo.service.twitch.StreamService;

import java.util.List;

@RestController
@RequiredArgsConstructor
@RequestMapping("/api/dashboard/stream")
public class StreamController {

    private final StreamService streamService;
    private final ActivityService activityService;

    @GetMapping("/info")
    public StreamInfoDto getStreamInfo(HttpSession session) {
        TwitchUser user = requireUser(session);
        return streamService.getStreamInfo(user.id());
    }

    @PatchMapping("/info")
    public StreamInfoDto updateStreamInfo(
            HttpSession session,
            @RequestBody UpdateStreamInfoRequest request
    ) {
        TwitchUser user = requireUser(session);
        return streamService.updateStreamInfo(
                user.id(),
                request.getTitle(),
                request.getGameId()
        );
    }

    @GetMapping("/live")
    public LiveStreamDto getLiveStatus(HttpSession session) {
        TwitchUser user = requireUser(session);
        return streamService.getLiveStatus(user.id());
    }

    @GetMapping("/stats")
    public DashboardStatsDto getStats(HttpSession session) {
        TwitchUser user = requireUser(session);
        return streamService.getDashboardStats(user.id());
    }

    @GetMapping("/activities")
    public List<ActivityDto> getRecentActivities(
            HttpSession session,
            @RequestParam(defaultValue = "50") int limit
    ) {
        TwitchUser user = requireUser(session);
        return activityService.getRecentActivities(user.id(), limit);
    }

    @GetMapping("/categories/search")
    public List<CategoryDto> searchCategories(
            HttpSession session,
            @RequestParam String query
    ) {
        TwitchUser user = requireUser(session);
        return streamService.searchCategories(user.id(), query);
    }

    private TwitchUser requireUser(HttpSession session) {
        TwitchUser user = (TwitchUser) session.getAttribute("user");
        if (user == null) {
            throw new ResponseStatusException(
                    HttpStatus.UNAUTHORIZED,
                    "Not logged in"
            );
        }
        return user;
    }
}