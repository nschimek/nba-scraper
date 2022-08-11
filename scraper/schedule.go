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
)

const (
	baseScheduleTableElement = "body #wrap #content #all_schedule #div_schedule table tbody" // targets the main schedule table
)

type ScheduleScraper struct {
	Config         *core.Config           `Inject:""`
	Colly          *colly.Collector       `Inject:""`
	ScheduleParser *parser.ScheduleParser `Inject:""`
	dateRange      *DateRange
	ScrapedData    []model.Schedule
	GameIds        map[string]struct{}
}

type DateRange struct {
	startDate, endDate time.Time
}

func (s *ScheduleScraper) ScrapeDateRange(startDate, endDate time.Time) {
	months, err := getMonths(startDate, endDate)

	if err != nil {
		core.Log.Fatal(err)
	}

	s.dateRange = &DateRange{startDate: startDate, endDate: endDate}
	s.Scrape(months...)
}

// Scraper interface methods
func (s *ScheduleScraper) GetData() interface{} {
	return s.ScrapedData
}

func (s *ScheduleScraper) Scrape(pageIds ...string) {
	s.GameIds = make(map[string]struct{})
	c := s.Colly.Clone()
	c.OnRequest(onRequestVisit)
	c.OnError(onError)

	s.Colly.OnHTML(baseScheduleTableElement, func(tbl *colly.HTMLElement) {
		for _, ps := range s.ScheduleParser.ScheduleTable(tbl, s.dateRange.startDate, s.dateRange.endDate) {
			s.ScrapedData = append(s.ScrapedData, ps)
			s.GameIds[ps.GameId] = exists
		}
	})

	for _, id := range pageIds {
		s.Colly.Visit(s.getUrl(id))
	}
}

func getMonths(startDate, endDate time.Time) ([]string, error) {
	if endDate.Before(startDate) {
		return nil, errors.New("end date is before start date")
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
