package streamer_website.demo.repository;

import streamer_website.demo.entity.About;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import streamer_website.demo.repository.AboutRepository;

import java.util.Optional;

import static org.assertj.core.api.Assertions.assertThat;

@SpringBootTest(classes = streamer_website.demo.DemoApplication.class)
public class AboutRepositoryTest {

    @Autowired
    private AboutRepository aboutRepository;

    @Test
    void shouldSaveAndFindAbout() {
        About about = new About();
        about.setHeadline("Test Headline");
        about.setDescription("Test Description");
        about.setProfileImageUrl("http://example.com/image.jpg");

        About saved = aboutRepository.save(about);
        Optional<About> found = aboutRepository.findById(saved.getId());

        assertThat(found).isPresent();
        assertThat(found.get().getHeadline()).isEqualTo("Test Headline");
        assertThat(found.get().getDescription()).isEqualTo("Test Description");
        assertThat(found.get().getProfileImageUrl()).isEqualTo("http://example.com/image.jpg");
    }
}
