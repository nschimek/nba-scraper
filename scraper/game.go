package scraper

import (
	"fmt"

	"github.com/gocolly/colly/v2"
	"github.com/nschimek/nba-scraper/parser"
)

const (
	gameBaseBodyElement         = "body div#wrap"
	baseContentElement          = gameBaseBodyElement + " div#content"
	scoreboxElements            = "div.scorebox > div"
	lineScoreTableElementBase   = gameBaseBodyElement + " .content_grid div:nth-child(1) div#all_line_score.table_wrapper"
	lineScoreTableElement       = "div#div_line_score.table_container > table > tbody"
	fourFactorsTableElementBase = gameBaseBodyElement + " .content_grid div:nth-child(2) div#all_four_factors.table_wrapper"
	fourFactorsTableElement     = "div#div_four_factors.table_container > table > tbody"
	basicBoxScoreTables         = "div.section_wrapper div.section_content div.table_wrapper div.table_container table"
	seasonLinkElement           = gameBaseBodyElement + " div#bottom_nav.section_wrapper div#bottom_nav_container.section_content ul li:nth-child(3) a"
	scoreboxFooterElement       = "div#content > div"
)

type GameScraper struct {
	colly       colly.Collector
	ScrapedData []parser.Game
	Errors      []error
	child       Scraper
	childUrls   map[string]string
}

func CreateGameScraper(c *colly.Collector) GameScraper {
	return GameScraper{
		colly:     *c,
		childUrls: make(map[string]string),
	}
}

// Scraper interface methods
func (s *GameScraper) GetData() interface{} {
	return s.ScrapedData
}

func (s *GameScraper) AttachChild(scraper *Scraper) {
	s.child = *scraper
}

func (s *GameScraper) GetChild() Scraper {
	return s.child
}

func (s *GameScraper) GetChildUrls() []string {
	return urlsMapToArray(s.childUrls)
}

func (s *GameScraper) Scrape(urls ...string) {

	for _, url := range urls {
		game := s.parseGamePage(url)
		s.ScrapedData = append(s.ScrapedData, game)
		s.childUrls[game.HomeTeam.TeamId] = game.HomeTeam.TeamUrl
		s.childUrls[game.AwayTeam.TeamId] = game.AwayTeam.TeamUrl
	}

	// fmt.Printf("%+v\n", s.ScrapedData)

	scrapeChild(s)
}

func (s *GameScraper) parseGamePage(url string) (game parser.Game) {
	c := s.colly.Clone()

	game.Id = parser.ParseLastId(url)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting ", r.URL.String())
	})

	c.OnHTML(baseContentElement, func(div *colly.HTMLElement) {
		game.GameTitle(div)

		fmt.Println(game.Type)

		div.ForEach(scoreboxElements, func(i int, box *colly.HTMLElement) {
			game.Scorebox(box, i)
		})
		game.SetResultAndAdjust()
	})

	c.OnHTML(lineScoreTableElementBase, func(div *colly.HTMLElement) {
		tbl, _ := transformHtmlElement(div, lineScoreTableElement, removeCommentsSyntax)
		game.LineScoreTable(tbl)
	})

	c.OnHTML(fourFactorsTableElementBase, func(div *colly.HTMLElement) {
		tbl, _ := transformHtmlElement(div, fourFactorsTableElement, removeCommentsSyntax)
		game.FourFactorsTable(tbl)
	})

	c.OnHTML(baseContentElement, func(div *colly.HTMLElement) {
		div.ForEach(basicBoxScoreTables, func(_ int, box *colly.HTMLElement) {
			game.ScoreboxStatTable(box)
		})
	})

	c.OnHTML(gameBaseBodyElement, func(div *colly.HTMLElement) {
		div.ForEach(scoreboxFooterElement, func(_ int, box *colly.HTMLElement) {
			game.InactivePlayersList(box)
		})
	})

	c.OnHTML(seasonLinkElement, func(li *colly.HTMLElement) {
		game.ScheduleLink(li)
	})

	c.Visit(url)

	return
}
