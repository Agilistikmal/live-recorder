package showroom

import (
	"fmt"
	"time"

	"github.com/agilistikmal/live-recorder/pkg/recorder"
)

// ShowroomResponses represents the Showroom API response structure
type ShowroomResponses struct {
	OnLives []ShowroomOnLive `json:"onlives"`
}

// ShowroomStreamingUrlResponses represents the Showroom streaming URL response
type ShowroomStreamingUrlResponses struct {
	StreamingUrlList []ShowroomStreamingUrl `json:"streaming_url_list"`
}

// ShowroomOnLive represents on-live information from Showroom
type ShowroomOnLive struct {
	GenreName string         `json:"genre_name"`
	Lives     []ShowroomLive `json:"lives"`
}

// ShowroomLive represents a live stream from Showroom
type ShowroomLive struct {
	RoomUrlKey       string                 `json:"room_url_key"`
	Telop            string                 `json:"telop"`
	FollowerNum      int                    `json:"follower_num"`
	StartedAt        int                    `json:"started_at"`
	ImageSquare      string                 `json:"image_square"`
	ViewNum          int                    `json:"view_num"`
	MainName         string                 `json:"main_name"`
	PremiumRoomType  int                    `json:"premium_room_type"`
	RoomID           int                    `json:"room_id"`
	StreamingUrlList []ShowroomStreamingUrl `json:"streaming_url_list"`
}

// ShowroomStreamingUrl represents streaming URL information
type ShowroomStreamingUrl struct {
	Label string `json:"label"`
	Url   string `json:"url"`
}

// ToLive converts ShowroomLive to recorder.Live
func (s *ShowroomLive) ToLive() *recorder.Live {
	startedAt := time.Unix(int64(s.StartedAt), 0)
	return &recorder.Live{
		ID: fmt.Sprintf("%v", s.RoomID),
		Streamer: &recorder.LiveStreamer{
			Username:      s.RoomUrlKey,
			Name:          s.MainName,
			FollowerCount: s.FollowerNum,
			ImageUrl:      s.ImageSquare,
		},
		Title:        s.Telop,
		Platform:     recorder.PlatformShowroom,
		PlatformUrl:  fmt.Sprintf("https://showroom-live.com/r/%v", s.RoomUrlKey),
		StreamingUrl: s.StreamingUrlList[0].Url,
		ImageUrl:     s.ImageSquare,
		ViewCount:    s.ViewNum,
		StartedAt:    &startedAt,
	}
}
