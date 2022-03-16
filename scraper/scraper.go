package scraper

import (
	"strings"

	"github.com/gocolly/colly/v2"
)

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

func parseGameId(link string) string {
	return strings.Replace(strings.Split(link, "/")[2], ".html", "", 1)
}

func parseTeamId(link string) string {
	return strings.Split(link, "/")[2]
}
