package models

import "time"

type IDNResponses struct {
	Data struct {
		GetLivestreams []IDNLive `json:"getLivestreams"`
	} `json:"data"`
}

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

type IDNCreator struct {
	Username      string `json:"username"`
	Name          string `json:"name"`
	FollowerCount int    `json:"follower_count"`
}

func (i *IDNLive) ToLive() *Live {
	return &Live{
		ID: i.Slug,
		Streamer: &LiveStreamer{
			Username:      i.Creator.Username,
			Name:          i.Creator.Name,
			FollowerCount: i.Creator.FollowerCount,
		},
		Title:        i.Title,
		Platform:     "IDNLive",
		PlatformUrl:  i.PlaybackUrl,
		StreamingUrl: i.PlaybackUrl,
		ImageUrl:     i.ImageUrl,
		ViewCount:    i.ViewCount,
		StartedAt:    i.LiveAt,
	}
}
