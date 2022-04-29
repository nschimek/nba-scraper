package scraper

import (
	"fmt"

	"github.com/gocolly/colly/v2"
	"github.com/nschimek/nba-scraper/parser"
)

const (
	basePlayerBodyElement = "body > div#wrap"
	playerInfoElement     = basePlayerBodyElement + " > div#info > div#meta > div:nth-child(2)"
)

type PlayerScraper struct {
	colly       colly.Collector
	season      int
	ScrapedData []parser.Player
	Errors      []error
	child       Scraper
	childUrls   map[string]string
}

func CreatePlayerScraper(c *colly.Collector) PlayerScraper {
	return PlayerScraper{
		colly:     *c,
		childUrls: make(map[string]string),
	}
}

func (s *PlayerScraper) GetData() interface{} {
	return s.ScrapedData
}

func (s *PlayerScraper) AttachChild(scraper *Scraper) {
	s.child = *scraper
}

func (s *PlayerScraper) GetChild() Scraper {
	return s.child
}

func (s *PlayerScraper) GetChildUrls() []string {
	return urlsMapToArray(s.childUrls)
}

func (s *PlayerScraper) Scrape(urls ...string) {

	for _, url := range urls {
		player := s.parsePlayerPage(url)
		s.ScrapedData = append(s.ScrapedData, player)
	}

	fmt.Printf("%+v\n", s.ScrapedData)

	scrapeChild(s)
}

func (s *PlayerScraper) parsePlayerPage(url string) (player parser.Player) {
	c := s.colly.Clone()
	c.OnRequest(onRequestVisit)

	player.Id = parser.ParseLastId(url)

	c.OnHTML(playerInfoElement, func(div *colly.HTMLElement) {
		player.PlayerInfoBox(div)
	})

	c.Visit(url)

	return
}
