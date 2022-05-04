package scraper

import (
	"github.com/gocolly/colly/v2"
	"github.com/nschimek/nba-scraper/parser"
)

type InjuriesScraper struct {
	colly       colly.Collector
	season      int
	ScrapedData []parser.Injuries
	urls        []string
	Errors      []error
	child       Scraper
	childUrls   map[string]string
}

const (
	injuriesUrl              = BaseHttp + "/friv/injuries.fcgi"
	injuriesTableBaseElement = "body > div#wrap > div#content > div#all_injuries.table_wrapper > div#div_injuries.table_container > table > tbody"
)

func CreateInjuriesScraper(c *colly.Collector, season int) InjuriesScraper {
	return InjuriesScraper{
		colly:     *c,
		season:    season,
		urls:      []string{injuriesUrl},
		childUrls: make(map[string]string),
	}
}

func (s *InjuriesScraper) GetData() interface{} {
	return s.ScrapedData
}

func (s *InjuriesScraper) AttachChild(scraper *Scraper) {
	s.child = *scraper
}

func (s *InjuriesScraper) GetChild() Scraper {
	return s.child
}

func (s *InjuriesScraper) GetChildUrls() []string {
	return urlsMapToArray(s.childUrls)
}

func (s *InjuriesScraper) Scrape(urls ...string) {
	c := s.colly.Clone()
	c.OnRequest(onRequestVisit)
	s.urls = append(s.urls, urls...)

	c.OnHTML(injuriesTableBaseElement, func(tbl *colly.HTMLElement) {
		for _, pi := range parser.InjuriesTable(tbl, s.season) {
			s.ScrapedData = append(s.ScrapedData, pi)
			// this scraper only captures the players as the child URLs, so Standings should be run first
			s.childUrls[pi.PlayerId] = tbl.Request.AbsoluteURL(pi.PlayerUrl)
		}
	})

	for _, url := range s.urls {
		c.Visit(url)
	}
}
