package streamer_website.demo.service;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import streamer_website.demo.entity.Video;
import streamer_website.demo.repository.VideoRepository;

import java.util.List;

@Service
@RequiredArgsConstructor
public class VideoService {

    private final VideoRepository videoRepository;

    public List<Video> getAllVideos() {
        return videoRepository.findAll();
    }

    public List<Video> getVideosByPlatform(String platform) {
        return videoRepository.findByPlatform(platform.toLowerCase());
    }

    public Video createVideo(Video video) {
        return videoRepository.save(video);
    }

    public Video updateVideo(Long id, Video updatedVideo) {
        return videoRepository.findById(id)
                .map(existing -> {
                    existing.setPlatform(updatedVideo.getPlatform());
                    existing.setTitle(updatedVideo.getTitle());
                    existing.setUrlId(updatedVideo.getUrlId());
                    existing.setApiKey(updatedVideo.getApiKey());
                    existing.setChannelId(updatedVideo.getChannelId());
                    return videoRepository.save(existing);
                })
                .orElseThrow(() -> new RuntimeException("Video not found"));
    }

    public void deleteVideo(Long id) {
        videoRepository.deleteById(id);
    }
}
