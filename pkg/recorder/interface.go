package recorder

import "time"

// Platform constants
const (
	PlatformShowroom = "showroom"
	PlatformIDN      = "idn"
	PlatformTiktok   = "tiktok"
)

// RecorderConfig holds configuration for recorders
type RecorderConfig struct {
	UserAgent string `json:"user_agent"`
	Referer   string `json:"referer"`
	Cookie    string `json:"cookie"`
}

// Live represents a live streaming session
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

// LiveStreamer represents information about a streamer
type LiveStreamer struct {
	Username      string `json:"username"`
	Name          string `json:"name"`
	FollowerCount int    `json:"follower_count"`
	ImageUrl      string `json:"image_url"`
}

// LiveQuery represents query parameters for filtering live streams
type LiveQuery struct {
	Platforms            []string `json:"platforms"`
	StreamerUsernameLike string   `json:"streamer_username_like"`
	TitleLike            string   `json:"title_like"`
}

// Recorder is the interface for recording live streams
type Recorder interface {
	GetLives() ([]*Live, error)
	GetLive(url string) (*Live, error)
	GetStreamingUrl(live *Live) (string, error)
	Record(live *Live, outputPath string) error
}
