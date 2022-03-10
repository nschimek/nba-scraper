package scraper

import (
	"errors"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

const (
	BasePath = "leagues"
)

type ScrapedSchedule struct {
	StartTime time.Time
	GameId    string
	GameUrl   string
}

func Schedule(season string, startDate, endDate time.Time) (map[string]string, error) {
	months, err := getMonths(startDate, endDate)

	if err != nil {
		return nil, err
	}

	c := colly.NewCollector()
	schedules := []ScrapedSchedule{}

	c.OnHTML("body #wrap #content #all_schedule #div_schedule table tbody", func(tbl *colly.HTMLElement) {
		tbl.ForEach("tr", func(_ int, tr *colly.HTMLElement) {
			schedules = append(schedules, parseScheduleRows(tr, startDate, endDate)...)
		})
	})

	for _, month := range months {
		c.Visit(getMonthUrl(month, season))
	}

	return buildGameMap(schedules), nil
}

func parseScheduleRows(tr *colly.HTMLElement, startDate, endDate time.Time) []ScrapedSchedule {
	schedules := []ScrapedSchedule{}

	if tr.Attr("class") != "thead" {
		schedule := ScrapedSchedule{}
		var parsedDate, parsedTime string

		parsedDate = tr.ChildText("th a")

		tr.ForEach("td", func(_ int, td *colly.HTMLElement) {
			switch td.Attr("data-stat") {
			case "box_score_text":
				schedule.GameUrl = td.ChildAttr("a", "href")
				schedule.GameId = strings.Replace(strings.Split(td.ChildAttr("a", "href"), "/")[2], ".html", "", 1)
			case "game_start_time":
				parsedTime = strings.Replace(td.Text, "p", " PM EST", 1)
			}
		})

		schedule.StartTime, _ = time.ParseInLocation("Mon, Jan 2, 2006 3:04 PM EST", parsedDate+" "+parsedTime, EST)

		if schedule.StartTime.After(startDate) && schedule.StartTime.Before(endDate) {
			schedules = append(schedules, schedule)
		}
	}

	return schedules
}

func buildGameMap(schedules []ScrapedSchedule) map[string]string {
	gameMap := make(map[string]string)
	for _, schedule := range schedules {
		gameMap[schedule.GameId] = schedule.GameUrl
	}
	return gameMap
}

func getMonthUrl(month time.Month, season string) string {
	monthString := strings.ToLower(month.String())
	return BaseHttp + "/" + BasePath + "/NBA_" + season + "_games-" + monthString + ".html"
}

func getMonths(startDate, endDate time.Time) ([]time.Month, error) {
	if endDate.Before(startDate) {
		return nil, errors.New("end date is before start date")
	}

	months := []time.Month{}
	startMonth := time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, time.Local)
	endMonth := time.Date(endDate.Year(), endDate.Month(), 1, 0, 0, 0, 0, time.Local)

	for startMonth.Before(endMonth) || startMonth.Equal(endMonth) {
		months = append(months, startMonth.Month())
		startMonth = startMonth.AddDate(0, 1, 0)
	}

	return months, nil
}
