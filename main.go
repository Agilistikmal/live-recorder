package main

import (
	"flag"
	"fmt"

	"github.com/agilistikmal/live-recorder/models"
	"github.com/agilistikmal/live-recorder/services"
	"github.com/agilistikmal/live-recorder/utils"
	"github.com/sirupsen/logrus"
)

func main() {
	liveService := services.NewLiveService()

	platform := flag.String("p", "", "Platform to record (showroom, idn)")
	query := flag.String("q", "", "Query to search for lives")
	url := flag.String("url", "", "URL to record (https://www.showroom-live.com/r/example)")

	flag.Parse()

	if *platform == "" {
		logrus.Fatalf("Platform is required")
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
		liveQuery.Platform = *platform
	}

	lives, err := liveService.GetLives(liveQuery)
	if err != nil {
		logrus.Fatalf("Failed to get lives: %v", err)
	}

	logrus.Infof("lives: %v", len(lives))
	if len(lives) < 50 {
		for _, live := range lives {
			streamingUrl, err := liveService.GetStreamingUrl(live)
			if err != nil {
				logrus.Fatalf("Failed to get streaming url: %v", err)
			}
			logrus.Infof("Starting to record: %s %s (%s)", live.Platform, live.Streamer.Username, streamingUrl)
			go utils.DownloadHLS(streamingUrl, fmt.Sprintf("./tmp/%v.mp4", live.Streamer.Username))
		}
	}

	select {}
}
