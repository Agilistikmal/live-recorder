package models

import (
	"fmt"
	"time"
)

type ShowroomResponses struct {
	OnLives []ShowroomOnLive `json:"onlives"`
}

type ShowroomStreamingUrlResponses struct {
	StreamingUrlList []ShowroomStreamingUrl `json:"streaming_url_list"`
}

type ShowroomOnLive struct {
	GenreName string         `json:"genre_name"`
	Lives     []ShowroomLive `json:"lives"`
}

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

type ShowroomStreamingUrl struct {
	Label string `json:"label"`
	Url   string `json:"url"`
}

func (s *ShowroomLive) ToLive() *Live {
	startedAt := time.Unix(int64(s.StartedAt), 0)
	return &Live{
		ID: s.RoomUrlKey,
		Streamer: &LiveStreamer{
			Username:      s.MainName,
			Name:          s.MainName,
			FollowerCount: s.FollowerNum,
			ImageUrl:      s.ImageSquare,
		},
		Title:        s.Telop,
		Platform:     "ShowroomLive",
		PlatformUrl:  fmt.Sprintf("https://showroom-live.com/r/%v", s.RoomUrlKey),
		StreamingUrl: s.StreamingUrlList[0].Url,
		ImageUrl:     s.ImageSquare,
		ViewCount:    s.ViewNum,
		StartedAt:    &startedAt,
	}
}
