package service

import (
	"context"
	"fmt"
	"time"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
	"demo/backend-go/internal/twitch"
)

// StreamDashboardService backs the Live-Dashboard home page: stream
// info/title/category, live status, aggregate stats, category search, and
// ad breaks. Every method takes an already-resolved twitchUserID - team
// delegation is resolved by the handler (resolveEffectiveTwitchID), this
// service doesn't know about teams, same split as GiveawayService.
type StreamDashboardService struct {
	authService    *AuthService
	channelClient  *twitch.ChannelClient
	appTokenClient *twitch.TwitchAppTokenClient
	activityRepo   repository.StreamActivityRepository
}

func NewStreamDashboardService(
	authService *AuthService,
	channelClient *twitch.ChannelClient,
	appTokenClient *twitch.TwitchAppTokenClient,
	activityRepo repository.StreamActivityRepository,
) *StreamDashboardService {
	return &StreamDashboardService{
		authService:    authService,
		channelClient:  channelClient,
		appTokenClient: appTokenClient,
		activityRepo:   activityRepo,
	}
}

// StreamInfo mirrors the frontend's StreamInfo type.
type StreamInfo struct {
	BroadcasterID   string `json:"broadcasterId"`
	BroadcasterName string `json:"broadcasterName"`
	Title           string `json:"title"`
	GameID          string `json:"gameId"`
	GameName        string `json:"gameName"`
}

func (s *StreamDashboardService) GetStreamInfo(ctx context.Context, twitchUserID string) (*StreamInfo, error) {
	accessToken, err := s.authService.GetFreshAccessToken(ctx, twitchUserID)
	if err != nil {
		return nil, err
	}

	info, err := s.channelClient.GetChannelInfo(ctx, twitchUserID, accessToken)
	if err != nil {
		return nil, err
	}

	return &StreamInfo{
		BroadcasterID:   info.BroadcasterID,
		BroadcasterName: info.BroadcasterLogin,
		Title:           info.Title,
		GameID:          info.GameID,
		GameName:        info.GameName,
	}, nil
}

func (s *StreamDashboardService) UpdateStreamInfo(ctx context.Context, twitchUserID string, title, gameID *string) (*StreamInfo, error) {
	accessToken, err := s.authService.GetFreshAccessToken(ctx, twitchUserID)
	if err != nil {
		return nil, err
	}

	if err := s.channelClient.ModifyChannelInfo(ctx, twitchUserID, accessToken, title, gameID); err != nil {
		return nil, err
	}

	return s.GetStreamInfo(ctx, twitchUserID)
}

// LiveStatus mirrors the frontend's LiveStream type.
type LiveStatus struct {
	IsLive       bool   `json:"isLive"`
	ViewerCount  int    `json:"viewerCount"`
	StartedAt    string `json:"startedAt,omitempty"`
	Title        string `json:"title,omitempty"`
	GameName     string `json:"gameName,omitempty"`
	ThumbnailURL string `json:"thumbnailUrl,omitempty"`
}

func (s *StreamDashboardService) GetLiveStatus(ctx context.Context, twitchUserID string) (*LiveStatus, error) {
	streams, err := s.appTokenClient.GetLiveStreams(ctx, []string{twitchUserID})
	if err != nil {
		return nil, err
	}

	live, isLive := streams[twitchUserID]
	if !isLive {
		return &LiveStatus{IsLive: false}, nil
	}

	return &LiveStatus{
		IsLive:      true,
		ViewerCount: live.ViewerCount,
		StartedAt:   live.StartedAt,
		Title:       live.Title,
		GameName:    live.GameName,
	}, nil
}

// DashboardStats mirrors the frontend's DashboardStats type.
type DashboardStats struct {
	IsLive          bool   `json:"isLive"`
	CurrentViewers  int    `json:"currentViewers"`
	FollowerCount   int    `json:"followerCount"`
	SubscriberCount int    `json:"subscriberCount"`
	Uptime          string `json:"uptime,omitempty"`
	FollowsToday    int    `json:"followsToday"`
	SubsThisWeek    int    `json:"subsThisWeek"`
	BitsToday       int    `json:"bitsToday"`
	AvgViewers      int    `json:"avgViewers"`
}

