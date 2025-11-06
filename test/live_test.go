package test

import (
	"strings"
	"testing"

	"github.com/agilistikmal/live-recorder/pkg/recorder"
	"github.com/agilistikmal/live-recorder/pkg/recorder/live"
	"github.com/stretchr/testify/assert"
)

func TestLiveService_GetLives(t *testing.T) {
	liveQuery := &recorder.LiveQuery{
		Platforms:            []string{recorder.PlatformIDN, recorder.PlatformShowroom},
		StreamerUsernameLike: "*4*",
	}
	liveRecorder := live.NewRecorder(liveQuery)
	lives, err := liveRecorder.GetLives()
	assert.NoError(t, err, "Failed to get lives")
	assert.Greater(t, len(lives), 0, "No lives found")

	usernames := make([]string, 0)
	for _, live := range lives {
		usernames = append(usernames, live.Streamer.Username)
	}
	assert.Contains(t, strings.Join(usernames, ","), "4")
	t.Logf("usernames: %v", strings.Join(usernames, ","))
}
