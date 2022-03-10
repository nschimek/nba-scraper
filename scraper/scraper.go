package scraper

type Scraper interface {
	Scrape(urls []string)
	getData() interface{}
	attachChild(scraper *Scraper)
}
