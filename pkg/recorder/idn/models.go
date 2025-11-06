package idn

import (
	"time"

	"github.com/agilistikmal/live-recorder/pkg/recorder"
)

// IDNResponses represents the IDN API response structure
type IDNResponses struct {
	Data struct {
		GetLivestreams []IDNLive `json:"getLivestreams"`
	} `json:"data"`
}

// IDNLive represents a live stream from IDN
type IDNLive struct {
	Slug        string      `json:"slug"`
	Title       string      `json:"title"`
	ImageUrl    string      `json:"image_url"`
	ViewCount   int         `json:"view_count"`
	PlaybackUrl string      `json:"playback_url"`
	Status      string      `json:"status"`
	LiveAt      *time.Time  `json:"live_at"`
	Creator     *IDNCreator `json:"creator"`
}

// IDNCreator represents creator information from IDN
type IDNCreator struct {
	Username      string `json:"username"`
	Name          string `json:"name"`
	FollowerCount int    `json:"follower_count"`
}

// ToLive converts IDNLive to recorder.Live
func (i *IDNLive) ToLive() *recorder.Live {
	return &recorder.Live{
		ID: i.Slug,
		Streamer: &recorder.LiveStreamer{
			Username:      i.Creator.Username,
			Name:          i.Creator.Name,
			FollowerCount: i.Creator.FollowerCount,
		},
		Title:        i.Title,
		Platform:     recorder.PlatformIDN,
		PlatformUrl:  i.PlaybackUrl,
		StreamingUrl: i.PlaybackUrl,
		ImageUrl:     i.ImageUrl,
		ViewCount:    i.ViewCount,
		StartedAt:    i.LiveAt,
	}
}
