package streamer_website.demo.controller;

import jakarta.validation.Valid;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import streamer_website.demo.dto.ContactRequestDTO;
import streamer_website.demo.entity.ContactRequest;
import streamer_website.demo.service.ContactService;

@RestController
@RequestMapping("/api/contact")
public class ContactController {

    private final ContactService service;

    public ContactController(ContactService service) {
        this.service = service;
    }

    @PostMapping
    public ResponseEntity<ContactRequest> submitContact(@Valid @RequestBody ContactRequestDTO dto) {
        ContactRequest saved = service.saveRequest(dto);
        return ResponseEntity.ok(saved);
    }
}
