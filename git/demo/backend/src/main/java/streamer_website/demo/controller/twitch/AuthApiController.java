package streamer_website.demo.controller.twitch;

import jakarta.servlet.http.HttpSession;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import streamer_website.demo.dto.TwitchUser;

@RestController
@RequestMapping("/api/auth")
public class AuthApiController {

    @GetMapping("/me")
    public ResponseEntity<?> getLoggedInUser(HttpSession session) {
        TwitchUser user = (TwitchUser) session.getAttribute("user");

        if (user == null) {
            return ResponseEntity.status(HttpStatus.UNAUTHORIZED).body("Not logged in");
        }
        return ResponseEntity.ok(user);
    }

    @PostMapping("/logout")
    public ResponseEntity<?> logout(HttpSession session) {
        session.invalidate();

        return ResponseEntity.ok("Logged out");
    }
}
