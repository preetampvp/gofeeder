package main

import (
	_ "github.com/pkg/browser"
	"github.com/preetampvp/gofeeder/feed"
	"github.com/preetampvp/gofeeder/ui"
)

func main() {
	hnScraper := feed.NewHackerNewsScraper()
	awsScraper := feed.NewAwsScraper()
	dcScraper := feed.NewDCScraper()
	k8sScrapper := feed.NewK8sScraper()
	scrapers := []feed.Scraper{hnScraper, awsScraper, dcScraper, k8sScrapper}
	ui := ui.NewFeedViewer(scrapers)
	ui.Show()
}
