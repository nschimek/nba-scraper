package scraper

import (
	"fmt"
	"strconv"

	"github.com/gocolly/colly/v2"
	"github.com/nschimek/nba-scraper/parser"
)

const (
	teamBaseBodyElement = "body div#wrap"
	teamInfoElement     = teamBaseBodyElement + " div#info div#meta"
)

type TeamScraper struct {
	colly       colly.Collector
	ScrapedData []parser.Team
	Errors      []error
	child       Scraper
	childUrls   map[string]string
}

func CreateTeamScraper(c *colly.Collector) TeamScraper {
	return TeamScraper{
		colly:     *c,
		childUrls: make(map[string]string),
	}
}

func (s *TeamScraper) GetData() interface{} {
	return s.ScrapedData
}

func (s *TeamScraper) AttachChild(scraper *Scraper) {
	s.child = *scraper
}

func (s *TeamScraper) GetChild() Scraper {
	return s.child
}

func (s *TeamScraper) GetChildUrls() []string {
	return urlsMapToArray(s.childUrls)
}

func (s *TeamScraper) Scrape(urls ...string) {

	for _, url := range urls {
		team := s.parseTeamPage(url)
		s.ScrapedData = append(s.ScrapedData, team)
	}

	fmt.Printf("%+v\n", s.ScrapedData)

	scrapeChild(s)
}

func (s *TeamScraper) parseTeamPage(url string) (team parser.Team) {
	c := s.colly.Clone()

	team.Id = parser.ParseTeamId(url)
	team.Season, _ = strconv.Atoi(parser.ParseLastId(url))

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Parsing: " + r.URL.String())
	})

	c.OnHTML(teamInfoElement, func(div *colly.HTMLElement) {
		team.Name = div.ChildText("div:nth-child(2) > h1 > span:nth-child(2)")
	})

	c.Visit(url)

	return team
}
