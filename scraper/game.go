package scraper

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

const (
	baseBodyElement         = "body #wrap #content"
	scoreboxElement         = baseBodyElement + " .scorebox .scorebox_meta"
	lineScoreTableElement   = baseBodyElement + " .content_grid div:first-child #all_line_score #div_line_score table tbody"
	fourFactorsTableElement = baseBodyElement + " .content_grid div #all_four_factors #div_four_factors table tbody"
)

type Game struct {
	HomeId, VisitorId, Location, WinnerId, LoserId string
	HomeScore, VisitorScore, Attendance            int
	StartTime                                      time.Time
	TimeOfGame                                     time.Duration
	HomeLineScore, VisitorLineScore                []GameLineScore // these will end up in their own table due to the possiblity of OT
	HomeFourFactors, VisitorFourFactors            GameFourFactors // also probably their own table
	GamePlayers                                    []GamePlayer
	GamePlayersBasicStats                          []GamePlayerBasicStats
	GamePlayersAdvancedStats                       []GamePlayerAdvancedStats
}

// return a new game object with all child objects and arrys initialized
func newGame() Game {
	return Game{
		HomeLineScore:            make([]GameLineScore, 0),
		VisitorLineScore:         make([]GameLineScore, 0),
		HomeFourFactors:          GameFourFactors{},
		VisitorFourFactors:       GameFourFactors{},
		GamePlayers:              make([]GamePlayer, 0),
		GamePlayersBasicStats:    make([]GamePlayerBasicStats, 0),
		GamePlayersAdvancedStats: make([]GamePlayerAdvancedStats, 0),
	}
}

type GameLineScore struct {
	TeamId         string
	Quarter, Score int
}

type GameFourFactors struct {
	TeamId                                                       string
	Pace, EffectiveFgPct, TurnoverPct, FtPerFga, OffensiveRating float64
}

type GamePlayer struct {
	GameId, TeamId, PlayerId, Status string
}

type GamePlayerBasicStats struct {
	GameId, TeamId, PlayerId                                                                                string
	Quarter                                                                                                 int
	TimePlayed                                                                                              time.Duration
	FieldGoals, FieldGoalsAttempted, ThreePointers, ThreePointersAttempted, FreeThrows, FreeThrowsAttempted int
	FieldGoalPct, ThreePointersPct, FreeThrowsPct                                                           float64
	OffensiveRB, DefensiveRB, TotalRB, Assists, Steals, Blocks, Turnovers, PersonalFouls, Points, PlusMinus int
}

type GamePlayerAdvancedStats struct {
	GameId, TeamId, PlayerId                                                                        string
	TrueShootingPct, EffectiveFgPct, ThreePtAttemptRate, OffensiveRbPct, DefensiveRbPct, TotalRbPct float64
	AssistPct, StealPct, BlockPct, TurnoverPct, UsagePct, BoxPlusMinus                              float64
	OffensiveRating, DefensiveRating                                                                int
}

type GameScraper struct {
	colly       colly.Collector
	ScrapedData []Game
	Errors      []error
	child       Scraper
	childUrls   map[string]string
}

func CreateGameScraper(c *colly.Collector) GameScraper {
	return GameScraper{
		colly: *c,
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

// Testing the approach for scraping Game pages, as these are more complex
func (s *GameScraper) Scrape(urls ...string) {

	for _, url := range urls {
		s.ScrapedData = append(s.ScrapedData, s.parseGamePage(url))
	}

	fmt.Println(s.ScrapedData)

	scrapeChild(s)
}

func (s *GameScraper) parseGamePage(url string) Game {
	game := Game{}
	c := s.colly.Clone()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting ", r.URL.String())
	})

	c.OnHTML(scoreboxElement, func(div *colly.HTMLElement) {
		game.StartTime, _ = time.ParseInLocation("3:04 PM, January 2, 2006", div.ChildText("div:first-child"), EST)
		game.Location = div.ChildText("div:nth-child(2)")
	})

	// unfortunately, the line score and four factor tables are commented out, so we have to get our hands dirty
	c.OnHTML("body #wrap #content .content_grid div:nth-child(1) div#all_line_score.table_wrapper", func(div *colly.HTMLElement) {
		html, _ := div.DOM.Html()
		doc, _ := goquery.NewDocumentFromReader(bytes.NewBufferString(strings.Replace(strings.Replace(html, "<!--", "", 1), "-->", "", 1)))
		sel := doc.Find("div#div_line_score.table_container table tbody")
		tbl := colly.NewHTMLElementFromSelectionNode(div.Response, sel, sel.Get(0), 0)

		game.HomeLineScore, game.VisitorLineScore = parseLineScoreTable(tbl)
	})

	c.Visit(url)

	return game
}

func parseLineScoreTable(tbl *colly.HTMLElement) (home []GameLineScore, visitor []GameLineScore) {
	tableMaps := ParseTable(tbl) // row 0 will be away, row 1 will be home

	visitor = lineScoreFromRow(tableMaps[0])
	home = lineScoreFromRow(tableMaps[1])

	fmt.Println(visitor)
	fmt.Println(home)

	return
}

func lineScoreFromRow(rowMap map[string]*colly.HTMLElement) (scores []GameLineScore) {
	teamId := parseTeamId(parseLink(rowMap["team"]))

	for key, cell := range rowMap {
		// loop through all non-team and total columns; those that remain are the quarters
		if key != "team" && key != "T" {
			score, _ := strconv.Atoi(cell.Text)
			scores = append(scores, GameLineScore{
				TeamId:  teamId,
				Quarter: lineScoreQuarter(key),
				Score:   score,
			})
		}
	}

	return
}

func lineScoreQuarter(c string) int {
	// if we can parse an int out, then it's quarters 1-4
	if i, err := strconv.Atoi(c); err == nil {
		return i
	} else {
		// remove OT, then parse what is left...an error indicates OT1 (5), as it will be blank
		ot, err := strconv.Atoi(strings.Replace(c, "OT", "", 1))
		if err != nil {
			return 5
		} else {
			return ot + 4
		}
	}
}
