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
	Config           *core.Config                             `Inject:""`
	Colly            *colly.Collector                         `Inject:""`
	TeamParser       *parser.TeamParser                       `Inject:""`
	TeamRepository   *repository.TeamRepository               `Inject:""`
	SimpleRepository *repository.SimpleRepository[model.Team] `Inject:""`
	ScrapedData      []model.Team
	PlayerIds        map[string]struct{}
}

func (s *TeamScraper) Scrape(idMap map[string]struct{}) {
	core.Log.WithField("ids", len(idMap)).Info("Got Team ID(s) to scrape")
	if s.Config.Suppression.Team > 0 {
		s.suppressRecent(idMap)
	}
	s.scrape(idMap)
}

func (s *TeamScraper) scrape(idMap map[string]struct{}) {
	core.Log.WithField("ids", len(idMap)).Info("Scraping Team ID(s)...")
	s.PlayerIds = make(map[string]struct{})

	for _, id := range core.IdMapToArray(idMap) {
		team := s.parseTeamPage(id)
		if !team.HasErrors() {
			s.ScrapedData = append(s.ScrapedData, team)
		} else {
			team.LogErrors()
		}
	}

	core.Log.WithField("teams", len(s.ScrapedData)).Info("Finished scraping Team page(s)!")
}

func (s *TeamScraper) Persist() {
	if len(s.ScrapedData) > 0 {
		s.TeamRepository.UpsertTeams(s.ScrapedData)
	} else {
		core.Log.Info("No Teams scraped to persist! Skipping...")
	}
}

func (s *TeamScraper) suppressRecent(idMap map[string]struct{}) {
	core.Log.Infof("Checking for teams updated in last %d days (for suppression)...", s.Config.Suppression.Team)
	ids, _ := s.SimpleRepository.GetRecentlyUpdated(s.Config.Suppression.Team, core.IdMapToArray(idMap), "Team")
	if len(ids) > 0 {
		core.SuppressIdMap(idMap, ids)
	}
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
			s.PlayerIds[p.PlayerId] = exists
		}
	})

	c.OnHTML(teamSalaryTableElementBase, func(div *colly.HTMLElement) {
		tbl, _ := transformHtmlElement(div, teamSalaryTableElement, removeCommentsSyntax)
		s.TeamParser.TeamSalariesTable(&team, tbl)

		for _, p := range team.TeamPlayerSalaries {
			s.PlayerIds[p.PlayerId] = exists
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		team.CaptureError(NewScraperError(err, r.Request.URL.String()))
	})

	c.Visit(s.getUrl(id))

	return team
}

func (t *TeamScraper) getUrl(id string) string {
	return fmt.Sprintf(teamIdPage, id, t.Config.Season)
}
