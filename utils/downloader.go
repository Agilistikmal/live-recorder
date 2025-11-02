package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func DownloadHLS(url string, outputPath string) map[string]interface{} {

	if _, err := os.Stat(filepath.Dir(outputPath)); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(outputPath), 0755)
	}

	cmd := exec.Command("ffmpeg",
		"-i", url,
		"-c", "copy",
		"-bsf:a", "aac_adtstoasc",
		"-movflags", "faststart",
		outputPath,
	)

	err := cmd.Run()
	if err != nil {
		return nil
	}

	fileInfo, err := os.Stat(outputPath)
	if err != nil {
		return nil
	}

	downloadInfo := map[string]interface{}{
		"url":          url,
		"output_path":  outputPath,
		"size":         fileInfo.Size(),
		"duration":     fileInfo.Size() / 1024 / 1024,
		"started_at":   time.Now(),
		"completed_at": time.Now(),
	}

	return downloadInfo
}
