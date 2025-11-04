package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/agilistikmal/live-recorder/models"
	"github.com/agilistikmal/live-recorder/services"
	"github.com/agilistikmal/live-recorder/utils"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	watchMode := flag.Bool("watch", false, "Watch for new lives")

	platforms := flag.String("p", "", "Platforms to record (showroom,idn)")
	query := flag.String("q", "", "Query to search for lives (streamer_username:*_JKT48,48_*;title:*JKT48*)")
	url := flag.String("url", "", "URL to record (https://www.showroom-live.com/r/example)")

	flag.Parse()

	if *platforms == "" {
		logrus.Fatalf("Platforms are required")
	}

	if *query == "" && *url == "" {
		logrus.Fatalf("Query or URL is required")
	}

	var liveQuery *models.LiveQuery
	if *query != "" {
		liveQuery = &models.LiveQuery{}
		err := utils.ParseLiveQuery(*query, liveQuery)
		if err != nil {
			logrus.Fatalf("Failed to parse live query: %v", err)
		}
		liveQuery.Platforms = strings.Split(*platforms, ",")
	}

	liveService := services.NewLiveService()
	if *watchMode {
		watchService := services.NewWatchService(liveService, liveQuery)
		watchService.StartWatchMode()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		logrus.Info("Application is running in Watch Mode. Waiting for signal to stop...")
		<-quit

		logrus.Info("Received stop signal. Exiting.")
	} else {
		runOnce(liveService, liveQuery)
	}
}

func runOnce(liveService *services.LiveService, liveQuery *models.LiveQuery) {
	logrus.Info("Once mode started")
	lives, err := liveService.GetLives(liveQuery)
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
		streamingUrl, err := liveService.GetStreamingUrl(live)
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
				logrus.Errorf("Failed to download HLS: %v", err)
				return
			}
			logrus.WithFields(downloadInfo).Infof("Download completed for %s", live.Streamer.Username)
		}()
	}
	wg.Wait()
	logrus.Info("All downloads completed")
}
