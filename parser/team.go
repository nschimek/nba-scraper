package parser

import (
	"strconv"

	"github.com/gocolly/colly/v2"
)

type Team struct {
	Id, Name           string
	Season             int
	TeamPlayers        []TeamPlayer
	TeamPlayerSalaries []TeamPlayerSalary
}

type TeamPlayer struct {
	PlayerId, Name, Position string
	Number                   int
}

type TeamPlayerSalary struct {
	PlayerId     string
	Salary, Rank int
}

func (t *Team) TeamInfoBox(box *colly.HTMLElement) {
	t.Name = box.ChildText("div:nth-child(2) > h1 > span:nth-child(2)")
}

func (t *Team) TeamPlayerTable(tbl *colly.HTMLElement) {
	t.TeamPlayers = parseTeamRosterTable(tbl)
}

func (t *Team) TeamSalariesTable(tbl *colly.HTMLElement) {
	t.TeamPlayerSalaries = parseTeamSalariesTable(tbl)
}

func parseTeamRosterTable(tbl *colly.HTMLElement) []TeamPlayer {
	roster := []TeamPlayer{}

	for _, rowMap := range Table(tbl) {
		roster = append(roster, *teamPlayerFromRow(rowMap))
	}

	return roster
}

func teamPlayerFromRow(rowMap map[string]*colly.HTMLElement) *TeamPlayer {
	tr := new(TeamPlayer)

	tr.Number, _ = strconv.Atoi(rowMap["number"].Text)
	tr.PlayerId = parsePlayerId(parseLink(rowMap["player"]))
	tr.Name = rowMap["player"].Text
	tr.Position = rowMap["pos"].Text

	return tr
}

func parseTeamSalariesTable(tbl *colly.HTMLElement) []TeamPlayerSalary {
	salaries := []TeamPlayerSalary{}

	for _, rowMap := range Table(tbl) {
		salaries = append(salaries, *teamSalaryFromRow(rowMap))
	}

	return salaries
}

func teamSalaryFromRow(rowMap map[string]*colly.HTMLElement) *TeamPlayerSalary {
	tps := new(TeamPlayerSalary)

	tps.PlayerId = parsePlayerId(parseLink(rowMap["player"]))
	tps.Salary, _ = strconv.Atoi(rowMap["salary"].Attr("csk"))
	tps.Rank, _ = strconv.Atoi(rowMap["ranker"].Text)

	return tps
}
