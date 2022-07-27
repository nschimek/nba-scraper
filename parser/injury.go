package parser

import (
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/nschimek/nba-scraper/core"
	"github.com/nschimek/nba-scraper/model"
)

type Injuries struct {
	PlayerId, PlayerUrl, TeamId string
	Season                      int
	UpdateDate                  time.Time
	Description                 string
}

type InjuryParser struct {
	Config *core.Config `Inject:""`
}

func (p *InjuryParser) InjuriesTable(tbl *colly.HTMLElement) []model.PlayerInjury {
	injuries := []model.PlayerInjury{}

	for _, rowMap := range Table(tbl) {
		inj := injuryFromRow(rowMap)
		inj.Season = p.Config.Season
		injuries = append(injuries, *inj)
	}

	return injuries
}

func injuryFromRow(rowMap map[string]*colly.HTMLElement) *model.PlayerInjury {
	inj := new(model.PlayerInjury)

	inj.PlayerId = ParseLastId(parseLink(rowMap["player"]))
	inj.TeamId = ParseTeamId(parseLink(rowMap["team_name"]))
	inj.SourceUpdateDate, _ = time.ParseInLocation("Mon, Jan 2, 2006", rowMap["date_update"].Text, CST)
	inj.Description = rowMap["note"].Text

	return inj
}
