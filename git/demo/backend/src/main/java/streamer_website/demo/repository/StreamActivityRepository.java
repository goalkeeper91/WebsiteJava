package streamer_website.demo.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import streamer_website.demo.entity.twitch.StreamActivity;

import java.time.Instant;
import java.util.List;

public interface StreamActivityRepository extends JpaRepository<StreamActivity, Long> {

    /**
     * Findet die letzten N Activities für einen User
     */
    @Query("""
        SELECT a FROM StreamActivity a\s
        WHERE a.twitchUserId = :userId\s
        ORDER BY a.timestamp DESC\s
        LIMIT :limit
   \s""")
    List<StreamActivity> findTopByTwitchUserIdOrderByTimestampDesc(
            @Param("userId") String twitchUserId,
            @Param("limit") int limit
    );

    /**
     * Findet Activities in einem Zeitraum
     */
    List<StreamActivity> findByTwitchUserIdAndTimestampBetweenOrderByTimestampDesc(
            String twitchUserId,
            Instant from,
            Instant to
    );

    /**
     * Zählt Activities nach Typ
     */
    @Query("""
        SELECT COUNT(a) FROM StreamActivity a\s
        WHERE a.twitchUserId = :userId\s
        AND a.type = :type\s
        AND a.timestamp > :since
   \s""")
    long countByTypeAndSince(
            @Param("userId") String twitchUserId,
            @Param("type") String type,
            @Param("since") Instant since
    );

    /**
     * Löscht alte Activities (Cleanup)
     */
    void deleteByTimestampBefore(Instant timestamp);
}
