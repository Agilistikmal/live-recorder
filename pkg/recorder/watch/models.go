package watch

import (
	"time"

	"github.com/agilistikmal/live-recorder/pkg/recorder"
)

// RecordingStatus represents the status of a recording
type RecordingStatus string

const (
	StatusInProgress RecordingStatus = "in_progress"
	StatusCompleted  RecordingStatus = "completed"
	StatusFailed     RecordingStatus = "failed"
)

// RecordingInfo contains information about a recording
type RecordingInfo struct {
	Live        *recorder.Live
	Status      RecordingStatus
	StartedAt   time.Time
	CompletedAt *time.Time
	FilePath    string
	FileSize    int64
	Error       error
}

// StatusUpdate represents a status update event
type StatusUpdate struct {
	StreamerID string
	Status     RecordingStatus
	Info       *RecordingInfo
}
