package service;

import entity.About;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import repository.AboutRepository;

import java.util.Optional;

@Service
public class AboutService {

    private final AboutRepository aboutRepository;

    public AboutService(AboutRepository aboutRepository) {
        this.aboutRepository = aboutRepository;
    }

    public About getAbout(Long id) {
        return aboutRepository.findById(id)
                .orElseThrow(() -> new IllegalArgumentException("About not found with id: " + id));
    }
}
