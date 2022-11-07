package parser

import (
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
		tp := teamPlayerFromRow(rowMap)
		tp.TeamId = teamId
		tp.Season = p.Config.Season
		roster = append(roster, *tp)
	}

	return roster
}

func teamPlayerFromRow(rowMap map[string]*colly.HTMLElement) *model.TeamPlayer {
	tp := new(model.TeamPlayer)

	tp.Number, _ = strconv.Atoi(rowMap["number"].Text)
	tp.PlayerId, _ = ParseLastId(parseLink(rowMap["player"]))
	tp.Position = rowMap["pos"].Text

	return tp
}

func (p *TeamParser) parseTeamSalariesTable(tbl *colly.HTMLElement, teamId string) []model.TeamPlayerSalary {
	salaries := []model.TeamPlayerSalary{}

	for _, rowMap := range Table(tbl) {
		tps := teamSalaryFromRow(rowMap)
		tps.TeamId = teamId
		tps.Season = p.Config.Season
		salaries = append(salaries, *tps)
	}

	return salaries
}

func teamSalaryFromRow(rowMap map[string]*colly.HTMLElement) *model.TeamPlayerSalary {
	tps := new(model.TeamPlayerSalary)

	tps.PlayerId, _ = ParseLastId(parseLink(rowMap["player"]))
	tps.Salary, _ = strconv.Atoi(rowMap["salary"].Attr("csk"))
	tps.Rank, _ = strconv.Atoi(rowMap["ranker"].Text)

	return tps
}
