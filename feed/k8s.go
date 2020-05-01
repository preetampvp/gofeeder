package feed

import (
	"fmt"

	"github.com/gocolly/colly"
)

func NewK8sScraper() Scraper {
	return &k8sScraper{
		feedBase: "https://kubernetes.io/blog/",
	}
}

type k8sScraper struct {
	feedBase  string
	nextPath  string
	prevPaths []string
	pageIndex int
}

func (s *k8sScraper) GetFeedName() string {
	return fmt.Sprintf("Kubernetes Blog")
}

func (s *k8sScraper) GetPageIndex() int {
	return s.pageIndex
}

func (s *k8sScraper) GetInitialFeed() chan Feed {
	s.pageIndex = 1
	s.prevPaths = make([]string, 0)
	s.prevPaths = append(s.prevPaths, s.feedBase)
	return s.scrapeFeed(fmt.Sprintf("%s", s.feedBase))
}

func (s *k8sScraper) GetNextFeed() chan Feed {
	if s.nextPath == "" {
		return s.scrapeFeed(fmt.Sprintf("%s", s.prevPaths[len(s.prevPaths)-1]))
	}

	s.pageIndex += 1
	if s.nextPath != "" && len(s.prevPaths) < s.pageIndex {
		s.prevPaths = append(s.prevPaths, s.nextPath)
	}
	return s.scrapeFeed(fmt.Sprintf("%s", s.nextPath))
}

func (s *k8sScraper) GetPrevFeed() chan Feed {
	if s.pageIndex > 1 {
		s.pageIndex -= 1
	}
	return s.scrapeFeed(fmt.Sprintf("%s", s.prevPaths[s.pageIndex-1]))
}

func (s *k8sScraper) scrapeFeed(url string) chan Feed {
	ch := make(chan Feed)

	go func() {
		collector := colly.NewCollector()
		defer close(ch)

		collector.OnHTML(".list-group-item", func(e *colly.HTMLElement) {
			title := e.ChildText("a")
			url := e.ChildAttr("a", "href")
			ch <- Feed{Title: title, Link: url}
		})

		// this will return an empty stirng as scraping a single page
		collector.OnHTML(".blog-pagination", func(e *colly.HTMLElement) {
			url := e.ChildAttr("a:first-child", "href")
			s.nextPath = url
		})

		collector.Visit(url)
	}()

	return ch
}
