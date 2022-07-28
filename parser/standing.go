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
		standing := standingFromRow(rowMap)
		standing.Season = p.Config.Season
		standings = append(standings, *standing)
	}

	return standings
}

func standingFromRow(rowMap map[string]*colly.HTMLElement) *model.TeamStanding {
	standing := new(model.TeamStanding)

	standing.Rank, _ = strconv.Atoi(rowMap["ranker"].Text)
	standing.TeamId = ParseTeamId(parseLink(rowMap["team_name"]))
	standing.Overall = parseWinLoss(rowMap["Overall"].Text)
	standing.Home = parseWinLoss(rowMap["Home"].Text)
	standing.Road = parseWinLoss(rowMap["Road"].Text)
	standing.East = parseWinLoss(rowMap["E"].Text)
	standing.West = parseWinLoss(rowMap["W"].Text)
	standing.Atlantic = parseWinLoss(rowMap["A"].Text)
	standing.Central = parseWinLoss(rowMap["C"].Text)
	standing.Southeast = parseWinLoss(rowMap["SE"].Text)
	standing.Northwest = parseWinLoss(rowMap["NW"].Text)
	standing.Pacific = parseWinLoss(rowMap["P"].Text)
	standing.Southwest = parseWinLoss(rowMap["SW"].Text)
	standing.PreAllStar = parseWinLoss(rowMap["Pre"].Text)
	standing.PostAllStar = parseWinLoss(rowMap["Post"].Text)
	standing.MarginLess3 = parseWinLoss(rowMap["3"].Text)
	standing.MarginGreater10 = parseWinLoss(rowMap["10"].Text)
	standing.October = parseWinLoss(rowMap["Oct"].Text)
	standing.November = parseWinLoss(rowMap["Nov"].Text)
	standing.December = parseWinLoss(rowMap["Dec"].Text)
	standing.January = parseWinLoss(rowMap["Jan"].Text)
	standing.February = parseWinLoss(rowMap["Feb"].Text)
	standing.March = parseWinLoss(rowMap["Mar"].Text)
	standing.April = parseWinLoss(rowMap["Apr"].Text)

	return standing
}

func parseWinLoss(s string) model.WinLoss {
	wl := new(model.WinLoss)
	p := strings.Split(s, "-")

	wl.Wins, _ = strconv.Atoi(p[0])
	wl.Losses, _ = strconv.Atoi(p[1])

	return *wl
}
