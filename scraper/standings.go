package scraper

import (
	"fmt"
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

const (
	expandedStandingsTableElementBase = "body > div#wrap > div#content > div#all_expanded_standings.table_wrapper"
	expandedStandingsTableElement     = "div#div_expanded_standings.table_container > table > tbody"
)

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

func (s *StandingsScraper) Scrape(urls ...string) {
	c := s.colly.Clone()
	c.OnRequest(onRequestVisit)
	s.urls = append(s.urls, urls...)

	c.OnHTML(expandedStandingsTableElementBase, func(div *colly.HTMLElement) {
		tbl, _ := transformHtmlElement(div, expandedStandingsTableElement, removeCommentsSyntax)
		for _, ps := range parser.StandingsTable(tbl, s.season) {
			s.ScrapedData = append(s.ScrapedData, ps)
			s.childUrls[ps.TeamId] = tbl.Request.AbsoluteURL(ps.TeamUrl)
		}
		fmt.Printf("%+v\n", s.ScrapedData)
		fmt.Println(s.childUrls)
	})

	for _, url := range s.urls {
		c.Visit(url)
	}
}

func getSeasonUrl(season int) string {
	return BaseHttp + "/" + baseLeaguesPath + "/NBA_" + strconv.Itoa(season) + "_standings.html"
}
