package services

import (
	"fmt"
	"strings"
	"sync"

	"github.com/agilistikmal/live-recorder/models"
	"github.com/sirupsen/logrus"
)

type LiveService struct {
	showroomLiveService *ShowroomLiveService
	idnLiveService      *IDNLiveService
}

func NewLiveService() *LiveService {
	return &LiveService{
		showroomLiveService: NewShowroomLiveService(),
		idnLiveService:      NewIDNLiveService(),
	}
}

func (s *LiveService) GetLives(liveQuery *models.LiveQuery) ([]*models.Live, error) {
	lives := make([]*models.Live, 0)
	wg := sync.WaitGroup{}
	for _, platform := range liveQuery.Platforms {
		switch platform {
		case models.PlatformShowroom:
			wg.Add(1)
			go func() {
				defer wg.Done()
				showroomLives, err := s.showroomLiveService.GetLives()
				if err != nil {
					logrus.Errorf("Failed to get showroom lives: %v", err)
					return
				}
				filteredShowroomLives, err := s.ApplyFilter(showroomLives, liveQuery)
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
				idnLives, err := s.idnLiveService.GetLives()
				if err != nil {
					logrus.Errorf("Failed to get idn lives: %v", err)
					return
				}
				filteredIdnLives, err := s.ApplyFilter(idnLives, liveQuery)
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
	return lives, nil
}

func (s *LiveService) GetStreamingUrl(live *models.Live) (string, error) {
	switch live.Platform {
	case models.PlatformShowroom:
		return s.showroomLiveService.GetStreamingUrl(live.ID)
	default:
		return live.StreamingUrl, nil
	}
}

func (s *LiveService) ApplyFilter(lives []*models.Live, query *models.LiveQuery) ([]*models.Live, error) {
	filteredList := make([]*models.Live, 0, len(lives))

	for _, live := range lives {
		if s.CheckFilters(live, query) {
			filteredList = append(filteredList, live)
		}
	}

	return filteredList, nil
}

func (s *LiveService) CheckFilters(live *models.Live, liveQuery *models.LiveQuery) bool {
	if live.Streamer == nil {
		if liveQuery.StreamerUsernameLike != "" {
			return false
		}
	}

	liveTitleLower := strings.ToLower(live.Title)
	liveUsernameLower := ""
	if live.Streamer != nil {
		liveUsernameLower = strings.ToLower(live.Streamer.Username)
	}

	// Filter Streamer Username LIKE (Wildcard *)
	liveQueryStreamerUsernames := strings.SplitSeq(liveUsernameLower, ",")
	streamerUsernameFilterPassed := false
	for liveQueryStreamerUsername := range liveQueryStreamerUsernames {
		streamerUsernameFilterPassed = streamerUsernameFilterPassed || s.CheckWildcardFilter(liveUsernameLower, liveQueryStreamerUsername)
	}

	// Filter Title LIKE (Wildcard *)
	liveQueryTitles := strings.SplitSeq(liveTitleLower, ",")
	titleFilterPassed := false
	for liveQueryTitle := range liveQueryTitles {
		titleFilterPassed = titleFilterPassed || s.CheckWildcardFilter(liveTitleLower, liveQueryTitle)
	}

	// Return true if both filters passed
	return streamerUsernameFilterPassed && titleFilterPassed
}

func (s *LiveService) CheckWildcardFilter(text, filter string) bool {
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
