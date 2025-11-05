package utils

import (
	"strings"

	"github.com/agilistikmal/live-recorder/pkg/recorder/models"
	"github.com/sirupsen/logrus"
)

func ParseLiveQuery(query string, liveQuery *models.LiveQuery) error {
	pairs := strings.Split(query, ",")

	for _, pair := range pairs {
		parts := strings.SplitN(strings.TrimSpace(pair), ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "title":
			liveQuery.TitleLike = value
		case "streamer_username":
			liveQuery.StreamerUsernameLike = value
		default:
			logrus.Warnf("Unknown query key: %s", key)
		}
	}
	return nil
}
