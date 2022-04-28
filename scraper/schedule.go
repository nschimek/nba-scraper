package scraper

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/nschimek/nba-scraper/parser"
)

const (
	baseTableElement = "body #wrap #content #all_schedule #div_schedule table tbody" // targets the main schedule table
)

type Schedule struct {
	StartTime                         time.Time
	GameId, VisitorTeamId, HomeTeamId string
	Played                            bool
}

type ScheduleScraper struct {
	colly       colly.Collector
	season      int
	dateRange   DateRange
	urls        []string
	ScrapedData []parser.Schedule
	Errors      []error
	child       Scraper
	childUrls   map[string]string
}

type DateRange struct {
	startDate, endDate time.Time
}

func CreateScheduleScraperWithDates(c *colly.Collector, season int, startDate, endDate time.Time) (ScheduleScraper, error) {
	dateRange := DateRange{startDate: startDate, endDate: endDate}
	urls, err := dateRangeToUrls(season, dateRange)

	if err != nil {
		return ScheduleScraper{}, err
	}

	scraper := CreateScheduleScraper(c, season)
	scraper.dateRange = dateRange
	scraper.urls = urls

	return scraper, nil
}

func CreateScheduleScraper(c *colly.Collector, season int) ScheduleScraper {
	return ScheduleScraper{
		colly:     *c,
		season:    season,
		childUrls: make(map[string]string),
	}
}

// Scraper interface methods
func (s *ScheduleScraper) GetData() interface{} {
	return s.ScrapedData
}

func (s *ScheduleScraper) AttachChild(scraper *Scraper) {
	s.child = *scraper
}

func (s *ScheduleScraper) GetChild() Scraper {
	return s.child
}

func (s *ScheduleScraper) GetChildUrls() []string {
	return urlsMapToArray(s.childUrls)
}

func (s *ScheduleScraper) Scrape(urls ...string) {
	s.urls = append(s.urls, urls...)

	s.colly.OnHTML(baseTableElement, func(tbl *colly.HTMLElement) {
		for _, ps := range parser.ScheduleTable(tbl, s.dateRange.startDate, s.dateRange.endDate) {
			s.ScrapedData = append(s.ScrapedData, ps)
			s.childUrls[ps.GameId] = tbl.Request.AbsoluteURL(ps.GameUrl)
		}
	})

	for _, url := range s.urls {
		s.colly.Visit(url)
	}

	scrapeChild(s)
}

func dateRangeToUrls(season int, dateRange DateRange) ([]string, error) {
	months, err := getMonths(dateRange.startDate, dateRange.endDate)

	if err != nil {
		return nil, err
	}

	urls := []string{}

	for _, month := range months {
		urls = append(urls, getMonthUrl(month, season))
	}

	return urls, nil
}

func getMonthUrl(month time.Month, season int) string {
	monthString := strings.ToLower(month.String())
	return BaseHttp + "/" + baseLeaguesPath + "/NBA_" + strconv.Itoa(season) + "_games-" + monthString + ".html"
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
