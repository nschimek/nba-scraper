package scraper

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

const (
	baseBodyElement             = "body #wrap #content"
	scoreboxMetaElement         = baseBodyElement + " .scorebox .scorebox_meta"
	lineScoreTableElementBase   = baseBodyElement + " .content_grid div:nth-child(1) div#all_line_score.table_wrapper"
	lineScoreTableElement       = "div#div_line_score.table_container table tbody"
	fourFactorsTableElementBase = baseBodyElement + " .content_grid div:nth-child(2) div#all_four_factors.table_wrapper"
	fourFactorsTableElement     = "div#div_four_factors.table_container table tbody"
	basicBoxScoreTables         = "div.section_wrapper div.section_content div.table_wrapper div.table_container table"
)

type Game struct {
	Id, HomeId, VisitorId, Location, WinnerId, LoserId    string
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
	GameId, TeamId, PlayerId                                                                                              string
	TrueShootingPct, EffectiveFgPct, ThreePtAttemptRate, FreeThrowAttemptRate, OffensiveRbPct, DefensiveRbPct, TotalRbPct float64
	AssistPct, StealPct, BlockPct, TurnoverPct, UsagePct, BoxPlusMinus                                                    float64
	OffensiveRating, DefensiveRating                                                                                      int
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
		game, homeUrl, visitorUrl := s.parseGamePage(url)
		s.ScrapedData = append(s.ScrapedData, game)
		s.childUrls[game.HomeId] = homeUrl
		s.childUrls[game.VisitorId] = visitorUrl
	}

	fmt.Println(s.ScrapedData)

	scrapeChild(s)
}

func (s *GameScraper) parseGamePage(url string) (game Game, homeUrl, visitorUrl string) {
	game = Game{}
	c := s.colly.Clone()

	game.Id = parseGameId(url)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting ", r.URL.String())
	})

	c.OnHTML(scoreboxMetaElement, func(div *colly.HTMLElement) {
		game.StartTime, _ = time.ParseInLocation("3:04 PM, January 2, 2006", div.ChildText("div:first-child"), EST)
		game.Location = div.ChildText("div:nth-child(2)")
	})

	c.OnHTML(lineScoreTableElementBase, func(div *colly.HTMLElement) {
		tbl, _ := transformHtmlElement(div, lineScoreTableElement, removeCommentsSyntax)
		game.HomeLineScore, game.VisitorLineScore, homeUrl, visitorUrl = parseLineScoreTable(tbl)
		game.setTotalsFromLineScore()
	})

	c.OnHTML(fourFactorsTableElementBase, func(div *colly.HTMLElement) {
		tbl, _ := transformHtmlElement(div, fourFactorsTableElement, removeCommentsSyntax)
		game.HomeFourFactors, game.VisitorFourFactors = parseFourFactorsTable(tbl)
	})

	c.OnHTML(baseBodyElement, func(div *colly.HTMLElement) {
		div.ForEach(basicBoxScoreTables, func(_ int, box *colly.HTMLElement) {
			teamId, boxType, quarter := parseBoxScoreTableProperties(box.Attr("id"))
			box.ForEach("tbody", func(_ int, tbl *colly.HTMLElement) {
				if boxType == "basic" && quarter > 0 && quarter < math.MaxInt {
					game.GamePlayersBasicStats = append(game.GamePlayersBasicStats, parseBasicBoxScoreTable(tbl, game.Id, teamId, quarter)...)
				} else if boxType == "advanced" {
					game.GamePlayersAdvancedStats = append(game.GamePlayersAdvancedStats, parseAdvancedBoxScoreTable(tbl, game.Id, teamId)...)
				}
			})
		})
	})

	c.Visit(url)

	return
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

func parseLineScoreTable(tbl *colly.HTMLElement) (home []GameLineScore, visitor []GameLineScore, homeUrl, visitorUrl string) {
	tableMaps := ParseTable(tbl) // row 0 will be away, row 1 will be home

	visitor, visitorUrl = lineScoreFromRow(tableMaps[0])
	home, homeUrl = lineScoreFromRow(tableMaps[1])

	return
}

