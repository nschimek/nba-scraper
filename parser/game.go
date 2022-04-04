package parser

import (
	"math"
	"time"

	"github.com/gocolly/colly/v2"
)

type Game struct {
	Id, HomeId, VisitorId, HomeUrl, VisitorUrl, Location, WinnerId, LoserId string
	Season, HomeScore, VisitorScore, Quarters                               int
	StartTime                                                               time.Time
	HomeLineScore, VisitorLineScore                                         []GameLineScore // these will end up in their own table due to the possiblity of OT
	HomeFourFactors, VisitorFourFactors                                     GameFourFactors // also probably their own table
	GamePlayers                                                             []GamePlayer
	GamePlayersBasicStats                                                   []GamePlayerBasicStats
	GamePlayersAdvancedStats                                                []GamePlayerAdvancedStats
}

func (g *Game) MetaScorebox(div *colly.HTMLElement) {
	g.StartTime, _ = time.ParseInLocation("3:04 PM, January 2, 2006", div.ChildText("div:first-child"), EST)
	g.Location = div.ChildText("div:nth-child(2)")
}

func (g *Game) LineScoreTable(tbl *colly.HTMLElement) {
	g.HomeLineScore, g.VisitorLineScore = ParseLineScoreTable(tbl)
	g.setTotalsFromLineScore()
}

func (g *Game) FourFactorsTable(tbl *colly.HTMLElement) {
	g.HomeFourFactors, g.VisitorFourFactors = ParseFourFactorsTable(tbl)
}

func (g *Game) ScoreboxStatTable(box *colly.HTMLElement) {
	teamId, boxType, quarter := parseBoxScoreTableProperties(box.Attr("id"))
	box.ForEach("tbody", func(_ int, tbl *colly.HTMLElement) {
		if boxType == "basic" && quarter > 0 && quarter < math.MaxInt {
			g.GamePlayersBasicStats = append(g.GamePlayersBasicStats, ParseBasicBoxScoreTable(tbl, g.Id, teamId, quarter)...)
		} else if boxType == "basic" && quarter == math.MaxInt {
			g.GamePlayers = append(g.GamePlayers, ParseBasicBoxScoreGameTable(tbl, g.Id, teamId)...)
		} else if boxType == "advanced" {
			g.GamePlayersAdvancedStats = append(g.GamePlayersAdvancedStats, ParseAdvancedBoxScoreTable(tbl, g.Id, teamId)...)
		}
	})
}

func (g *Game) setTotalsFromLineScore() {

	g.HomeId = g.HomeLineScore[0].TeamId
	g.HomeUrl = g.HomeLineScore[0].TeamUrl
	g.VisitorId = g.VisitorLineScore[0].TeamId
	g.VisitorUrl = g.VisitorLineScore[0].TeamUrl

	g.Quarters = len(g.HomeLineScore)

	for _, ls := range g.HomeLineScore {
		g.HomeScore = g.HomeScore + ls.Score
	}
	for _, ls := range g.VisitorLineScore {
		g.VisitorScore = g.VisitorScore + ls.Score
	}

	g.WinnerId, g.LoserId = g.getResult()
}

// I looked this up: there can be no ties in the NBA!
func (g *Game) getResult() (winnerId, loserId string) {
	if g.HomeScore > g.VisitorScore {
		return g.HomeId, g.VisitorId
	} else {
		return g.VisitorId, g.HomeId
	}
}
