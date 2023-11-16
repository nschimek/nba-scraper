package parser

import (
	"errors"
	"strconv"

	"github.com/gocolly/colly/v2"
	"github.com/nschimek/nba-scraper/core"
	"github.com/nschimek/nba-scraper/model"
)

type TeamParser struct {
	Config *core.Config `Inject:""`
}

func (p *TeamParser) TeamInfoBox(t *model.Team, box *colly.HTMLElement) {
	t.Name = box.ChildText("div:nth-child(2) > h1 > span:nth-child(2)")
	if t.Name == "" {
		t.CaptureError(errors.New("could not parse team name"))
	}
}

func (p *TeamParser) TeamPlayerTable(t *model.Team, tbl *colly.HTMLElement) {
	t.TeamPlayers = p.parseTeamRosterTable(tbl, t.ID)
}

func (p *TeamParser) TeamSalariesTable(t *model.Team, tbl *colly.HTMLElement) {
	t.TeamPlayerSalaries = p.parseTeamSalariesTable(tbl, t.ID)
}

func (p *TeamParser) parseTeamRosterTable(tbl *colly.HTMLElement, teamId string) []model.TeamPlayer {
	roster := []model.TeamPlayer{}

	for _, rowMap := range Table(tbl) {
		tp, err := teamPlayerFromRow(rowMap)
		if err == nil {
			tp.TeamId = teamId
			tp.Season = p.Config.Season
			roster = append(roster, tp)
		} else {
			core.Log.Warnf("error encountered parsing team player list: %v", err)
		}
	}

	if len(roster) == 0 {
		core.Log.Warn("did not parse any team players")
	}

	return roster
}

func teamPlayerFromRow(rowMap map[string]*colly.HTMLElement) (model.TeamPlayer, error) {
	var err error
	tp := new(model.TeamPlayer)

	tp.PlayerId, err = ParseLastId(parseLink(rowMap["player"]))

	if tp.PlayerId == "" {
		return model.TeamPlayer{}, errors.New("player ID is blank")
	}

	tp.Number, _ = strconv.Atoi(getColumnText(rowMap, "number"))
	tp.Position = getColumnText(rowMap, "pos")

	return *tp, err
}

func (p *TeamParser) parseTeamSalariesTable(tbl *colly.HTMLElement, teamId string) []model.TeamPlayerSalary {
	salaries := []model.TeamPlayerSalary{}

	for _, rowMap := range Table(tbl) {
		tps, err := teamSalaryFromRow(rowMap)
		if err == nil {
			tps.TeamId = teamId
			tps.Season = p.Config.Season
			salaries = append(salaries, tps)
		} else {
			core.Log.Warnf("error encountered parsing team player salary list: %v", err)
		}
	}

	if len(salaries) == 0 {
		core.Log.Warn("did not parse any team player salaries")
	}

	return salaries
}

func teamSalaryFromRow(rowMap map[string]*colly.HTMLElement) (model.TeamPlayerSalary, error) {
	var err error
	tps := new(model.TeamPlayerSalary)

	tps.PlayerId, err = ParseLastId(parseLink(rowMap["player"]))

	if tps.PlayerId == "" {
		return model.TeamPlayerSalary{}, errors.New("player ID is blank")
	}

	tps.Salary, _ = strconv.Atoi(rowMap["salary"].Attr("csk"))
	tps.Rank, _ = strconv.Atoi(getColumnText(rowMap, "ranker"))

	return *tps, err
}
