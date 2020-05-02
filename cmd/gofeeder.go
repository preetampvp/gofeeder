package main

import (
	"github.com/preetampvp/gofeeder/config"
	"github.com/preetampvp/gofeeder/feed"
	"github.com/preetampvp/gofeeder/ui"
)

func main() {
	config := config.GetConfig()
	scrapers := feed.GetScrapers(config)
	ui := ui.NewFeedViewer(scrapers)
	ui.Show()
}
