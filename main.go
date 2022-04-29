package main

import (
	"github.com/gocolly/colly/v2"

	"github.com/nschimek/nba-scraper/scraper"
)

func main() {
	c := colly.NewCollector(colly.AllowedDomains(scraper.AllowedDomain))
	c.Limit(&scraper.LimitRule)

	// startDate, _ := time.Parse("2006-01-02", "2021-10-20")
	// endDate, _ := time.Parse("2006-01-02", "2021-10-25")

	// scheduleScraper, _ := scraper.CreateScheduleScraperWithDates(c, "2022", startDate, endDate)

	// scheduleScraper.Scrape()
	// fmt.Println(scheduleScraper.GetData())
	// fmt.Println(scheduleScraper.GetChildUrls())

	// gameScraper := scraper.CreateGameScraper(c)
	// gameScraper.Scrape("https://www.basketball-reference.com/boxscores/202110300WAS.html", "https://www.basketball-reference.com/boxscores/202204180GSW.html")

	// teamScraper := scraper.CreateTeamScraper(c)
	// teamScraper.Scrape("https://www.basketball-reference.com/teams/TOR/2022.html")

	// playerScraper := scraper.CreatePlayerScraper(c)
	// playerScraper.Scrape("https://www.basketball-reference.com/players/v/vandeja01.html")

	standingsScraper := scraper.CreateStandingsScraper(c, 2022)
	standingsScraper.Scrape()
}
