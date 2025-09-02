package streamer_website.demo.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import streamer_website.demo.entity.ContactRequest;

import java.time.LocalDateTime;

public interface ContactRequestRepository extends JpaRepository<ContactRequest, Long> {

    void deleteByCreatedAtBefore(LocalDateTime cutoff);

}
