package services

import (
	"errors"
	"log"
	"sync"

	"github.com/heyrovsky/rsscurator/pkg/content"
	"github.com/heyrovsky/rsscurator/pkg/rss"
	"go.uber.org/zap"
)

func FetchNewsItems() ([]content.NewsItemHashed, error) {

	if len(RssFeedReaders) == 0 {
		return nil, errors.New("no RSS feed readers provided")
	}

	newsItems := []content.NewsItemHashed{}
	var newsItemsChain = make(chan content.NewsItemHashed)
	var wg sync.WaitGroup

	for _, reader := range RssFeedReaders {
		wg.Add(1)
		go func(reader rss.RssReader) {
			defer wg.Done()
			items, err := reader.ReadNews()
			if err != nil {
				reader.Logger.Error("Error reading news", zap.Error(err))
				return
			}

			for _, item := range items {
				newsItemsChain <- item
			}
		}(reader)
	}

	go func() {
		wg.Wait()
		close(newsItemsChain)
	}()

	var mutex sync.Mutex
	for hashedItem := range newsItemsChain {
		if hashedItem.Hash == "" {
			log.Println("Skipping empty hashed news item")
			continue
		}

		mutex.Lock()
		newsItems = append(newsItems, hashedItem)
		mutex.Unlock()
	}

	return newsItems, nil
}
