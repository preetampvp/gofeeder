package ui

import (
	ui "github.com/gizak/termui/v3"
)

func (f *feedViewer) initEventsPolling() {
	uiEvents := ui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			{
				switch e.ID {
				case "q", "<C-c>":
					return
				case "k":
					f.feedList.ScrollUp()
				case "j":
					f.feedList.ScrollDown()
				case "K":
					f.sourceList.ScrollUp()
				case "J":
					f.sourceList.ScrollDown()
				case "l", "L":
					f.currentScraper = f.sourceList.SelectedRow
					f.loadFeed(f.scrapers[f.currentScraper].GetInitialFeed)
				case "n":
					f.loadFeed(f.scrapers[f.currentScraper].GetNextFeed)
				case "p":
					f.loadFeed(f.scrapers[f.currentScraper].GetPrevFeed)
				case "r":
					f.loadFeed(f.scrapers[f.currentScraper].GetInitialFeed)
				case "<Resize>":
					payload := e.Payload.(ui.Resize)
					f.grid.SetRect(0, 0, payload.Width, payload.Height-1)
				case "<Enter>":
					f.openArticle()
				}

				f.render()
			}
		}
	}
}
