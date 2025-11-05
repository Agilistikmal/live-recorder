package recorder

import "github.com/agilistikmal/live-recorder/pkg/recorder/models"

type RecorderConfig struct {
	UserAgent string `json:"user_agent"`
	Referer   string `json:"referer"`
	Cookie    string `json:"cookie"`
}

type Recorder interface {
	GetLives() ([]*models.Live, error)
	GetStreamingUrl(live *models.Live) (string, error)
	Record(live *models.Live, outputPath string) error
}
