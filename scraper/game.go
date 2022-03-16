package scraper

import (
	"fmt"
	"time"

	"github.com/gocolly/colly/v2"
)

const (
	baseBodyElement         = "body #wrap #content"
	scoreboxElement         = baseBodyElement + " .scorebox .scorebox_meta"
	lineScoreTableElement   = baseBodyElement + " .content_grid div #all_line_score #div_line_score table tbody"
	fourFactorsTableElement = baseBodyElement + " .content_grid div #all_four_factors #div_four_factors table tbody"
)

var overtimes = []string{"OT", "OT1", "OT2", "OT3", "OT4", "OT5", "OT6", "OT7", "OT8", "OT9", "OT10"}

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
	TeamId, Quarter string
	Score           int
}

type GameFourFactors struct {
	TeamId                                                       string
	Pace, EffectiveFgPct, TurnoverPct, FtPerFga, OffensiveRating float64
}

type GamePlayer struct {
	GameId, TeamId, PlayerId, Status string
}

type GamePlayerBasicStats struct {
	GameId, TeamId, PlayerId, Quarter                                                                       string
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

	c.OnHTML(lineScoreTableElement, func(tbl *colly.HTMLElement) {
		game.HomeLineScore, game.VisitorLineScore = parseLineScoreTable(tbl)
	})

	c.Visit(url)

	return game
}

func parseLineScoreTable(tbl *colly.HTMLElement) (home []GameLineScore, visitor []GameLineScore) {
	tableMaps := ParseTable(tbl) // row 0 will be away, row 1 will be home

	visitor = lineScoreFromRow(tableMaps[0])
	home = lineScoreFromRow(tableMaps[1])

	return
}

func lineScoreFromRow(rowMap map[string]*colly.HTMLElement) (scores []GameLineScore) {
	// TODO: create one row per quarter or OT column in the table map and return

	return
}
