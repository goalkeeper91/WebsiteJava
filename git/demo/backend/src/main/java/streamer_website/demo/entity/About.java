package streamer_website.demo.entity;

import jakarta.persistence.Entity;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;

@Entity
public class About {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    private String headline;
    private String description;
    private String profileImageUrl;

    // Konstruktoren

    public About() {
    }

    public About(Long id, String headline, String description, String profileImageUrl) {
        this.id = id;
        this.headline = headline;
        this.description = description;
        this.profileImageUrl = profileImageUrl;
    }

    // Getter & Setter

    public Long getId() {
        return id;
    }

    public String getHeadline() {
        return headline;
    }

    public void setHeadline(String headline) {
        this.headline = headline;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public String getProfileImageUrl() {
        return profileImageUrl;
    }

    public void setProfileImageUrl(String profileImageUrl) {
        this.profileImageUrl = profileImageUrl;
    }
}

