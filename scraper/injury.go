package scraper

import (
	"github.com/gocolly/colly/v2"
	"github.com/nschimek/nba-scraper/core"
	"github.com/nschimek/nba-scraper/model"
	"github.com/nschimek/nba-scraper/parser"
	"github.com/nschimek/nba-scraper/repository"
)

type InjuryScraper struct {
	Config       *core.Config                                     `Inject:""`
	Colly        *colly.Collector                                 `Inject:""`
	InjuryParser *parser.InjuryParser                             `Inject:""`
	Repository   *repository.SimpleRepository[model.PlayerInjury] `Inject:""`
	ScrapedData  []model.PlayerInjury
	PlayerIds    map[string]struct{}
	TeamIds      map[string]struct{}
}

const (
	injuriesUrl              = BaseHttp + "/friv/injuries.fcgi"
	injuriesTableBaseElement = "body > div#wrap > div#content > div#all_injuries.table_wrapper > div#div_injuries.table_container > table > tbody"
)

func (s *InjuryScraper) Scrape() {
	c := core.CloneColly(s.Colly)
	s.PlayerIds = make(map[string]struct{})
	s.TeamIds = make(map[string]struct{})

	c.OnHTML(injuriesTableBaseElement, func(tbl *colly.HTMLElement) {
		for _, pi := range s.InjuryParser.InjuriesTable(tbl) {
			if !pi.HasErrors() {
				s.ScrapedData = append(s.ScrapedData, pi)
				s.PlayerIds[pi.PlayerId] = exists
				s.TeamIds[pi.TeamId] = exists
			} else {
				pi.LogErrors()
			}
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		core.Log.Error(NewScraperError(err, r.Request.URL.String()))
	})

	c.Visit(injuriesUrl)

	core.Log.WithField("injuries", len(s.ScrapedData)).Info("Finished scraping Player Injuries page!")
}

func (s *InjuryScraper) Persist() {
	if len(s.ScrapedData) > 0 {
		s.Repository.Upsert(s.ScrapedData, "Player Injuries")
	} else {
		core.Log.Warn("No Player Injuries scraped to persist! Skipping...")
	}
}
