package scraper

import (
	"fmt"

	"github.com/gocolly/colly/v2"
	"github.com/nschimek/nba-scraper/core"
	"github.com/nschimek/nba-scraper/model"
	"github.com/nschimek/nba-scraper/parser"
	"github.com/nschimek/nba-scraper/repository"
)

const (
	teams                      = "teams"
	teamBaseBodyElement        = "body div#wrap"
	teamInfoElement            = teamBaseBodyElement + " div#info div#meta"
	teamRosterTableElement     = teamBaseBodyElement + " div#all_roster > div#div_roster > table#roster > tbody"
	teamSalaryTableElementBase = teamBaseBodyElement + " div#all_salaries2"
	teamSalaryTableElement     = "div#div_salaries2 > table#salaries2 > tbody"
)

type TeamScraper struct {
	Config      *core.Config               `Inject:""`
	Colly       *colly.Collector           `Inject:""`
	TeamParser  *parser.TeamParser         `Inject:""`
	Repository  *repository.TeamRepository `Inject:""`
	ScrapedData []model.Team
	PlayerIds   map[string]bool
}

func (s *TeamScraper) Scrape(idMap map[string]bool) {
	s.PlayerIds = make(map[string]bool)

	for _, id := range core.IdMapToArray(idMap) {
		team := s.parseTeamPage(id)
		s.ScrapedData = append(s.ScrapedData, team)
	}

	core.Log.WithField("teams", len(s.ScrapedData)).Info("Successfully scraped Team page(s)!")
}

func (s *TeamScraper) Persist() {
	s.Repository.UpsertTeams(s.ScrapedData)
}

func (s *TeamScraper) parseTeamPage(id string) (team model.Team) {
	c := core.CloneColly(s.Colly)

	team.ID = id
	team.Season = s.Config.Season

	c.OnHTML(teamInfoElement, func(div *colly.HTMLElement) {
		s.TeamParser.TeamInfoBox(&team, div)
	})

	c.OnHTML(teamRosterTableElement, func(tbl *colly.HTMLElement) {
		s.TeamParser.TeamPlayerTable(&team, tbl)

		for _, p := range team.TeamPlayers {
			s.PlayerIds[p.PlayerId] = true
		}
	})

	c.OnHTML(teamSalaryTableElementBase, func(div *colly.HTMLElement) {
		tbl, _ := transformHtmlElement(div, teamSalaryTableElement, removeCommentsSyntax)
		s.TeamParser.TeamSalariesTable(&team, tbl)
	})

	c.Visit(s.getUrl(id))

	return team
}

func (t *TeamScraper) getUrl(id string) string {
	return fmt.Sprintf(teamIdPage, id, t.Config.Season)
}
