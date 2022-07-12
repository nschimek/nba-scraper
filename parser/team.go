package parser

import (
	"strconv"

	"github.com/gocolly/colly/v2"
	"github.com/nschimek/nba-scraper/model"
)

type TeamParser struct{}

func (*TeamParser) TeamInfoBox(t *model.Team, box *colly.HTMLElement) {
	t.Name = box.ChildText("div:nth-child(2) > h1 > span:nth-child(2)")
}

func (*TeamParser) TeamPlayerTable(t *model.Team, tbl *colly.HTMLElement) {
	t.TeamPlayers = parseTeamRosterTable(tbl, t.ID)
}

func (*TeamParser) TeamSalariesTable(t *model.Team, tbl *colly.HTMLElement) {
	t.TeamPlayerSalaries = parseTeamSalariesTable(tbl, t.ID)
}

func parseTeamRosterTable(tbl *colly.HTMLElement, teamId string) []model.TeamPlayer {
	roster := []model.TeamPlayer{}

	for _, rowMap := range Table(tbl) {
		tp := teamPlayerFromRow(rowMap)
		tp.TeamId = teamId
		roster = append(roster, *tp)
	}

	return roster
}

func teamPlayerFromRow(rowMap map[string]*colly.HTMLElement) *model.TeamPlayer {
	tr := new(model.TeamPlayer)

	tr.Number, _ = strconv.Atoi(rowMap["number"].Text)
	tr.PlayerId = ParseLastId(parseLink(rowMap["player"]))
	tr.Position = rowMap["pos"].Text

	return tr
}

func parseTeamSalariesTable(tbl *colly.HTMLElement, teamId string) []model.TeamPlayerSalary {
	salaries := []model.TeamPlayerSalary{}

	for _, rowMap := range Table(tbl) {
		tps := teamSalaryFromRow(rowMap)
		tps.TeamId = teamId
		salaries = append(salaries, *tps)
	}

	return salaries
}

func teamSalaryFromRow(rowMap map[string]*colly.HTMLElement) *model.TeamPlayerSalary {
	tps := new(model.TeamPlayerSalary)

	tps.PlayerId = ParseLastId(parseLink(rowMap["player"]))
	tps.Salary, _ = strconv.Atoi(rowMap["salary"].Attr("csk"))
	tps.Rank, _ = strconv.Atoi(rowMap["ranker"].Text)

	return tps
}