func (s *StreamDashboardService) GetDashboardStats(ctx context.Context, twitchUserID string) (*DashboardStats, error) {
	accessToken, err := s.authService.GetFreshAccessToken(ctx, twitchUserID)
	if err != nil {
		return nil, err
	}

	live, err := s.GetLiveStatus(ctx, twitchUserID)
	if err != nil {
		return nil, err
	}

	followerCount, err := s.channelClient.GetFollowerCount(ctx, twitchUserID, accessToken)
	if err != nil {
		return nil, err
	}

	subscriberCount, err := s.channelClient.GetSubscriberCount(ctx, twitchUserID, accessToken)
	if err != nil {
		return nil, err
	}

	followsToday, err := s.countActivitiesSince(ctx, twitchUserID, domain.ActivityFollow, startOfToday())
	if err != nil {
		return nil, err
	}

	subsThisWeek, err := s.countSubsSince(ctx, twitchUserID, time.Now().AddDate(0, 0, -7))
	if err != nil {
		return nil, err
	}

	bitsToday, err := s.sumBitsSince(ctx, twitchUserID, startOfToday())
	if err != nil {
		return nil, err
	}

	uptime := ""
	if live.IsLive && live.StartedAt != "" {
		if startedAt, err := time.Parse(time.RFC3339, live.StartedAt); err == nil {
			uptime = time.Since(startedAt).Round(time.Minute).String()
		}
	}

	// avgViewers is deliberately stubbed to the current viewer count - a
	// real rolling average needs a periodic viewer-count sampler and a new
	// table, out of scope here (see plan's "Nicht Teil dieser Phase").
	return &DashboardStats{
		IsLive:          live.IsLive,
		CurrentViewers:  live.ViewerCount,
		FollowerCount:   followerCount,
		SubscriberCount: subscriberCount,
		Uptime:          uptime,
		FollowsToday:    followsToday,
		SubsThisWeek:    subsThisWeek,
		BitsToday:       bitsToday,
		AvgViewers:      live.ViewerCount,
	}, nil
}

const statsHistoryLimit = 500

func (s *StreamDashboardService) countActivitiesSince(ctx context.Context, twitchUserID string, activityType domain.ActivityType, since time.Time) (int, error) {
	activities, err := s.activityRepo.GetByType(ctx, twitchUserID, activityType, statsHistoryLimit)
	if err != nil {
		return 0, fmt.Errorf("fehler beim Zählen der Activities: %w", err)
	}

	count := 0
	for _, a := range activities {
		if !a.Timestamp.Before(since) {
			count++
		}
	}
	return count, nil
}

func (s *StreamDashboardService) countSubsSince(ctx context.Context, twitchUserID string, since time.Time) (int, error) {
	count := 0
	for _, t := range []domain.ActivityType{domain.ActivitySubscribe, domain.ActivityGiftSub, domain.ActivityResubscribe} {
		n, err := s.countActivitiesSince(ctx, twitchUserID, t, since)
		if err != nil {
			return 0, err
		}
		count += n
	}
	return count, nil
}

func (s *StreamDashboardService) sumBitsSince(ctx context.Context, twitchUserID string, since time.Time) (int, error) {
	activities, err := s.activityRepo.GetByType(ctx, twitchUserID, domain.ActivityCheer, statsHistoryLimit)
	if err != nil {
		return 0, fmt.Errorf("fehler beim Summieren der Bits: %w", err)
	}

	total := 0
	for _, a := range activities {
		if a.Bits != nil && !a.Timestamp.Before(since) {
			total += *a.Bits
		}
	}
	return total, nil
}

func startOfToday() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

func (s *StreamDashboardService) SearchCategories(ctx context.Context, query string) ([]twitch.Category, error) {
	return s.appTokenClient.SearchCategories(ctx, query)
}

func (s *StreamDashboardService) StartCommercial(ctx context.Context, twitchUserID string, lengthSeconds int) (*twitch.CommercialResult, error) {
	accessToken, err := s.authService.GetFreshAccessToken(ctx, twitchUserID)
	if err != nil {
		return nil, err
	}

	return s.channelClient.StartCommercial(ctx, twitchUserID, accessToken, lengthSeconds)
}
