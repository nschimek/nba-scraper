package scraper

import (
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
	Colly            *colly.Collector             `Inject:""`
	PlayerParser     *parser.PlayerParser         `Inject:""`
	PlayerRepository *repository.PlayerRepository `Inject:""`
	ScrapedData      []model.Player
}

func (s *PlayerScraper) GetData() interface{} {
	return s.ScrapedData
}

func (s *PlayerScraper) Scrape(urls ...string) {
	for _, url := range urls {
		player := s.parsePlayerPage(url)
		s.ScrapedData = append(s.ScrapedData, player)
	}

	s.PlayerRepository.CreateBatch(s.ScrapedData)
}

func (s *PlayerScraper) parsePlayerPage(url string) (player model.Player) {
	c := s.Colly.Clone()
	c.OnRequest(onRequestVisit)

	player.ID = parser.ParseLastId(url)

	c.OnHTML(playerInfoElement, func(div *colly.HTMLElement) {
		s.PlayerParser.PlayerInfoBox(&player, div)
	})

	c.Visit(url)

	return
}
