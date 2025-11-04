package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/agilistikmal/live-recorder/configs"
	"github.com/agilistikmal/live-recorder/models"
)

type ShowroomLiveService struct {
	requestConfig configs.RequestConfig
	httpClient    *http.Client
}

func NewShowroomLiveService() *ShowroomLiveService {
	requestConfig := configs.NewRequestConfig(
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"https://www.showroom-live.com/",
		"",
	)
	httpClient := &http.Client{}

	return &ShowroomLiveService{
		requestConfig: requestConfig,
		httpClient:    httpClient,
	}
}

func (s *ShowroomLiveService) GetLives() ([]*models.Live, error) {
	req, err := http.NewRequest("GET", "https://www.showroom-live.com/api/live/onlives", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", s.requestConfig.UserAgent)
	req.Header.Set("Referer", s.requestConfig.Referer)
	req.Header.Set("Cookie", s.requestConfig.Cookie)

	resp, err := s.httpClient.Do(req)
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
	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.showroom-live.com/api/live/streaming_url?abr_available=1&room_id=%v", roomId), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", s.requestConfig.UserAgent)
	req.Header.Set("Referer", s.requestConfig.Referer)
	req.Header.Set("Cookie", s.requestConfig.Cookie)

	resp, err := s.httpClient.Do(req)
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
