package feed

import (
	"fmt"

	"github.com/gocolly/colly"
	"github.com/preetampvp/gofeeder/config"
)

func NewScraper(config config.FeedConfig) Scraper {
	return &scraper{config: config}
}

type scraper struct {
	config    config.FeedConfig
	nextPath  string
	prevPaths []string
	pageIndex int
}

func (s *scraper) GetFeedName() string {
	return fmt.Sprintf("%s", s.config.Name)
}

func (s *scraper) GetPageIndex() int {
	return s.pageIndex
}

func (s *scraper) GetInitialFeed() chan Feed {
	s.pageIndex = 1
	s.prevPaths = make([]string, 0)
	initialPath := fmt.Sprintf("%s%s", s.config.UrlBase, s.config.UrlPath)
	s.prevPaths = append(s.prevPaths, initialPath)
	return s.scrapeFeed(initialPath)
}

func (s *scraper) GetNextFeed() chan Feed {
	if s.nextPath == "" {
		return s.scrapeFeed(fmt.Sprintf("%s", s.prevPaths[len(s.prevPaths)-1]))
	}

	s.pageIndex += 1
	if s.nextPath != "" && len(s.prevPaths) < s.pageIndex {
		s.prevPaths = append(s.prevPaths, s.nextPath)
	}

	if s.config.NextLinkIsRelative {
		return s.scrapeFeed(fmt.Sprintf("%s%s", s.config.UrlBase, s.nextPath))
	}

	return s.scrapeFeed(fmt.Sprintf("%s", s.nextPath))
}

func (s *scraper) GetPrevFeed() chan Feed {
	if s.pageIndex > 1 {
		s.pageIndex -= 1
	}

	if s.config.NextLinkIsRelative {
		return s.scrapeFeed(fmt.Sprintf("%s%s", s.config.UrlBase, s.prevPaths[s.pageIndex-1]))
	}

	return s.scrapeFeed(fmt.Sprintf("%s", s.prevPaths[s.pageIndex-1]))
}

func (s *scraper) scrapeFeed(url string) chan Feed {
	ch := make(chan Feed)

	go func() {
		collector := colly.NewCollector()
		defer close(ch)

		collector.OnHTML(s.config.LinkSelector, func(e *colly.HTMLElement) {

			title := e.Text
			url := e.Attr("href")

			ch <- Feed{Title: title, Link: url}
		})

		if s.config.NextPageSelector != "" {
			collector.OnHTML(s.config.NextPageSelector, func(e *colly.HTMLElement) {

				url := e.Attr("href")
				s.nextPath = url
			})
		}

		collector.Visit(url)
	}()

	return ch
}
