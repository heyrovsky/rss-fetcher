package rss

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/heyrovsky/rsscurator/common/reader"
	"github.com/heyrovsky/rsscurator/pkg/content"
	"go.uber.org/zap"
)

type RssReader struct {
	Url      string
	Category string
	Logger   *zap.Logger
}

func (r *RssReader) ReadNews() ([]content.NewsItemHashed, error) {
	newsItems := []content.NewsItemHashed{}

	feed, err := reader.Url.ParseURL(r.Url)
	if err != nil {
		return newsItems, err
	}

	for _, item := range feed.Items {
		newsItem := content.NewsItem{
			Title:       item.Title,
			Description: item.Description,
			Links:       item.Links,
			Categories:  item.Categories,
			FeedSource:  feed.Title,
			Authors:     item.Authors,
			Image:       item.Image,
			Updated:     item.UpdatedParsed,
			Published:   item.PublishedParsed,
		}

		hash, err := generateHash(newsItem)
		if err != nil {
			return newsItems, err
		}

		hashedItem := content.NewsItemHashed{
			Hash: hash,
			Item: newsItem,
		}

		newsItems = append(newsItems, hashedItem)
	}

	return newsItems, nil
}

func generateHash(item content.NewsItem) (string, error) {
	jsonData, err := json.Marshal(item)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(jsonData)

	hashString := fmt.Sprintf("%x", hash)

	return hashString, nil
}
