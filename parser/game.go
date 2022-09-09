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
	g.HomeLineScores, g.AwayLineScores = p.GS.parseLineScoreTable(tbl, g.ID)
	g.Quarters = len(g.HomeLineScores)
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

func (p *GameParser) CheckScheduleLinkSeason(a *colly.HTMLElement) {
	// link format: /leagues/NBA_2022_games.html (we want the 2022 obviously)
	var season, _ = strconv.Atoi(strings.Split(a.Attr("href"), "_")[1])
	if season != p.Config.Season {
		core.Log.Fatalf("Scraped season (%d) is different from configured season (%d)!  This is a bad idea.", season, p.Config.Season)
	}
}

func parseMetaScorebox(box *colly.HTMLElement) (startTime time.Time, location string) {
	startTime, _ = time.ParseInLocation("3:04 PM, January 2, 2006", box.ChildText("div:first-child"), EST)
	location = box.ChildText("div:nth-child(2)")
	return
}

func parseTypeFromTitle(title string) string {
	// play-in games are weird...giving them a different type so they can be queried/excluded
	if strings.HasPrefix(title, "Play-In Game") {
		return "I"
	} else if strings.Contains(title, "NBA") {
		return "P"
	} else {
		return "R"
	}
}

// I looked this up: there can be no ties in the NBA!
// Also, the W-L we scraped are including this game; we want them as of before this game.  So we simply adjust.
func (*GameParser) SetResultAndAdjust(g *model.Game) {
	if g.Home.Score > g.Away.Score {
		g.Home.Result = "W"
		g.Away.Result = "L"
		g.Home.Wins--
		g.Away.Losses--
	} else {
		g.Home.Result = "L"
		g.Away.Result = "W"
		g.Home.Losses--
		g.Away.Wins--
	}
}
