package parser

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/nschimek/nba-scraper/core"
	"github.com/nschimek/nba-scraper/model"
)

type GameStatsParser struct{}

// Parse the Scorebox - if we can't get a team ID or score, this will return an error.
// Any issues with parsing the win-loss record will instead log warnings.
func (*GameStatsParser) parseScorebox(box *colly.HTMLElement) (model.GameTeam, error) {
	var err error
	gt := &model.GameTeam{}

	gt.TeamId, err = ParseTeamId(box.ChildAttr("div:first-child strong a", "href"))
	if err != nil {
		return model.GameTeam{}, err
	}

	gt.Score, err = strconv.Atoi(box.ChildText("div.scores div.score"))
	if err != nil {
		return model.GameTeam{}, errors.New("could not parse score from scorebox")
	}

	wl := strings.Split(box.ChildText("div:nth-child(3)"), "-")
	if len(wl) == 2 {
		gt.Wins, err = strconv.Atoi(wl[0])
		if err != nil {
			core.Log.Warnf("could not parse wins from scorebox for team %s", gt.TeamId)
		}
		gt.Losses, err = strconv.Atoi(wl[1])
		if err != nil {
			core.Log.Warnf("could not parse losses from scorebox for team %s", gt.TeamId)
		}
	} else {
		core.Log.Warnf("could not parse win-loss record from scorebox for team %s", gt.TeamId)
	}

	return *gt, nil
}

func (*GameStatsParser) parseLineScoreTable(tbl *colly.HTMLElement, gameId string) (home, visitor []model.GameLineScore, err error) {
	tableMaps := Table(tbl) // row 0 will be away, row 1 will be home

	visitor, err = lineScoreFromRow(tableMaps[0], gameId)

	if err != nil {
		return nil, nil, err
	}

	home, err = lineScoreFromRow(tableMaps[1], gameId)

	return
}

func (*GameStatsParser) parseFourFactorsTable(tbl *colly.HTMLElement, gameId string) (home, visitor model.GameFourFactor, err error) {
	tableMaps := Table(tbl) // row 0 will be away, row 1 will be home

	if len(tableMaps) != 2 {
		return model.GameFourFactor{}, model.GameFourFactor{}, errors.New("unexpected number of four factor rows")
	}

	visitor, err = gameFourFactorsFromRow(tableMaps[0])

	if err != nil {
		return model.GameFourFactor{}, model.GameFourFactor{}, err
	}

	visitor.GameId = gameId
	home, err = gameFourFactorsFromRow(tableMaps[1])
	home.GameId = gameId

	return
}

func lineScoreFromRow(rowMap map[string]*colly.HTMLElement, gameId string) (scores []model.GameLineScore, err error) {
	teamId, err := ParseTeamId(parseLink(rowMap["team"]))

	for key, cell := range rowMap {
		// loop through all non-team and total columns; those that remain are the quarters
		if key != "team" && key != "T" {
			score, e := strconv.Atoi(cell.Text)
			if e == nil {
				scores = append(scores, model.GameLineScore{
					GameId:  gameId,
					TeamId:  teamId,
					Quarter: lineScoreQuarter(key),
					Score:   score,
				})
			}
		}
	}

	if len(scores) == 0 {
		err = errors.New("did not parse any scores from line score row")
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

func gameFourFactorsFromRow(rowMap map[string]*colly.HTMLElement) (model.GameFourFactor, error) {
	var err error
	gff := new(model.GameFourFactor)

	gff.TeamId, err = ParseTeamId(parseLink(rowMap["team_id"]))

	if err != nil {
		return model.GameFourFactor{}, err
	}

	gff.Pace, _ = strconv.ParseFloat(rowMap["pace"].Text, 64)
	gff.EffectiveFgPct, _ = strconv.ParseFloat(getColumnText(rowMap, "efg_pct"), 64)
	gff.TurnoverPct, _ = strconv.ParseFloat(getColumnText(rowMap, "tov_pct"), 64)
	gff.OffensiveRbPct, _ = strconv.ParseFloat(getColumnText(rowMap, "orb_pct"), 64)
	gff.FtPerFga, _ = strconv.ParseFloat(getColumnText(rowMap, "ft_rate"), 64)
	gff.OffensiveRating, _ = strconv.ParseFloat(getColumnText(rowMap, "off_rtg"), 64)

	return *gff, nil
}
