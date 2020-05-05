package feed

import (
	"fmt"
	"net/url"
	"path"
	"strings"

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

type buildParam struct {
	link       string
	forceBuild bool
}

func (s *scraper) buildUrl(param buildParam) string {

	if !strings.HasPrefix(param.link, "http") || param.forceBuild {
		url, _ := url.Parse(s.config.UrlBase)
		url.Path = path.Join(url.Path, param.link)
		url, _ = url.Parse(url.Path)
		return url.String()
	}

	return param.link
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
	initialPath := s.buildUrl(buildParam{link: s.config.UrlPath, forceBuild: true})
	s.prevPaths = append(s.prevPaths, initialPath)
	return s.scrapeFeed(initialPath)
}

func (s *scraper) GetNextFeed() chan Feed {
	if s.nextPath == "" {
		return s.scrapeFeed(fmt.Sprintf("%s", s.prevPaths[len(s.prevPaths)-1]))
	}

	s.pageIndex += 1
	if len(s.prevPaths) < s.pageIndex {
		s.prevPaths = append(s.prevPaths, s.nextPath)
	}

	return s.scrapeFeed(s.nextPath)
}

func (s *scraper) GetPrevFeed() chan Feed {
	if s.pageIndex > 1 {
		s.pageIndex -= 1
	}

	return s.scrapeFeed(fmt.Sprintf("%s", s.prevPaths[s.pageIndex-1]))
}

func (s *scraper) scrapeFeed(url string) chan Feed {
	ch := make(chan Feed)
	s.nextPath = ""

	go func() {
		collector := colly.NewCollector()
		defer close(ch)

		collector.OnHTML(s.config.LinkSelector, func(e *colly.HTMLElement) {
			title := e.Text
			url := s.buildUrl(buildParam{link: e.Attr("href")})
			ch <- Feed{Title: title, Link: url}
		})

		if s.config.NextPageSelector != "" {
			collector.OnHTML(s.config.NextPageSelector, func(e *colly.HTMLElement) {
				url := e.Attr("href")
				s.nextPath = url
				if s.nextPath != "" {
					s.nextPath = s.buildUrl(buildParam{link: url})
				}
			})
		}

		collector.Visit(url)
	}()

	return ch
}
