package services

import (
	"errors"
	"fmt"

	"github.com/heyrovsky/rsscurator/common/constants"
	"github.com/heyrovsky/rsscurator/config"
	"github.com/heyrovsky/rsscurator/pkg/rss"
	"go.uber.org/zap"
)

var RssFeedReaders []rss.RssReader

func InitServices(logger *zap.Logger) error {
	if logger == nil {
		return errors.New("logger cannot be nil")
	}
	fmt.Println(config.FEEDS)
	for _, feedurl := range config.FEEDS {
		if feedurl == "" {
			logger.Warn("skipping empty URL in CYBERSEC config")
			continue
		}

		reader := rss.RssReader{
			Url:      feedurl,
			Category: constants.CAT_CYBERSEC,
			Logger: logger.With(
				zap.String("type", constants.FROM_CONCURRENCY),
				zap.String("category", constants.CAT_CYBERSEC),
				zap.String("url", feedurl),
			),
		}

		RssFeedReaders = append(RssFeedReaders, reader)
		logger.Info(fmt.Sprintf("Initialized RSS feed reader for URL: %s", feedurl))
	}

	if len(RssFeedReaders) == 0 {
		return errors.New("no valid RSS feed URLs found in config")
	}

	return nil
}
