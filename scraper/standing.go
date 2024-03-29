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
	Config      *core.Config                                     `Inject:""`
	Colly       *colly.Collector                                 `Inject:""`
	Parser      *parser.StandingParser                           `Inject:""`
	Repository  *repository.SimpleRepository[model.TeamStanding] `Inject:""`
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

func (s *StandingScraper) Scrape() {
	c := core.CloneColly(s.Colly)
	s.TeamIds = make(map[string]struct{})

	c.OnHTML(expandedStandingsTableElementBase, func(div *colly.HTMLElement) {
		tbl, _ := transformHtmlElement(div, expandedStandingsTableElement, removeCommentsSyntax)
		for _, ts := range s.Parser.StandingsTable(tbl) {
			if !ts.HasErrors() {
				s.ScrapedData = append(s.ScrapedData, ts)
				s.TeamIds[ts.TeamId] = exists
			} else {
				ts.LogErrors()
			}
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		core.Log.Error(NewScraperError(err, r.Request.URL.String()))
	})

	c.Visit(s.getUrl())

	core.Log.WithField("standings", len(s.ScrapedData)).Info("Finished scraping Team Standings page!")
}

func (s *StandingScraper) Persist() {
	if len(s.ScrapedData) > 0 {
		s.Repository.Upsert(s.ScrapedData, "Team Standing")
	} else {
		core.Log.Warn("No Team Standings scraped to persist, skipping...")
	}
}

func (s *StandingScraper) getUrl() string {
	return fmt.Sprintf(standingPage, s.Config.Season)
}
