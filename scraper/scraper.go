package scraper

import "github.com/gocolly/colly/v2"

type Scraper interface {
	Scrape(urls ...string)
	GetData() interface{}
	AttachChild(scraper *Scraper)
	GetChild() Scraper
	GetChildUrls() []string
}

func scrapeChild(s Scraper) {
	if s.GetChild() != nil && len(s.GetChildUrls()) > 0 {
		s.GetChild().Scrape(s.GetChildUrls()...)
	}
}

func urlsMapToArray(urlsMap map[string]string) (urlsArray []string) {
	for _, url := range urlsMap {
		urlsArray = append(urlsArray, url)
	}
	return
}

func parseLink(e *colly.HTMLElement) string {
	return e.ChildAttr("a", "href")
}
