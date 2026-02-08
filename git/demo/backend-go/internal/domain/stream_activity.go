package domain

import (
	"time"
)

type ActivityType string

const (
	ActivityFollow      ActivityType = "FOLLOW"
	ActivitySubscribe   ActivityType = "SUBSCRIBE"
	ActivityRaid        ActivityType = "RAID"
	ActivityCheer       ActivityType = "CHEER"
	ActivityGiftSub     ActivityType = "GIFT_SUB"
	ActivityHosting     ActivityType = "HOSTING"
	ActivityResubscribe ActivityType = "RESUBSCRIBE"
)

type StreamActivity struct {
	ID           int64        `json:"id"`
	TwitchUserID string       `json:"twitch_user_id"`
	Type         ActivityType `json:"type"`
	Username     string       `json:"username"`
	DisplayName  string       `json:"display_name"`
	Viewers      *int         `json:"viewers,omitempty"`
	Bits         *int         `json:"bits,omitempty"`
	Tier         *string      `json:"tier,omitempty"`
	Message      *string      `json:"message,omitempty"`
	Timestamp    time.Time    `json:"timestamp"`
}

func NewStreamActivity(
	twitchUserID string,
	activityType ActivityType,
	username, displayName string,
) *StreamActivity {
	return &StreamActivity{
		TwitchUserID: twitchUserID,
		Type:         activityType,
		Username:     username,
		DisplayName:  displayName,
		Timestamp:    time.Now(),
	}
}

func (a *StreamActivity) SetViewers(viewers int) {
	a.Viewers = &viewers
}

func (a *StreamActivity) SetBits(bits int) {
	a.Bits = &bits
}

func (a *StreamActivity) SetTier(tier string) {
	a.Tier = &tier
}

func (a *StreamActivity) SetMessage(message string) {
	a.Message = &message
}