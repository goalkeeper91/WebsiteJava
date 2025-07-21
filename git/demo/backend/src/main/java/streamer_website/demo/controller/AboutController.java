package streamer_website.demo.controller;

import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.*;
import streamer_website.demo.entity.About;
import streamer_website.demo.service.AboutService;

@RestController
@RequestMapping("/api/about")
public class AboutController {

    private final AboutService aboutService;

    public AboutController(AboutService aboutService) {
        this.aboutService = aboutService;
    }

    @GetMapping("/{id}")
    public About getAbout(@PathVariable Long id) {
        return aboutService.getAbout(id);
    }

    @PostMapping
    @ResponseStatus(HttpStatus.CREATED)
    public About createAbout(@RequestBody About about) {
        About saved = aboutService.saveAbout(about);
        System.out.println("Saved About with ID: " + saved.getId());
        return saved;
    }
}

