package scraper

import (
	"fmt"

	"github.com/gocolly/colly/v2"
	"github.com/nschimek/nba-scraper/model"
	"github.com/nschimek/nba-scraper/parser"
)

const (
	basePlayerBodyElement = "body > div#wrap"
	playerInfoElement     = basePlayerBodyElement + " > div#info > div#meta > div:nth-child(2)"
)

type PlayerScraper struct {
	Colly        *colly.Collector     `Inject:""`
	PlayerParser *parser.PlayerParser `Inject:""`
	ScrapedData  []model.Player
}

func (s *PlayerScraper) GetData() interface{} {
	return s.ScrapedData
}

func (s *PlayerScraper) Scrape(urls ...string) {
	for _, url := range urls {
		player := s.parsePlayerPage(url)
		s.ScrapedData = append(s.ScrapedData, player)
	}

	fmt.Printf("%+v\n", s.ScrapedData)
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
