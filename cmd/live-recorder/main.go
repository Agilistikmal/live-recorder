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
		watchService := watch.NewWatchLive(liveRecorder)
		watchService.StartWatchMode()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		logrus.Info("Application is running in Watch Mode. Waiting for signal to stop...")
		<-quit

		logrus.Info("Received stop signal. Exiting.")
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
