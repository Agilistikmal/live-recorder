package models

import "time"

const (
	PlatformShowroom = "showroom"
	PlatformIDN      = "idn"
)

type Live struct {
	ID           string        `json:"id"`
	Streamer     *LiveStreamer `json:"streamer"`
	Title        string        `json:"title"`
	Platform     string        `json:"platform"`
	PlatformUrl  string        `json:"platform_url"`
	StreamingUrl string        `json:"streaming_url"`
	ImageUrl     string        `json:"image_url"`
	ViewCount    int           `json:"view_count"`
	StartedAt    *time.Time    `json:"started_at"`
}

type LiveStreamer struct {
	Username      string `json:"username"`
	Name          string `json:"name"`
	FollowerCount int    `json:"follower_count"`
	ImageUrl      string `json:"image_url"`
}

type LiveQuery struct {
	Platform             string `json:"platform"`
	StreamerUsernameLike string `json:"streamer_username_like"`
	TitleLike            string `json:"title_like"`
}
