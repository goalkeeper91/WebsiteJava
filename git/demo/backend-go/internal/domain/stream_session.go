package domain

import "time"

// StreamSession is one continuous stream (detected via Twitch's started_at
// timestamp), used to compute a real average/peak viewer count instead of
// just reporting the latest live.ViewerCount. It holds a running sum rather
// than individual viewer-count samples - the average is just
// ViewerSum/SampleCount, no aggregation query over many rows needed.
type StreamSession struct {
	ID           int64      `json:"id"`
	TwitchUserID string     `json:"twitch_user_id"`
	StartedAt    time.Time  `json:"started_at"`
	EndedAt      *time.Time `json:"ended_at,omitempty"`
	ViewerSum    int64      `json:"viewer_sum"`
	SampleCount  int        `json:"sample_count"`
	PeakViewers  int        `json:"peak_viewers"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// AverageViewers returns 0 if no sample has been recorded yet (e.g. the
// sampler hasn't ticked since this session started).
func (s *StreamSession) AverageViewers() int {
	if s.SampleCount == 0 {
		return 0
	}
	return int(s.ViewerSum / int64(s.SampleCount))
}
