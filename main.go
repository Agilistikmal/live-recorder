package main

import (
	"github.com/agilistikmal/live-recorder/services"
	"github.com/agilistikmal/live-recorder/utils"
	"github.com/sirupsen/logrus"
)

func main() {
	showroomLiveService := services.NewShowroomLiveService()
	showroomLives, err := showroomLiveService.GetLives()
	if err != nil {
		logrus.Fatalf("Failed to get showroom lives: %v", err)
	}

	logrus.Infof("showroom lives: %v", len(showroomLives))
	if len(showroomLives) < 1 {
		logrus.Fatalf("No showroom lives found")
	}

	streamingUrl, err := showroomLiveService.GetStreamingUrl(showroomLives[0].ID)
	if err != nil {
		logrus.Fatalf("Failed to get streaming url: %v", err)
	}

	err = utils.DownloadHLS(streamingUrl, "output.mp4")
	if err != nil {
		logrus.Fatalf("Failed to download HLS: %v", err)
	}
}
