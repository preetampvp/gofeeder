package ui

import (
	"fmt"

	"github.com/pkg/browser"
	"github.com/preetampvp/gofeeder/feed"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type FeedViewer interface {
	Show()
}

type feedLoader func() chan feed.Feed

// NewUi - description
func NewFeedViewer(scrapers []feed.Scraper) FeedViewer {
	viewer := &feedViewer{scrapers: scrapers, currentScraper: 0}
	return viewer
}

type feedViewer struct {
	scrapers       []feed.Scraper
	feed           []feed.Feed
	grid           *ui.Grid
	sourceList     *widgets.List
	feedList       *widgets.List
	infoText       *widgets.Paragraph
	currentScraper int
}

// Show - Show ui
func (f *feedViewer) Show() {
	if err := ui.Init(); err != nil {
		fmt.Printf("failed to initialize termui: %v", err)
		return
	}
	defer ui.Close()

	f.initSourceView()
	f.initFeedView()
	f.initGrid()
	f.render()
	f.loadFeed(f.scrapers[f.currentScraper].GetInitialFeed)
	f.initEventsPolling()
}

func (f *feedViewer) shortcutsText() string {
	return "[ Shortcuts   ](fg:white,bg:black) [ Enter ](fg:black)[ Open ](fg:black,bg:green) " +
		"[ k ](fg:black)[ Feed Up ](fg:black,bg:green) " +
		"[ j ](fg:black)[ Feed Down ](fg:black,bg:green) " +
		"[ K ](fg:black)[  Source Up ](fg:black,bg:green) " +
		"[ J ](fg:black)[ Source Down ](fg:black,bg:green) " +
		"[ l ](fg:black)[ Load Feed ](fg:black,bg:green) " +
		"[ n ](fg:black)[ Next ](fg:black,bg:green) " +
		"[ p ](fg:black)[ Prev ](fg:black,bg:green) " +
		"[ r ](fg:black)[ Refresh ](fg:black,bg:green) " +
		"[ q ](fg:black)[ Quit ](fg:black,bg:green) "
}

func (f *feedViewer) loadFeed(loader feedLoader) {
	f.infoText.Text = "loading feed...."
	f.render()
	f.feed = make([]feed.Feed, 0)
	f.feedList.Rows = make([]string, 0)
	for item := range loader() {
		f.feed = append(f.feed, item)
		f.feedList.Rows = append(f.feedList.Rows, item.Title)
	}
	f.feedList.Title = fmt.Sprintf("  %s [%d]  ", f.scrapers[f.currentScraper].GetFeedName(), f.scrapers[f.currentScraper].GetPageIndex())
	f.infoText.Text = f.shortcutsText()
	f.feedList.SelectedRow = 0
	f.render()
}

func (f *feedViewer) render() {
	ui.Clear()
	ui.Render(f.grid)
}

func (f *feedViewer) initFeedView() {
	f.feedList = widgets.NewList()
	f.feedList.TextStyle = ui.NewStyle(ui.ColorWhite)
	f.feedList.BorderStyle.Fg = ui.ColorMagenta
	f.feedList.SelectedRowStyle.Fg = ui.ColorBlack
	f.feedList.SelectedRowStyle.Bg = ui.ColorWhite
}

func (f *feedViewer) initSourceView() {
	f.sourceList = widgets.NewList()
	f.sourceList.TextStyle = ui.NewStyle(ui.ColorWhite)
	f.sourceList.BorderStyle.Fg = ui.ColorMagenta
	f.sourceList.SelectedRowStyle.Fg = ui.ColorBlack
	f.sourceList.SelectedRowStyle.Bg = ui.ColorWhite
	f.sourceList.Title = "  Sources  "
	f.sourceList.Rows = make([]string, 0)
	for _, s := range f.scrapers {
		f.sourceList.Rows = append(f.sourceList.Rows, s.GetFeedName())
	}
}

func (f *feedViewer) initGrid() {
	f.grid = ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	f.grid.SetRect(0, 0, termWidth, termHeight-1)

	f.infoText = widgets.NewParagraph()
	f.infoText.WrapText = true
	f.infoText.Border = false
	f.infoText.TextStyle = ui.Style{Modifier: ui.ModifierBold, Bg: ui.ColorWhite}
	f.infoText.Text = "initiating..."

	f.grid.Set(ui.NewRow(0.9, ui.NewCol(0.2, f.sourceList), ui.NewCol(0.8, f.feedList)), ui.NewRow(0.1, ui.NewCol(1.0, f.infoText)))
}

func (f *feedViewer) openArticle() {
	index := f.feedList.SelectedRow
	if len(f.feed) > index {
		browser.Stderr = nil
		browser.Stdout = nil
		_ = browser.OpenURL(f.feed[index].Link)
		f.render()
	}
}
