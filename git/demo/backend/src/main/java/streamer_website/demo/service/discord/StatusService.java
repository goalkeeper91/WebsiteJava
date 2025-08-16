package streamer_website.demo.service.discord;

import lombok.Getter;
import org.springframework.stereotype.Service;

import java.time.Instant;

@Getter
@Service
public class StatusService {

    private boolean running = false;
    private Instant startTime;

    public void setRunning(boolean running) {
        this.running = running;
        if (running) {
            this.startTime = Instant.now();
        }
    }

}
