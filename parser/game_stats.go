package parser

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/nschimek/nba-scraper/model"
)

type GameStatsParser struct{}

func (*GameStatsParser) parseScorebox(box *colly.HTMLElement) (model.GameTeam, error) {
	var err error
	gt := &model.GameTeam{}

	gt.TeamId = ParseTeamId(box.ChildAttr("div:first-child strong a", "href"))
	gt.Score, err = strconv.Atoi(box.ChildText("div.scores div.score"))

	if err != nil {
		return model.GameTeam{}, errors.New("could not parse score from scorebox")
	}

	wl := strings.Split(box.ChildText("div:nth-child(3)"), "-")
	if len(wl) == 2 {
		gt.Wins, err = strconv.Atoi(wl[0])
		if err != nil {
			return model.GameTeam{}, errors.New("could not parse wins from scorebox")
		}
		gt.Losses, err = strconv.Atoi(wl[1])
		if err != nil {
			return model.GameTeam{}, errors.New("could not parse losses from scorebox")
		}
	} else {
		return model.GameTeam{}, errors.New("could not parse win-loss record from scorebox")
	}

	return *gt, nil
}

func (*GameStatsParser) parseLineScoreTable(tbl *colly.HTMLElement, gameId string) (home, visitor []model.GameLineScore) {
	tableMaps := Table(tbl) // row 0 will be away, row 1 will be home

	visitor = lineScoreFromRow(tableMaps[0], gameId)
	home = lineScoreFromRow(tableMaps[1], gameId)

	return
}

func (*GameStatsParser) parseFourFactorsTable(tbl *colly.HTMLElement, gameId string) (home, visitor model.GameFourFactor) {
	tableMaps := Table(tbl) // row 0 will be away, row 1 will be home

	visitor = gameFourFactorsFromRow(tableMaps[0])
	visitor.GameId = gameId
	home = gameFourFactorsFromRow(tableMaps[1])
	home.GameId = gameId

	return
}

func lineScoreFromRow(rowMap map[string]*colly.HTMLElement, gameId string) (scores []model.GameLineScore) {
	teamId := ParseTeamId(parseLink(rowMap["team"]))

	for key, cell := range rowMap {
		// loop through all non-team and total columns; those that remain are the quarters
		if key != "team" && key != "T" {
			score, _ := strconv.Atoi(cell.Text)
			scores = append(scores, model.GameLineScore{
				GameId:  gameId,
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

func gameFourFactorsFromRow(rowMap map[string]*colly.HTMLElement) model.GameFourFactor {
	gff := new(model.GameFourFactor)

	gff.TeamId = ParseTeamId(parseLink(rowMap["team_id"]))
	gff.Pace, _ = strconv.ParseFloat(rowMap["pace"].Text, 64)
	gff.EffectiveFgPct, _ = strconv.ParseFloat(rowMap["efg_pct"].Text, 64)
	gff.TurnoverPct, _ = strconv.ParseFloat(rowMap["tov_pct"].Text, 64)
	gff.OffensiveRbPct, _ = strconv.ParseFloat(rowMap["orb_pct"].Text, 64)
	gff.FtPerFga, _ = strconv.ParseFloat(rowMap["ft_rate"].Text, 64)
	gff.OffensiveRating, _ = strconv.ParseFloat(rowMap["off_rtg"].Text, 64)

	return *gff
}
