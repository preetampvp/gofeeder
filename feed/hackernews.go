package feed

import (
	"fmt"

	"github.com/gocolly/colly"
)

func NewHackerNewsScraper() Scraper {
	return &hnScraper{
		feedBase:   "https://news.ycombinator.com/",
		newestPath: "newest",
	}
}

type hnScraper struct {
	feedBase   string
	newestPath string
	nextPath   string
	prevPaths  []string
	pageIndex  int
}

func (s *hnScraper) GetFeedName() string {
	return fmt.Sprintf("Hacker News Feed")
}

func (s *hnScraper) GetInitialFeed() chan Feed {
	s.pageIndex = 1
	s.prevPaths = make([]string, 0)
	s.prevPaths = append(s.prevPaths, s.newestPath)
	return s.scrapeFeed(fmt.Sprintf("%s%s", s.feedBase, s.newestPath))
}

func (s *hnScraper) GetNextFeed() chan Feed {
	s.pageIndex += 1
	if s.nextPath != "" && len(s.prevPaths) < s.pageIndex {
		s.prevPaths = append(s.prevPaths, s.nextPath)
	}
	return s.scrapeFeed(fmt.Sprintf("%s%s", s.feedBase, s.nextPath))
}

func (s *hnScraper) GetPrevFeed() chan Feed {
	if s.pageIndex > 1 {
		s.pageIndex -= 1
	}
	return s.scrapeFeed(fmt.Sprintf("%s%s", s.feedBase, s.prevPaths[s.pageIndex-1]))
}

func (s *hnScraper) scrapeFeed(url string) chan Feed {
	ch := make(chan Feed)

	go func() {
		collector := colly.NewCollector()
		defer close(ch)

		collector.OnHTML(".storylink", func(e *colly.HTMLElement) {

			title := e.Text
			url := e.Attr("href")

			ch <- Feed{Title: title, Link: url}
		})

		collector.OnHTML(".morelink", func(e *colly.HTMLElement) {

			url := e.Attr("href")
			s.nextPath = url
		})

		collector.Visit(url)
	}()

	return ch
}
