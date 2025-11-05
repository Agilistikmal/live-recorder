package test

import (
	"os"
	"testing"

	"github.com/agilistikmal/live-recorder/pkg/recorder/showroom"
	"github.com/agilistikmal/live-recorder/utils"
)

func TestShowroomLiveService_Download(t *testing.T) {
	outputPath := "showroom_test_result.mp4"
	defer os.Remove(outputPath)
	showroomRecorder := showroom.NewRecorder()
	showroomLives, err := showroomRecorder.GetLives()
	if err != nil {
		t.Fatalf("Failed to get showroom lives: %v", err)
	}

	t.Logf("showroom lives: %v", len(showroomLives))
	if len(showroomLives) < 1 {
		t.Fatalf("No showroom lives found")
	}

	streamingUrl, err := showroomRecorder.GetStreamingUrl(showroomLives[0])
	if err != nil {
		t.Fatalf("Failed to get streaming url: %v", err)
	}

	t.Logf("streaming url: %v", streamingUrl)

	t.Logf("Downloading 5 seconds HLS...")
	go func() {
		downloadInfo := utils.DownloadHLS(streamingUrl, &outputPath)
		if downloadInfo == nil {
			t.Log("Failed to download HLS")
		}
		t.Logf("Download completed: %v", downloadInfo)
	}()
}
