package watch

import (
	"fmt"
	"sync"
	"time"

	"math/rand"

	"github.com/agilistikmal/live-recorder/pkg/recorder"
	"github.com/agilistikmal/live-recorder/utils"
	"github.com/sirupsen/logrus"
)

type WatchLive struct {
	liveRecorder    recorder.Recorder
	recordedStreams map[string]bool
	mu              sync.Mutex
	wg              sync.WaitGroup
}

func NewWatchLive(ls recorder.Recorder) *WatchLive {
	return &WatchLive{
		liveRecorder:    ls,
		recordedStreams: make(map[string]bool),
	}
}

func (ws *WatchLive) StartWatchMode() {
	logrus.Info("Watch mode started")

	ws.CheckAndStartRecording()

	tickerDuration := time.Duration(rand.Intn(15)+15) * time.Second
	ticker := time.NewTicker(tickerDuration)
	defer ticker.Stop()
	for range ticker.C {
		logrus.Debug("Checking for new live streams...")
		ws.CheckAndStartRecording()
	}
}

func (ws *WatchLive) CheckAndStartRecording() {
	lives, err := ws.liveRecorder.GetLives()
	if err != nil {
		logrus.Errorf("Failed to get lives: %v", err)
		return
	}

	for _, live := range lives {
		streamerID := live.Streamer.Username

		ws.mu.Lock()
		isRecording := ws.recordedStreams[streamerID]
		ws.mu.Unlock()

		if isRecording {
			continue
		}

		streamingUrl, err := ws.liveRecorder.GetStreamingUrl(live)
		if err != nil {
			logrus.Errorf("Failed to get streaming url: %v", err)
			return
		}

		logrus.Infof("Recording started for %s", live.Streamer.Username)

		ws.mu.Lock()
		ws.recordedStreams[streamerID] = true
		ws.mu.Unlock()

		ws.wg.Add(1)
		go func(l recorder.Live) {
			defer ws.wg.Done()
			defer func() {
				ws.mu.Lock()
				delete(ws.recordedStreams, streamerID)
				ws.mu.Unlock()
			}()

			filename := fmt.Sprintf("./tmp/%s/%s.mp4", l.Platform, l.Streamer.Username)
			downloadInfo := utils.DownloadHLS(streamingUrl, &filename)

			if downloadInfo == nil {
				logrus.Errorf("Recording failed for %s", l.Streamer.Username)
				return
			}
			logrus.WithFields(downloadInfo).Infof("Download completed for %s", l.Streamer.Username)
		}(*live)
	}
}
