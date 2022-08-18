package scraper

import (
	"github.com/gocolly/colly/v2"
	"github.com/nschimek/nba-scraper/core"
	"github.com/nschimek/nba-scraper/model"
	"github.com/nschimek/nba-scraper/parser"
	"github.com/nschimek/nba-scraper/repository"
)

type InjuryScraper struct {
	Config       *core.Config                 `Inject:""`
	Colly        *colly.Collector             `Inject:""`
	InjuryParser *parser.InjuryParser         `Inject:""`
	Repository   *repository.SimpleRepository `Inject:""`
	ScrapedData  []model.PlayerInjury
	PlayerIds    map[string]struct{}
}

const (
	injuriesUrl              = BaseHttp + "/friv/injuries.fcgi"
	injuriesTableBaseElement = "body > div#wrap > div#content > div#all_injuries.table_wrapper > div#div_injuries.table_container > table > tbody"
)

// we this has a static URL, so we have no use for the IDs...but leaving it for a future interface
func (s *InjuryScraper) Scrape(pageIds ...string) {
	s.PlayerIds = make(map[string]struct{})
	c := s.Colly.Clone()
	c.OnRequest(onRequestVisit)
	c.OnError(onError)

	c.OnHTML(injuriesTableBaseElement, func(tbl *colly.HTMLElement) {
		for _, pi := range s.InjuryParser.InjuriesTable(tbl) {
			s.ScrapedData = append(s.ScrapedData, pi)
			s.PlayerIds[pi.PlayerId] = exists
		}
	})

	c.Visit(injuriesUrl)

	core.Log.WithField("injuries", len(s.ScrapedData)).Info("Successfully scraped Player Injuries page!")
}

func (s *InjuryScraper) Persist() {
	s.Repository.Upsert(s.ScrapedData, "player_injuries")
}
