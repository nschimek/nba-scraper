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
	Config      *core.Config               `Inject:""`
	Colly       *colly.Collector           `Inject:""`
	GameParser  *parser.GameParser         `Inject:""`
	Repository  *repository.GameRepository `Inject:""`
	ScrapedData []model.Game
	PlayerIds   map[string]struct{}
	TeamIds     map[string]struct{}
}

func (s *GameScraper) Scrape(idMap map[string]struct{}) {
	s.PlayerIds = make(map[string]struct{})
	s.TeamIds = make(map[string]struct{})

	for _, id := range core.IdMapToArray(idMap) {
		game := s.parseGamePage(id)
		if !game.HasErrors() {
			s.ScrapedData = append(s.ScrapedData, game)
			for _, gp := range game.GamePlayers {
				s.PlayerIds[gp.PlayerId] = exists
			}
			s.TeamIds[game.Home.TeamId] = exists
			s.TeamIds[game.Away.TeamId] = exists
		} else {
			game.LogErrors()
		}
	}

	core.Log.WithField("games", len(s.ScrapedData)).Info("Finished scraping Game page(s)!")
}

func (s *GameScraper) Persist() {
	if len(s.ScrapedData) > 0 {
		s.Repository.UpsertGames(s.ScrapedData)
	} else {
		core.Log.Warn("No Games scraped to persist! Skipping...")
	}
}

func (s *GameScraper) parseGamePage(id string) (game model.Game) {
	c := core.CloneColly(s.Colly)

	game.ID = id
	game.Season = s.Config.Season

	c.OnHTML(baseContentElement, func(div *colly.HTMLElement) {
		s.GameParser.GameTitle(&game, div)

		div.ForEach(scoreboxElements, func(i int, box *colly.HTMLElement) {
			s.GameParser.Scorebox(&game, box, i)
		})
		s.GameParser.SetResultAndAdjust(&game)
	})

	c.OnHTML(lineScoreTableElementBase, func(div *colly.HTMLElement) {
		tbl, _ := transformHtmlElement(div, lineScoreTableElement, removeCommentsSyntax)
		s.GameParser.LineScoreTable(&game, tbl)
	})

	c.OnHTML(fourFactorsTableElementBase, func(div *colly.HTMLElement) {
		tbl, _ := transformHtmlElement(div, fourFactorsTableElement, removeCommentsSyntax)
		s.GameParser.FourFactorsTable(&game, tbl)
	})

	c.OnHTML(baseContentElement, func(div *colly.HTMLElement) {
		div.ForEach(basicBoxScoreTables, func(_ int, box *colly.HTMLElement) {
			s.GameParser.ScoreboxStatTable(&game, box)
		})
	})

	c.OnHTML(gameBaseBodyElement, func(div *colly.HTMLElement) {
		div.ForEach(scoreboxFooterElement, func(_ int, box *colly.HTMLElement) {
			s.GameParser.InactivePlayersList(&game, box)
		})
	})

	c.OnHTML(seasonLinkElement, func(li *colly.HTMLElement) {
		if err := s.GameParser.CheckScheduleLinkSeason(li); err != nil {
			game.CaptureError(err)
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		game.CaptureError(NewScraperError(err, r.Request.URL.String()))
	})

	c.Visit(s.getUrl(id))

	return
}

func (s *GameScraper) getUrl(id string) string {
	return fmt.Sprintf(gameIdPage, id)
}
