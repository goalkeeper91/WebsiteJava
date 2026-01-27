package streamer_website.demo.controller.twitch;

import jakarta.servlet.http.HttpSession;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import streamer_website.demo.dto.twitch.TwitchUser;

import java.util.Map;

@RestController
@RequestMapping("/api/auth")
public class AuthApiController {

    @GetMapping("/me")
    public ResponseEntity<?> getLoggedInUser(HttpSession session) {
        // Hol das ganze Objekt, so wie wir es im AuthController gespeichert haben
        TwitchUser user = (TwitchUser) session.getAttribute("user");

        if (user == null) {
            return ResponseEntity
                    .status(HttpStatus.UNAUTHORIZED)
                    .body("Not logged in");
        }

        // Wir geben die Daten so zurück, dass dein React-Frontend (data.username) sie versteht
        return ResponseEntity.ok(
                Map.of(
                        "userId", user.id(),
                        "username", user.login(), // TwitchUser.login() ist meist der Anzeigename
                        "profileImageUrl", user.profileImageUrl()
                )
        );
    }

    @PostMapping("/logout")
    public ResponseEntity<?> logout(HttpSession session) {
        session.invalidate();

        return ResponseEntity.ok("Logged out");
    }
}
