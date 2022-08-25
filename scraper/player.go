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
	basePlayerBodyElement = "body > div#wrap"
	playerInfoElement     = basePlayerBodyElement + " > div#info > div#meta > div:nth-child(2)"
)

type PlayerScraper struct {
	Colly        *colly.Collector             `Inject:""`
	PlayerParser *parser.PlayerParser         `Inject:""`
	Repository   *repository.PlayerRepository `Inject:""`
	ScrapedData  []model.Player
}

func (s *PlayerScraper) Scrape(idMap map[string]bool) {
	core.Log.WithField("ids", len(idMap)).Info("Got Player IDs to Scrape, suppressing recent...")
	s.Repository.SuppressRecentlyUpdated(30, idMap)

	for _, id := range IdMapToArray(idMap) {
		player := s.parsePlayerPage(id)
		s.ScrapedData = append(s.ScrapedData, player)
	}

	core.Log.WithField("players", len(s.ScrapedData)).Info("Successfully scraped Player page(s)!")
}

func (s *PlayerScraper) Persist() {
	s.Repository.Upsert(s.ScrapedData, "players")
}

func (s *PlayerScraper) suppressRecent(idMap map[string]bool) {

}

func (s *PlayerScraper) parsePlayerPage(id string) (player model.Player) {
	c := core.CloneColly(s.Colly)

	player.ID = id

	c.OnHTML(playerInfoElement, func(div *colly.HTMLElement) {
		s.PlayerParser.PlayerInfoBox(&player, div)
	})

	c.Visit(s.getUrl(id))

	return
}

func (*PlayerScraper) getUrl(id string) string {
	return fmt.Sprintf(playerIdPage, id[0:1], id)
}
