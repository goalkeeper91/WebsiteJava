package controller;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest;
import org.springframework.http.MediaType;
import org.springframework.test.context.ContextConfiguration;
import org.springframework.test.context.bean.override.mockito.MockitoBean;
import org.springframework.test.web.servlet.MockMvc;
import streamer_website.demo.controller.AboutController;
import streamer_website.demo.entity.About;
import streamer_website.demo.service.AboutService;

import static org.mockito.ArgumentMatchers.any;
import static org.mockito.BDDMockito.given;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath;


@WebMvcTest(AboutController.class)
@ContextConfiguration(classes = streamer_website.demo.DemoApplication.class)
public class AboutControllerTest {

    @Autowired
    private MockMvc mockMvc;

    @MockitoBean
    private AboutService aboutService;

    @Test
    void shouldReturnAboutInfo() throws Exception {
        About about = new About(1L, "My App", "This is a cool app.", "http://image.url/profile.png");
        given(aboutService.getAbout(1L)).willReturn(about);

        mockMvc.perform(get("/api/about/1"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.headline").value("My App"))
                .andExpect(jsonPath("$.description").value("This is a cool app."))
                .andExpect(jsonPath("$.profileImageUrl").value("http://image.url/profile.png"));
    }

    @Test
    void shouldCreateNewAbout() throws Exception {
        About savedAbout = new About(1L, "New Headline", "New Description", "http://image.url");

        given(aboutService.saveAbout(any(About.class))).willReturn(savedAbout);

        mockMvc.perform(post("/api/about")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content("""
                {
                    "headline": "New Headline",
                    "description": "New Description",
                    "profileImageUrl": "http://image.url"
                }
            """))
                .andExpect(status().isCreated())
                .andExpect(jsonPath("$.id").value(1L))
                .andExpect(jsonPath("$.headline").value("New Headline"))
                .andExpect(jsonPath("$.description").value("New Description"))
                .andExpect(jsonPath("$.profileImageUrl").value("http://image.url"));
    }

}
