package parser

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

type Game struct {
	Id, Location                     string
	Season, Quarters                 int
	StartTime                        time.Time
	HomeTeam, AwayTeam               GameTeam
	HomeLineScore, AwayLineScore     []GameLineScore // these will end up in their own table due to the possiblity of OT
	HomeFourFactors, AwayFourFactors GameFourFactors // also probably their own table
	GamePlayers                      []GamePlayer
	GamePlayersBasicStats            []GamePlayerBasicStats
	GamePlayersAdvancedStats         []GamePlayerAdvancedStats
}

func (g *Game) Scorebox(box *colly.HTMLElement, index int) {
	if index == 0 { // away team is first
		g.AwayTeam = parseScorebox(box)
	} else if index == 1 { // home team is second
		g.HomeTeam = parseScorebox(box)
	} else if index == 2 && box.Attr("class") == "scorebox_meta" {
		g.StartTime, g.Location = parseMetaScorebox(box)
	}
}

func (g *Game) LineScoreTable(tbl *colly.HTMLElement) {
	g.HomeLineScore, g.AwayLineScore = parseLineScoreTable(tbl)
	g.Quarters = len(g.HomeLineScore)
}

func (g *Game) FourFactorsTable(tbl *colly.HTMLElement) {
	g.HomeFourFactors, g.AwayFourFactors = parseFourFactorsTable(tbl)
}

func (g *Game) ScoreboxStatTable(box *colly.HTMLElement) {
	teamId, boxType, quarter := parseBoxScoreTableProperties(box.Attr("id"))
	box.ForEach("tbody", func(_ int, tbl *colly.HTMLElement) {
		if boxType == "basic" && quarter > 0 && quarter < math.MaxInt {
			g.GamePlayersBasicStats = append(g.GamePlayersBasicStats, parseBasicBoxScoreTable(tbl, teamId, quarter)...)
		} else if boxType == "basic" && quarter == math.MaxInt {
			g.GamePlayers = append(g.GamePlayers, parseBasicBoxScoreGameTable(tbl, teamId)...)
		} else if boxType == "advanced" {
			g.GamePlayersAdvancedStats = append(g.GamePlayersAdvancedStats, parseAdvancedBoxScoreTable(tbl, teamId)...)
		}
	})
}

func (g *Game) ScheduleLink(a *colly.HTMLElement) {
	// link format: /leagues/NBA_2022_games.html (we want the 2022 obviously)
	g.Season, _ = strconv.Atoi(strings.Split(a.Attr("href"), "_")[1])
}

func parseMetaScorebox(box *colly.HTMLElement) (startTime time.Time, location string) {
	startTime, _ = time.ParseInLocation("3:04 PM, January 2, 2006", box.ChildText("div:first-child"), EST)
	location = box.ChildText("div:nth-child(2)")
	return
}

// I looked this up: there can be no ties in the NBA!
// Also, the W-L we scraped are including this game; we want them as of before this game.  So we simply adjust.
func (g *Game) SetResultAndAdjust() {
	if g.HomeTeam.Score > g.AwayTeam.Score {
		g.HomeTeam.Result = "W"
		g.AwayTeam.Result = "L"
		g.HomeTeam.Wins--
		g.AwayTeam.Losses--
	} else {
		g.HomeTeam.Result = "L"
		g.AwayTeam.Result = "W"
		g.HomeTeam.Losses--
		g.AwayTeam.Wins--
	}
}
