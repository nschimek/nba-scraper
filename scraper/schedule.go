package scraper

import (
	"errors"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

const (
	BasePath         = "leagues"
	baseTableElement = "body #wrap #content #all_schedule #div_schedule table tbody" // targets the main schedule table
)

type Schedule struct {
	StartTime                         time.Time
	GameId, VisitorTeamId, HomeTeamId string
	Played                            bool
}

type ScheduleScraper struct {
	colly       colly.Collector
	season      string
	dateRange   DateRange
	urls        []string
	ScrapedData []Schedule
	Errors      []error
	Child       Scraper
	childUrls   map[string]string
}

type DateRange struct {
	startDate, endDate time.Time
}

func CreateScheduleScraperWithDates(c *colly.Collector, season string, startDate, endDate time.Time) (ScheduleScraper, error) {
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

func CreateScheduleScraper(c *colly.Collector, season string) ScheduleScraper {
	return ScheduleScraper{
		colly:     *c,
		season:    season,
		childUrls: make(map[string]string),
	}
}

func (s *ScheduleScraper) GetData() []Schedule {
	return s.ScrapedData
}

func (s *ScheduleScraper) AttachChild(scraper *Scraper) {
	s.Child = *scraper
}

func (s *ScheduleScraper) GetChild() Scraper {
	return s.Child
}

func (s *ScheduleScraper) Scrape(urls ...string) {
	s.urls = append(s.urls, urls...)

	s.colly.OnHTML(baseTableElement, func(tbl *colly.HTMLElement) {
		tbl.ForEach("tr", func(_ int, tr *colly.HTMLElement) {
			s.parseRow(tr)
		})
	})

	for _, url := range s.urls {
		s.colly.Visit(url)
	}

	if s.GetChild() != nil && len(s.childUrls) > 0 {
		urls := []string{}
		for _, url := range s.childUrls {
			urls = append(urls, url)
		}
		s.GetChild().Scrape(urls...)
	}
}

func dateRangeToUrls(season string, dateRange DateRange) ([]string, error) {
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

func (s *ScheduleScraper) parseRow(tr *colly.HTMLElement) {
	if tr.Attr("class") != "thead" {
		schedule := Schedule{}
		var parsedDate, parsedTime string

		parsedDate = tr.ChildText("th a")

		tr.ForEach("td", parseColumnCallback(schedule, parsedTime, s.childUrls))

		schedule.StartTime, _ = time.ParseInLocation("Mon, Jan 2, 2006 3:04 PM EST", parsedDate+" "+parsedTime, EST)

		if schedule.StartTime.After(s.dateRange.startDate) && schedule.StartTime.Before(s.dateRange.endDate) {
			s.ScrapedData = append(s.ScrapedData, schedule)
		}
	}
}

func parseColumnCallback(schedule Schedule, parsedTime string, gameUrls map[string]string) func(int, *colly.HTMLElement) {
	return func(_ int, td *colly.HTMLElement) {
		switch td.Attr("data-stat") {
		case "box_score_text":
			gameUrl := td.ChildAttr("a", "href")
			schedule.GameId = strings.Replace(strings.Split(td.ChildAttr("a", "href"), "/")[2], ".html", "", 1)
			gameUrls[schedule.GameId] = gameUrl
		case "game_start_time":
			parsedTime = strings.Replace(td.Text, "p", " PM EST", 1)
		}
	}
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
