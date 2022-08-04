package parser

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/nschimek/nba-scraper/core"
	"github.com/nschimek/nba-scraper/model"
)

type GameParser struct {
	Config *core.Config           `Inject:""`
	GS     *GameStatsParser       `Inject:""`
	GPS    *GamePlayerStatsParser `Inject:""`
}

type Game struct {
	Id, Location, Type               string
	Season, Quarters                 int
	StartTime                        time.Time
	HomeTeam, AwayTeam               GameTeam
	HomeLineScore, AwayLineScore     []GameLineScore // these will end up in their own table due to the possiblity of OT
	HomeFourFactors, AwayFourFactors GameFourFactors // also probably their own table
	GamePlayers                      []GamePlayer
	GamePlayersBasicStats            []GamePlayerBasicStats
	GamePlayersAdvancedStats         []GamePlayerAdvancedStats
}

func (*GameParser) GameTitle(g *model.Game, div *colly.HTMLElement) {
	g.Type = parseTypeFromTitle(div.ChildText("h1"))
}

func (p *GameParser) Scorebox(g *model.Game, box *colly.HTMLElement, index int) {
	if index == 0 { // away team is first
		g.Away = *p.GS.parseScorebox(box)
	} else if index == 1 { // home team is second
		g.Home = *p.GS.parseScorebox(box)
	} else if index == 2 && box.Attr("class") == "scorebox_meta" {
		g.StartTime, g.Location = parseMetaScorebox(box)
	}
}

func (p *GameParser) LineScoreTable(g *model.Game, tbl *colly.HTMLElement) {
	g.HomeLineScore, g.AwayLineScore = p.GS.parseLineScoreTable(tbl, g.ID)
	g.Quarters = len(g.HomeLineScore)
}

func (p *GameParser) FourFactorsTable(g *model.Game, tbl *colly.HTMLElement) {
	g.HomeFourFactors, g.AwayFourFactors = p.GS.parseFourFactorsTable(tbl, g.ID)
}

func (p *GameParser) ScoreboxStatTable(g *model.Game, box *colly.HTMLElement) {
	teamId, boxType, quarter := parseBoxScoreTableProperties(box.Attr("id"))
	box.ForEach("tbody", func(_ int, tbl *colly.HTMLElement) {
		if boxType == "basic" && quarter > 0 && quarter < math.MaxInt {
			g.GamePlayersBasicStats = append(g.GamePlayersBasicStats, p.GPS.parseBasicBoxScoreTable(tbl, g.ID, teamId, quarter)...)
		} else if boxType == "basic" && quarter == math.MaxInt {
			g.GamePlayers = append(g.GamePlayers, p.GPS.parseBasicBoxScoreGameTable(tbl, g.ID, teamId)...)
		} else if boxType == "advanced" {
			g.GamePlayersAdvancedStats = append(g.GamePlayersAdvancedStats, p.GPS.parseAdvancedBoxScoreTable(tbl, g.ID, teamId)...)
		}
	})
}

func (p *GameParser) InactivePlayersList(g *model.Game, box *colly.HTMLElement) {
	if box.Attr("id") == "" && box.Attr("class") == "" {
		g.GamePlayers = append(g.GamePlayers, p.GPS.parseInactivePlayersList(box, g.ID)...)
	}
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

func parseTypeFromTitle(title string) string {
	if strings.Contains(title, "NBA") {
		return "P"
	} else {
		return "R"
	}
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
