package watch

import (
	"fmt"
	"os"
	"sync"
	"time"

	"math/rand"

	"github.com/agilistikmal/live-recorder/pkg/recorder"
	"github.com/agilistikmal/live-recorder/utils"
	"github.com/sirupsen/logrus"
)

type WatchLive struct {
	liveRecorder recorder.Recorder
	recordings   map[string]*RecordingInfo
	mu           sync.RWMutex
	wg           sync.WaitGroup
	liveChan     chan *recorder.Live
	statusChan   chan *StatusUpdate
	outputDir    string
}

func NewWatchLive(ls recorder.Recorder, outputDir string) *WatchLive {
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.MkdirAll(outputDir, 0755)
	}

	return &WatchLive{
		liveRecorder: ls,
		outputDir:    outputDir,
		recordings:   make(map[string]*RecordingInfo),
		liveChan:     nil,
		statusChan:   nil,
	}
}

// SetLiveChannel sets the channel for receiving new live stream data.
// If channel is nil, no data will be sent.
// Channel should be buffered to avoid blocking.
func (ws *WatchLive) SetLiveChannel(ch chan *recorder.Live) {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	ws.liveChan = ch
}

// SetStatusChannel sets the channel for receiving status updates.
// If channel is nil, no status updates will be sent.
// Channel should be buffered to avoid blocking.
func (ws *WatchLive) SetStatusChannel(ch chan *StatusUpdate) {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	ws.statusChan = ch
}

// GetStatus returns the recording info for a given streamer ID.
// Returns nil if not found.
func (ws *WatchLive) GetStatus(streamerID string) (*RecordingInfo, bool) {
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	info, exists := ws.recordings[streamerID]
	return info, exists
}

// GetAllStatuses returns all recording statuses.
func (ws *WatchLive) GetAllStatuses() map[string]*RecordingInfo {
	ws.mu.RLock()
	defer ws.mu.RUnlock()

	result := make(map[string]*RecordingInfo)
	for k, v := range ws.recordings {
		result[k] = v
	}
	return result
}

// GetStatusesByStatus returns all recordings with a specific status.
func (ws *WatchLive) GetStatusesByStatus(status RecordingStatus) []*RecordingInfo {
	ws.mu.RLock()
	defer ws.mu.RUnlock()

	var result []*RecordingInfo
	for _, info := range ws.recordings {
		if info.Status == status {
			result = append(result, info)
		}
	}
	return result
}

// sendStatusUpdate sends a status update to the status channel (non-blocking)
func (ws *WatchLive) sendStatusUpdate(streamerID string, status RecordingStatus, info *RecordingInfo) {
	ws.mu.RLock()
	ch := ws.statusChan
	ws.mu.RUnlock()

	if ch != nil {
		select {
		case ch <- &StatusUpdate{
			StreamerID: streamerID,
			Status:     status,
			Info:       info,
		}:
			// Successfully sent
		default:
			// Channel is full, log warning but don't block
			logrus.Warnf("Status channel is full, dropping status update for %s", streamerID)
		}
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

		ws.mu.RLock()
		_, exists := ws.recordings[streamerID]
		ws.mu.RUnlock()

		if exists {
			continue
		}

		streamingUrl, err := ws.liveRecorder.GetStreamingUrl(live)
		if err != nil {
			logrus.Errorf("Failed to get streaming url: %v", err)
			return
		}

		// Create recording info with InProgress status
		recordingInfo := &RecordingInfo{
			Live:      live,
			Status:    StatusInProgress,
			StartedAt: time.Now(),
		}

		// Store recording info
		ws.mu.Lock()
		ws.recordings[streamerID] = recordingInfo
		liveCh := ws.liveChan
		ws.mu.Unlock()

		// Send live data to channel if available (non-blocking)
		if liveCh != nil {
			select {
			case liveCh <- live:
				// Successfully sent to channel
			default:
				// Channel is full, log warning but don't block
				logrus.Warnf("Live channel is full, dropping live data for %s", live.Streamer.Username)
			}
		}

		// Send status update for InProgress
		ws.sendStatusUpdate(streamerID, StatusInProgress, recordingInfo)

		ws.wg.Add(1)
		go func(l recorder.Live, streamID string) {
			defer ws.wg.Done()

			filename := fmt.Sprintf("%s/%s/%s.mp4", ws.outputDir, l.Platform, l.Streamer.Username)
			downloadInfo := utils.DownloadHLS(streamingUrl, &filename)

			// Update status based on result
			ws.mu.Lock()
			recordingInfo := ws.recordings[streamID]
			if recordingInfo == nil {
				ws.mu.Unlock()
				return
			}

			if downloadInfo == nil {
				// Recording failed
				recordingInfo.Status = StatusFailed
				recordingInfo.Error = fmt.Errorf("recording failed")
				ws.mu.Unlock()

				logrus.Errorf("Recording failed for %s", l.Streamer.Username)
				ws.sendStatusUpdate(streamID, StatusFailed, recordingInfo)
			} else {
				// Recording completed
				now := time.Now()
				recordingInfo.Status = StatusCompleted
				recordingInfo.CompletedAt = &now
				recordingInfo.FilePath = filename

				// Extract file size if available
				if size, ok := downloadInfo["size"].(int64); ok {
					recordingInfo.FileSize = size
				}

				ws.mu.Unlock()

				ws.sendStatusUpdate(streamID, StatusCompleted, recordingInfo)
			}
		}(*live, streamerID)
	}
}
