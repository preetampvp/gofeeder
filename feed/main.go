package feed

import (
	"github.com/preetampvp/gofeeder/config"
)

type Scraper interface {
	GetInitialFeed() chan Feed
	GetNextFeed() chan Feed
	GetPrevFeed() chan Feed
	GetFeedName() string
	GetPageIndex() int
}

type Feed struct {
	Title string
	Link  string
}

// GetScrappers - Builds a list of scrapers based on config
func GetScrapers(config *config.Config) []Scraper {
	scrapers := make([]Scraper, 0)
	for _, feedItem := range config.Feeds {
		scrapers = append(scrapers, NewScraper(feedItem))
	}

	return scrapers
}
