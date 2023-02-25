package scraper

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/nschimek/nba-scraper/core"
	"github.com/nschimek/nba-scraper/model"
	"github.com/nschimek/nba-scraper/parser"
	"github.com/nschimek/nba-scraper/repository"
)

const (
	baseScheduleTableElement = "body #wrap #content #all_schedule #div_schedule table tbody" // targets the main schedule table
)

type ScheduleScraper struct {
	Config         *core.Config               `Inject:""`
	Colly          *colly.Collector           `Inject:""`
	ScheduleParser *parser.ScheduleParser     `Inject:""`
	GameRepository *repository.GameRepository `Inject:""`
	dateRange      *DateRange
	ScrapedData    []model.Schedule
	GameIds        map[string]struct{}
}

type DateRange struct {
	startDate, endDate time.Time
}

func (s *ScheduleScraper) ScrapeDateRange(startDate, endDate time.Time) {
	if startDate.IsZero() {
		startDate, _ = s.GameRepository.GetMostRecentGame()
	}
	if endDate.IsZero() {
		endDate = time.Now()
	} else {
		endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, parser.EST)
	}

	core.Log.Infof("Schedule Scraper started with date range: %s - %s",
		startDate.Format(core.DateRangeFormat), endDate.Format(core.DateRangeFormat))

	months, err := getMonths(startDate, endDate)

	if err != nil {
		core.Log.Error(err)
	} else {
		s.dateRange = &DateRange{startDate: startDate, endDate: endDate}
		s.Scrape(months...)
	}
}

func (s *ScheduleScraper) Scrape(pageIds ...string) {
	c := core.CloneColly(s.Colly)
	s.GameIds = make(map[string]struct{})

	c.OnHTML(baseScheduleTableElement, func(tbl *colly.HTMLElement) {
		for _, ps := range s.ScheduleParser.ScheduleTable(tbl, s.dateRange.startDate, s.dateRange.endDate) {
			if !ps.HasErrors() {
				s.ScrapedData = append(s.ScrapedData, ps)
				s.GameIds[ps.GameId] = exists
			} else {
				ps.LogErrors()
			}
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		core.Log.Error(NewScraperError(err, r.Request.URL.String()))
	})

	for _, id := range pageIds {
		url := s.getUrl(id)
		// really ugly covid hack...the NBA played games in both October 2019 and October 2020 this season!
		if s.Config.Season == 2020 && id == "october" {
			url = strings.Replace(url, ".html", fmt.Sprintf("-%d.html", s.dateRange.endDate.Year()), 1)
		}
		c.Visit(url)
	}

	if len(s.GameIds) > 0 {
		core.Log.WithField("gameIds", len(s.GameIds)).Info("Successfully scraped game IDs from Schedule!")
	} else {
		core.Log.Warn("No game IDs scraped from given date range...")
	}
}

func getMonths(startDate, endDate time.Time) ([]string, error) {
	if endDate.Before(startDate) {
		return nil, errors.New("end date is before start date or invalid")
	}

	months := []string{}
	startMonth := time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, time.Local)
	endMonth := time.Date(endDate.Year(), endDate.Month(), 1, 0, 0, 0, 0, time.Local)

	for startMonth.Before(endMonth) || startMonth.Equal(endMonth) {
		months = append(months, strings.ToLower(startMonth.Month().String()))
		startMonth = startMonth.AddDate(0, 1, 0)
	}

	return months, nil
}

func (s *ScheduleScraper) getUrl(month string) string {
	return fmt.Sprintf(scheduleIdPage, s.Config.Season, month)
}
