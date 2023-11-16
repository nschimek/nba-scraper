package parser

import (
	"errors"
	"fmt"
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

func (p *GameParser) GameTitle(g *model.Game, div *colly.HTMLElement) {
	var err error
	g.Type, err = parseTypeFromTitle(div.ChildText("h1"))
	g.CaptureError(err)
}

func (p *GameParser) Scorebox(g *model.Game, box *colly.HTMLElement, index int) {
	var err error
	if index == 0 { // away team is first
		g.Away, err = p.GS.parseScorebox(box)
	} else if index == 1 { // home team is second
		g.Home, err = p.GS.parseScorebox(box)
	} else if index == 2 && box.Attr("class") == "scorebox_meta" {
		g.StartTime, g.Location, err = parseMetaScorebox(box)
	}
	g.CaptureError(err)
}

func (p *GameParser) LineScoreTable(g *model.Game, tbl *colly.HTMLElement) {
	var err error
	g.HomeLineScores, g.AwayLineScores, err = p.GS.parseLineScoreTable(tbl, g.ID)
	g.Quarters = len(g.HomeLineScores)
	g.CaptureError(err)
}

func (p *GameParser) FourFactorsTable(g *model.Game, tbl *colly.HTMLElement) {
	var err error
	g.HomeFourFactors, g.AwayFourFactors, err = p.GS.parseFourFactorsTable(tbl, g.ID)
	g.CaptureError(err)
}

func (p *GameParser) ScoreboxStatTable(g *model.Game, box *colly.HTMLElement) {
	teamId, boxType, quarter, err := parseBoxScoreTableProperties(box.Attr("id"))
	if err == nil {
		box.ForEach("tbody", func(_ int, tbl *colly.HTMLElement) {
			if boxType == "basic" && quarter > 0 && quarter < math.MaxInt {
				g.GamePlayersBasicStats = append(g.GamePlayersBasicStats, p.GPS.parseBasicBoxScoreTable(tbl, g.ID, teamId, quarter)...)
			} else if boxType == "basic" && quarter == math.MaxInt {
				g.GamePlayers = append(g.GamePlayers, p.GPS.parseBasicBoxScoreGameTable(tbl, g.ID, teamId)...)
			} else if boxType == "advanced" {
				g.GamePlayersAdvancedStats = append(g.GamePlayersAdvancedStats, p.GPS.parseAdvancedBoxScoreTable(tbl, g.ID, teamId)...)
			}
		})
	} else {
		core.Log.Warnf("error encountered while parsing box scores: %s", err.Error())
	}
}

func (p *GameParser) InactivePlayersList(g *model.Game, box *colly.HTMLElement) {
	if box.Attr("id") == "" && box.Attr("class") == "" {
		g.GamePlayers = append(g.GamePlayers, p.GPS.parseInactivePlayersList(box, g.ID)...)
	}
}

func (p *GameParser) CheckScheduleLinkSeason(a *colly.HTMLElement) error {
	var err error
	var season int
	// link format: /leagues/NBA_2022_games.html (we want the 2022 obviously)
	parts := strings.Split(a.Attr("href"), "_")

	if len(parts) < 2 {
		err = errors.New("unexpected format of league season link, so could not validate season")
		return err
	}

	season, err = strconv.Atoi(strings.Split(a.Attr("href"), "_")[1])
	if season != p.Config.Season && err != nil {
		err = errors.New("scraped season is different from configured season")
		return err
	}

	return nil
}

func parseMetaScorebox(box *colly.HTMLElement) (startTime time.Time, location string, err error) {
	// BR now puts "In-Season Tournament" as the top DIV for this box for tournament games (but not regular games)
	box.ForEachWithBreak("div", func(i int, div *colly.HTMLElement) bool {
		// look the start date; once we find it, the next div will be the location
		startTime, err = time.ParseInLocation("3:04 PM, January 2, 2006", div.Text, EST)
		if err == nil {
			location = box.ChildText(fmt.Sprintf("div:nth-child(%d)", i + 2))
			return false
		} else {
			return true
		}
	})
	if err != nil  {
		err = errors.New("could not parse game start time from scorebox")
	}
	return
}

func parseTypeFromTitle(title string) (string, error) {
	if title != "" { // play-in games are weird...giving them a different type so they can be queried/excluded
		if strings.HasPrefix(title, "Play-In Game") {
			return "I", nil
		} else if strings.Contains(title, "NBA") {
			return "P", nil
		} else {
			return "R", nil
		}
	} else {
		return "", errors.New("could not get game title, so could not determine game type")
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
