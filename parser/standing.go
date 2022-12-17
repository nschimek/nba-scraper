package parser

import (
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/nschimek/nba-scraper/core"
	"github.com/nschimek/nba-scraper/model"
)

type StandingParser struct {
	Config *core.Config `Inject:""`
}

type WinLoss struct {
	Wins, Losses int
}

func (p *StandingParser) StandingsTable(tbl *colly.HTMLElement) []model.TeamStanding {
	standings := []model.TeamStanding{}

	for _, rowMap := range Table(tbl) {
		standing, err := standingFromRow(rowMap)
		standing.CaptureError(err)
		standing.Season = p.Config.Season
		standings = append(standings, standing)
	}

	return standings
}

func standingFromRow(rowMap map[string]*colly.HTMLElement) (model.TeamStanding, error) {
	var err error
	standing := new(model.TeamStanding)

	standing.Rank, _ = strconv.Atoi(getColumnText(rowMap, "ranker"))
	standing.TeamId, err = ParseTeamId(parseLink(rowMap["team_name"]))
	standing.Overall = parseWinLoss(getColumnText(rowMap, "Overall"))
	standing.Home = parseWinLoss(getColumnText(rowMap, "Home"))
	standing.Road = parseWinLoss(getColumnText(rowMap, "Road"))
	standing.East = parseWinLoss(getColumnText(rowMap, "E"))
	standing.West = parseWinLoss(getColumnText(rowMap, "W"))
	standing.Atlantic = parseWinLoss(getColumnText(rowMap, "A"))
	standing.Central = parseWinLoss(getColumnText(rowMap, "C"))
	standing.Southeast = parseWinLoss(getColumnText(rowMap, "SE"))
	standing.Northwest = parseWinLoss(getColumnText(rowMap, "NW"))
	standing.Pacific = parseWinLoss(getColumnText(rowMap, "P"))
	standing.Southwest = parseWinLoss(getColumnText(rowMap, "SW"))
	standing.PreAllStar = parseWinLoss(getColumnText(rowMap, "Pre"))
	standing.PostAllStar = parseWinLoss(getColumnText(rowMap, "Post"))
	standing.MarginLess3 = parseWinLoss(getColumnText(rowMap, "3"))
	standing.MarginGreater10 = parseWinLoss(getColumnText(rowMap, "10"))
	standing.October = parseWinLoss(getColumnText(rowMap, "Oct"))
	standing.November = parseWinLoss(getColumnText(rowMap, "Nov"))
	standing.December = parseWinLoss(getColumnText(rowMap, "Dec"))
	standing.January = parseWinLoss(getColumnText(rowMap, "Jan"))
	standing.February = parseWinLoss(getColumnText(rowMap, "Feb"))
	standing.March = parseWinLoss(getColumnText(rowMap, "Mar"))
	standing.April = parseWinLoss(getColumnText(rowMap, "Apr"))

	return *standing, err
}

func parseWinLoss(s string) model.WinLoss {
	wl := new(model.WinLoss)

	if p := strings.Split(s, "-"); s != "" && len(p) == 2 {
		wl.Wins, _ = strconv.Atoi(p[0])
		wl.Losses, _ = strconv.Atoi(p[1])
	} else {
		core.Log.Warn("could not parse win/loss record, using 0-0")
		wl.Wins, wl.Losses = 0, 0
	}

	return *wl
}
