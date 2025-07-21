package streamer_website.demo.service;

import streamer_website.demo.entity.About;
import org.springframework.stereotype.Service;
import streamer_website.demo.repository.AboutRepository;

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

    public About saveAbout(About about) {
        return aboutRepository.save(about);
    }
}
