package parser

import (
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/nschimek/nba-scraper/core"
	"github.com/nschimek/nba-scraper/model"
)

type InjuryParser struct {
	Config *core.Config `Inject:""`
}

func (p *InjuryParser) InjuriesTable(tbl *colly.HTMLElement) []model.PlayerInjury {
	injuries := []model.PlayerInjury{}

	for _, rowMap := range Table(tbl) {
		inj, err := injuryFromRow(rowMap)
		inj.CaptureError(err)
		inj.Season = p.Config.Season
		injuries = append(injuries, inj)
	}

	if len(injuries) == 0 {
		core.Log.Warn("did not scrape any player injuries...")
	}

	return injuries
}

func injuryFromRow(rowMap map[string]*colly.HTMLElement) (model.PlayerInjury, error) {
	var err error
	inj := new(model.PlayerInjury)

	inj.PlayerId, err = ParseLastId(parseLink(rowMap["player"]))
	inj.TeamId, err = ParseTeamId(parseLink(rowMap["team_name"]))
	inj.SourceUpdateDate, _ = time.ParseInLocation("Mon, Jan 2, 2006", rowMap["date_update"].Text, CST)
	inj.Description = rowMap["note"].Text

	return *inj, err
}
