package showroom

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/agilistikmal/live-recorder/pkg/recorder"
	"github.com/agilistikmal/live-recorder/utils"
)

type ShowroomRecorder struct {
	recorderConfig recorder.RecorderConfig
	httpClient     *http.Client
}

func NewRecorder() recorder.Recorder {
	recorderConfig := recorder.RecorderConfig{
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		Referer:   "https://www.showroom-live.com/",
		Cookie:    "",
	}
	httpClient := &http.Client{}

	return &ShowroomRecorder{
		recorderConfig: recorderConfig,
		httpClient:     httpClient,
	}
}

func (s *ShowroomRecorder) GetLives() ([]*recorder.Live, error) {
	req, err := http.NewRequest("GET", "https://www.showroom-live.com/api/live/onlives", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", s.recorderConfig.UserAgent)
	req.Header.Set("Referer", s.recorderConfig.Referer)
	req.Header.Set("Cookie", s.recorderConfig.Cookie)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var showroomResponses ShowroomResponses
	err = json.NewDecoder(resp.Body).Decode(&showroomResponses)
	if err != nil {
		return nil, err
	}

	liveList := make([]*recorder.Live, 0)
	for _, onLive := range showroomResponses.OnLives {
		for _, showroomLive := range onLive.Lives {
			liveList = append(liveList, showroomLive.ToLive())
		}
	}

	return liveList, nil
}

func (s *ShowroomRecorder) GetLive(url string) (*recorder.Live, error) {
	return nil, nil
}

func (s *ShowroomRecorder) GetStreamingUrl(live *recorder.Live) (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.showroom-live.com/api/live/streaming_url?abr_available=1&room_id=%v", live.ID), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", s.recorderConfig.UserAgent)
	req.Header.Set("Referer", s.recorderConfig.Referer)
	req.Header.Set("Cookie", s.recorderConfig.Cookie)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var srStreamingUrlResponses ShowroomStreamingUrlResponses
	err = json.NewDecoder(resp.Body).Decode(&srStreamingUrlResponses)
	if err != nil {
		return "", err
	}

	if len(srStreamingUrlResponses.StreamingUrlList) < 1 {
		return "", fmt.Errorf("showroom streaming url not found for room id: %s", live.ID)
	}

	return srStreamingUrlResponses.StreamingUrlList[1].Url, nil
}

func (s *ShowroomRecorder) Record(live *recorder.Live, outputPath string) error {
	downloadInfo := utils.DownloadHLS(live.StreamingUrl, &outputPath)
	if downloadInfo == nil {
		return fmt.Errorf("failed to download hls: %v", live.StreamingUrl)
	}
	return nil
}
