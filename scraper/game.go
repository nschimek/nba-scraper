package scraper

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

const (
	baseBodyElement             = "body #wrap #content"
	scoreboxElement             = baseBodyElement + " .scorebox .scorebox_meta"
	lineScoreTableElementBase   = baseBodyElement + " .content_grid div:nth-child(1) div#all_line_score.table_wrapper"
	lineScoreTableElement       = "div#div_line_score.table_container table tbody"
	fourFactorsTableElementBase = baseBodyElement + " .content_grid div:nth-child(2) div#all_four_factors.table_wrapper"
	fourFactorsTableElement     = "div#div_four_factors.table_container table tbody"
	basicBoxScoreTables         = "div.section_wrapper.toggleable div.section_content div.table_wrapper div.table_container table"
)

type Game struct {
	HomeId, VisitorId, Location, WinnerId, LoserId        string
	Season, HomeScore, VisitorScore, Quarters, Attendance int
	StartTime                                             time.Time
	TimeOfGame                                            time.Duration
	HomeLineScore, VisitorLineScore                       []GameLineScore // these will end up in their own table due to the possiblity of OT
	HomeFourFactors, VisitorFourFactors                   GameFourFactors // also probably their own table
	GamePlayers                                           []GamePlayer
	GamePlayersBasicStats                                 []GamePlayerBasicStats
	GamePlayersAdvancedStats                              []GamePlayerAdvancedStats
}

type GameLineScore struct {
	TeamId         string
	Quarter, Score int
}

type GameFourFactors struct {
	TeamId                                                                       string
	Pace, EffectiveFgPct, TurnoverPct, OffensiveRbPct, FtPerFga, OffensiveRating float64
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

	c.OnHTML(lineScoreTableElementBase, func(div *colly.HTMLElement) {
		tbl, _ := transformHtmlElement(div, lineScoreTableElement, removeCommentsSyntax)
		game.HomeLineScore, game.VisitorLineScore = parseLineScoreTable(tbl)
		game.setTotalsFromLineScore()
	})

	c.OnHTML(fourFactorsTableElementBase, func(div *colly.HTMLElement) {
		tbl, _ := transformHtmlElement(div, fourFactorsTableElement, removeCommentsSyntax)
		game.HomeFourFactors, game.VisitorFourFactors = parseFourFactorsTable(tbl)
	})

	c.OnHTML(baseBodyElement, func(div *colly.HTMLElement) {
		div.ForEach(basicBoxScoreTables, func(_ int, box *colly.HTMLElement) {
			// TODO: handle box score tables
		})
	})

	c.Visit(url)

	return game
}

func (g *Game) setTotalsFromLineScore() {

	g.HomeId = g.HomeLineScore[0].TeamId
	g.VisitorId = g.VisitorLineScore[0].TeamId

	g.Quarters = len(g.HomeLineScore)

	for _, ls := range g.HomeLineScore {
		g.HomeScore = g.HomeScore + ls.Score
	}
	for _, ls := range g.VisitorLineScore {
		g.VisitorScore = g.VisitorScore + ls.Score
	}

	// I looked this up: there can be no ties in the NBA!
	if g.HomeScore > g.VisitorScore {
		g.WinnerId = g.HomeId
		g.LoserId = g.VisitorId
	} else {
		g.WinnerId = g.VisitorId
		g.LoserId = g.HomeId
	}

}

func (g *Game) setBoxScoreCallbacks(c *colly.Collector) {
	q := 1
	for q <= g.Quarters {
		// fmt.Println(basicBoxScoreTableSelector(g.HomeId, q))

		c.OnHTML(baseBodyElement, func(tbl *colly.HTMLElement) {
			fmt.Println(tbl.DOM.Html())
		})
		q++
	}
}

func parseLineScoreTable(tbl *colly.HTMLElement) (home []GameLineScore, visitor []GameLineScore) {
	tableMaps := ParseTable(tbl) // row 0 will be away, row 1 will be home

	visitor = lineScoreFromRow(tableMaps[0])
	home = lineScoreFromRow(tableMaps[1])

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

func parseFourFactorsTable(tbl *colly.HTMLElement) (home GameFourFactors, visitor GameFourFactors) {
	tableMaps := ParseTable(tbl) // row 0 will be away, row 1 will be home

	visitor = gameFourFactorsFromRow(tableMaps[0])
	home = gameFourFactorsFromRow(tableMaps[1])

	return
}

func gameFourFactorsFromRow(rowMap map[string]*colly.HTMLElement) (factors GameFourFactors) {
	factors.TeamId = parseTeamId(parseLink(rowMap["team_id"]))
	factors.Pace, _ = strconv.ParseFloat(rowMap["pace"].Text, 64)
	factors.EffectiveFgPct, _ = strconv.ParseFloat(rowMap["efg_pct"].Text, 64)
	factors.TurnoverPct, _ = strconv.ParseFloat(rowMap["tov_pct"].Text, 64)
	factors.OffensiveRbPct, _ = strconv.ParseFloat(rowMap["orb_pct"].Text, 64)
	factors.FtPerFga, _ = strconv.ParseFloat(rowMap["ft_rate"].Text, 64)
	factors.OffensiveRating, _ = strconv.ParseFloat(rowMap["off_rtg"].Text, 64)

	return
}

func basicBoxScoreTableSelector(teamId string, quarter int) string {
	var q string

	if quarter <= 4 {
		q = fmt.Sprintf("q%d", quarter)
	} else {
		q = fmt.Sprintf("ot%d", quarter-4)
	}

	return "div#all_box-" + teamId + "-" + q + "-basic div.section_content div.table_wrapper div.table_container table tbody"
}
