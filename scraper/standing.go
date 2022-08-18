package scraper

import (
	"fmt"

	"github.com/gocolly/colly/v2"
	"github.com/nschimek/nba-scraper/core"
	"github.com/nschimek/nba-scraper/model"
	"github.com/nschimek/nba-scraper/parser"
	"github.com/nschimek/nba-scraper/repository"
)

type StandingScraper struct {
	Config      *core.Config                 `Inject:""`
	Colly       *colly.Collector             `Inject:""`
	Parser      *parser.StandingParser       `Inject:""`
	Repository  *repository.SimpleRepository `Inject:""`
	ScrapedData []model.TeamStanding
	TeamIds     map[string]struct{}
}

const (
	expandedStandingsTableElementBase = "body > div#wrap > div#content > div#all_expanded_standings.table_wrapper"
	expandedStandingsTableElement     = "div#div_expanded_standings.table_container > table > tbody"
)

func (s *StandingScraper) GetData() []model.TeamStanding {
	return s.ScrapedData
}

// we this has a static URL, so we have no use for the IDs...but leaving it for a future interface
func (s *StandingScraper) Scrape(pageIds ...string) {
	s.TeamIds = make(map[string]struct{})
	c := s.Colly.Clone()
	c.OnRequest(onRequestVisit)
	c.OnError(onError)

	c.OnHTML(expandedStandingsTableElementBase, func(div *colly.HTMLElement) {
		tbl, _ := transformHtmlElement(div, expandedStandingsTableElement, removeCommentsSyntax)
		for _, ts := range s.Parser.StandingsTable(tbl) {
			s.ScrapedData = append(s.ScrapedData, ts)
			s.TeamIds[ts.TeamId] = exists
		}
	})

	c.Visit(s.getUrl())

	core.Log.WithField("standings", len(s.ScrapedData)).Info("Successfully scraped Team Standings page!")
}

func (s *StandingScraper) Persist() {
	s.Repository.Upsert(s.ScrapedData, "player_standings")
}

func (t *StandingScraper) getUrl() string {
	return fmt.Sprintf(standingPage, t.Config.Season)
}
