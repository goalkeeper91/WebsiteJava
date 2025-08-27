package streamer_website.demo.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import streamer_website.demo.entity.ContactRequest;

public interface ContactRequestRepository extends JpaRepository<ContactRequest, Long> {
}
