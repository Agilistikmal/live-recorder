package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/agilistikmal/live-recorder/models"
)

type ShowroomLiveService struct {
}

func NewShowroomLiveService() *ShowroomLiveService {
	return &ShowroomLiveService{}
}

func (s *ShowroomLiveService) GetLives() ([]*models.Live, error) {
	resp, err := http.Get("https://www.showroom-live.com/api/live/onlives")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var showroomResponses models.ShowroomResponses
	err = json.NewDecoder(resp.Body).Decode(&showroomResponses)
	if err != nil {
		return nil, err
	}

	liveList := make([]*models.Live, 0)
	for _, onLive := range showroomResponses.OnLives {
		for _, showroomLive := range onLive.Lives {
			liveList = append(liveList, showroomLive.ToLive())
		}
	}

	return liveList, nil
}

func (s *ShowroomLiveService) GetStreamingUrl(roomId string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://www.showroom-live.com/api/live/streaming_url?abr_available=1&room_id=%v", roomId))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var srStreamingUrlResponses models.ShowroomStreamingUrlResponses
	err = json.NewDecoder(resp.Body).Decode(&srStreamingUrlResponses)
	if err != nil {
		return "", err
	}

	if len(srStreamingUrlResponses.StreamingUrlList) < 1 {
		return "", errors.New("streaming url not found")
	}

	return srStreamingUrlResponses.StreamingUrlList[1].Url, nil
}
