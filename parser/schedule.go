package parser

import (
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/nschimek/nba-scraper/model"
)

type ScheduleParser struct{}

func (*ScheduleParser) ScheduleTable(tbl *colly.HTMLElement, startDate, endDate time.Time) []model.Schedule {
	schedules := []model.Schedule{}

	for _, row := range Table(tbl) {
		s, err := mapScheduleRow(row)
		s.CaptureError(err)
		if s.Played && s.StartTime.After(startDate) && s.StartTime.Before(endDate) {
			schedules = append(schedules, s)
		}
	}

	return schedules
}

// the row map here has the data-stat attribute as the key and the colly HTML Element (cell) as the value
func mapScheduleRow(r map[string]*colly.HTMLElement) (model.Schedule, error) {
	var err error
	s := new(model.Schedule)

	parsedDate := r["date_game"].ChildText("a")
	parsedTime := strings.Replace(r["game_start_time"].Text, "p", " PM EST", 1)

	s.StartTime, err = time.ParseInLocation("Mon, Jan 2, 2006 3:04 PM EST", parsedDate+" "+parsedTime, EST)
	s.VisitorTeamId, err = ParseTeamId(parseLink(r["visitor_team_name"]))
	s.HomeTeamId, err = ParseTeamId(parseLink(r["home_team_name"]))

	gameUrl := parseLink(r["box_score_text"])

	if gameUrl != "" {
		s.GameId, err = ParseLastId(gameUrl)
		s.Played = true
	}

	return *s, err
}
