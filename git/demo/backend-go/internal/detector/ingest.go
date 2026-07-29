package detector

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
	"time"

	"demo/backend-go/internal/infrastructure/redis"
)

const (
	audioSampleRate = 16000
	audioChunkBytes = audioSampleRate * 2 // 1 Sekunde, 16-bit mono PCM

	videoWidth  = 160
	videoHeight = 90
	videoFPS    = 2
	frameBytes  = videoWidth * videoHeight * 3 // rgb24

	triggerCooldown   = 90 * time.Second
	combinedThreshold = 4.0 // Summe der Spike-Scores, die einen Trigger auslöst
	scoreDecayInterval = 2 * time.Second
	scoreDecayFactor   = 0.7 // Score verfällt mit der Zeit, wenn keine neuen Spikes reinkommen
)

// StreamIngestor beobachtet einen einzelnen Live-Stream (Audio + Video
// Heuristiken) und meldet über Redis, wenn ein clip-würdiger Moment erkannt wird.
type StreamIngestor struct {
	userTwitchID string
	channelLogin string
	redisService *redis.RedisService
	cancel       context.CancelFunc
}

// StartIngestor startet die Beobachtung eines Streams im Hintergrund.
func StartIngestor(ctx context.Context, userTwitchID, channelLogin string, redisService *redis.RedisService) *StreamIngestor {
	ctx, cancel := context.WithCancel(ctx)
	ing := &StreamIngestor{
		userTwitchID: userTwitchID,
		channelLogin: channelLogin,
		redisService: redisService,
		cancel:       cancel,
	}

	go ing.run(ctx)

	return ing
}

// Stop beendet die Beobachtung dieses Streams (Subprocesses werden über den
// Context abgebrochen).
func (ing *StreamIngestor) Stop() {
	ing.cancel()
}

// UserTwitchID und ChannelLogin geben Auskunft über den beobachteten Kanal -
// genutzt vom Clip-Detector, um seinen aktuellen Status für die Admin-Kunden-
// übersicht (Phase 5) als Redis-Snapshot zu veröffentlichen.
func (ing *StreamIngestor) UserTwitchID() string {
	return ing.userTwitchID
}

func (ing *StreamIngestor) ChannelLogin() string {
	return ing.channelLogin
}

func (ing *StreamIngestor) run(ctx context.Context) {
	log.Printf("🎬 Starte Ingestion für %s", ing.channelLogin)
	defer log.Printf("🛑 Ingestion beendet für %s", ing.channelLogin)

	streamURL, err := resolveStreamURL(ctx, ing.channelLogin)
	if err != nil {
		log.Printf("⚠️ Konnte Stream-URL für %s nicht auflösen: %v", ing.channelLogin, err)
		return
	}

	scoreCh := make(chan float64, 16)
	go ing.runAudio(ctx, streamURL, scoreCh)
	go ing.runVideo(ctx, streamURL, scoreCh)

	var lastTrigger time.Time
	combined := 0.0

	decay := time.NewTicker(scoreDecayInterval)
	defer decay.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case s := <-scoreCh:
			combined += s
			if combined >= combinedThreshold && time.Since(lastTrigger) > triggerCooldown {
				lastTrigger = time.Now()
				combined = 0
				if err := ing.redisService.SendClipTrigger(ing.userTwitchID, "audio_video_spike"); err != nil {
					log.Printf("❌ Fehler beim Senden des Clip-Triggers für %s: %v", ing.channelLogin, err)
				} else {
					log.Printf("🎯 Clip-Trigger ausgelöst für %s", ing.channelLogin)
				}
			}

		case <-decay.C:
			combined *= scoreDecayFactor
		}
	}
}

func (ing *StreamIngestor) runAudio(ctx context.Context, streamURL string, scoreCh chan<- float64) {
	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-i", streamURL,
		"-vn",
		"-f", "s16le",
		"-ar", fmt.Sprintf("%d", audioSampleRate),
		"-ac", "1",
		"-loglevel", "error",
		"pipe:1",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("⚠️ Audio-ffmpeg Setup fehlgeschlagen für %s: %v", ing.channelLogin, err)
		return
	}
	if err := cmd.Start(); err != nil {
		log.Printf("⚠️ Audio-ffmpeg Start fehlgeschlagen für %s: %v", ing.channelLogin, err)
		return
	}
	defer cmd.Wait()

	detector := NewAudioSpikeDetector()
	reader := bufio.NewReaderSize(stdout, audioChunkBytes*2)
	buf := make([]byte, audioChunkBytes)

	for {
		n, err := io.ReadFull(reader, buf)
		if n > 0 {
			if isSpike, score := detector.Observe(RMS(buf[:n])); isSpike {
				select {
				case scoreCh <- score:
				default:
				}
			}
		}
		if err != nil {
			return
		}
	}
}

func (ing *StreamIngestor) runVideo(ctx context.Context, streamURL string, scoreCh chan<- float64) {
	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-i", streamURL,
		"-an",
		"-vf", fmt.Sprintf("fps=%d,scale=%d:%d", videoFPS, videoWidth, videoHeight),
		"-f", "rawvideo",
		"-pix_fmt", "rgb24",
		"-loglevel", "error",
		"pipe:1",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("⚠️ Video-ffmpeg Setup fehlgeschlagen für %s: %v", ing.channelLogin, err)
		return
	}
	if err := cmd.Start(); err != nil {
		log.Printf("⚠️ Video-ffmpeg Start fehlgeschlagen für %s: %v", ing.channelLogin, err)
		return
	}
	defer cmd.Wait()

	detector := NewMotionSpikeDetector()
	reader := bufio.NewReaderSize(stdout, frameBytes*2)
	prev := make([]byte, frameBytes)
	curr := make([]byte, frameBytes)
	haveFrame := false

	for {
		n, err := io.ReadFull(reader, curr)
		if n == frameBytes {
			if haveFrame {
				if isSpike, score := detector.Observe(FrameDiff(prev, curr)); isSpike {
					select {
					case scoreCh <- score:
					default:
					}
				}
			}
			haveFrame = true
			prev, curr = curr, prev
		}
		if err != nil {
			return
		}
	}
}

// resolveStreamURL fragt streamlink nach der direkten HLS-Playlist-URL eines
// öffentlichen Twitch-Kanals — Twitch bietet dafür keine offizielle API.
func resolveStreamURL(ctx context.Context, channelLogin string) (string, error) {
	cmd := exec.CommandContext(ctx, "streamlink",
		fmt.Sprintf("https://twitch.tv/%s", channelLogin),
		"best",
		"--stream-url",
	)

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("streamlink fehlgeschlagen: %w", err)
	}

	url := strings.TrimSpace(string(out))
	if url == "" {
		return "", fmt.Errorf("streamlink lieferte keine URL")
	}

	return url, nil
}
