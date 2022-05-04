package parser

import (
	"time"

	"github.com/gocolly/colly/v2"
)

type Injuries struct {
	PlayerId, PlayerUrl, TeamId string
	Season                      int
	UpdateDate                  time.Time
	Description                 string
}

func InjuriesTable(tbl *colly.HTMLElement, season int) []Injuries {
	injuries := []Injuries{}

	for _, rowMap := range Table(tbl) {
		injuries = append(injuries, injuryFromRow(rowMap, season))
	}

	return injuries
}

func injuryFromRow(rowMap map[string]*colly.HTMLElement, season int) (injury Injuries) {
	injury.Season = season
	injury.PlayerUrl = parseLink(rowMap["player"])
	injury.PlayerId = ParseLastId(injury.PlayerUrl)
	injury.TeamId = ParseTeamId(parseLink(rowMap["team_name"]))
	injury.UpdateDate, _ = time.Parse("Mon, Jan 2, 2006", rowMap["date_update"].Text)
	injury.Description = rowMap["note"].Text
	return
}
