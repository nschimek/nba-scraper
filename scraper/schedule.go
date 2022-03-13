package scraper

import (
	"errors"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

const (
	basePath         = "leagues"
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
	child       Scraper
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
		s.parseTable(tbl)
	})

	for _, url := range s.urls {
		s.colly.Visit(url)
	}

	scrapeChild(s)
}

func (s *ScheduleScraper) parseTable(tbl *colly.HTMLElement) {
	tbl.ForEach("tr", func(_ int, tr *colly.HTMLElement) {
		row := parseRow(tr)
		schedule := row.toSchedule()

		if schedule.Played && schedule.StartTime.After(s.dateRange.startDate) && schedule.StartTime.Before(s.dateRange.endDate) {
			s.ScrapedData = append(s.ScrapedData, schedule)
			s.childUrls[schedule.GameId] = row.gameUrl
		}
	})
}

type ScheduleRow struct {
	date, time, visitorUrl, homeUrl, gameUrl string
}

func parseRow(tr *colly.HTMLElement) (sr ScheduleRow) {
	if tr.Attr("class") != "thead" {
		sr.date = tr.ChildText("th a")

		tr.ForEach("td", func(_ int, td *colly.HTMLElement) {
			sr.parseColumn(td)
		})
	}
	return
}

func (sr *ScheduleRow) parseColumn(td *colly.HTMLElement) {
	switch td.Attr("data-stat") {
	case "game_start_time":
		sr.time = strings.Replace(td.Text, "p", " PM EST", 1)
	case "visitor_team_name":
		sr.visitorUrl = td.ChildAttr("a", "href")
	case "home_team_name":
		sr.homeUrl = td.ChildAttr("a", "href")
	case "box_score_text":
		sr.gameUrl = td.ChildAttr("a", "href")
	}
}

func (sr *ScheduleRow) toSchedule() (schedule Schedule) {
	schedule.HomeTeamId = parseTeamId(sr.homeUrl)
	schedule.VisitorTeamId = parseTeamId(sr.visitorUrl)
	if sr.gameUrl != "" {
		schedule.GameId = parseGameId(sr.gameUrl)
		schedule.Played = true
	}
	schedule.StartTime, _ = time.ParseInLocation("Mon, Jan 2, 2006 3:04 PM EST", sr.date+" "+sr.time, EST)
	return
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

func getMonthUrl(month time.Month, season string) string {
	monthString := strings.ToLower(month.String())
	return BaseHttp + "/" + basePath + "/NBA_" + season + "_games-" + monthString + ".html"
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

func parseGameId(link string) string {
	return strings.Replace(strings.Split(link, "/")[2], ".html", "", 1)
}

func parseTeamId(link string) string {
	return strings.Split(link, "/")[2]
}
