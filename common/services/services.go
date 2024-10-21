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

	for _, cyberurl := range config.CYBERSEC {
		if cyberurl == "" {
			logger.Warn("skipping empty URL in CYBERSEC config")
			continue
		}

		reader := rss.RssReader{
			Url:      cyberurl,
			Category: constants.CAT_CYBERSEC,
			Logger: logger.With(
				zap.String("type", constants.FROM_CONCURRENCY),
				zap.String("category", constants.CAT_CYBERSEC),
				zap.String("url", cyberurl),
			),
		}

		RssFeedReaders = append(RssFeedReaders, reader)
		logger.Info(fmt.Sprintf("Initialized RSS feed reader for URL: %s", cyberurl))
	}

	for _, blockchainurl := range config.BLOCKCHAIN {
		if blockchainurl == "" {
			logger.Warn("skipping empty URL in BLOCKCHAIN config")
			continue
		}

		reader := rss.RssReader{
			Url:      blockchainurl,
			Category: constants.CAT_BLOCKCHAIN,
			Logger: logger.With(
				zap.String("type", constants.FROM_CONCURRENCY),
				zap.String("category", constants.CAT_BLOCKCHAIN),
				zap.String("url", blockchainurl),
			),
		}

		RssFeedReaders = append(RssFeedReaders, reader)
		logger.Info(fmt.Sprintf("Initialized RSS feed reader for URL: %s", blockchainurl))
	}

	for _, technologyurl := range config.TECHNOLOGY {
		if technologyurl == "" {
			logger.Warn("skipping empty URL in TECHNOLOGY config")
			continue
		}

		reader := rss.RssReader{
			Url:      technologyurl,
			Category: constants.CAT_TECHNOLOGY,
			Logger: logger.With(
				zap.String("type", constants.FROM_CONCURRENCY),
				zap.String("category", constants.CAT_TECHNOLOGY),
				zap.String("url", technologyurl),
			),
		}

		RssFeedReaders = append(RssFeedReaders, reader)
		logger.Info(fmt.Sprintf("Initialized RSS feed reader for URL: %s", technologyurl))
	}

	if len(RssFeedReaders) == 0 {
		return errors.New("no valid RSS feed URLs found in config")
	}

	return nil
}
