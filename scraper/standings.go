package scraper

import (
	"strconv"

	"github.com/gocolly/colly/v2"
	"github.com/nschimek/nba-scraper/parser"
)

type StandingsScraper struct {
	colly       colly.Collector
	season      int
	ScrapedData []parser.Standings
	urls        []string
	Errors      []error
	child       Scraper
	childUrls   map[string]string
}

func CreateStandingsScraper(c *colly.Collector, season int) StandingsScraper {
	urls := []string{getSeasonUrl(season)}
	return StandingsScraper{
		colly:     *c,
		season:    season,
		urls:      urls,
		childUrls: make(map[string]string),
	}
}

func (s *StandingsScraper) GetData() interface{} {
	return s.ScrapedData
}

func (s *StandingsScraper) AttachChild(scraper *Scraper) {
	s.child = *scraper
}

func (s *StandingsScraper) GetChild() Scraper {
	return s.child
}

func (s *StandingsScraper) GetChildUrls() []string {
	return urlsMapToArray(s.childUrls)
}

func getSeasonUrl(season int) string {
	return BaseHttp + "/" + baseLeaguesPath + "/NBA_" + strconv.Itoa(season) + "_standings.html"
}
