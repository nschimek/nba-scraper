package parser

import (
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

type Schedule struct {
	StartTime                                  time.Time
	GameId, GameUrl, VisitorTeamId, HomeTeamId string
	Played                                     bool
}

func ScheduleTable(tbl *colly.HTMLElement, startDate, endDate time.Time) []Schedule {
	schedules := []Schedule{}

	for _, row := range Table(tbl) {
		s := mapScheduleRow(row)
		if s.Played && s.StartTime.After(startDate) && s.StartTime.Before(endDate) {
			schedules = append(schedules, s)
		}
	}

	return schedules
}

// the row map here has the data-stat attribute as the key and the colly HTML Element (cell) as the value
func mapScheduleRow(r map[string]*colly.HTMLElement) (s Schedule) {
	parsedDate := r["date_game"].ChildText("a")
	parsedTime := strings.Replace(r["game_start_time"].Text, "p", " PM EST", 1)

	s.StartTime, _ = time.ParseInLocation("Mon, Jan 2, 2006 3:04 PM EST", parsedDate+" "+parsedTime, EST)
	s.VisitorTeamId = ParseTeamId(parseLink(r["visitor_team_name"]))
	s.HomeTeamId = ParseTeamId(parseLink(r["home_team_name"]))

	s.GameUrl = parseLink(r["box_score_text"])

	if s.GameUrl != "" {
		s.GameId = ParseLastId(s.GameUrl)
		s.Played = true
	}

	return
}
