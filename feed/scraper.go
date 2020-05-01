package feed

type Scraper interface {
	GetInitialFeed() chan Feed
	GetNextFeed() chan Feed
	GetPrevFeed() chan Feed
	GetFeedName() string
	GetPageIndex() int
}

// todo: add sitestr
type Feed struct {
	Title string
	Link  string
}
