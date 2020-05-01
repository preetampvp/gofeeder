package feed

import (
	"fmt"

	"github.com/gocolly/colly"
)

func NewAwsScraper() Scraper {
	return &awsScraper{
		feedBase: "https://aws.amazon.com/blogs/aws/",
	}
}

type awsScraper struct {
	feedBase  string
	nextPath  string
	prevPaths []string
	pageIndex int
}

func (s *awsScraper) GetFeedName() string {
	return fmt.Sprintf("AWS Blog")
}

func (s *awsScraper) GetPageIndex() int {
	return s.pageIndex
}

func (s *awsScraper) GetInitialFeed() chan Feed {
	s.pageIndex = 1
	s.prevPaths = make([]string, 0)
	s.prevPaths = append(s.prevPaths, s.feedBase)
	return s.scrapeFeed(fmt.Sprintf("%s", s.feedBase))
}

func (s *awsScraper) GetNextFeed() chan Feed {
	if s.nextPath == "" {
		return s.scrapeFeed(fmt.Sprintf("%s", s.prevPaths[len(s.prevPaths)-1]))
	}

	s.pageIndex += 1
	if s.nextPath != "" && len(s.prevPaths) < s.pageIndex {
		s.prevPaths = append(s.prevPaths, s.nextPath)
	}
	return s.scrapeFeed(fmt.Sprintf("%s", s.nextPath))
}

func (s *awsScraper) GetPrevFeed() chan Feed {
	if s.pageIndex > 1 {
		s.pageIndex -= 1
	}
	return s.scrapeFeed(fmt.Sprintf("%s", s.prevPaths[s.pageIndex-1]))
}

func (s *awsScraper) scrapeFeed(url string) chan Feed {
	ch := make(chan Feed)

	go func() {
		collector := colly.NewCollector()
		defer close(ch)

		collector.OnHTML(".blog-post-title", func(e *colly.HTMLElement) {
			title := e.ChildText("a")
			url := e.ChildAttr("a", "href")
			ch <- Feed{Title: title, Link: url}
		})

		collector.OnHTML(".blog-pagination", func(e *colly.HTMLElement) {
			url := e.ChildAttr("a:first-child", "href")
			s.nextPath = url
		})

		collector.Visit(url)
	}()

	return ch
}