func lineScoreFromRow(rowMap map[string]*colly.HTMLElement) (scores []GameLineScore, teamUrl string) {
	teamUrl = parseLink(rowMap["team"])
	teamId := parseTeamId(teamUrl)

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

func parseBoxScoreTableProperties(id string) (team, boxType string, quarter int) {
	parts := strings.Split(id, "-") // expected format: box-TEAM-q/ot#-basic
	team = parts[1]
	boxType = parts[3]

	if boxType == "basic" {
		if q := parts[2]; strings.HasPrefix(q, "q") {
			quarter, _ = strconv.Atoi(strings.ReplaceAll(q, "q", ""))
		} else if strings.HasPrefix(q, "ot") {
			quarter, _ = strconv.Atoi(strings.ReplaceAll(q, "ot", ""))
			quarter = quarter + 4 // adjust for OT
		} else if q == "game" {
			quarter = math.MaxInt
		}
	}

	// quarter will be 0 for other box types, such as advanced

	return
}

func parseBasicBoxScoreTable(tbl *colly.HTMLElement, gameId, teamId string, quarter int) []GamePlayerBasicStats {
	tableMaps := ParseTable(tbl)
	stats := make([]GamePlayerBasicStats, 12)

	for _, rowMap := range tableMaps {
		stats = append(stats, gamePlayerBasicStatsFromRow(rowMap, gameId, teamId, quarter))
	}

	return stats
}

func gamePlayerBasicStatsFromRow(rowMap map[string]*colly.HTMLElement, gameId, teamId string, quarter int) (gpbs GamePlayerBasicStats) {
	gpbs.GameId = gameId
	gpbs.TeamId = teamId
	gpbs.Quarter = quarter
	gpbs.PlayerId = parsePlayerId(parseLink(rowMap["player"]))

	if _, ok := rowMap["reason"]; !ok { // a "reason" column indicates the player did not play
		gpbs.TimePlayed, _ = parseDuration(rowMap["mp"].Text)
		gpbs.FieldGoals, _ = strconv.Atoi(rowMap["fg"].Text)
		gpbs.FieldGoalsAttempted, _ = strconv.Atoi(rowMap["fga"].Text)
		gpbs.FieldGoalPct, _ = parseFloatStat(rowMap["fg_pct"].Text)
		gpbs.ThreePointers, _ = strconv.Atoi(rowMap["fg3"].Text)
		gpbs.ThreePointersAttempted, _ = strconv.Atoi(rowMap["fg3a"].Text)
		gpbs.ThreePointersPct, _ = parseFloatStat(rowMap["fg3_pct"].Text)
		gpbs.FreeThrows, _ = strconv.Atoi(rowMap["ft"].Text)
		gpbs.FreeThrowsAttempted, _ = strconv.Atoi(rowMap["fta"].Text)
		gpbs.FreeThrowsPct, _ = parseFloatStat(rowMap["ft_pct"].Text)
		gpbs.OffensiveRB, _ = strconv.Atoi(rowMap["orb"].Text)
		gpbs.DefensiveRB, _ = strconv.Atoi(rowMap["drb"].Text)
		gpbs.TotalRB, _ = strconv.Atoi(rowMap["trb"].Text)
		gpbs.Assists, _ = strconv.Atoi(rowMap["ast"].Text)
		gpbs.Steals, _ = strconv.Atoi(rowMap["stl"].Text)
		gpbs.Blocks, _ = strconv.Atoi(rowMap["blk"].Text)
		gpbs.PersonalFouls, _ = strconv.Atoi(rowMap["pf"].Text)
		gpbs.Points, _ = strconv.Atoi(rowMap["pts"].Text)
		gpbs.PlusMinus, _ = strconv.Atoi(strings.Replace(rowMap["plus_minus"].Text, "+", "", 1))
	}

	return
}

func parseAdvancedBoxScoreTable(tbl *colly.HTMLElement, gameId, teamId string) []GamePlayerAdvancedStats {
	tableMaps := ParseTable(tbl)
	stats := make([]GamePlayerAdvancedStats, 12)

	for _, rowMap := range tableMaps {
		stats = append(stats, gamePlayerAdvancedStatsFromRow(rowMap, gameId, teamId))
	}

	return stats
}

func gamePlayerAdvancedStatsFromRow(rowMap map[string]*colly.HTMLElement, gameId, teamId string) (gpas GamePlayerAdvancedStats) {
	gpas.GameId = gameId
	gpas.TeamId = teamId
	gpas.PlayerId = parsePlayerId(parseLink(rowMap["player"]))

	if _, ok := rowMap["reason"]; !ok { // a "reason" column indicates the player did not play
		gpas.TrueShootingPct, _ = parseFloatStat(rowMap["ts_pct"].Text)
		gpas.EffectiveFgPct, _ = parseFloatStat(rowMap["efg_pct"].Text)
		gpas.ThreePtAttemptRate, _ = parseFloatStat(rowMap["fg3a_per_fga_pct"].Text)
		gpas.FreeThrowAttemptRate, _ = parseFloatStat(rowMap["fta_per_fga_pct"].Text)
		gpas.OffensiveRbPct, _ = parseFloatStat(rowMap["orb_pct"].Text)
		gpas.DefensiveRbPct, _ = parseFloatStat(rowMap["drb_pct"].Text)
		gpas.TotalRbPct, _ = parseFloatStat(rowMap["trb_pct"].Text)
		gpas.AssistPct, _ = parseFloatStat(rowMap["ast_pct"].Text)
		gpas.StealPct, _ = parseFloatStat(rowMap["stl_pct"].Text)
		gpas.BlockPct, _ = parseFloatStat(rowMap["blk_pct"].Text)
		gpas.TurnoverPct, _ = parseFloatStat(rowMap["tov_pct"].Text)
		gpas.UsagePct, _ = parseFloatStat(rowMap["usg_pct"].Text)
		gpas.OffensiveRating, _ = strconv.Atoi(rowMap["off_rtg"].Text)
		gpas.DefensiveRating, _ = strconv.Atoi(rowMap["def_rtg"].Text)
		gpas.BoxPlusMinus, _ = parseFloatStat(rowMap["bpm"].Text)
	}

	return
}
