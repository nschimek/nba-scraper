package parser

import (
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

type Standings struct {
	TeamId, TeamUrl                                              string
	Season, Rank                                                 int
	Overall, Home, Road, East, West                              WinLoss
	Atlantic, Central, Southeast, Northwest, Pacific, Southwest  WinLoss
	PreAllStar, PostAllStar, MarginLess3, MarginGreater10        WinLoss
	October, November, December, January, February, March, April WinLoss
}

type WinLoss struct {
	Wins, Losses int
}

func StandingsTable(tbl *colly.HTMLElement, season int) []Standings {
	standings := []Standings{}

	for _, rowMap := range Table(tbl) {
		standings = append(standings, standingsFromRow(rowMap, season))
	}

	return standings
}

func standingsFromRow(rowMap map[string]*colly.HTMLElement, season int) (standings Standings) {
	standings.Season = season
	standings.Rank, _ = strconv.Atoi(rowMap["ranker"].Text)
	standings.TeamUrl = parseLink(rowMap["team_name"])
	standings.TeamId = ParseTeamId(standings.TeamUrl)
	standings.Overall = parseWinLoss(rowMap["Overall"].Text)
	standings.Home = parseWinLoss(rowMap["Home"].Text)
	standings.Road = parseWinLoss(rowMap["Road"].Text)
	standings.East = parseWinLoss(rowMap["E"].Text)
	standings.West = parseWinLoss(rowMap["W"].Text)
	standings.Atlantic = parseWinLoss(rowMap["A"].Text)
	standings.Central = parseWinLoss(rowMap["C"].Text)
	standings.Southeast = parseWinLoss(rowMap["SE"].Text)
	standings.Northwest = parseWinLoss(rowMap["NW"].Text)
	standings.Pacific = parseWinLoss(rowMap["P"].Text)
	standings.Southwest = parseWinLoss(rowMap["SW"].Text)
	standings.PreAllStar = parseWinLoss(rowMap["Pre"].Text)
	standings.PostAllStar = parseWinLoss(rowMap["Post"].Text)
	standings.MarginLess3 = parseWinLoss(rowMap["3"].Text)
	standings.MarginGreater10 = parseWinLoss(rowMap["10"].Text)
	standings.October = parseWinLoss(rowMap["Oct"].Text)
	standings.November = parseWinLoss(rowMap["Nov"].Text)
	standings.December = parseWinLoss(rowMap["Dec"].Text)
	standings.January = parseWinLoss(rowMap["Jan"].Text)
	standings.February = parseWinLoss(rowMap["Feb"].Text)
	standings.March = parseWinLoss(rowMap["Mar"].Text)
	standings.April = parseWinLoss(rowMap["Apr"].Text)
	return
}

func parseWinLoss(s string) (wl WinLoss) {
	p := strings.Split(s, "-")
	wl.Wins, _ = strconv.Atoi(p[0])
	wl.Losses, _ = strconv.Atoi(p[1])
	return
}
