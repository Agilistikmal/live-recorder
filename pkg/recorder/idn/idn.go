package idn

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/agilistikmal/live-recorder/pkg/recorder"
	"github.com/agilistikmal/live-recorder/pkg/recorder/models"
	"github.com/agilistikmal/live-recorder/utils"
)

type IDNRecorder struct {
	recorderConfig recorder.RecorderConfig
	httpClient     *http.Client
}

func NewRecorder() recorder.Recorder {
	recorderConfig := recorder.RecorderConfig{
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		Referer:   "https://www.idnlive.com/",
		Cookie:    "",
	}
	httpClient := &http.Client{}
	return &IDNRecorder{
		recorderConfig: recorderConfig,
		httpClient:     httpClient,
	}
}

func (s *IDNRecorder) GetLives() ([]*models.Live, error) {
	lives := make([]*models.Live, 0)

	page := 1
	for {
		query, err := json.Marshal(map[string]any{
			"query": fmt.Sprintf(`
				query GetLivestreams {
					getLivestreams(page: %v) {
						slug
						title
						image_url
						view_count
						playback_url
						status
						live_at
						gift_icon_url
						creator {
								username
								name
								follower_count
						}
					}
				}
				`, page),
		})
		if err != nil {
			return nil, err
		}

		gReq, err := http.NewRequest("POST", "https://api.idn.app/graphql", bytes.NewBuffer(query))
		if err != nil {
			return nil, err
		}
		gReq.Header.Set("Content-Type", "application/json")
		gReq.Header.Set("User-Agent", s.recorderConfig.UserAgent)
		gReq.Header.Set("Referer", s.recorderConfig.Referer)
		gReq.Header.Set("Cookie", s.recorderConfig.Cookie)
		resp, err := s.httpClient.Do(gReq)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var idnResponses models.IDNResponses
		err = json.Unmarshal(body, &idnResponses)
		if err != nil {
			return nil, err
		}

		if idnResponses.Data.GetLivestreams == nil {
			return nil, errors.New("idn response is nil")
		}

		if len(idnResponses.Data.GetLivestreams) == 0 {
			break
		}

		for _, l := range idnResponses.Data.GetLivestreams {
			if l.Status != "live" {
				continue
			}
			lives = append(lives, l.ToLive())
		}
		page++
	}
	return lives, nil
}

func (s *IDNRecorder) GetStreamingUrl(live *models.Live) (string, error) {
	return live.StreamingUrl, nil
}

func (s *IDNRecorder) Record(live *models.Live, outputPath string) error {
	downloadInfo := utils.DownloadHLS(live.StreamingUrl, &outputPath)
	if downloadInfo == nil {
		return fmt.Errorf("failed to download hls: %v", live.StreamingUrl)
	}
	return nil
}
