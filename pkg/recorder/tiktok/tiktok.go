package tiktok

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/agilistikmal/live-recorder/pkg/recorder"
	"github.com/agilistikmal/live-recorder/utils"
)

type TiktokRecorder struct {
	recorderConfig recorder.RecorderConfig
	httpClient     *http.Client
}

func NewRecorder() recorder.Recorder {
	recorderConfig := recorder.RecorderConfig{
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		Referer:   "https://www.tiktok.com/",
		Cookie:    "1%7Cz7FKki38aKyy7i-BC9rEDwcrVvjcLcFEL6QIeqldoy4%7C1761302831%7C6c1461e9f1f980cbe0404c5190",
	}
	httpClient := &http.Client{}
	return &TiktokRecorder{
		recorderConfig: recorderConfig,
		httpClient:     httpClient,
	}
}

func (s *TiktokRecorder) GetLives() ([]*recorder.Live, error) {
	lives := make([]*recorder.Live, 0)

	return lives, nil
}

func (s *TiktokRecorder) GetLive(url string) (*recorder.Live, error) {
	req, err := http.NewRequest("GET", url, nil)
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`<script id="SIGI_STATE" type="application/json">(.*?)</script>`)
	matches := re.FindStringSubmatch(string(body))
	if len(matches) < 2 {
		return nil, fmt.Errorf("SIGI_STATE not found in response")
	}

	jsonData := matches[1]
	var tiktokResponses TiktokResponses
	err = json.Unmarshal([]byte(jsonData), &tiktokResponses)
	if err != nil {
		return nil, err
	}

	liveRoom := tiktokResponses.LiveRoom.LiveRoomUserInfo
	user := liveRoom.User
	startedAt := time.Unix(liveRoom.LiveRoom.StartTime, 0)

	if user.Status != 2 {
		return nil, fmt.Errorf("user is not live")
	}

	live := &recorder.Live{
		ID:          liveRoom.LiveRoom.StreamId,
		Title:       liveRoom.LiveRoom.Title,
		Platform:    recorder.PlatformTiktok,
		PlatformUrl: fmt.Sprintf("https://www.tiktok.com/@%s/live", user.UniqueId),
		ImageUrl:    liveRoom.LiveRoom.CoverUrl,
		ViewCount:   liveRoom.LiveRoom.LiveRoomStats.UserCount,
		StartedAt:   &startedAt,
		Streamer: &recorder.LiveStreamer{
			Username: user.UniqueId,
			Name:     user.Nickname,
			ImageUrl: user.AvatarLarger,
		},
	}

	streamDataStr := liveRoom.LiveRoom.StreamData.PullData.StreamData
	m3u8UrlList, err := getVideoQualityUrl(streamDataStr, "hls")
	if err != nil {
		return nil, err
	}
	live.StreamingUrl = m3u8UrlList[0].URL

	return live, nil
}

func (s *TiktokRecorder) GetStreamingUrl(live *recorder.Live) (string, error) {
	return live.StreamingUrl, nil
}

func (s *TiktokRecorder) Record(live *recorder.Live, outputPath string) error {
	downloadInfo := utils.DownloadHLS(live.StreamingUrl, &outputPath)
	if downloadInfo == nil {
		return fmt.Errorf("failed to download hls: %v", live.StreamingUrl)
	}
	return nil
}

// qKey: hls, flv
func getVideoQualityUrl(streamDataStr string, qKey string) ([]TiktokVideoQualityInfo, error) {
	var streamDataMap map[string]any
	if err := json.Unmarshal([]byte(streamDataStr), &streamDataMap); err != nil {
		return nil, fmt.Errorf("failed to parse stream_data: %w", err)
	}

	dataSection, ok := streamDataMap["data"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("data section not found in stream_data")
	}

	playList := make([]TiktokVideoQualityInfo, 0)

	for _, value := range dataSection {
		qualityData, ok := value.(map[string]any)
		if !ok {
			continue
		}

		mainData, ok := qualityData["main"].(map[string]any)
		if !ok {
			continue
		}

		sdkParamsStr, ok := mainData["sdk_params"].(string)
		if !ok {
			continue
		}

		var sdkParams map[string]any
		if err := json.Unmarshal([]byte(sdkParamsStr), &sdkParams); err != nil {
			continue
		}

		vbitrate := 0
		if vbitrateVal, ok := sdkParams["vbitrate"].(float64); ok {
			vbitrate = int(vbitrateVal)
		}

		vCodec := ""
		if vCodecVal, ok := sdkParams["VCodec"].(string); ok {
			vCodec = vCodecVal
		}

		// Get the URL for the requested quality key (flv or hls)
		playUrl := ""
		if urlVal, exists := mainData[qKey]; exists {
			if urlStr, ok := urlVal.(string); ok && urlStr != "" {
				if strings.HasSuffix(urlStr, ".flv") || strings.HasSuffix(urlStr, ".m3u8") {
					playUrl = urlStr + "?codec=" + vCodec
				} else {
					playUrl = urlStr + "&codec=" + vCodec
				}
			}
		}

		resolution, ok := sdkParams["resolution"].(string)
		if vbitrate != 0 && ok && resolution != "" {
			parts := strings.Split(resolution, "x")
			if len(parts) == 2 {
				width, err1 := strconv.Atoi(parts[0])
				height, err2 := strconv.Atoi(parts[1])
				if err1 == nil && err2 == nil {
					playList = append(playList, TiktokVideoQualityInfo{
						URL:        playUrl,
						VBitrate:   vbitrate,
						Resolution: [2]int{width, height},
					})
				}
			}
		}
	}

	sort.Slice(playList, func(i, j int) bool {
		if playList[i].VBitrate != playList[j].VBitrate {
			return playList[i].VBitrate > playList[j].VBitrate
		}
		if playList[i].Resolution[0] != playList[j].Resolution[0] {
			return playList[i].Resolution[0] > playList[j].Resolution[0]
		}
		return playList[i].Resolution[1] > playList[j].Resolution[1]
	})

	return playList, nil
}
