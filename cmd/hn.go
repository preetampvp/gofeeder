package main

import (
	_ "github.com/pkg/browser"
	"github.com/preetampvp/gofeeder/feed"
	"github.com/preetampvp/gofeeder/ui"
)

func main() {
	feed := feed.NewHackerNewsScraper()
	ui := ui.NewFeedViewer(feed)
	ui.Show()
}
