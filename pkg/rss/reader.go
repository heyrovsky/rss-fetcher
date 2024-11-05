package rss

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/heyrovsky/rsscurator/pkg/content"
	"github.com/mmcdole/gofeed"
	"go.uber.org/zap"
)

// RssReader reads news items from an RSS feed and generates unique hashed news items.
type RssReader struct {
	Url      string
	Category string
	Logger   *zap.Logger
}

// ReadNews fetches and processes news items from the configured RSS feed URL, returning a list of hashed news items.
func (r *RssReader) ReadNews() ([]content.NewsItemHashed, error) {
	if r.Logger == nil {
		return nil, fmt.Errorf("logger is not initialized")
	}
	if r.Url == "" {
		r.Logger.Error("empty RSS feed URL")
		return nil, fmt.Errorf("RSS feed URL cannot be empty")
	}

	// Initialize parser and parse the feed URL
	r.Logger.Info("Starting to read RSS feed", zap.String("url", r.Url))
	parser := gofeed.NewParser()
	feed, err := parser.ParseURL(r.Url)
	if err != nil {
		r.Logger.Error("failed to parse RSS feed", zap.String("url", r.Url), zap.Error(err))
		return nil, fmt.Errorf("failed to parse RSS feed: %w", err)
	}

	var newsItems []content.NewsItemHashed
	for _, item := range feed.Items {
		newsItem, err := r.createNewsItem(feed, item)
		if err != nil {
			r.Logger.Warn("failed to create news item", zap.String("title", item.Title), zap.Error(err))
			continue
		}

		hashedItem, err := r.hashNewsItem(newsItem)
		if err != nil {
			r.Logger.Warn("failed to generate hash for news item", zap.String("title", item.Title), zap.Error(err))
			continue
		}

		newsItems = append(newsItems, hashedItem)
		r.Logger.Debug("news item processed", zap.String("title", newsItem.Title), zap.String("hash", hashedItem.Hash))
	}

	r.Logger.Info("RSS feed reading completed", zap.Int("items_processed", len(newsItems)))
	return newsItems, nil
}

// createNewsItem constructs a NewsItem from the RSS feed and item details.
func (r *RssReader) createNewsItem(feed *gofeed.Feed, item *gofeed.Item) (content.NewsItem, error) {
	if item == nil {
		return content.NewsItem{}, fmt.Errorf("nil item cannot be processed")
	}
	return content.NewsItem{
		Title:       item.Title,
		Description: item.Description,
		Links:       item.Links,
		Categories:  item.Categories,
		FeedSource:  feed.Title,
		Authors:     item.Authors,
		Image:       item.Image,
		Updated:     item.UpdatedParsed,
		Published:   item.PublishedParsed,
	}, nil
}

// hashNewsItem generates a unique SHA-256 hash for the given news item.
func (r *RssReader) hashNewsItem(newsItem content.NewsItem) (content.NewsItemHashed, error) {
	hash, err := generateHash(newsItem)
	if err != nil {
		return content.NewsItemHashed{}, fmt.Errorf("failed to generate hash: %w", err)
	}
	return content.NewsItemHashed{
		Hash: hash,
		Item: newsItem,
	}, nil
}

// generateHash serializes a NewsItem to JSON and computes its SHA-256 hash.
func generateHash(item content.NewsItem) (string, error) {
	jsonData, err := json.Marshal(item)
	if err != nil {
		return "", fmt.Errorf("failed to marshal news item for hashing: %w", err)
	}

	hash := sha256.Sum256(jsonData)
	return fmt.Sprintf("%x", hash), nil
}
