package test

import (
	"testing"

	"github.com/agilistikmal/live-recorder/pkg/recorder/tiktok"
	"github.com/stretchr/testify/assert"
)

func TestTiktokLiveService_GetLive(t *testing.T) {
	tiktokRecorder := tiktok.NewRecorder()
	live, err := tiktokRecorder.GetLive("https://www.tiktok.com/@kucing_lusuh/live")
	assert.NoError(t, err, "Failed to get tiktok live")
	assert.NotNil(t, live, "Live is nil")

	assert.NotEmpty(t, live.Title, "Live title is empty")
	assert.NotEmpty(t, live.Platform, "Live platform is empty")
	assert.NotEmpty(t, live.StreamingUrl, "Live streaming url is empty")

	t.Logf("Live Title: %s", live.Title)
	t.Logf("Live Platform: %s", live.Platform)
	t.Logf("Live Streaming Url: %s", live.StreamingUrl)
}

func TestTiktokLiveService_Record(t *testing.T) {
	tiktokRecorder := tiktok.NewRecorder()
	live, err := tiktokRecorder.GetLive("https://www.tiktok.com/@bossdikha/live")
	assert.NoError(t, err, "Failed to get tiktok live")
	assert.NotNil(t, live, "Live is nil")

	assert.NotEmpty(t, live.Title, "Live title is empty")
	assert.NotEmpty(t, live.Platform, "Live platform is empty")
	assert.NotEmpty(t, live.StreamingUrl, "Live streaming url is empty")

	t.Logf("Live Title: %s", live.Title)
	t.Logf("Live Platform: %s", live.Platform)
	t.Logf("Live Streaming Url: %s", live.StreamingUrl)

	outputPath := "./tmp/tiktok/tiktok_test_result.mp4"
	err = tiktokRecorder.Record(live, outputPath)
	assert.NoError(t, err, "Failed to record tiktok live")
	t.Logf("Live recorded to: %s", outputPath)
}
