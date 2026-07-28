package service

import (
	"context"
	"log"
	"time"

	"demo/backend-go/internal/repository"
	"demo/backend-go/internal/twitch"
)

// maxLiveStreamsPerCall mirrors Twitch's Get Streams limit of 100 user_id
// query params per call.
const maxLiveStreamsPerCall = 100

// StreamSessionService tracks per-stream viewer-count sessions so the
// dashboard can show a real average instead of just the latest reading.
type StreamSessionService struct {
	sessionRepo    repository.StreamSessionRepository
	appTokenClient *twitch.TwitchAppTokenClient
	userRepo       repository.UserRepository
}

func NewStreamSessionService(
	sessionRepo repository.StreamSessionRepository,
	appTokenClient *twitch.TwitchAppTokenClient,
	userRepo repository.UserRepository,
) *StreamSessionService {
	return &StreamSessionService{
		sessionRepo:    sessionRepo,
		appTokenClient: appTokenClient,
		userRepo:       userRepo,
	}
}

// RunSampleTick polls every channel's live status in one batch and records
// one more viewer-count sample for each currently live channel's open
// session - creating a new session when a stream just started (detected via
// a changed started_at), and closing a channel's open session once it's no
// longer live.
func (s *StreamSessionService) RunSampleTick(ctx context.Context) {
	users, err := s.userRepo.GetAllNonBotUsers(ctx)
	if err != nil {
		log.Printf("Stream-Session-Sampler: Fehler beim Laden der Nutzer: %v", err)
		return
	}
	if len(users) == 0 {
		return
	}

	ids := make([]string, len(users))
	for i, u := range users {
		ids[i] = u.TwitchID
	}

	liveStreams := make(map[string]twitch.LiveStream)
	for i := 0; i < len(ids); i += maxLiveStreamsPerCall {
		end := i + maxLiveStreamsPerCall
		if end > len(ids) {
			end = len(ids)
		}
		chunk, err := s.appTokenClient.GetLiveStreams(ctx, ids[i:end])
		if err != nil {
			log.Printf("Stream-Session-Sampler: Fehler beim Abrufen der Live-Streams: %v", err)
			continue
		}
		for id, stream := range chunk {
			liveStreams[id] = stream
		}
	}

	openSessions, err := s.sessionRepo.GetOpenSessionsByUsers(ctx, ids)
	if err != nil {
		log.Printf("Stream-Session-Sampler: Fehler beim Laden offener Sessions: %v", err)
		return
	}

	for _, id := range ids {
		live, isLive := liveStreams[id]
		open := openSessions[id]

		if !isLive {
			if open != nil {
				if err := s.sessionRepo.CloseSession(ctx, open.ID); err != nil {
					log.Printf("Stream-Session-Sampler: Fehler beim Schließen der Session fuer %s: %v", id, err)
				}
			}
			continue
		}

		startedAt, err := time.Parse(time.RFC3339, live.StartedAt)
		if err != nil {
			log.Printf("Stream-Session-Sampler: ungueltiges started_at fuer %s: %v", id, err)
			continue
		}

		if open != nil && !open.StartedAt.Equal(startedAt) {
			// A different started_at means a new stream began since the
			// last tick (the sampler missed the actual offline moment, e.g.
			// a short disconnect) - close the stale session before opening
			// a fresh one so its average isn't polluted across streams.
			if err := s.sessionRepo.CloseSession(ctx, open.ID); err != nil {
				log.Printf("Stream-Session-Sampler: Fehler beim Schliessen der alten Session fuer %s: %v", id, err)
			}
			open = nil
		}

		if open == nil {
			created, err := s.sessionRepo.CreateSession(ctx, id, startedAt)
			if err != nil {
				log.Printf("Stream-Session-Sampler: Fehler beim Anlegen der Session fuer %s: %v", id, err)
				continue
			}
			open = created
		}

		if err := s.sessionRepo.RecordSample(ctx, open.ID, live.ViewerCount); err != nil {
			log.Printf("Stream-Session-Sampler: Fehler beim Speichern des Samples fuer %s: %v", id, err)
		}
	}
}

// GetAverageViewers returns the running average for the channel's currently
// open session while live, or its last completed session's average while
// offline. Returns 0 if no matching session exists yet (e.g. the sampler
// hasn't ticked since the stream started, or the channel has never gone
// live since this feature was deployed).
func (s *StreamSessionService) GetAverageViewers(ctx context.Context, twitchUserID string, isLive bool) (int, error) {
	if isLive {
		sessions, err := s.sessionRepo.GetOpenSessionsByUsers(ctx, []string{twitchUserID})
		if err != nil {
			return 0, err
		}
		open := sessions[twitchUserID]
		if open == nil {
			return 0, nil
		}
		return open.AverageViewers(), nil
	}

	last, err := s.sessionRepo.GetLastClosedSession(ctx, twitchUserID)
	if err != nil {
		return 0, err
	}
	if last == nil {
		return 0, nil
	}
	return last.AverageViewers(), nil
}
