package live

import (
	"fmt"
	"strings"
	"sync"

	"github.com/agilistikmal/live-recorder/pkg/recorder"
	"github.com/agilistikmal/live-recorder/pkg/recorder/idn"
	"github.com/agilistikmal/live-recorder/pkg/recorder/models"
	"github.com/agilistikmal/live-recorder/pkg/recorder/showroom"
	"github.com/agilistikmal/live-recorder/pkg/recorder/tiktok"
	"github.com/sirupsen/logrus"
)

type LiveRecorder struct {
	showroomRecorder recorder.Recorder
	idnRecorder      recorder.Recorder
	tiktokRecorder   recorder.Recorder
	liveQuery        *models.LiveQuery
}

func NewRecorder(liveQuery *models.LiveQuery) recorder.Recorder {
	return &LiveRecorder{
		showroomRecorder: showroom.NewRecorder(),
		idnRecorder:      idn.NewRecorder(),
		tiktokRecorder:   tiktok.NewRecorder(),
		liveQuery:        liveQuery,
	}
}

func (s *LiveRecorder) GetLives() ([]*models.Live, error) {
	lives := make([]*models.Live, 0)
	wg := sync.WaitGroup{}
	for _, platform := range s.liveQuery.Platforms {
		switch platform {
		case models.PlatformShowroom:
			wg.Add(1)
			go func() {
				defer wg.Done()
				showroomLives, err := s.showroomRecorder.GetLives()
				if err != nil {
					logrus.Errorf("Failed to get showroom lives: %v", err)
					return
				}
				filteredShowroomLives, err := s.ApplyFilter(showroomLives, s.liveQuery)
				if err != nil {
					logrus.Errorf("Failed to apply filter to showroom lives: %v", err)
					return
				}
				lives = append(lives, filteredShowroomLives...)
			}()
		case models.PlatformIDN:
			wg.Add(1)
			go func() {
				defer wg.Done()
				idnLives, err := s.idnRecorder.GetLives()
				if err != nil {
					logrus.Errorf("Failed to get idn lives: %v", err)
					return
				}
				filteredIdnLives, err := s.ApplyFilter(idnLives, s.liveQuery)
				if err != nil {
					logrus.Errorf("Failed to apply filter to idn lives: %v", err)
					return
				}
				lives = append(lives, filteredIdnLives...)
			}()
		default:
			logrus.Errorf("Invalid platform: %s", platform)
			return nil, fmt.Errorf("invalid platform: %s", platform)
		}
	}
	wg.Wait()

	// Remove duplicates
	uniqueLives := make([]*models.Live, 0)
	seen := make(map[string]bool)
	for _, live := range lives {
		if !seen[live.ID] {
			seen[live.ID] = true
			uniqueLives = append(uniqueLives, live)
		}
	}
	lives = uniqueLives
	return lives, nil
}

func (s *LiveRecorder) GetLive(url string) (*models.Live, error) {
	if len(s.liveQuery.Platforms) < 1 {
		return nil, fmt.Errorf("no platforms provided")
	}

	switch s.liveQuery.Platforms[0] {
	case models.PlatformShowroom:
		return s.showroomRecorder.GetLive(url)
	case models.PlatformIDN:
		return s.idnRecorder.GetLive(url)
	case models.PlatformTiktok:
		return s.tiktokRecorder.GetLive(url)
	default:
		return nil, fmt.Errorf("invalid platform: %s", s.liveQuery.Platforms[0])
	}
}

func (s *LiveRecorder) GetStreamingUrl(live *models.Live) (string, error) {
	switch live.Platform {
	case models.PlatformShowroom:
		return s.showroomRecorder.GetStreamingUrl(live)
	default:
		return live.StreamingUrl, nil
	}
}

func (s *LiveRecorder) Record(live *models.Live, outputPath string) error {
	switch live.Platform {
	case models.PlatformShowroom:
		return s.showroomRecorder.Record(live, outputPath)
	default:
		return s.idnRecorder.Record(live, outputPath)
	}
}

func (s *LiveRecorder) ApplyFilter(lives []*models.Live, query *models.LiveQuery) ([]*models.Live, error) {
	filteredList := make([]*models.Live, 0, len(lives))

	for _, live := range lives {
		if s.CheckFilters(live, query) {
			filteredList = append(filteredList, live)
		}
	}

	return filteredList, nil
}

func (s *LiveRecorder) CheckFilters(live *models.Live, liveQuery *models.LiveQuery) bool {
	if live.Streamer == nil {
		if liveQuery.StreamerUsernameLike != "" {
			return false
		}
	}

	liveQuery.StreamerUsernameLike = strings.ToLower(liveQuery.StreamerUsernameLike)
	liveQuery.TitleLike = strings.ToLower(liveQuery.TitleLike)

	liveTitleLower := strings.ToLower(live.Title)
	liveUsernameLower := ""
	if live.Streamer != nil {
		liveUsernameLower = strings.ToLower(live.Streamer.Username)
	}

	// Filter Streamer Username LIKE (Wildcard *)
	liveQueryStreamerUsernames := strings.SplitSeq(liveQuery.StreamerUsernameLike, ",")
	streamerUsernameFilterPassed := false
	for liveQueryStreamerUsername := range liveQueryStreamerUsernames {
		streamerUsernameFilterPassed = streamerUsernameFilterPassed || s.CheckWildcardFilter(liveUsernameLower, liveQueryStreamerUsername)
	}

	// Filter Title LIKE (Wildcard *)
	liveQueryTitles := strings.SplitSeq(liveQuery.TitleLike, ",")
	titleFilterPassed := false
	for liveQueryTitle := range liveQueryTitles {
		titleFilterPassed = titleFilterPassed || s.CheckWildcardFilter(liveTitleLower, liveQueryTitle)
	}

	// Return true if both filters passed
	return streamerUsernameFilterPassed && titleFilterPassed
}

func (s *LiveRecorder) CheckWildcardFilter(text, filter string) bool {
	if strings.HasPrefix(filter, "*") && !strings.HasSuffix(filter, "*") && len(filter) > 1 {
		suffix := strings.TrimPrefix(filter, "*")
		return strings.HasSuffix(text, suffix)
	}

	if strings.HasSuffix(filter, "*") && !strings.HasPrefix(filter, "*") && len(filter) > 1 {
		prefix := strings.TrimSuffix(filter, "*")
		return strings.HasPrefix(text, prefix)
	}

	cleanedFilter := strings.ReplaceAll(filter, "*", "")
	if cleanedFilter == "" {
		return true
	}

	return strings.Contains(text, cleanedFilter)
}
