package feed

import (
	"fmt"

	"github.com/gocolly/colly"
)

func NewDCScraper() Scraper {
	return &dcScraper{
		feedBase: "https://dave.cheney.net/category/golang",
	}
}

type dcScraper struct {
	feedBase  string
	nextPath  string
	prevPaths []string
	pageIndex int
}

func (s *dcScraper) GetFeedName() string {
	return fmt.Sprintf("Dave Cheney Blog")
}

func (s *dcScraper) GetPageIndex() int {
	return s.pageIndex
}

func (s *dcScraper) GetInitialFeed() chan Feed {
	s.pageIndex = 1
	s.prevPaths = make([]string, 0)
	s.prevPaths = append(s.prevPaths, s.feedBase)
	return s.scrapeFeed(fmt.Sprintf("%s", s.feedBase))
}

func (s *dcScraper) GetNextFeed() chan Feed {
	if s.nextPath == "" {
		return s.scrapeFeed(fmt.Sprintf("%s", s.prevPaths[len(s.prevPaths)-1]))
	}

	s.pageIndex += 1
	if len(s.prevPaths) < s.pageIndex {
		s.prevPaths = append(s.prevPaths, s.nextPath)
	}
	return s.scrapeFeed(fmt.Sprintf("%s", s.nextPath))
}

func (s *dcScraper) GetPrevFeed() chan Feed {
	if s.pageIndex > 1 {
		s.pageIndex -= 1
	}
	return s.scrapeFeed(fmt.Sprintf("%s", s.prevPaths[s.pageIndex-1]))
}

func (s *dcScraper) scrapeFeed(url string) chan Feed {
	ch := make(chan Feed)

	go func() {
		collector := colly.NewCollector()
		defer close(ch)

		collector.OnHTML(".entry-title", func(e *colly.HTMLElement) {
			title := e.ChildText("a")
			url := e.ChildAttr("a", "href")
			ch <- Feed{Title: title, Link: url}
		})

		collector.OnHTML(".nav-previous", func(e *colly.HTMLElement) {
			url := e.ChildAttr("a:first-child", "href")
			s.nextPath = url
		})

		collector.Visit(url)
	}()

	return ch
}
