package streamer_website.demo.service;

import org.springframework.stereotype.Service;
import streamer_website.demo.dto.ContactRequestDTO;
import streamer_website.demo.entity.ContactRequest;
import streamer_website.demo.repository.ContactRequestRepository;
import streamer_website.demo.service.discord.DiscordNotificationService;

@Service
public class ContactService {
    private final ContactRequestRepository repository;
    private final DiscordNotificationService discordService;

    public ContactService(ContactRequestRepository repository,
                          DiscordNotificationService discordService) {
        this.repository = repository;
        this.discordService = discordService;
    }

    public ContactRequest saveRequest(ContactRequestDTO dto) {
        ContactRequest request = new ContactRequest();
        request.setName(dto.getName());
        request.setEmail(dto.getEmail());
        request.setPhone(dto.getPhone());
        request.setSubject(dto.getSubject());
        request.setMessage(dto.getMessage());

        ContactRequest saved = repository.save(request);

        discordService.notifyNewContactRequest(saved).subscribe();

        return saved;
    }
}
