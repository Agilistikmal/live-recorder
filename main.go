package main

import (
	"fmt"

	"github.com/agilistikmal/live-recorder/services"
	"github.com/sirupsen/logrus"
)

func main() {
	showroomLiveService := services.NewShowroomLiveService()
	showroomLives, err := showroomLiveService.GetLives()
	if err != nil {
		logrus.Fatalf("Failed to get showroom lives: %v", err)
	}

	for _, showroomLive := range showroomLives {
		fmt.Println(showroomLive.PlatformUrl)
	}
}
