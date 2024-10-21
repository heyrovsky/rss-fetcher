package content

import (
	"time"

	"github.com/mmcdole/gofeed"
)

// NewsFeed represents a collection of news items with a timestamp
type NewsFeed struct {
	Items       []NewsItemHashed `json:"items"`
	LastUpdated time.Time        `json:"time"`
}

type NewsItemHashed struct {
	Hash string   `json:"hash"`
	Item NewsItem `json:"item"`
}

// NewsItem represents an individual news article or item
type NewsItem struct {
	Title       string           `json:"title,omitempty"`
	Description string           `json:"description,omitempty"`
	Links       []string         `json:"links,omitempty"`
	Categories  []string         `json:"categories,omitempty"`
	FeedSource  string           `json:"source,omitempty"`
	Authors     []*gofeed.Person `json:"authors,omitempty"`
	Image       *gofeed.Image    `json:"image,omitempty"`
	Updated     *time.Time       `json:"updated,omitempty"`
	Published   *time.Time       `json:"published,omitempty"`
}
