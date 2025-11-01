package services

import (
	"errors"
	"strings"

	"github.com/agilistikmal/live-recorder/models"
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
	switch liveQuery.Platform {
	case models.PlatformShowroom:
		lives, err := s.showroomLiveService.GetLives()
		if err != nil {
			return nil, err
		}
		return s.ApplyFilter(lives, liveQuery)
	case models.PlatformIDN:
		// TODO: Implement IDN live service
		return nil, errors.New("idn platform not implemented")
	default:
		return nil, errors.New("invalid platform")
	}
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
	if liveQuery.StreamerUsernameLike != "" {
		filterValue := strings.ToLower(liveQuery.StreamerUsernameLike)

		if !s.CheckWildcardFilter(liveUsernameLower, filterValue) {
			return false
		}
	}

	// Filter Title LIKE (Wildcard *)
	if liveQuery.TitleLike != "" {
		filterValue := strings.ToLower(liveQuery.TitleLike)

		if !s.CheckWildcardFilter(liveTitleLower, filterValue) {
			return false
		}
	}

	return true
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
