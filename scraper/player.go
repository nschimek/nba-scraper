package scraper

import (
	"fmt"

	"github.com/gocolly/colly/v2"
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
	Repository   *repository.SimpleRepository `Inject:""`
	ScrapedData  []model.Player
}

func (s *PlayerScraper) GetData() interface{} {
	return s.ScrapedData
}

func (s *PlayerScraper) Scrape(ids ...string) {
	for _, id := range ids {
		player := s.parsePlayerPage(id)
		s.ScrapedData = append(s.ScrapedData, player)
	}

	s.Repository.Upsert(s.ScrapedData, "players")
}

func (s *PlayerScraper) parsePlayerPage(id string) (player model.Player) {
	c := s.Colly.Clone()
	c.OnRequest(onRequestVisit)
	c.OnError(onError)

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
