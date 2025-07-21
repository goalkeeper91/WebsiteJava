package service;

import entity.About;
import org.junit.jupiter.api.Test;
import repository.AboutRepository;

import java.util.Optional;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

public class AboutServiceTest {

    @Test
    void shouldReturnAboutText() {
        AboutRepository mockRepo = mock(AboutRepository.class);
        AboutService service = new AboutService(mockRepo);

        About mockAbout = new About(1L, "Über mich", "Ich bin Java-Entwickler & Streamer", "url.jpg");
        when(mockRepo.findById(1L)).thenReturn(Optional.of(mockAbout));

        About result = service.getAbout(1L);

        assertEquals("Über mich", result.getHeadline());
        assertEquals("Ich bin Java-Entwickler & Streamer", result.getDescription());
        assertEquals("url.jpg", result.getProfileImageUrl());
    }
}
