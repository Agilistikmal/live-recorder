package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func DownloadHLS(url string, outputPath *string) map[string]interface{} {
	if _, err := os.Stat(filepath.Dir(*outputPath)); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(*outputPath), 0755)
	}

	ext := filepath.Ext(*outputPath)
	outputPathWithoutExt := strings.TrimSuffix(*outputPath, ext)

	timestamp := time.Now().Unix()
	outputPathTemp := fmt.Sprintf("%s_%d.tmp%s", outputPathWithoutExt, timestamp, ext)

	cmd := exec.Command("ffmpeg",
		// "-t", "30", // For testing purposes (recording 30 seconds)
		"-i", url,
		"-y",
		"-c", "copy",
		"-bsf:a", "aac_adtstoasc",
		"-movflags", "faststart",
		outputPathTemp,
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		logrus.Errorf("Failed to download HLS using ffmpeg: %v, stderr: %s", err, stderr.String())
		return nil
	}

	tempFiles, err := filepath.Glob(fmt.Sprintf("%s_*.tmp%s", outputPathWithoutExt, ext))
	if err != nil {
		logrus.Errorf("Failed to get temp files: %v", err)
		return nil
	}

	sort.Strings(tempFiles)

	// File list for FFmpeg concat demuxer
	listFilePath := outputPathWithoutExt + ".list"
	listContent := ""
	for _, tempFile := range tempFiles {
		listContent += fmt.Sprintf("file '%s'\n", filepath.Base(tempFile))
	}
	err = os.WriteFile(listFilePath, []byte(listContent), 0644)
	if err != nil {
		logrus.Errorf("Failed to write list file: %v", err)
		return nil
	}

	outputPathFinal := fmt.Sprintf("%s_%d%s", outputPathWithoutExt, timestamp, ext)
	*outputPath = outputPathFinal

	// Joining Files to Output
	cmd = exec.Command("ffmpeg",
		"-f", "concat",
		"-safe", "0",
		"-i", listFilePath,
		"-y",
		"-c", "copy",
		"-bsf:a", "aac_adtstoasc",
		"-movflags", "faststart",
		*outputPath,
	)

	err = cmd.Run()

	// Cleanup temporary files
	os.Remove(listFilePath)
	for _, tempFile := range tempFiles {
		os.Remove(tempFile)
	}

	if err != nil {
		logrus.Errorf("Failed to join files: %v", err)
		return nil
	}

	fileInfo, err := os.Stat(*outputPath)
	if err != nil {
		logrus.Errorf("Failed to get file info: %v", err)
		return nil
	}

	downloadInfo := map[string]interface{}{
		"url":          url,
		"output_path":  *outputPath,
		"size":         fileInfo.Size(),
		"duration":     fileInfo.Size() / 1024 / 1024,
		"started_at":   time.Now(),
		"completed_at": time.Now(),
	}

	return downloadInfo
}
