package utils

import "os/exec"

func DownloadHLS(url string, outputPath string) error {
	cmd := exec.Command("ffmpeg",
		"-i", url,
		"-c", "copy",
		"-bsf:a", "aac_adtstoasc",
		"-movflags", "faststart",
		outputPath,
	)

	return cmd.Run()
}
