package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/agilistikmal/live-recorder/pkg/recorder"
	"github.com/agilistikmal/live-recorder/pkg/recorder/live"
	"github.com/agilistikmal/live-recorder/pkg/recorder/watch"
	"github.com/agilistikmal/live-recorder/utils"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	watchMode := flag.Bool("watch", false, "Watch for new lives")

	platforms := flag.String("p", "", "Platforms to record (showroom,idn)")
	query := flag.String("q", "", "Query to search for lives (streamer_username:*_JKT48,48_*;title:*JKT48*)")
	url := flag.String("url", "", "URL to record (https://www.tiktok.com/@user/live)")

	flag.Parse()

	if *platforms == "" {
		logrus.Fatalf("Platforms are required")
	}

	if *query == "" && *url == "" {
		logrus.Fatalf("Query or URL is required")
	}

	if *url != "" {
		liveRecorder := live.NewRecorder(&recorder.LiveQuery{
			Platforms: strings.Split(*platforms, ","),
		})
		live, err := liveRecorder.GetLive(*url)
		if err != nil {
			logrus.Fatalf("Failed to get live: %v", err)
		}
		filename := fmt.Sprintf("./tmp/%s/%s.mp4", live.Platform, live.Streamer.Username)
		liveRecorder.Record(live, filename)
		logrus.Infof("Download completed: %v", filename)
		return
	}

	var liveQuery *recorder.LiveQuery
	if *query != "" {
		liveQuery = &recorder.LiveQuery{}
		err := utils.ParseLiveQuery(*query, liveQuery)
		if err != nil {
			logrus.Fatalf("Failed to parse live query: %v", err)
		}
		liveQuery.Platforms = strings.Split(*platforms, ",")
	}

	liveRecorder := live.NewRecorder(liveQuery)
	if *watchMode {
		// Create buffered channels for events
		liveChan := make(chan *recorder.Live, 100)        // Buffer 100 events
		statusChan := make(chan *watch.StatusUpdate, 100) // Buffer 100 status updates

		watchService := watch.NewWatchLive(liveRecorder, "./tmp")
		watchService.SetLiveChannel(liveChan)
		watchService.SetStatusChannel(statusChan)

		// Start goroutine to consume live events from channel
		go func() {
			for live := range liveChan {
				// Process live event here
				// Contoh: logging, notification, webhook, dll
				logrus.WithFields(logrus.Fields{
					"platform":      live.Platform,
					"streamer":      live.Streamer.Username,
					"title":         live.Title,
					"view_count":    live.ViewCount,
					"platform_url":  live.PlatformUrl,
					"streaming_url": live.StreamingUrl,
				}).Info("New live stream detected")
			}
		}()

		// Start goroutine to consume status updates from channel
		go func() {
			for update := range statusChan {
				// Process status updates here
				// Contoh: logging, notification, webhook, database update, dll
				logrus.WithFields(logrus.Fields{
					"streamer_id":  update.StreamerID,
					"status":       update.Status,
					"platform":     update.Info.Live.Platform,
					"title":        update.Info.Live.Title,
					"started_at":   update.Info.StartedAt,
					"completed_at": update.Info.CompletedAt,
					"file_path":    update.Info.FilePath,
					"file_size":    update.Info.FileSize,
					"error":        update.Info.Error,
				}).Info("Recording status update")

				// Contoh: Handle berdasarkan status
				switch update.Status {
				case watch.StatusInProgress:
					logrus.Infof("Recording started for %s", update.StreamerID)
				case watch.StatusCompleted:
					logrus.Infof("Recording completed for %s: %s (Size: %d bytes)",
						update.StreamerID, update.Info.FilePath, update.Info.FileSize)
				case watch.StatusFailed:
					logrus.Errorf("Recording failed for %s: %v", update.StreamerID, update.Info.Error)
				}
			}
		}()

		// Start watch mode in goroutine so it doesn't block
		go watchService.StartWatchMode()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		logrus.Info("Application is running in Watch Mode. Waiting for signal to stop...")
		<-quit

		// Cleanup: close channels
		close(liveChan)
		close(statusChan)
		logrus.Info("Received stop signal. Exiting.")

		// Contoh: Print final status summary
		allStatuses := watchService.GetAllStatuses()
		logrus.Infof("Final status summary: %d recordings", len(allStatuses))
		for streamerID, info := range allStatuses {
			logrus.Infof("  %s: %s", streamerID, info.Status)
		}
	} else {
		runOnce(liveRecorder, liveQuery)
	}
}

func runOnce(liveRecorder recorder.Recorder, liveQuery *recorder.LiveQuery) {
	logrus.Info("Once mode started")
	lives, err := liveRecorder.GetLives()
	if err != nil {
		logrus.Errorf("Failed to get lives: %v", err)
		return
	}

	logrus.Infof("Found %d lives", len(lives))

	if len(lives) == 0 {
		logrus.Info("No lives found")
		return
	}

	wg := sync.WaitGroup{}
	for _, live := range lives {
		wg.Add(1)
		streamingUrl, err := liveRecorder.GetStreamingUrl(live)
		if err != nil {
			logrus.Errorf("Failed to get streaming url: %v", err)
			return
		}

		go func() {
			defer wg.Done()
			logrus.Infof("Recording started for %s", live.Streamer.Username)

			filename := fmt.Sprintf("./tmp/%s/%s.mp4", live.Platform, live.Streamer.Username)
			downloadInfo := utils.DownloadHLS(streamingUrl, &filename)
			if downloadInfo == nil {
				return
			}
			logrus.WithFields(downloadInfo).Infof("Download completed for %s", live.Streamer.Username)
		}()
	}
	wg.Wait()
	logrus.Info("All downloads completed")
}
