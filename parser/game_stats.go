package parser

import (
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

type GameTeam struct {
	TeamId, TeamUrl, Result string
	Wins, Losses, Score     int
}

type GameFourFactors struct {
	TeamId, TeamUrl                                                              string
	Pace, EffectiveFgPct, TurnoverPct, OffensiveRbPct, FtPerFga, OffensiveRating float64
}

type GameLineScore struct {
	TeamId, TeamUrl string
	Quarter, Score  int
}

func parseScorebox(box *colly.HTMLElement) (gt GameTeam) {
	gt.TeamUrl = box.Request.AbsoluteURL(box.ChildAttr("div:first-child strong a", "href"))
	gt.TeamId = parseTeamId(gt.TeamUrl)
	gt.Score, _ = strconv.Atoi(box.ChildText("div.scores div.score"))

	wl := strings.Split(box.ChildText("div:nth-child(3)"), "-")
	gt.Wins, _ = strconv.Atoi(wl[0])
	gt.Losses, _ = strconv.Atoi(wl[1])

	return
}

func parseLineScoreTable(tbl *colly.HTMLElement) (home, visitor []GameLineScore) {
	tableMaps := Table(tbl) // row 0 will be away, row 1 will be home

	visitor = lineScoreFromRow(tableMaps[0])
	home = lineScoreFromRow(tableMaps[1])

	return
}

func parseFourFactorsTable(tbl *colly.HTMLElement) (home, visitor GameFourFactors) {
	tableMaps := Table(tbl) // row 0 will be away, row 1 will be home

	visitor = gameFourFactorsFromRow(tableMaps[0])
	home = gameFourFactorsFromRow(tableMaps[1])

	return
}

func lineScoreFromRow(rowMap map[string]*colly.HTMLElement) (scores []GameLineScore) {
	teamUrl := parseLink(rowMap["team"])
	teamId := parseTeamId(teamUrl)

	for key, cell := range rowMap {
		// loop through all non-team and total columns; those that remain are the quarters
		if key != "team" && key != "T" {
			score, _ := strconv.Atoi(cell.Text)
			scores = append(scores, GameLineScore{
				TeamId:  teamId,
				TeamUrl: teamUrl,
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
